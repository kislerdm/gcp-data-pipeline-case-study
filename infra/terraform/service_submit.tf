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
  topics_submit = toset([for k, v in local.topics : k if v.groud == "submission"])
}

resource "google_service_account" "submit" {
  project      = local.project
  account_id   = "submitter"
  display_name = "submitter"
  description  = "Account for the submission service."
}

resource "google_pubsub_topic_iam_member" "submit" {
  for_each = local.topics_submit
  project  = local.project
  topic    = each.key
  role     = "roles/pubsub.publisher"
  member   = "serviceAccount:${google_service_account.submit.email}"
}

resource "google_storage_bucket_iam_member" "submit" {
  bucket = google_storage_bucket.data.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.submit.email}"
}

locals {
  path_code_services        = "${path.module}/../../services"
  path_code_lib             = "${local.path_code_services}/lib"
  path_code_services_submit = "${local.path_code_services}/submit"

  hash_codeabse_libs = flatten([
    [for f in fileset(local.path_code_lib, "**") : filebase64sha256("${local.path_code_lib}/${f}")],
    ["${local.path_code_services}/Dockerfile", "${local.path_code_services}/.dockerignore"],
  ])

  hash_codeabse_submit = [for f in fileset(local.path_code_services_submit, "**") : filebase64sha256("${local.path_code_services_submit}/${f}")]

  hash_submit = sha256(join(",", flatten([local.hash_codeabse_libs, local.hash_codeabse_submit])))

  image_tag_submit = substr(local.hash_submit, 0, 5)
}

resource "null_resource" "submit_docker_rebuild" {
  triggers = {
    trigger = local.image_tag_submit
  }
  provisioner "local-exec" {
    command = "cd ${path.module}/../.. && make service.rebuild SERVICE_NAME=submit PROJECT_ID=${local.project} IMAGE_TAG=${local.image_tag_submit} && cd ${path.module}"
  }
}

resource "google_cloud_run_service" "submit" {
  project  = local.project
  location = local.region
  name     = "submit"

  template {
    spec {
      containers {
        resources {
          limits = {
            cpu    = "1"
            memory = "1Gi"
          }
        }
        image = "eu.gcr.io/${local.project}/submit:${local.image_tag_submit}"
        ports {
          container_port = 9000
        }

        env {
          name  = "COLD_STORAGE_BUCKET"
          value = google_storage_bucket.data.name
        }
        env {
          name  = "NOTIFICATION_TOPIC"
          value = google_pubsub_topic._["submission"].name
        }
        env {
          name  = "NOTIFICATION_TOPIC_FAIL"
          value = google_pubsub_topic._["submission-fail"].name
        }
      }
      container_concurrency = 20
      timeout_seconds       = 30
      service_account_name  = google_service_account.submit.email
    }
    metadata {
      annotations = {
        "autoscaling.knative.dev/minScale" = "1"
        "autoscaling.knative.dev/maxScale" = "500"
        "run.googleapis.com/client-name"   = "submit"
      }
      labels = {
        layer = "ingress"
        type  = "interface"
      }
      namespace = "submit"
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
  autogenerate_revision_name = true

  depends_on = [
    google_service_account.submit,
    null_resource.submit_docker_rebuild,
  ]
}
