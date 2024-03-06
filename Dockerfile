FROM golang:1.21.4-alpine as builder

ENV CGO_ENABLED=1

RUN apk add --no-cache \
  #required for go-sqlite3
  gcc \
  #required for Alpine
  musl-dev

WORKDIR /workdir

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY . ./

RUN go build -o /out/main ./main.go

# DEPLOY
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /usr/src

COPY --from=builder /out/ .

EXPOSE 3000

ENTRYPOINT [ "./main" ]