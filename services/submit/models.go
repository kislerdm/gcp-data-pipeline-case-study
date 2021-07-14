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
*/

package main

import (
	"platform/lib/utils"
	"time"

	"github.com/goccy/go-json"
)

type PayloadHotstorage struct {
	SubmitterID     string `json:"submitter_id"`
	SubmissionID    string `json:"submission_id"`
	SubmissionEpoch int64  `json:"submission_epoch"`
	Valid           bool   `json:"valid"`
	Payload         []byte `json:"payload"`
}

func NewPayloadHotStorage(submitterID string, payload []byte, valid bool) *PayloadHotstorage {
	return &PayloadHotstorage{
		SubmitterID:     submitterID,
		SubmissionID:    utils.GenerateUUID4(),
		SubmissionEpoch: time.Now().Unix(),
		Valid:           valid,
		Payload:         payload,
	}
}

// payloadLocation defines the cold storage object location.
type payloadLocation struct {
	SubmitterID  string `json:"submitter_id"`
	SubmissionID string `json:"submission_id"`
	Bucket       string `json:"bucket"`
	Obj          string `json:"key"`
}

func (p *payloadLocation) MustSerialize() []byte {
	o, _ := json.Marshal(p)
	return o
}

type response struct {
	SubmissionID string   `json:"submission_id"`
	Errors       []string `json:"errors,omitempty"`
}

func (r *response) MustSerialize() []byte {
	o, _ := json.Marshal(r)
	return o
}
