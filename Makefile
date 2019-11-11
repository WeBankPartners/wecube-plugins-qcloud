current_dir=$(shell pwd)
version=${PLUGIN_VERSION}
project_name=$(shell basename "${current_dir}")


APP_HOME=src/github.com/WeBankPartners/wecube-plugins-qcloud
PORT_BINDINGS={{host_port}}:8081

ifndef RUN_MODE
  RUN_MODE=dev
endif

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
	chmod +x ./build/register.xml.tpl
	sed -i 's/{{PLUGIN_VERSION}}/$(version)/' ./build/register.xml.tpl > ./register.xml
	sed -i 's/{{PORTBINDINGS}}/$(PORT_BINDINGS)/' ./register.xml
	sed -i 's/{{IMAGENAME}}/$(project_name):$(version)/' ./register.xml 
	docker save -o  image.tar $(project_name):$(version)
	zip  $(project_name)-$(version).zip image.tar register.xml
	rm -rf ./*.tar
	rm -f register.xml
	docker rmi $(project_name):$(version)
	rm -rf $(project_name)
