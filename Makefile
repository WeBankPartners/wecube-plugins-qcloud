export GOPATH=$(PWD)

current_dir=$(shell pwd)
version=${PLUGIN_VERSION}
project_name=$(shell basename "${current_dir}" )


APP_HOME=src/github.com/WeBankPartners/wecube-plugins-qcloud

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
	cp start.sh stop.sh docker_run.sh docker_stop.sh makefile dockerfile register.xml target
	cd target && chmod 755 *.sh
	cp $(APP_HOME)/wecube-plugins-qcloud target
	cd target && tar cvfz $(PKG_NAME) *

clean:
	rm -rf $(project_name)
	rm -rf  ./*.tar
	rm -rf ./*.zip
fmt:
	docker run --rm -v $(current_dir):/go/src/github.com/WeBankPartners/$(project_name) --name build_$(project_name) -w /go/src/github.com/WeBankPartners/$(project_name)/  golang:1.12.5 go fmt ./...


build: clean
	chmod +x ./build/*.sh
	docker run --rm -v $(current_dir):/go/src/github.com/WeBankPartners/$(project_name) --name build_$(project_name) golang:1.12.5 /bin/bash /go/src/github.com/WeBankPartners/$(project_name)/build/build.sh 

image: build
	docker build -t $(project_name):$(version) .
     
package: image 
	sed 's/{{IMAGE_TAG}}/$(version)/' ./build/register.xml.tpl > ./register.xml
	sed -i 's/{{PLUGIN_VERSION}}/$(version)/' ./register.xml 
	docker save -o  image.tar $(project_name):$(version)
	zip  $(project_name)_$(version).zip image.tar register.xml
	rm -rf ./*.tar
	docker rmi $(project_name):$(version)
	rm -rf $(project_name)
