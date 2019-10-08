package securitygroup

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const ENV_SECRET_ID = "SECRET_ID"
const ENV_SECRET_KEY = "SECRET_KEY"
const ENV_SUPPORT_REGIONS = "REGIONS" //用分号隔开多个地域

func getProviderParams(region string) (string, error) {
	secretId := os.Getenv(ENV_SECRET_ID)
	secretKey := os.Getenv(ENV_SECRET_KEY)

	if secretId == "" {
		err := errors.New("can't get secretId from env")

		logrus.Errorf("getProviderParams meet error=%v", err)
		return "", err
	}
	if secretKey == "" {
		err := errors.New("can't get secretKey from env")

		logrus.Errorf("getProviderParams meet error=%v", err)
		return "", err
	}
	if region == "" {
		err := errors.New("input region is empty")

		logrus.Errorf("getProviderParams meet error=%v", err)
		return "", err
	}

	return fmt.Sprintf("Region=%s;SecretID=%s;SecretKey=%s", region, secretId, secretKey), nil
}

func getRegions() ([]string, error) {
	regions := strings.Split(os.Getenv(ENV_SUPPORT_REGIONS), ";")
	if len(regions) == 0 {
		err := errors.New("can't get region from env")

		logrus.Errorf("getRegions meet error=%v", err)
		return regions, err
	}

	return regions, nil
}
