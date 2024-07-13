FROM golang:1.22 AS build

RUN mkdir /build
COPY . /build
WORKDIR /build

RUN go mod download

ENV GOCACHE=/root/.cache/go-build

RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 go build -o hugespaceship ./cmd/hugespaceship
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 go build -o resource_server ./cmd/resource_server

FROM scratch AS hugespaceship
COPY --from=build /build/hugespaceship /
# this is to stop HS from trying to write a log file, which is currently the default behaviour
ENV HS_LOG_FILE_LOGGING=false
ENTRYPOINT [ "/hugespaceship" ]

FROM scratch AS resource_server
COPY --from=build /build/resource_server /
ENV HS_LOG_FILE_LOGGING=false
ENTRYPOINT [ "/resource_server" ]
