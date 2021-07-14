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
	_ "embed"
	"time"

	"platform/lib/jsonschema"

	"github.com/goccy/go-json"
)

type queryTimestamp struct {
	Min *time.Time `json:"min,omitempty"`
	Max *time.Time `json:"max,omitempty"`
}

type queryFloat struct {
	Min *float64 `json:"min,omitempty"`
	Max *float64 `json:"max,omitempty"`
}

// Query defines the query format.
type Query struct {
	PayloadTimestamp *queryTimestamp `json:"timestamp,omitempty"`
	PayloadMean      *queryFloat     `json:"mean,omitempty"`
	PayloadStddev    *queryFloat     `json:"standard_deviation,omitempty"`
}

// DeserializeQuery deserializes the data.
func DeserializeQuery(data []byte) (q *Query) {
	json.Unmarshal(data, &q)
	return
}

//go:embed request_query.json
var schema []byte

var s, _ = jsonschema.NewSchema(schema)

// ValidateQuery validates the object.
func ValidateQuery(data []byte) error {
	return s.ValidateBytes(data)
}
