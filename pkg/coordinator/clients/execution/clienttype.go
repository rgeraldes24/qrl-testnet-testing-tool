package execution

import (
	"fmt"
	"regexp"
)

type ClientType int8

var (
	AnyClient     ClientType
	UnknownClient ClientType = -1
	GzondClient   ClientType = 4
)
var clientTypePatterns = map[ClientType]*regexp.Regexp{
	GzondClient: regexp.MustCompile("(?i)^Gzond/.*"),
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
	case "gzond":
		return GzondClient
	default:
		return UnknownClient
	}
}

func (client *Client) GetClientType() ClientType {
	return client.clientType
}

func (clientType ClientType) String() string {
	switch clientType {
	case GzondClient:
		return "gzond"
	default:
		return fmt.Sprintf("unknown: %d", clientType)
	}
}
