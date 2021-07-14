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

SHELL := /bin/bash

.PHONY: service.rebuild

PROJECT_ID :=
SERVICE_NAME :=

REPO := eu.gcr.io/$(PROJECT_ID)
IMAGE_TAG := v1.0
IMAGE := $(REPO)/$(SERVICE_NAME):$(IMAGE_TAG)

service.rebuild: service.image.build service.image.push

service.image.build:
	@ cd services \
	&& docker build --build-arg SERVICE=$(SERVICE_NAME) -t $(IMAGE) .

service.image.push:
	@ docker push $(IMAGE)

tf.apply:
	@ cd ./infra/terraform \
	&& export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/.gcp/ginkgo/terraform.json \
	&& terraform init \
	&& terraform plan \
	&& terraform apply -auto-approve
