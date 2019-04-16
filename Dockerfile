FROM alpine-wecube:v1
LABEL maintainer = "Webank CTB Team"

ENV APP_HOME=/home/app/wecube-plugins
ENV APP_CONF=$APP_HOME/conf
ENV LOG_PATH=$APP_HOME/logs

RUN mkdir -p $APP_HOME $APP_CONF $LOG_PATH

ADD wecube-plugins $APP_HOME/
ADD *.sh $APP_HOME/
ADD conf $APP_CONF/

RUN chmod +x $APP_HOME/*.*

WORKDIR $APP_HOME

ENTRYPOINT ["/bin/sh", "start.sh"]