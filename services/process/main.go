/* Copyright 2021 Dmitry Kisler dkisler.com

Licensed under the Apache License,Version 2.0 (the "License");
you may not use this file except in compliance with the License. You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE,
AND NONINFRINGEMENT. IN NO EVENT WILL THE LICENSOR OR OTHER CONTRIBUTORS BE LIABLE FOR ANY CLAIM, DAMAGES,
OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF,
OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

See the License for the specific language governing permissions and limitations under the License.

Package defines the service which provides the interface for the data submission into the data platform.

Modus operandi:

1. Transforms raw data sample:
	- Converts input timestamp timezone to UTC
	- Calculates Mean and Std dev of submitted data distribution.
2. Stores data to the service store (GCP Datastore).
3. Pushes notification message to the message bus (GCP PubSub).
*/

package main

import (
	"log"
	"platform/lib/api/http"
	"platform/lib/io/bus/pubsub"
	"platform/lib/io/meta"
	"platform/lib/io/store/gcs"
	"platform/lib/utils"
	"platform/process/store"
)

type runner struct {
	Success     *pubsub.Publisher
	Fail        *pubsub.Publisher
	ColdStorage *gcs.Client
	HotStorage  *store.Client
}

var (
	r *runner = &runner{}
	s *http.Server
)

func setServer() {
	endpoints := map[string]*http.HandlerEndpoint{
		"/":      http.NewHandlerEndpoint(process(r), []string{"POST"}),
		"/query": http.NewHandlerEndpoint(query(r), []string{"POST"}),
		"/fetch": http.NewHandlerEndpoint(fetch(r), []string{"GET"}),
	}
	handlers := http.NewRequestHandlers(endpoints).WithDefaultHeaders(
		map[string]string{
			"tag-layer":  "process",
			"tag-branch": "fast",
		})

	s = http.NewServer(handlers)
	s.SetName("process")
}

func init() {
	projectID := utils.GetEnv("GCP_PROJECT", "")
	if projectID == "" {
		projectID = meta.GetProjectID()
		if projectID == "" {
			log.Fatalln("specify gcp project, e.g. as envvar 'GCP_PROJECT'")
		}
	}
	topic := utils.GetEnv("NOTIFICATION_TOPIC", "")
	if topic == "" {
		log.Fatalln("specify the the message bus notification topic by setting envvar 'NOTIFICATION_TOPIC'")
	}
	topicFail := utils.GetEnv("NOTIFICATION_TOPIC_FAIL", "")
	if topic == "" {
		log.Fatalln("specify the the message bus notification for fail topic by setting envvar 'NOTIFICATION_TOPIC_FAIL'")
	}
	c, err := pubsub.NewClient(projectID)
	if err != nil {
		log.Fatalln(err)
	}
	r.Success = c.GetPublisher(topic).WithCCLimit(1)
	r.Fail = c.GetPublisher(topicFail).WithCCLimit(1)

	r.ColdStorage, err = gcs.NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	r.HotStorage, err = store.NewClient(projectID)
	if err != nil {
		log.Fatalln(err)
	}

	setServer()
}

func main() {
	s.Start(utils.GetEnv("PORT", "9000"))
}
