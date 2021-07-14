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

package http

import (
	"regexp"
	"strings"
)

func split(s string) []string {
	return strings.Split(s[1:], "/")
}

var re = regexp.MustCompile("\\{:(.*?)\\}")

func parse(s string) string {
	o := re.FindStringSubmatch(s)
	if len(o) == 0 {
		return ""
	}
	return o[1]
}

type PositionalQueryElement struct {
	// Route element
	Index int
	// Query element
	Key, Value string
}

func InitPositionalQueryElement(i int, s string) *PositionalQueryElement {
	k := parse(s)
	if k == "" {
		return nil
	}
	return &PositionalQueryElement{Index: i, Key: k}
}

type PositionalQueryElements []*PositionalQueryElement

func (p *PositionalQueryElements) setValues(el []string) {
	for _, q := range *p {
		q.Value = el[q.Index]
	}
}

func (p *PositionalQueryElements) getPositionalQuery() map[string]string {
	o := map[string]string{}
	for _, e := range *p {
		o[e.Key] = e.Value
	}
	return o
}

type Route struct {
	// Original route
	Source string
	// Number of route elements
	Length int
	// Keys for lookup to identify correct route
	Keys map[int]string
	// Query route query parameters
	Query *PositionalQueryElements
}

// matchByKey checks is the request route matches the handlers pointers.
func (r *Route) matchByKey(el []string) bool {
	for i, k := range r.Keys {
		if el[i] != k {
			return false
		}
	}
	return true
}

func (r *Route) GetPositionalQuery() map[string]string {
	return r.Query.getPositionalQuery()
}

func NewRouteElement(s string) *Route {
	el := split(s)
	keys := map[int]string{}
	q := PositionalQueryElements{}
	for i, e := range el {
		if qK := InitPositionalQueryElement(i, e); qK != nil {
			q = append(q, qK)
		} else {
			keys[i] = e
		}
	}
	return &Route{
		Source: s,
		Length: len(el),
		Keys:   keys,
		Query:  &q,
	}
}

// Routes route elements to search for the handler.
type Routes []*Route

func (rs *Routes) filterByLength(l int) *Routes {
	o := Routes{}
	for _, r := range *rs {
		if r.Length == l {
			o = append(o, r)
		}
	}
	if len(o) == 0 {
		return nil
	}
	return &o
}

// ParseRaw parses the request path string.
func (r *Routes) ParseRaw(s string) *Route {
	el := split(s)
	candidates := r.filterByLength(len(el))
	if candidates == nil {
		return nil
	}
	for _, c := range *candidates {
		if c.matchByKey(el) {
			c.Query.setValues(el)
			return c
		}
	}
	return nil
}
