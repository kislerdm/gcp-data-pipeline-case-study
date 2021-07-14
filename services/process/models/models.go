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

package models

import (
	"platform/process/transformation"
	"time"

	"github.com/goccy/go-json"
)

type payloadLocation struct {
	SubmitterID  string `json:"submitter_id"`
	SubmissionID string `json:"submission_id"`
	Bucket       string `json:"bucket"`
	Obj          string `json:"key"`
}

func DeserializePayloadLocation(data []byte) (p *payloadLocation, err error) {
	err = json.Unmarshal(data, &p)
	return
}

type Notification struct {
	SubmitterID  string `json:"submitter_id"`
	SubmissionID string `json:"submission_id"`
	Error        string `json:"error,omitempty"`
}

func (p *Notification) MustSerialize() []byte {
	o, _ := json.Marshal(p)
	return o
}

type input struct {
	Time time.Time `json:"time_stamp"`
	Data []float64 `json:"data"`
}

type payload struct {
	Time   time.Time `json:"timestamp"`
	Mean   float64   `json:"mean"`
	Stddev float64   `json:"standard_deviation"`
}

type outputProcessing struct {
	SubmitterID         string   `json:"submitter_id"`
	SubmissionID        string   `json:"submission_id"`
	TransformationEpoch int64    `json:"transformation_epoch"`
	Payload             *payload `json:"payload" datastore:",flatten"`
}

func (o *outputProcessing) MustSerialize() []byte {
	out, _ := json.Marshal(o)
	return out
}

func DeserializeInput(data []byte) (i *input, err error) {
	err = json.Unmarshal(data, &i)
	return
}

func (i *input) Transform() *outputProcessing {
	dstrStats := transformation.NewDistribution(i.Data).Stats()
	ts := transformation.ConvertTimestampUTC(i.Time)
	return &outputProcessing{
		Payload: &payload{
			Time:   ts,
			Mean:   dstrStats.Mean,
			Stddev: dstrStats.Stddev,
		},
	}
}

type QueryResults []outputProcessing

type outputElement struct {
	SubmissionID string   `json:"submission_id"`
	Payload      *payload `json:"payload"`
}

type output []*outputElement

func (i *QueryResults) Transform() *output {
	var out output
	for _, o := range *i {
		out = append(out, &outputElement{
			SubmissionID: o.SubmissionID,
			Payload:      o.Payload,
		})
	}
	return &out
}

func (o *output) MustSerialize() []byte {
	out, _ := json.Marshal(&o)
	return out
}
