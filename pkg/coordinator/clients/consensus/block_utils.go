package consensus

import (
	"errors"

	"github.com/rgeraldes24/go-qrl-beacon-client/spec"
)

func GetExecutionExtraData(v *spec.VersionedSignedBeaconBlock) ([]byte, error) {
	//nolint:exhaustive // ignore
	switch v.Version {
	case spec.DataVersionZond:
		if v.Zond == nil || v.Zond.Message == nil || v.Zond.Message.Body == nil || v.Zond.Message.Body.ExecutionPayload == nil {
			return nil, errors.New("no capella block")
		}

		return v.Zond.Message.Body.ExecutionPayload.ExtraData, nil
	default:
		return nil, errors.New("unknown version")
	}
}

func GetBlockBody(v *spec.VersionedSignedBeaconBlock) any {
	//nolint:exhaustive // ignore
	switch v.Version {
	case spec.DataVersionZond:
		return v.Zond
	default:
		return nil
	}
}
