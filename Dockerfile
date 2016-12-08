FROM alpine:latest

RUN apk --no-cache update && \
  apk --no-cache upgrade && \
  apk --no-cache add ca-certificates && \
  adduser -D kaffeefee

COPY bin/kaffeefee /home/kaffeefee/
COPY static/ /home/kaffeefee/static/
COPY tpl/ /home/kaffeefee/tpl/

RUN chown -R kaffeefee:kaffeefee /home/kaffeefee && chmod +x /home/kaffeefee/kaffeefee

ENV DB="postgres://kaffeefee:kaffeefee@db/kaffeefee?sslmode=disable"
EXPOSE 8080
VOLUME /home/kaffeefee/static

USER kaffeefee
WORKDIR /home/kaffeefee

CMD ./kaffeefee
