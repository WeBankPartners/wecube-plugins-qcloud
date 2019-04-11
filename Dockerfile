FROM alpine:latest
LABEL maintainer = "Webank CTB Team"

RUN mkdir -p /home/app
ENV LOG_PATH=/home/app/logs

COPY wecube-plugins /home/app
ADD conf /home/app/conf
RUN cd /home/app && chmod +x wecube-plugins

WORKDIR /home/app

ENTRYPOINT ["./wecube-plugins"]

