# Copyright 2021 Dmitry Kisler dkisler.com
# Licensed under the Apache License,Version 2.0 (the "License");
# you may not use this file except in compliance with the License. You may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
# INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE,
# AND NONINFRINGEMENT. IN NO EVENT WILL THE LICENSOR OR OTHER CONTRIBUTORS BE LIABLE FOR ANY CLAIM, DAMAGES,
# OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF,
# OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
# See the License for the specific language governing permissions and limitations under the License.

FROM golang:1.16.2-alpine3.12 AS build

RUN mkdir /go/src/app
WORKDIR /go/src/app

ARG SERVICE

COPY lib lib
COPY ${SERVICE} app

RUN cd app \
    && go mod tidy \
    && CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -gcflags=all="-l -B -C" -ldflags="-w -s" -o /root/runner *.go

RUN echo "ginkgo:x:1000:1000:ginkgo,,,::/bin/false" > /user.txt

FROM scratch AS run

# adds x509 cert
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /user.txt /etc/passwd
COPY --from=build /root/runner /runner

USER ginkgo

ENV GCP_PROJECT ""
ENV NOTIFICATION_TOPIC ""
ENV NOTIFICATION_TOPIC_FAIL ""

ENV PORT 9000
EXPOSE ${PORT}

ENTRYPOINT ["./runner"]
