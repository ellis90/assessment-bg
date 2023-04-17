FROM golang:1.18.7-alpine3.16 as builder

RUN apk add --no-cache \
    build-base \
    libmediainfo-dev \
     openssl

WORKDIR /app

ENV GO111MODULE=on
ENV GOPATH /go


COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .


COPY create-env.sh .
RUN chmod +x create-env.sh


RUN go install github.com/cosmtrek/air@latest


CMD ["air", "-c", ".air.toml"]

#FROM scratch
#
#WORKDIR /
#
#COPY --from=builder /app /app
#
#

##RUN set -eux; \
##	export GOROOT="$(go env GOROOT)"; \
##	./air \