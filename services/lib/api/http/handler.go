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

import http "github.com/valyala/fasthttp"

// Request defines the HandlerEndpoint's request object.
type Request struct {
	Method          string
	RouteParameters map[string]string
	Headers         map[string]string
	Query           map[string]string
	Body            []byte
}

// Response defines the HandlerEndpoint's response object.
type Response struct {
	// Body defines the response payload.
	Body []byte
	// ContentType defines the response MIME content type.
	ContentType string
	// Response status code
	StatusCode int
	// Headers defines the response headers.
	Headers map[string]string
}

const defaultContentType = "application/json"

// NewResponse defines the response object.
// The default content type application/json is being set.
// To set a custom content type, use the method SetContentType
func NewResponse(b []byte, statusCode int) *Response {
	return &Response{
		Body:        b,
		StatusCode:  statusCode,
		ContentType: defaultContentType,
	}
}

// SetContentType sets the custom response content type.
func (r *Response) SetContentType(s string) {
	r.ContentType = s
}

// SetHeader sets the custom response header.
func (r *Response) SetHeader(key, val string) {
	r.Headers[key] = val
}

// AddHeaders adds response headers.
func (r *Response) AddHeaders(h map[string]string) {
	r.Headers = h
}

// Action defines the requests HandlerEndpoint's action.
type Action func(r *Request) (*Response, error)

// ActionHealthcheck defines the function for the status healthcheck.
func ActionHealthcheck(r *Request) (*Response, error) {
	return NewResponse([]byte{}, http.StatusOK), nil
}

// AllowedMethods defines the allowed HTTP request methods.
type AllowedMethods []string

// IsIn checks if the HTTP method is allowed.
func (m *AllowedMethods) IsIn(t string) bool {
	for _, method := range *m {
		if t == method {
			return true
		}
	}
	return false
}

// Add adds a method to the list of allowed methods.
func (m *AllowedMethods) Add(method string) {
	*m = append(*m, method)
}

// HandlerEndpoint defines the logic to handle requests comming to an endpoint.
type HandlerEndpoint struct {
	// Action defines the action of the request HandlerEndpoint.
	Action Action
	// AllowedMethods defines the allowed HTTP method.
	AllowedMethods *AllowedMethods
	// ContentType MIME content type.
	ContentType string
}

// NewHandlerEndpoint initiates a new HandlerEndpoint.
// Example:
// submitEndpoint := map[string]*HandlerEndpoint{
// 	"submit": NewHandlerEndpoint(submitAction, []string{"POST"})
// }
func NewHandlerEndpoint(action Action, allowedMethods []string) *HandlerEndpoint {
	return &HandlerEndpoint{
		Action:         action,
		AllowedMethods: (*AllowedMethods)(&allowedMethods),
		ContentType:    defaultContentType,
	}
}

// ProcessRequest process the request incoming to the endpoint.
func (h *HandlerEndpoint) ProcessRequest(r *Request) (*Response, error) {
	return h.Action(r)
}

// WithOPTION sets the OPTION method as allowed.
func (h *HandlerEndpoint) WithOPTION() *HandlerEndpoint {
	h.AllowedMethods.Add("OPTIONS")
	return h
}

// SetAllowedMethods sets allowed .
func (h *HandlerEndpoint) SetAllowedMethods(methods []string) {
	h.AllowedMethods = (*AllowedMethods)(&methods)
}

// AddAllowedMethods sets the OPTION method as allowed.
func (h *HandlerEndpoint) AddAllowedMethods(methods []string) {
	for _, method := range methods {
		h.AllowedMethods.Add(method)
	}
}

// IsAllowedMethod checks if the request method is allowed.
func (h *HandlerEndpoint) IsAllowedMethod(method string) bool {
	return h.AllowedMethods.IsIn(method)
}

// HealthcheckHandler defines the handler for the status healthckeck endpoint.
var HealthcheckHandler = NewHandlerEndpoint(ActionHealthcheck, []string{"GET"})
