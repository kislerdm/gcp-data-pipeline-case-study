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

Package wraps the interface to GCP Storage.
*/

package gcs

import (
	"io"
	bg "platform/lib/io/context"

	"cloud.google.com/go/storage"
)

// Client defines the client to interact with bigquery
type Client struct {
	*storage.Client
}

// NewClient init a new Bigquery client.
func NewClient() (*Client, error) {
	c, err := storage.NewClient(bg.CtxBG)
	return &Client{c}, err
}

// Write writes object to the bucket.
func (c *Client) Write(bucket, path string, obj []byte) error {
	writer := c.Bucket(bucket).Object(path).NewWriter(bg.CtxBG)
	defer writer.Close()
	_, err := writer.Write(obj)
	return err
}

// Read reads object from bucket.
func (c *Client) Read(bucket, path string) (data []byte, err error) {
	r, err := c.Bucket(bucket).Object(path).NewReader(bg.CtxBG)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}
