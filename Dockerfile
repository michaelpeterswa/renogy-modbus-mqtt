# -=-=-=-=-=-=- Compile Image -=-=-=-=-=-=-

FROM --platform=$BUILDPLATFORM golang:1-alpine AS stage-compile

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./... && CGO_ENABLED=0 GOOS=linux go build ./cmd/renogy-modbus-mqtt

# -=-=-=-=- Final Distroless Image -=-=-=-=-

FROM debian:bookworm-slim as stage-final

COPY --from=stage-compile /go/src/app/renogy-modbus-mqtt /
CMD ["/renogy-modbus-mqtt"]
