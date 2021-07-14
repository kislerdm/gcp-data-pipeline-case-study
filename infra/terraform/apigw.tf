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

resource "google_api_gateway_api" "_" {
  provider     = google-beta
  project      = local.project
  api_id       = "api"
  display_name = "api"
  labels = {
    layer = "gateway"
  }
}

resource "google_service_account" "apigw" {
  project      = local.project
  account_id   = "gateway"
  display_name = "gateway"
  description  = "Account for the API GW to access backend services."
}

resource "google_cloud_run_service_iam_member" "cloudrun" {
  project  = local.project
  for_each = toset([google_cloud_run_service.submit.name, google_cloud_run_service.process.name])
  location = local.region
  service  = each.key
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.apigw.email}"
}

resource "google_api_gateway_api_config" "_" {
  provider      = google-beta
  project       = local.project
  api           = google_api_gateway_api._.api_id
  api_config_id = "cfg"

  openapi_documents {
    document {
      path     = "spec.yaml"
      contents = base64encode(local.api_config)
    }
  }

  gateway_config {
    backend_config {
      google_service_account = google_service_account.apigw.email
    }
  }

  lifecycle {
    create_before_destroy = false
  }

  labels = {
    layer = "gateway"
  }

  depends_on = [google_service_account.apigw]
}

resource "google_api_gateway_gateway" "_" {
  project    = local.project
  region     = local.region
  provider   = google-beta
  api_config = google_api_gateway_api_config._.id
  gateway_id = "api-gw"

  labels = {
    layer = "gateway"
  }
}
