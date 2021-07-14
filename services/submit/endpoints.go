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
	"fmt"
	httpStatus "net/http"
	"path"
	"platform/lib/api/http"
	"platform/submit/models"
	"sync"
)

var wg sync.WaitGroup

const submitterID = "test"

func submit(runner *runner, bucket *string) http.Action {
	return func(r *http.Request) (*http.Response, error) {
		errOut := []string{}

		validErrs := models.ValidatePayload(r.Body)
		if validErrs != nil {
			errOut = append(errOut, validErrs.Error())
		}

		payloadToDispatch := NewPayloadHotStorage(submitterID, r.Body, validErrs == nil)

		keyColdStorage := path.Join(
			payloadToDispatch.SubmitterID,
			payloadToDispatch.SubmissionID,
			fmt.Sprintf("%s.json", payloadToDispatch.SubmissionID),
		)

		func() {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := runner.ColdStorage.Write(*bucket, keyColdStorage, r.Body); err != nil {
					errOut = append(errOut, err.Error())
				}
			}()
			wg.Add(1)
			go func() {
				defer wg.Done()
				notification := payloadLocation{
					SubmitterID:  payloadToDispatch.SubmitterID,
					SubmissionID: payloadToDispatch.SubmissionID,
					Bucket:       *bucket,
					Obj:          keyColdStorage,
				}
				var err error
				if payloadToDispatch.Valid {
					_, err = runner.Success.Push(notification.MustSerialize())
				} else {
					_, err = runner.Fail.Push(notification.MustSerialize())
				}
				if err != nil {
					errOut = append(errOut, err.Error())
				}
			}()
		}()
		wg.Wait()

		resp := &response{
			SubmissionID: payloadToDispatch.SubmissionID,
			Errors:       errOut,
		}
		status := httpStatus.StatusOK
		if len(resp.Errors) > 0 {
			status = httpStatus.StatusInternalServerError
			if !payloadToDispatch.Valid {
				status = httpStatus.StatusBadRequest
			}
		}
		return http.NewResponse(resp.MustSerialize(), status), nil
	}
}

func read(runner *runner, bucket *string) http.Action {
	return func(r *http.Request) (*http.Response, error) {
		submissionID, ok := r.RouteParameters["submission_id"]
		if !ok {
			return http.NewResponse([]byte(`{"error": "missing submission_id"}`), httpStatus.StatusBadRequest), nil
		}
		keyColdStorage := path.Join(submitterID, submissionID, fmt.Sprintf("%s.json", submissionID))
		data, err := runner.ColdStorage.Read(*bucket, keyColdStorage)
		if err != nil {
			if err.Error() == "storage: object doesn't exist" {
				return http.NewResponse([]byte(`{"error": "data not found"}`), httpStatus.StatusNotFound), nil
			}
			return nil, err
		}
		return http.NewResponse(data, httpStatus.StatusOK), nil
	}
}
