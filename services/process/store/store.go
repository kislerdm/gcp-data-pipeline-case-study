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

Package defines the service store interface by wrapping around the logic of cloud.google.com/go/datastore.
*/

package store

import (
	"context"
	"platform/process/models"
	"time"

	"cloud.google.com/go/datastore"
)

var bg = context.Background()

const timeoutDefault = 20 * time.Second

// Config contains configuration for the Datastore query runner.
type Config struct {
	timeout time.Duration
}

// NewConfig return configuration for the Datastore client.
func NewConfig() *Config {
	return &Config{timeout: timeoutDefault}
}

// WithTimeout sets the operation timeout.
func (c *Config) WithTimeout(t time.Duration) *Config {
	c.timeout = t
	return c
}

// Client defines the client to interact with datastore
type Client struct {
	c   *datastore.Client
	cfg *Config
}

// NewClient init a new datastore client.
func NewClient(projectID string) (*Client, error) {
	c, err := datastore.NewClient(bg, projectID)
	return &Client{c: c, cfg: NewConfig()}, err
}

// WithCfg configures the client.
func (c *Client) WithCfg(cfg *Config) *Client {
	c.cfg = cfg
	return c
}

// Write writes object to the store.
func (c *Client) Write(collection string, obj interface{}) error {
	ctx, cancel := context.WithTimeout(bg, c.cfg.timeout)
	defer cancel()
	newKey := datastore.IncompleteKey(collection, nil)
	_, err := c.c.Put(ctx, newKey, obj)
	return err
}

// Read read object(s) from the store according to the query.
// The method requires the collection and the query objects to identify the output
// If collection is left empty, the query runs across all collections
// - limit defines the number of results to be returned
// - offset defines how many query results to be jumped over
// - out is the pointer to the object expected to be returned from db
func (c *Client) Read(collection string, query *models.Query, limit, offset int, out interface{}) error {
	ctx, cancel := context.WithTimeout(bg, c.cfg.timeout)
	defer cancel()
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	q := datastore.NewQuery(collection).Offset(offset).Limit(limit)
	if query != nil {
		if query.PayloadTimestamp != nil {
			if query.PayloadTimestamp.Min != nil {
				q = q.Filter("Payload.Time >=", query.PayloadTimestamp.Min)
			}
			if query.PayloadTimestamp.Max != nil {
				q = q.Filter("Payload.Time <=", query.PayloadTimestamp.Max)
			}
		} else if query.PayloadMean != nil {
			if query.PayloadMean.Min != nil {
				q = q.Filter("Payload.Mean >=", query.PayloadMean.Min)
			}
			if query.PayloadMean.Max != nil {
				q = q.Filter("Payload.Mean <=", query.PayloadMean.Max)
			}
		} else if query.PayloadStddev != nil {
			if query.PayloadStddev.Min != nil {
				q = q.Filter("Payload.Stddev >=", query.PayloadStddev.Min)
			}
			if query.PayloadStddev.Max != nil {
				q = q.Filter("Payload.Stddev <=", query.PayloadStddev.Max)
			}
		}
	}
	_, err := c.c.GetAll(ctx, q, out)
	return err
}
