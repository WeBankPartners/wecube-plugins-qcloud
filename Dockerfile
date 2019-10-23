FROM alpine
LABEL maintainer = "Webank CTB Team"

ENV APP_HOME=/home/app/wecube-plugins-qcloud
ENV APP_CONF=$APP_HOME/conf
ENV LOG_PATH=$APP_HOME/logs

RUN apk add ca-certificates
RUN mkdir -p $APP_HOME $APP_CONF $LOG_PATH

ADD scripts  $APP_HOME/
ADD wecube-plugins-qcloud $APP_HOME/
ADD build/start.sh $APP_HOME/
ADD build/stop.sh $APP_HOME/
ADD conf $APP_CONF/

RUN chmod +x $APP_HOME/*.*

WORKDIR $APP_HOME

ENTRYPOINT ["/bin/sh", "start.sh"]
