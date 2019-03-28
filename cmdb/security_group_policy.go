package cmdb

const (
	SECURITY_GROUP_POLICY_CI_NAME = "wb_security_group_policy"
)

type UpdateSecurityGroupPolicyCiEntry struct {
	Guid              string `json:"guid,omitempty"`
	State             string `json:"state,omitempty"`
	SecurityGroupId   string `json:"security_group_id,omitempty"`
	Type              string `json:"rule_type,omitempty"`
	PolicyIndex       int64  `json:"index,omitempty"`
	Protocol          string `json:"ip_protocol,omitempty"`
	Port              string `json:"port_range,omitempty"`
	CidrBlock         string `json:"cidr_ip,omitempty"`
	Action            string `json:"policy,omitempty"`
	PolicyDescription string `json:"description,omitempty"`
}
