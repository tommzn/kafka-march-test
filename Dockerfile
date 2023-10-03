FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.19 as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

ARG Version
ARG GitCommit

ENV CGO_ENABLED=1
ENV GO111MODULE=on
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

WORKDIR /go/build

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY .  .

RUN go build -tags musl --ldflags "-extldflags -static" -v -o build_artifact_bin

FROM --platform=${BUILDPLATFORM:-linux/amd64} gcr.io/distroless/static:nonroot

LABEL org.opencontainers.image.source=https://github.com/tommzn/hdb-api

WORKDIR /go

COPY --from=builder /go/build/build_artifact_bin kafka-march-bin
USER nonroot:nonroot

ENTRYPOINT ["/go/kafka-march-bin"]
