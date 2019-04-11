FROM alpine:latest
LABEL maintainer = "Webank CTB Team"

ENV LOG_PATH=/home/app/logs

ADD wecube-plugins /home
ADD conf /home
RUN cd /home && chmod +x wecube-plugins

WORKDIR /home

ENTRYPOINT ["./wecube-plugins"]

