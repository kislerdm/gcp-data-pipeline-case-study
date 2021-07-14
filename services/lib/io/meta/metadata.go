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

Package defines the logic to identify GCP project ID.
*/

package meta

import (
	"io"
	"net/http"
	"os"
	"platform/lib/io/fs"

	"github.com/goccy/go-json"
)

const url = "http://metadata.google.internal/computeMetadata/v1/project/project-id"

func callMetadataServer() string {
	c := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	req.Header.Add("Metadata-Flavor", "Google")
	res, err := c.Do(req)
	if err != nil {
		return ""
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	return string(b)
}

type key struct {
	ProjectID string `json:"project_id"`
}

func mustDeserializeKey(data []byte) (k *key) {
	json.Unmarshal(data, &k)
	return
}

func readKeyFromDisk(p string) *key {
	f, err := fs.FRead(p)
	if err != nil {
		return nil
	}
	return mustDeserializeKey(f)
}

// GetProjectID fetches the project ID
// The order to check:
// 1. secret key specified in GOOGLE_APPLICATION_CREDENTIALS envvar
// 2. GOOGLE_CLOUD_PROJECT envvar
// 3. GCLOUD_PROJECT envvar
// 4. Check metadata server
func GetProjectID() string {
	keyPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if keyPath != "" {
		key := readKeyFromDisk(keyPath)
		if key != nil {
			return key.ProjectID
		}
	}
	if o := os.Getenv("GOOGLE_CLOUD_PROJECT"); o != "" {
		return o
	}
	if o := os.Getenv("GCLOUD_PROJECT"); o != "" {
		return o
	}
	return callMetadataServer()
}
