package consensus

import (
	"reflect"
	"time"

	"github.com/rgeraldes24/go-qrl-beacon-client/spec/zond"
)

type ForkVersion struct {
	Epoch           uint64
	CurrentVersion  []byte
	PreviousVersion []byte
}

// https://github.com/ethereum/consensus-specs/blob/dev/configs/mainnet.yaml
type ChainSpec struct {
	PresetBase           string        `yaml:"PRESET_BASE"`
	ConfigName           string        `yaml:"CONFIG_NAME"`
	MinGenesisTime       time.Time     `yaml:"MIN_GENESIS_TIME"`
	GenesisForkVersion   zond.Version  `yaml:"GENESIS_FORK_VERSION"`
	SecondsPerSlot       time.Duration `yaml:"SECONDS_PER_SLOT"`
	SlotsPerEpoch        uint64        `yaml:"SLOTS_PER_EPOCH"`
	MaxCommitteesPerSlot uint64        `yaml:"MAX_COMMITTEES_PER_SLOT"`
}

func (chain *ChainSpec) CheckMismatch(chain2 *ChainSpec) []string {
	mismatches := []string{}

	chainT := reflect.ValueOf(chain).Elem()
	chain2T := reflect.ValueOf(chain2).Elem()

	for i := 0; i < chainT.NumField(); i++ {
		if chainT.Field(i).Interface() != chain2T.Field(i).Interface() {
			mismatches = append(mismatches, chainT.Type().Field(i).Name)
		}
	}

	return mismatches
}
