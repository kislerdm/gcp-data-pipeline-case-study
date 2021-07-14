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

variable "jsonschema" {
  type        = string
  description = "JSON schema content."
  default     = ""
}

variable "file_jsonschema" {
  type        = string
  description = "Path to JSON schema file."
  default     = ""
}

locals {
  schema_str               = var.file_jsonschema == "" ? var.jsonschema : file(var.file_jsonschema)
  submit_post_schema_input = jsondecode(local.schema_str)
  submit_post_schema_openapi = {
    for k, v in local.submit_post_schema_input :
    k => v if !contains(["$id", "$schema", "additionalItems"], k)
  }
}

output "obj" {
  description = "JSON schema parsed to the object."
  value       = local.submit_post_schema_openapi
}

output "yaml" {
  description = "JSON schema parsed and encoded as YAML following open API standards."
  value       = yamlencode(local.submit_post_schema_openapi)
}
