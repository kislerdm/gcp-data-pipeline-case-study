# Copyright 2021 Dmitry Kisler dkisler.com

# Licensed under the Apache License,Version 2.0 (the "License");
# you may not use this file except in compliance with the License. You may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
# INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE,
# AND NONINFRINGEMENT. IN NO EVENT WILL THE LICENSOR OR OTHER CONTRIBUTORS BE LIABLE FOR ANY CLAIM, DAMAGES,
# OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF,
# OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

# See the License for the specific language governing permissions and limitations under the License.

locals {
  topics_process = toset([for k, v in local.topics : k if v.groud == "process"])
}

resource "google_service_account" "process" {
  project      = local.project
  account_id   = "processor"
  display_name = "processor"
  description  = "Account for the processing service."
}

resource "google_pubsub_topic_iam_member" "process" {
  for_each = local.topics_process
  project  = local.project
  topic    = each.key
  role     = "roles/pubsub.publisher"
  member   = "serviceAccount:${google_service_account.process.email}"
}

resource "google_storage_bucket_iam_member" "process" {
  bucket = google_storage_bucket.data.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.process.email}"
}

resource "google_project_iam_member" "process" {
  project = local.project
  role    = "roles/datastore.user"
  member  = "serviceAccount:${google_service_account.process.email}"
}

locals {
  path_code_services_process = "${local.path_code_services}/process"
  hash_codeabse_process      = [for f in fileset(local.path_code_services_process, "**") : filebase64sha256("${local.path_code_services_process}/${f}")]
  hash_process               = sha256(join(",", flatten([local.hash_codeabse_libs, local.hash_codeabse_process])))
  image_tag_process          = substr(local.hash_process, 0, 5)
}

resource "null_resource" "process_docker_rebuild" {
  triggers = {
    trigger = local.image_tag_process
  }
  provisioner "local-exec" {
    command = "cd ${path.module}/../.. && make service.rebuild SERVICE_NAME=process PROJECT_ID=${local.project} IMAGE_TAG=${local.image_tag_process} && cd ${path.module}"
  }
}

resource "google_cloud_run_service" "process" {
  project  = local.project
  location = local.region
  name     = "process"

  template {
    spec {
      containers {
        resources {
          limits = {
            cpu    = "1"
            memory = "1Gi"
          }
        }
        image = "eu.gcr.io/${local.project}/process:${local.image_tag_process}"
        ports {
          container_port = 9000
        }

        env {
          name  = "COLD_STORAGE_BUCKET"
          value = google_storage_bucket.data.name
        }
        env {
          name  = "NOTIFICATION_TOPIC"
          value = google_pubsub_topic._["process"].name
        }
        env {
          name  = "NOTIFICATION_TOPIC_FAIL"
          value = google_pubsub_topic._["process-fail"].name
        }
      }
      container_concurrency = 20
      timeout_seconds       = 30
      service_account_name  = google_service_account.process.email
    }
    metadata {
      annotations = {
        "autoscaling.knative.dev/minScale" = "1"
        "autoscaling.knative.dev/maxScale" = "500"
        "run.googleapis.com/client-name"   = "process"
      }
      labels = {
        layer = "ingress"
        type  = "interface"
      }
      namespace = "process"
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
  autogenerate_revision_name = true

  depends_on = [
    google_service_account.process,
    null_resource.process_docker_rebuild,
  ]
}

resource "google_service_account" "trigger_process" {
  project      = local.project
  account_id   = "trigger-processing"
  display_name = "trigger-processing"
  description  = "Account to invoke data processing service."
}

resource "google_cloud_run_service_iam_member" "trigger_process" {
  project  = local.project
  location = local.region
  service  = google_cloud_run_service.process.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.trigger_process.email}"
}

resource "google_pubsub_subscription" "process" {
  project = local.project
  name    = "trigger-data_processing"
  topic   = google_pubsub_topic._["submission"].name
  push_config {
    push_endpoint = "${google_cloud_run_service.process.status[0].url}/"
    oidc_token {
      service_account_email = google_service_account.trigger_process.email
    }
  }

  ack_deadline_seconds       = 600
  message_retention_duration = "600s"
  retain_acked_messages      = false

  expiration_policy {
    ttl = ""
  }

  retry_policy {
    minimum_backoff = "600s"
    maximum_backoff = "600s"
  }

  depends_on = [
    google_cloud_run_service.process,
    google_pubsub_topic._,
    google_service_account.trigger_process,
  ]
}
