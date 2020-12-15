FROM alpine:3.7

WORKDIR /app

ADD ./dist/smtpd /app/smtpd
ADD ./cmd/config.yml /app/config.yml

EXPOSE 2525/tcp

CMD ["/app/smtpd"]
