package securitygroup
import (
	"os"
	"errors"
	"strings"
)
const ENV_SECRET_ID ="SECRET_ID"
const ENV_SECRET_KEY="SECRET_KEY" 
const ENV_SUPPORT_REGIONS="REGIONS"  //用分号隔开多个地域

func getProviderParams(region string)(string,error) {
	secretId:=os.Getenv(ENV_SECRET_ID)
	secretKey:=os.Getenv(ENV_SECRET_KEY)

	if secretId == "" {
		return "",errors.New("can't get secretId from env")
	}
	if secretKey == "" {
		return "",errors.New("can't get secretKey from env")
	}
	if region =="" {
		return "",errors.New("input region is empty")
	}

    return fmt.Sprintf("Region=%s;SecretID=%s;SecretKey=%s",region,secretId,secretKey),nil
}

func getRegions()([]string,error) {
	regions:=strings.Split(os.Getenv(ENV_SUPPORT_REGIONS),";")
	if len(regions) == 0 {
		return regions,errors.New("can't get region from env")
	}
	return regions,nil 
}