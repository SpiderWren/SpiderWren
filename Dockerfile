FROM golang:1.16-alpine AS build
WORKDIR /build
COPY . .
RUN apk add build-base
RUN go mod download
RUN go build -o wren-web

FROM alpine
COPY --from=build /build/wren-web /usr/bin/
ENV PATH="/usr/bin:${PATH}"