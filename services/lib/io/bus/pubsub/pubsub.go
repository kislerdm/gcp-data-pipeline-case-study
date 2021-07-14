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

Package wraps the interface to GCP PubSub.
*/

package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
)

var ctx = context.Background()

// Client defines the PubSub client.
type Client struct {
	projectID string
	Handler   *pubsub.Client
}

// NewClient init a new PubSub client.
func NewClient(projectID string) (*Client, error) {
	c, err := pubsub.NewClient(ctx, projectID)
	return &Client{
		projectID: projectID,
		Handler:   c,
	}, err
}

// Config contains configuration for the PubSub publisher client.
type Config struct {
	pubsub.PublishSettings
}

// NewConfig return configuration for the PubSub client.
func NewConfig() *Config {
	return &Config{}
}

// WithCCLimit overwrites default number of goroutines to controll concurrency.
// If n is not positive, the number concurrency is set to 25 * runtime.GOMAXPROCS(0)
// See: cloud.google.com/go/pubsub/topic for details
func (c *Config) WithCCLimit(n int) *Config {
	if n < 0 {
		n = 0
	}
	c.NumGoroutines = n
	return c
}

// GetPublisher init a publisher for a topic.
func (c *Client) GetPublisher(topic string) *Publisher {
	if topic == "" {
		return nil
	}
	t := c.Handler.Topic(topic)
	return &Publisher{t}
}

// GetPublisherWithConfig init a publisher for a topic with configuration.
func (c *Client) GetPublisherWithConfig(topic string, cfg *Config) *Publisher {
	if topic == "" {
		return nil
	}
	t := c.Handler.Topic(topic)
	t.PublishSettings = cfg.PublishSettings
	return &Publisher{t}
}

// Publisher defines the client to publish to a topic.
type Publisher struct {
	t *pubsub.Topic
}

// WithCCLimit overwrites the CCLimit specified in the client cfg.
// If n is not positive, the number concurrency is set to 25 * runtime.GOMAXPROCS(0)
// See: cloud.google.com/go/pubsub/topic for details
func (p *Publisher) WithCCLimit(n int) *Publisher {
	if n < 0 {
		n = 0
	}
	p.t.PublishSettings.NumGoroutines = n
	return p
}

// Push pushes data to the topic.
func (p *Publisher) Push(data []byte) (id string, err error) {
	return p.t.Publish(ctx, &pubsub.Message{Data: data}).Get(ctx)
}
