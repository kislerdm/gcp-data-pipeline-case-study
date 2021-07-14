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
	"time"
)

// Config defines the HTTP client configuration.
type Config struct {
	// Timeout to get response
	ReadTimeout time.Duration
	// Timeout to send request
	WriteTimeout time.Duration
	// MaxRequestBodySize sets the max body size in bytes
	MaxRequestBodySize int
	// MaxHeaderSize sets the max reader buffer size in bytes
	ReaderBufferSize int
	// MaxWriterSize sets the max writer buffer size in bytes
	WriterBufferSize int
	// Concurrency max concurent connections to be served by the server
	Concurrency int
}

const (
	defaultTimeoutRead  = 60 * time.Second
	defaultTimeoutWrite = 600 * time.Second
	defaultSize         = 10 * 1 << 22
	defaultBufferSize   = 4 * 1 << 10
	defaultConcurrency  = 1e6
)

// NewConfig initiates HTTP server config.
//
// Default settings:
//
// timeout for read: 60 sec
//
// timeout for write: 600 sec
//
// max reader body size: 10 Mb
//
// reader buffer size: 10 Mb
//
// writer buffer size: 10 Mb
//
// max connections concurrency: 1 million connections
func NewConfig() *Config {
	return &Config{
		ReadTimeout:        defaultTimeoutRead,
		WriteTimeout:       defaultTimeoutWrite,
		MaxRequestBodySize: defaultSize,
		ReaderBufferSize:   defaultBufferSize,
		WriterBufferSize:   defaultBufferSize,
		Concurrency:        defaultConcurrency,
	}
}

// SetReadTimeout sets the server read timeout.
func (c *Config) SetReadTimeout(t time.Duration) {
	c.ReadTimeout = t
}

// SetWriteTimeout sets the server write timeout.
func (c *Config) SetWriteTimeout(t time.Duration) {
	c.WriteTimeout = t
}

// SetMaxRequestBodySize sets the max request body size in bytes.
func (c *Config) SetMaxRequestBodySize(size int) {
	c.MaxRequestBodySize = size
}

// SetReaderBufferSize sets the reader buffer size in bytes.
func (c *Config) SetReaderBufferSize(size int) {
	c.ReaderBufferSize = size
}

// SetWriterBufferSize sets the reader buffer size in bytes.
func (c *Config) SetWriterBufferSize(size int) {
	c.WriterBufferSize = size
}

// SetMaxCCConnections sets the max number of cc connection the server could serve.
func (c *Config) SetMaxCCConnections(n int) {
	c.Concurrency = n
}
