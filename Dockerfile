FROM golang:latest as binder
LABEL stage=builder

COPY . /api
WORKDIR /api

ARG TAG

# Input the origin package path
ENV GO111MODULE=on

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-X main.VERSION=$TAG" \
	-o api

# 
FROM alpine:latest 

EXPOSE 5000

COPY --from=binder /api/api /api

WORKDIR /
ENTRYPOINT ["./api"]
CMD  ["version"]