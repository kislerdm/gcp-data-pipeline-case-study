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

module "schema_error" {
  source          = "./modules/jsonschema_openapi"
  file_jsonschema = "${local.path_code_services}/models/error.json"
}

module "schema_submission_data_req" {
  source          = "./modules/jsonschema_openapi"
  file_jsonschema = "${local.path_code_services_submit}/models/request.json"
}

module "schema_submit_post_resp_ok" {
  source          = "./modules/jsonschema_openapi"
  file_jsonschema = "${local.path_code_services_submit}/models/response_ok.json"
}

module "schema_submit_post_resp_fail" {
  source          = "./modules/jsonschema_openapi"
  file_jsonschema = "${local.path_code_services_submit}/models/response_fail.json"
}

module "schema_process_resp" {
  source          = "./modules/jsonschema_openapi"
  file_jsonschema = "${local.path_code_services_process}/models/response.json"
}

module "schema_process_req_query" {
  source          = "./modules/jsonschema_openapi"
  file_jsonschema = "${local.path_code_services_process}/models/request_query_openapi2.json"
}

locals {
  api_config_template = templatefile("${path.module}/openapi.yaml",
    {
      # gw backends
      submit_service_url  = google_cloud_run_service.submit.status[0].url,
      process_service_url = google_cloud_run_service.process.status[0].url,
      # schema definitions keys mapping
      error                = "error"
      submission_data_req  = "submission_data_req"
      submission_resp_ok   = "submission_resp_ok"
      submission_resp_fail = "submission_resp_fail"
      process_resp         = "process_resp"
      process_query_req    = "process_query_req"
    },
  )
  api_config_obj_base = yamldecode(local.api_config_template)
  api_config_obj = merge(local.api_config_obj_base,
    {
      definitions = {
        error                = module.schema_error.obj
        submission_data_req  = module.schema_submission_data_req.obj
        submission_resp_ok   = module.schema_submit_post_resp_ok.obj
        submission_resp_fail = module.schema_submit_post_resp_fail.obj
        process_resp         = module.schema_process_resp.obj
        process_query_req    = module.schema_process_req_query.obj
      }
  })
  api_config = yamlencode(local.api_config_obj)
}
