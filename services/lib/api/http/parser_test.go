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

package http_test

import (
	"platform/lib/api/http"
	"reflect"
	"testing"
)

func TestNewRouteElement(t *testing.T) {
	tests := []struct {
		in   string
		want *http.Route
	}{
		{
			in: "/foo/{:id}",
			want: &http.Route{
				Source: "/foo/{:id}",
				Length: 2,
				Keys:   map[int]string{0: "foo"},
				Query: &http.PositionalQueryElements{
					http.InitPositionalQueryElement(1, "{:id}"),
				},
			},
		},
		{
			in: "/foo/bar/{:id}/zip/{:guid}",
			want: &http.Route{
				Source: "/foo/bar/{:id}/zip/{:guid}",
				Length: 5,
				Keys: map[int]string{
					0: "foo",
					1: "bar",
					3: "zip",
				},
				Query: &http.PositionalQueryElements{
					http.InitPositionalQueryElement(2, "{:id}"),
					http.InitPositionalQueryElement(4, "{:guid}"),
				},
			},
		},
	}
	for _, test := range tests {
		got := http.NewRouteElement(test.in)
		if !reflect.DeepEqual(got, test.want) {
			t.Fatalf("error!\nwant: %v\ngot: %v", test.want, got)
		}
	}
}

func TestParseRaw(t *testing.T) {
	r := http.Routes{
		http.NewRouteElement("/test"),
		http.NewRouteElement("/foo/{:id}"),
		http.NewRouteElement("/foo/bar/{:id}/zip/{:guid}"),
	}
	tests := []struct {
		in   string
		want *http.Route
	}{
		{
			in: "/foo/1",
			want: &http.Route{
				Source: "/foo/{:id}",
				Length: 2,
				Keys:   map[int]string{0: "foo"},
				Query: &http.PositionalQueryElements{
					&http.PositionalQueryElement{
						Index: 1,
						Key:   "id",
						Value: "1",
					},
				},
			},
		},
		{
			in:   "/",
			want: nil,
		},
		{
			in:   "/a/b/c/d/e/f/g/f",
			want: nil,
		},
	}
	for _, test := range tests {
		got := r.ParseRaw(test.in)
		if !reflect.DeepEqual(got, test.want) {
			t.Fatalf("error!\nwant: %v\ngot: %v", test.want, got)
		}
	}
}

func TestGetPositionalQuery(t *testing.T) {
	r := http.Route{
		Source: "/foo/{:id}",
		Length: 2,
		Keys:   map[int]string{0: "foo"},
		Query: &http.PositionalQueryElements{
			&http.PositionalQueryElement{
				Index: 1,
				Key:   "id",
				Value: "1",
			},
		},
	}
	want := map[string]string{"id": "1"}
	got := r.GetPositionalQuery()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("error!\nwant: %v\ngot: %v", want, got)
	}
}
