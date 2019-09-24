package securitygroup

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	INGRESS_RULE = "ingress" //入站规则
	EGRESS_RULE  = "egress"  //出栈规则

	TCP_PROTOCOL  = "TCP"
	UDP_PROTOCOL  = "UDP"
	ICMP_PROTOCOL = "ICMP"

	POLICY_ACTION_ACCEPT = "ACCEPT"
	POLICY_ACTION_DROP   = "DROP"

	ALL_PORT = "ALL"
)

func isContainInList(input string, list []string) bool {
	for _, item := range list {
		if item == input {
			return true
		}
	}
	return false
}

func isValidValueIgnoreCase(inputValue string, validValues []string) error {
	for _, validValue := range validValues {
		if strings.EqualFold(validValue, inputValue) {
			return nil
		}
	}
	return fmt.Errorf("%s is not valid value in(%++v)", inputValue, validValues)
}

func isValidProtocol(protocol string) error {
	validProtocols := []string{TCP_PROTOCOL, UDP_PROTOCOL, ICMP_PROTOCOL}
	if err := isValidValueIgnoreCase(protocol, validProtocols); err != nil {
		return fmt.Errorf("protocol(%s) is invalid", protocol)
	}
	return nil
}

func isValidIp(ip string) error {
	ipaddr := net.ParseIP(ip)
	if ipaddr == nil {
		return fmt.Errorf("ip(%s) is invalid", ip)
	}
	return nil
}

func isValidPort(port string) (int, error) {
	portInt, err := strconv.Atoi(port)
	if err != nil || portInt >= 65535 {
		return 0, fmt.Errorf("port(%s) is invalid", port)
	}
	return portInt, nil
}

func isValidAction(action string) error {
	validActions := []string{POLICY_ACTION_ACCEPT, POLICY_ACTION_DROP}
	if err := isValidValueIgnoreCase(action, validActions); err != nil {
		return fmt.Errorf("action(%s) is invalid", action)
	}
	return nil
}

func isValidDirection(direction string) error {
	validDirections := []string{INGRESS_RULE, EGRESS_RULE}
	if err := isValidValueIgnoreCase(direction, validDirections); err != nil {
		return fmt.Errorf("direction(%s) is invalid", direction)
	}

	return nil
}

//"8090;8080;80-7000;ALL"这种格式
func getPortsByPolicyFormat(portStr string) ([]string, error) {
	allPorts := []string{}
	singlePorts := []string{}
	rangePorts := []string{}
	rtnPorts := []string{}

	ports := strings.Split(portStr, ";")
	for _, port := range ports {
		//all
		port = strings.TrimSpace(port)
		if strings.EqualFold(port, ALL_PORT) {
			allPorts = append(allPorts, ALL_PORT)
			break
		}

		//single port
		portInt, err := strconv.Atoi(port)
		if err == nil && portInt <= 65535 {
			singlePorts = append(singlePorts, port)
			continue
		}

		//range port
		portRange := strings.Split(port, "-")
		if len(portRange) == 2 {
			firstPort, firstErr := isValidPort(portRange[0])
			lastPort, lastErr := isValidPort(portRange[1])
			if firstErr == nil && lastErr == nil && firstPort < lastPort {
				rangePorts = append(rangePorts, port)
				continue
			}
		}

		return rtnPorts, fmt.Errorf("port(%s) is invalid", port)
	}
	if len(allPorts) > 0 {
		return allPorts, nil
	}

	if len(singlePorts) > 0 {
		rtnPorts = append(rtnPorts, strings.Join(singlePorts, ","))
	}

	if len(rangePorts) > 0 {
		rtnPorts = append(rtnPorts, rangePorts...)
	}
	return rtnPorts, nil
}
