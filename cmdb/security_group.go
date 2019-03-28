package cmdb

const (
	SECURITY_GROUP_CI_NAME = "wb_security_group"
)

type UpdateSecurityGroupCiEntry struct {
	Guid              string `json:"guid,omitempty"`
	SecurityGroupId   string `json:"security_group_id,omitempty"`
	SecurityGroupName string `json:"security_group_name,omitempty"`
	SecurityGroupDesc string `json:"security_group_desc,omitempty"`
	State             string `json:"state,omitempty"`
}
