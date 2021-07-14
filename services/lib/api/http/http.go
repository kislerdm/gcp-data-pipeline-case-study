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

Package defines the web-server (HTTP interface).
It wraps around github.com/valyala/fasthttp for base logic plus defines a custom routing algorithm.
*/

package http

import (
	"errors"
	"fmt"
	"log"
	"os"

	http "github.com/valyala/fasthttp"
)

var logger = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds|log.Lmsgprefix|log.LUTC|log.Llongfile)

// Server defines the API interface.
type Server struct {
	// HTTP server
	http *http.Server
}

// Start starts the HTTP service listening on the given port.
func (s *Server) Start(port interface{}) error {
	return s.http.ListenAndServe(fmt.Sprintf(":%v", port))
}

// NewServerDummy instantiates new HTTP server for test purposes.
func NewServerDummy() *Server {
	cfg := NewConfig()
	return &Server{
		&http.Server{
			Logger:               logger,
			ReadTimeout:          cfg.ReadTimeout,
			WriteTimeout:         cfg.WriteTimeout,
			MaxRequestBodySize:   cfg.MaxRequestBodySize,
			ReadBufferSize:       cfg.ReaderBufferSize,
			WriteBufferSize:      cfg.WriterBufferSize,
			Concurrency:          cfg.Concurrency,
			NoDefaultContentType: true,
			Handler:              NewRequestHandlersDummy().Router(),
		},
	}
}

// NewServer instantiates new HTTP server.
func NewServer(handlers *Handlers) *Server {
	s := NewServerDummy()
	s.http.Handler = handlers.Router()
	return s
}

// WithConfig adds custom configuration to the server.
func (s *Server) WithConfig(cfg *Config) *Server {
	s.http.ReadTimeout = cfg.ReadTimeout
	s.http.WriteTimeout = cfg.WriteTimeout
	s.http.MaxRequestBodySize = cfg.MaxRequestBodySize
	s.http.ReadBufferSize = cfg.ReaderBufferSize
	s.http.WriteBufferSize = cfg.WriterBufferSize
	s.http.Concurrency = cfg.Concurrency
	return s
}

// SetName sets the server name header
func (s *Server) SetName(n string) {
	s.http.Name = n
}

const healthcheckEndpointRoute = "/healthcheck"

// Handlers defines the request handlers for the server endpoint(s).
// It includes the endpoint to resolve GET requests to the "/healthcheck" route to check the services status.
type Handlers struct {
	// Endpoints defines the map of endpoint handlers.
	Endpoints map[string]*HandlerEndpoint
	// DefaultHeaders define the handlers to be append to all responses for all endpoints.
	// It may be useful for CORS settings.
	DefaultHeaders map[string]string
	// Routing defines the routes mapping
	Routing *Routes
}

// NewRequestHandlersDummy defines dummy endpoints handler.
// It includes only "/healthcheck" endpoint for test purposes.
func NewRequestHandlersDummy() *Handlers {
	return &Handlers{
		Endpoints:      map[string]*HandlerEndpoint{healthcheckEndpointRoute: HealthcheckHandler},
		DefaultHeaders: map[string]string{},
	}
}

// NewRequestHandlers defines the handlers for all endpoints.
func NewRequestHandlers(endpoints map[string]*HandlerEndpoint) *Handlers {
	if _, ok := endpoints[healthcheckEndpointRoute]; !ok {
		endpoints[healthcheckEndpointRoute] = HealthcheckHandler
	}
	r := Routes{}
	for route := range endpoints {
		r = append(r, NewRouteElement(route))
	}
	return &Handlers{
		Endpoints:      endpoints,
		DefaultHeaders: map[string]string{},
		Routing:        &r,
	}
}

// WithDefaultHeaders add default response headers to the endpoints handlers.
func (h *Handlers) WithDefaultHeaders(headers map[string]string) *Handlers {
	h.DefaultHeaders = headers
	return h
}

// WithoutHealthCheck deactivates the default healthcheck endpoint.
func (h *Handlers) WithoutHealthCheck() *Handlers {
	delete(h.Endpoints, healthcheckEndpointRoute)
	return h
}

// AddEndpointHandler add the handler to resolve requests to the route endpoint.
func (h *Handlers) AddEndpointHandler(route string, handler *HandlerEndpoint) {
	h.Endpoints[route] = handler
}

func (h *Handlers) addDefaultHeaders(ctx *http.RequestCtx) {
	for k, v := range h.DefaultHeaders {
		ctx.Response.Header.Add(k, v)
	}
}

func (h *Handlers) reply(ctx *http.RequestCtx, actionResp *Response) {
	ctx.SetBody(actionResp.Body)
	ctx.SetStatusCode(actionResp.StatusCode)
	ctx.SetContentType(actionResp.ContentType)
	for k, v := range actionResp.Headers {
		ctx.Response.Header.Add(k, v)
	}
}

func (h *Handlers) router(ctx *http.RequestCtx) {
	h.addDefaultHeaders(ctx)
	p := string(ctx.Path())
	hdlr, ok := h.Endpoints[p]

	missingRoute := func() {
		e := "unsupported path"
		logger.Println(e)
		h.reply(ctx, NewResponse([]byte(fmt.Sprintf(`{"error": "%s"}`, e)), http.StatusNotFound))
	}

	routeParameters := map[string]string{}

	if !ok {
		r := h.Routing.ParseRaw(p)
		if r == nil {
			missingRoute()
			return
		}
		hdlr, ok = h.Endpoints[r.Source]
		if !ok {
			missingRoute()
			return
		}
		routeParameters = r.GetPositionalQuery()
	}

	method := string(ctx.Method())
	if !hdlr.IsAllowedMethod(method) {
		e := "unsupported method"
		logger.Println(e)
		h.reply(ctx, NewResponse([]byte(fmt.Sprintf(`{"error": "%s"}`, e)), http.StatusMethodNotAllowed))
		return
	}
	resp, err := hdlr.ProcessRequest(&Request{
		Method:          method,
		RouteParameters: routeParameters,
		Query:           ParseRequestKV(ctx.QueryArgs().VisitAll),
		Headers:         ParseRequestKV(ctx.Request.Header.VisitAll),
		Body:            ctx.Request.Body(),
	})
	if err != nil {
		h.reply(ctx,
			NewResponse(
				[]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())),
				http.StatusInternalServerError,
			),
		)
	}
	h.reply(ctx, resp)
}

// Router defines the request routing to corresponding handler.
func (h *Handlers) Router() http.RequestHandler {
	return func(ctx *http.RequestCtx) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("unknown error")
				}
				logger.Println(err)
			}
		}()
		h.router(ctx)
	}
}

type parser struct {
	Content map[string]string
}

func (h *parser) Parse(key, value []byte) {
	h.Content[string(key)] = string(value)
}

func ParseRequestKV(f func(func(key, value []byte))) map[string]string {
	o := parser{
		Content: map[string]string{},
	}
	f(o.Parse)
	return o.Content
}
