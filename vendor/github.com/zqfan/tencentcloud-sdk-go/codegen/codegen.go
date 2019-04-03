package main

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// goFormat returns the Go formated string of the input.
func goFormat(s string) string {
	buf, err := format.Source([]byte(s))
	if err != nil {
		log.Panicf("Failed to format %s because %s", s, err.Error())
	}
	return string(buf)
}

func main() {
	tmpl := template.Must(template.New("api").Parse(`package {{.PkgName}}

import (
    "github.com/zqfan/tencentcloud-sdk-go/common"
)

const APIVersion = "{{.APIVersion}}"
{{range .Actions}}
func New{{.Action}}Request() (request *{{.Action}}Request) {
    request = &{{.Action}}Request{
        BaseRequest: &common.BaseRequest{},
    }
    request.Init().WithApiInfo("{{.SubDomain}}", APIVersion, "{{.Action}}")
    return
}

func New{{.Action}}Response() (response *{{.Action}}Response) {
    response = &{{.Action}}Response{
        BaseResponse: &common.BaseResponse{},
    }
    return
}

func (c *Client) {{.Action}}(request *{{.Action}}Request) (response *{{.Action}}Response, err error) {
    if request == nil {
        request = New{{.Action}}Request()
    }
    response = New{{.Action}}Response()
    err = c.Send(request, response)
    return
}
{{end}}`))
	type Action struct {
		Action    string
		SubDomain string
	}
	type Service struct {
		PkgName    string
		APIVersion string
		Actions    []Action
	}
	services := []Service{
		Service{
			PkgName:    "cvm",
			APIVersion: "2017-03-12",
			Actions: []Action{
				{"DescribeAddresses", "eip"},
				{"ReleaseAddresses", "eip"},
				{"ModifyAddressAttribute", "eip"},
				{"AllocateAddresses", "eip"},
				{"AssociateAddress", "eip"},
				{"DisassociateAddress", "eip"},
				{"RunInstances", "cvm"},
				{"DescribeInstances", "cvm"},
				{"TerminateInstances", "cvm"},
			},
		},
		Service{
			PkgName:    "cbs",
			APIVersion: "",
			Actions: []Action{
				{"CreateSnapshot", "snapshot"},
				{"DeleteSnapshot", "snapshot"},
				{"ModifySnapshot", "snapshot"},
				{"DescribeSnapshots", "snapshot"},
				{"DescribeCbsStorages", "snapshot"},
			},
		},
		Service{
			PkgName:    "vpc",
			APIVersion: "",
			Actions: []Action{
				{"AddDnaptRule", "vpc"},
				{"CreateNatGateway", "vpc"},
				{"DeleteDnaptRule", "vpc"},
				{"DeleteNatGateway", "vpc"},
				{"DescribeNatGateway", "vpc"},
				{"DescribeNetworkInterfaces", "vpc"},
				{"DescribeVpcEx", "vpc"},
				{"DescribeVpcTaskResult", "vpc"},
				{"EipBindNatGateway", "vpc"},
				{"EipUnBindNatGateway", "vpc"},
				{"GetDnaptRule", "vpc"},
				{"ModifyDnaptRule", "vpc"},
				{"ModifyNatGateway", "vpc"},
				{"QueryNatGatewayProductionStatus", "vpc"},
				{"UpgradeNatGateway", "vpc"},
			},
		},
		Service{
			PkgName:    "ccs",
			APIVersion: "",
			Actions: []Action{
				{"AddClusterInstances", "ccs"},
				{"AddClusterInstancesFromExistedCvm", "ccs"},
				{"CreateCluster", "ccs"},
				{"DeleteCluster", "ccs"},
				{"DeleteClusterInstances", "ccs"},
				{"DescribeCluster", "ccs"},
				{"DescribeClusterInstances", "ccs"},
				{"DescribeClusterSecurityInfo", "ccs"},
				{"DescribeClusterTaskResult", "ccs"},
				{"OperateClusterVip", "ccs"},
			},
		},
		Service{
			PkgName:    "lb",
			APIVersion: "",
			Actions: []Action{
				{"CreateLoadBalancer", "lb"},
				{"DeleteLoadBalancers", "lb"},
				{"DescribeForwardLBBackends", "lb"},
				{"DescribeForwardLBListeners", "lb"},
				{"DescribeLoadBalancers", "lb"},
				{"DescribeLoadBalancersTaskResult", "lb"},
				{"DeregisterInstancesFromForwardLB", "lb"},
				{"ModifyForwardLBName", "lb"},
				{"ModifyLoadBalancerAttributes", "lb"},
				{"RegisterInstancesWithForwardLBSeventhListener", "lb"},
			},
		},
	}
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	serviceDir := filepath.Join(dir, "..", "services")
	for _, service := range services {
		var buffer bytes.Buffer
		tmpl.Execute(&buffer, service)
		var versionPath string
		if service.APIVersion == "" {
			versionPath = "unversioned"
		} else {
			versionPath = "v" + strings.Replace(service.APIVersion, "-", "", -1)
		}
		apiFilePath := filepath.Join(serviceDir, service.PkgName, versionPath, "api.go")
		err := ioutil.WriteFile(apiFilePath, []byte(goFormat(buffer.String())), 0644)
		if err != nil {
			log.Println(err)
		}
	}
}
