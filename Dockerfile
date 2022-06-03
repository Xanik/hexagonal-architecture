FROM golang:latest AS builder

RUN mkdir -p /go/src/study

WORKDIR /go/src/study

COPY . .

RUN GIT_TERMINAL_PROMPT=1 \
    GOARCH=amd64 \
    GOOS=linux \
    CGO_ENABLED=0 \
    go build -v --installsuffix cgo --ldflags="-s" -o study

FROM alpine:latest

RUN mkdir -p /app/study
COPY --from=builder /go/src/study/study /app/study/

WORKDIR /app/study

RUN apk add --no-cache mailcap
RUN apk add --no-cache tzdata

CMD ["./study"]