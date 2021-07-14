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
	"log"
	httpStatus "net/http"
	"platform/lib/api/http"
	"platform/lib/utils"
	"platform/process/models"
	"time"
)

const (
	hotStorageCollection = "processed"
)

func defaultReturn() (*http.Response, error) {
	return http.NewResponse([]byte{}, httpStatus.StatusOK), nil
}

func process(runner *runner) http.Action {
	return func(r *http.Request) (*http.Response, error) {
		locationDataRaw, err := models.DeserializePayloadLocation(r.Body)
		if err != nil {
			runner.Fail.Push([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
			return defaultReturn()
		}
		n := models.Notification{
			SubmitterID:  locationDataRaw.SubmitterID,
			SubmissionID: locationDataRaw.SubmissionID,
			Error:        "",
		}
		sendFail := func(err error) {
			n.Error = err.Error()
			runner.Fail.Push(n.MustSerialize())
		}

		if err != nil {
			sendFail(err)
			return defaultReturn()
		}

		data, err := runner.ColdStorage.Read(locationDataRaw.Bucket, locationDataRaw.Obj)
		if err != nil {
			sendFail(err)
			return defaultReturn()
		}
		inpt, err := models.DeserializeInput(data)
		if err != nil {
			sendFail(err)
			return defaultReturn()
		}
		o := inpt.Transform()
		o.SubmitterID = locationDataRaw.SubmitterID
		o.SubmissionID = locationDataRaw.SubmissionID
		o.TransformationEpoch = time.Now().Unix()

		if err := runner.HotStorage.Write(hotStorageCollection, o); err != nil {
			log.Println(err)
			sendFail(err)
			return defaultReturn()
		}

		if _, err := runner.Success.Push(o.MustSerialize()); err != nil {
			log.Println(err)
		}
		return defaultReturn()
	}
}

func query(runner *runner) http.Action {
	return func(r *http.Request) (*http.Response, error) {
		var q *models.Query
		err := models.ValidateQuery(r.Body)
		if err != nil {
			return http.NewResponse(
				[]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())),
				httpStatus.StatusBadRequest,
			), err
		}
		q = models.DeserializeQuery(r.Body)
		l := utils.MustAtoi(r.Query["limit"])
		offset := utils.MustAtoi(r.Query["offset"])
		var qRes models.QueryResults
		if err := runner.HotStorage.Read(hotStorageCollection, q, l, offset, &qRes); err != nil {
			return nil, err
		}
		out := qRes.Transform()
		return http.NewResponse(out.MustSerialize(), httpStatus.StatusOK), nil
	}
}

func fetch(runner *runner) http.Action {
	return func(r *http.Request) (*http.Response, error) {
		l := utils.MustAtoi(r.Query["limit"])
		offset := utils.MustAtoi(r.Query["offset"])
		var qRes models.QueryResults
		if err := runner.HotStorage.Read(hotStorageCollection, nil, l, offset, &qRes); err != nil {
			return nil, err
		}
		out := qRes.Transform()
		return http.NewResponse(out.MustSerialize(), httpStatus.StatusOK), nil
	}
}
