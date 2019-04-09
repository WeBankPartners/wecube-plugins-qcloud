FROM alpine:latest
LABEL maintainer = "Webank CTB Team"

ARG DEPLOY_PATH=/home/app/wecube-plugins

ENV LOG_PATH=$DEPLOY_PATH/logs
RUN mkdir -p $DEPLOY_PATH/

ADD wecube-plugins $DEPLOY_PATH/
ADD conf $DEPLOY_PATH/

WORKDIR $DEPLOY_PATH
ENTRYPOINT ["wecube-plugins"]
