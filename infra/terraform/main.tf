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

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "3.69.0"
    }
  }
  backend "gcs" {
    bucket = "data-case-dmitry-tf"
    prefix = "terraform/state/prod"
  }
}

provider "google" {
  project = "data-case-dmitry"
  region  = "europe-west1"
}

locals {
  project = "data-case-dmitry"
  region  = "europe-west1"
}

resource "google_project_service" "_" {
  project = local.project
  for_each = toset([
    "iam",
    "storage",
    "logging",
    "run",
    "servicemanagement",
    "serviceconsumermanagement",
    "servicecontrol",
    "cloudresourcemanager",
    "compute",
    "apigateway",
    "pubsub",
    "datastore",
  ])
  service                    = "${each.value}.googleapis.com"
  disable_dependent_services = true
}
