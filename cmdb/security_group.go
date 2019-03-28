package cmdb

const (
	SECURITY_GROUP_CI_NAME = "wb_security_group"
)

type UpdateSecurityGroupCiEntry struct {
	Guid            string `json:"guid,omitempty"`
	SecurityGroupId string `json:"id,omitempty"`
	State           string `json:"state,omitempty"`
}
