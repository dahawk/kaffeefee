FROM golang:1.16-alpine3.12 AS build

WORKDIR /go/src

COPY . .

RUN go build -o kaffeefee

FROM alpine:3.12

RUN apk add ca-certificates && \
  adduser -D kaffeefee

EXPOSE 8080
VOLUME /home/kaffeefee/static

USER kaffeefee
WORKDIR /home/kaffeefee

CMD ["/home/kaffeefee/kaffeefee"]

COPY --from=build --chown=kaffeefee:kaffeefee /go/src/kaffeefee /home/kaffeefee/
COPY static/ /home/kaffeefee/static/
COPY tpl/ /home/kaffeefee/tpl/