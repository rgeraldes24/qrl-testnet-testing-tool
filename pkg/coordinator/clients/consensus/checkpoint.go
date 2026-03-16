package consensus

import "github.com/rgeraldes24/go-qrl-beacon-client/spec/zond"

type FinalizedCheckpoint struct {
	Epoch zond.Epoch
	Root  zond.Root
}
