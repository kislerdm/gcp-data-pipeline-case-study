// Copyright 2021 Dmitry Kisler dkisler.com

// Licensed under the Apache License,Version 2.0 (the "License");
// you may not use this file except in compliance with the License. You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE,
// AND NONINFRINGEMENT. IN NO EVENT WILL THE LICENSOR OR OTHER CONTRIBUTORS BE LIABLE FOR ANY CLAIM, DAMAGES,
// OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF,
// OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// See the License for the specific language governing permissions and limitations under the License.

module platform/process

go 1.16

replace platform/lib => ../lib

require (
	cloud.google.com/go/datastore v1.5.0
	github.com/goccy/go-json v0.7.4
	google.golang.org/genproto v0.0.0-20210713002101-d411969a0d9a // indirect
	platform/lib v0.0.0-00010101000000-000000000000
)
