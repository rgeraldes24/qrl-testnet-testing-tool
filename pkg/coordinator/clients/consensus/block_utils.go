package consensus

import (
	"errors"

	"github.com/attestantio/go-eth2-client/spec"
)

func GetExecutionExtraData(v *spec.VersionedSignedBeaconBlock) ([]byte, error) {
	//nolint:exhaustive // ignore
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil || v.Capella.Message.Body.ExecutionPayload == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.ExecutionPayload.ExtraData, nil
	default:
		return nil, errors.New("unknown version")
	}
}

func GetBlockBody(v *spec.VersionedSignedBeaconBlock) any {
	//nolint:exhaustive // ignore
	switch v.Version {
	case spec.DataVersionCapella:
		return v.Capella
	default:
		return nil
	}
}
