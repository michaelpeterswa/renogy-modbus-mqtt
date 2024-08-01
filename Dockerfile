# -=-=-=-=-=-=- Compile Image -=-=-=-=-=-=-

FROM --platform=$BUILDPLATFORM golang:1 AS stage-compile

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./... && CGO_ENABLED=0 GOOS=linux go build ./cmd/renogy-modbus-mqtt

# -=-=-=-=- Final Distroless Image -=-=-=-=-

FROM gcr.io/distroless/static-debian12:latest as stage-final

COPY --from=stage-compile /go/src/app/renogy-modbus-mqtt /
CMD ["/renogy-modbus-mqtt"]
