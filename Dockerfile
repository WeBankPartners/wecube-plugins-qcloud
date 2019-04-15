FROM alpine:latest
LABEL maintainer = "Webank CTB Team"

ARG DEPLOY_PATH=/home/app/wecube-plugins

RUN mkdir -p /home/app
ENV LOG_PATH=/home/app/logs

COPY wecube-plugins /home/app
ADD conf /home/app/conf

RUN cd /home/app && chmod +x wecube-plugins

RUN apk upgrade && apk add --no-cache ca-certificates
RUN apk add --update curl && rm -rf /var/cache/apk/*

WORKDIR /home/app

ENTRYPOINT ["./wecube-plugins"]

