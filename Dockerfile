FROM golang:alpine AS build
WORKDIR /go/src/currency-converter
COPY ./ ./
ARG GIT_COMMIT
ARG GOPROXY

RUN go build -ldflags maing.go -o /go/bin/currency-converter

FROM gcr.io/distroless/base-debian11
COPY --from=build /go/bin/currency-converter /usr/local/bin/currency-converter

ENTRYPOINT [ "/usr/local/bin/dcurrency-converter" ]
