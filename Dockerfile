#
# Copyright 2021 The Sigstore Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.23.3@sha256:73f06be4578c9987ce560087e2e2ea6485fb605e3910542cadd8fa09fc5f3e31 as builder
ENV APP_ROOT=/opt/app-root
ENV GOPATH=$APP_ROOT

WORKDIR $APP_ROOT/src/
ADD go.mod go.sum $APP_ROOT/src/
RUN go mod download

# Add source code
ADD ./cmd/ $APP_ROOT/src/cmd/
ADD ./pkg/ $APP_ROOT/src/pkg/

RUN go build ./cmd/verifier
RUN CGO_ENABLED=0 go build -gcflags "all=-N -l"  -o verifier_debug ./cmd/verifier

# Multi-Stage build
FROM golang:1.23.3@sha256:73f06be4578c9987ce560087e2e2ea6485fb605e3910542cadd8fa09fc5f3e31 as deploy

# Retrieve the binary from the previous stage
COPY --from=builder /opt/app-root/src/verifier /usr/local/bin/verifier

# Set the binary as the entrypoint of the container
CMD ["verifier"]

# debug compile options & debugger
FROM deploy as debug
RUN go install github.com/go-delve/delve/cmd/dlv@v1.23.1

# overwrite utility and include debugger
COPY --from=builder /opt/app-root/src/verifier_debug /usr/local/bin/verifier
