FROM golang:1.23 AS build

RUN mkdir /build
COPY . /build
WORKDIR /build

RUN go mod download

ENV GOCACHE=/root/.cache/go-build

RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 go build -o hugespaceship ./cmd/hugespaceship

FROM scratch AS hugespaceship
COPY --from=build /build/hugespaceship /
ENTRYPOINT [ "/hugespaceship" ]