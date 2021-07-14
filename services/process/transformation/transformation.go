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

Package defines the logic to transform raw data accroding to requirements.
*/

package transformation

import (
	"math"
	"time"
)

// distribution defines the distribution of floats.
type distribution []float64

// NewDistribution init distribution.
func NewDistribution(data []float64) distribution {
	return distribution(data)
}

// distributionStats defines basic distribution stats.
type DistributionStats struct {
	Mean   float64
	Stddev float64
}

// Stats calculates the distribution stats.
func (d distribution) Stats() *DistributionStats {
	var sum float64
	var sumSq float64
	total := float64(len(d))
	for _, el := range d {
		sum += el
		sumSq += el * el
	}
	meanSq := sumSq / total
	mean := sum / total
	return &DistributionStats{
		Mean:   mean,
		Stddev: math.Sqrt(meanSq - mean*mean),
	}
}

var loc, _ = time.LoadLocation("UTC")

// ConvertTimestampUTC converts timestamp timezone to UTC.
func ConvertTimestampUTC(t time.Time) time.Time {
	return t.In(loc)
}
