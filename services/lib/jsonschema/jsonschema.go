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

Package wraps around github.com/xeipuuv/gojsonschema.
*/

package jsonschema

import (
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// ValidationError defines JSON validation errors.
type ValidationError struct {
	o []gojsonschema.ResultError
}

func (e *ValidationError) Error() string {
	o := []string{}
	for _, i := range e.o {
		o = append(o, fmt.Sprintf("Field %s: %s", i.Field(), i.Description()))
	}
	return strings.Join(o, "\n")
}

// Schema defines the schema object.
type Schema struct {
	*gojsonschema.Schema
}

// NewSchema defines the Schema object.
func NewSchema(def []byte) (s *Schema, err error) {
	sch, err := gojsonschema.NewSchemaLoader().Compile(gojsonschema.NewBytesLoader(def))
	if err != nil {
		return nil, err
	}
	return &Schema{sch}, nil
}

// ValidateBytes validates bytes object.
func (s *Schema) ValidateBytes(data []byte) error {
	res, err := s.Validate(gojsonschema.NewBytesLoader(data))
	if err != nil {
		return err
	}
	if res.Valid() {
		return nil
	}
	return &ValidationError{res.Errors()}
}

// ValidateObject validates GO object.
func (s *Schema) ValidateObject(data interface{}) error {
	res, err := s.Validate(gojsonschema.NewGoLoader(data))
	if err != nil {
		return err
	}
	if res.Valid() {
		return nil
	}
	return &ValidationError{res.Errors()}
}
