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

package transformation_test

import (
	"platform/process/transformation"
	"reflect"
	"testing"
	"time"
)

func toTS(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func TestConvertTimestampUTC(t *testing.T) {
	tests := []struct {
		in   time.Time
		want time.Time
	}{
		{
			in:   toTS("2019-05-01T06:00:00-04:00"),
			want: toTS("2019-05-01T10:00:00Z"),
		},
		{
			in:   toTS("2021-07-08T06:00:00-04:00"),
			want: toTS("2021-07-08T10:00:00Z"),
		},
		{
			in:   toTS("2021-07-08T00:00:00-04:00"),
			want: toTS("2021-07-08T04:00:00Z"),
		},
	}
	for _, test := range tests {
		got := transformation.ConvertTimestampUTC(test.in)
		if got != test.want {
			t.Fatalf("converter fail!\nwant:%v\ngot: %v\n", test.want, got)
		}
	}
}

func TestStats(t *testing.T) {
	tests := []struct {
		in   []float64
		want *transformation.DistributionStats
	}{
		{
			in:   []float64{0, 0, 0, 0},
			want: &transformation.DistributionStats{0, 0},
		},
		{
			in:   []float64{3, 3, 3},
			want: &transformation.DistributionStats{3., 0.},
		},
	}
	for _, test := range tests {
		got := transformation.NewDistribution(test.in).Stats()
		if !reflect.DeepEqual(got, test.want) {
			t.Fatalf("dist stats fail!\nwant:%v\ngot: %v\n", test.want, got)
		}
	}
}
