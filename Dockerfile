FROM golang:1.12-alpine AS development

ENV GO111MODULE on

WORKDIR /go/src/github.com/IcaliaLabs/docker-image-manager

ENV HOME /go/src/github.com/IcaliaLabs/docker-image-manager

RUN apk add --no-cache git

# Stage II: Testing ============================================================

FROM development AS testing

COPY ./go.mod ./go.sum /go/src/github.com/IcaliaLabs/docker-image-manager/

RUN go mod download

COPY . /go/src/github.com/IcaliaLabs/docker-image-manager

# Stage III: Linux amd64 Builder ===============================================

FROM testing AS linux-amd64-builder

ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0

RUN go build \
  -a \
  -installsuffix cgo \
  -ldflags="-s -w" \
  -o /builds/docker-image-manager \
  main.go

# Stage IV: Release Linux amd64 ================================================

FROM scratch AS linux-amd64-release

COPY --from=linux-amd64-builder /builds/docker-image-manager /

ENTRYPOINT ["/docker-image-manager"]


