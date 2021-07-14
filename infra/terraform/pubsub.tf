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
  statuses = ["success", "fail"]
  groups = [
    {
      name  = "submission"
      layer = "ingress"
    },
    {
      name  = "process"
      layer = "processing"
    },
  ]
  topics_obj = flatten([
    for group in local.groups : [
      for status in local.statuses : {
        name  = status == "fail" ? "${group.name}-${status}" : group.name
        groud = group.name
        layer = group.layer
      }
    ]
  ])
  topics = { for t in local.topics_obj : t.name => t }
}

resource "google_pubsub_topic" "_" {
  for_each = local.topics
  project  = local.project
  name     = each.key
  labels = {
    layer = each.value.layer
  }
}
