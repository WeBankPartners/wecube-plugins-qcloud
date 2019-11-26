FROM webankpartners/alpine-ca-certificate
LABEL maintainer = "Webank CTB Team"

ENV APP_HOME=/home/app/wecube-plugins-qcloud
ENV APP_CONF=$APP_HOME/conf
ENV LOG_PATH=$APP_HOME/logs

RUN mkdir -p $APP_HOME $APP_CONF $LOG_PATH

COPY build/start.sh $APP_HOME/ 
COPY build/stop.sh $APP_HOME/ 
RUN chmod +x $APP_HOME/*.*

COPY scripts $APP_HOME/scripts/
COPY wecube-plugins-qcloud $APP_HOME/
COPY conf $APP_CONF/

WORKDIR $APP_HOME

ENTRYPOINT ["/bin/sh", "start.sh"]
