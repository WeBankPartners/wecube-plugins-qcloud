export GOPATH=$(PWD)

APP_HOME=src/git.webank.io/wecube-plugins

ifndef RUN_MODE
  RUN_MODE=dev
endif

archive:
	tar cvfz source.tar.gz *
	rm -rf src
	mkdir -p $(APP_HOME)
	rm -rf target
	mkdir target
	tar zxvf source.tar.gz -C $(APP_HOME)
	rm -rf source.tar.gz
	cd $(APP_HOME) && CGO_ENABLED=0 GOOS=linux go build
	cp -R $(APP_HOME)/conf target
	cp start.sh stop.sh makefile target
	cd target && chmod 755 *.sh
	cp $(APP_HOME)/wecube-plugins target
	cd target && tar cvfz $(PKG_NAME) *
