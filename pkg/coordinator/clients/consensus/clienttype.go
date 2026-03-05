package consensus

import (
	"fmt"
	"regexp"
)

type ClientType int8

var (
	AnyClient     ClientType
	UnknownClient ClientType = -1
	QrysmClient   ClientType = 4
)

var clientTypePatterns = map[ClientType]*regexp.Regexp{
	QrysmClient: regexp.MustCompile("(?i)^Qrysm/.*"), // TODO(rgeraldes24)
}

func (client *Client) parseClientVersion(version string) {
	for clientType, versionPattern := range clientTypePatterns {
		if versionPattern.MatchString(version) {
			client.clientType = clientType
			return
		}
	}

	client.clientType = UnknownClient
}

func ParseClientType(name string) ClientType {
	switch name {
	case "qrysm":
		return QrysmClient
	default:
		return UnknownClient
	}
}

func (client *Client) GetClientType() ClientType {
	return client.clientType
}

func (clientType ClientType) String() string {
	switch clientType {
	case QrysmClient:
		return "qrysm"
	default:
		return fmt.Sprintf("unknown: %d", clientType)
	}
}
