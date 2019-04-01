package cmdb

import (
	"github.com/sirupsen/logrus"
)

const (
	SECURITY_GROUP_INPUT_NAME = "SECURITY-GROUP-IDC"
)

type SecurityGroupInput struct {
	Guid              string `json:"guid,omitempty"`
	ProviderParams    string `json:"provider_params,omitempty"`
	Name              string `json:"name,omitempty"`
	Id                string `json:"id,omitempty"`
	Description       string `json:"description,omitempty"`
	State             string `json:"state,omitempty"`
	ProcessInstanceId string `json:"process_instance_id,omitempty"`
	RulePriority      string `json:"rule_priority,omitempty"`
	RuleType          string `json:"rule_type,omitempty"`
	RuleCidrIp        string `json:"rule_cidr_ip,omitempty"`
	RuleIpProtocol    string `json:"rule_ip_protocol,omitempty"`
	RulePortRange     string `json:"rule_port_range,omitempty"`
	RuleDescription   string `json:"rule_description,omitempty"`
}

func GetSecurityGroupInputsByProcessInstanceId(queryParam *CmdbCiQueryParam) ([]SecurityGroupInput, int, error) {
	results := []SecurityGroupInput{}
	queryParam.ResultColumn = ExtractColumnFromStruct(SecurityGroupInput{})

	total, err := ListIntegrateEntries(SECURITY_GROUP_INPUT_NAME, queryParam, &results)
	if err != nil {
		logrus.Errorf("GetSecurityGroupInputsByProcessInstanceId meet error err=%v,queryParam=%v", err, queryParam)
	}

	return results, total, err
}

const (
	SECURITY_GROUP_CI_NAME = "wb_security_group"
)

type SecurityGroupOutput struct {
	SecurityGroupId string `json:"id,omitempty"`
	State           string `json:"state,omitempty"`
}
