FROM golang:1.23.2 AS builder

WORKDIR /boombox-web-runner

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./dist/app 

FROM ghcr.io/go-rod/rod AS runner

RUN apt update -y && apt upgrade -y
RUN apt install -y ffmpeg

RUN mkdir -p /go/pkg/mod/github.com/navicstein

COPY --from=builder /boombox-web-runner/dist/app /usr/bin/app
COPY --from=builder /go/pkg/mod/github.com/navicstein /go/pkg/mod/github.com/navicstein

CMD ["/usr/bin/app"]
