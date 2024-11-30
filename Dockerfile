FROM --platform=$BUILDPLATFORM golang:1.23 AS build

RUN mkdir /build
COPY . /build
WORKDIR /build

RUN go mod download

ENV GOCACHE=/root/.cache/go-build

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target="/root/.cache/go-build" \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 \
    go build -o hugespaceship-${TARGETOS}-${TARGETARCH} ./cmd/hugespaceship

FROM scratch AS hugespaceship

ARG TARGETOS
ARG TARGETARCH

COPY --from=build /build/hugespaceship-${TARGETOS}-${TARGETARCH} /
ENTRYPOINT [ "/hugespaceship" ]