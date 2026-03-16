package logger

import "github.com/theQRL/qrl-testnet-testing-tool/pkg/coordinator/db"

type LogReader interface {
	GetLogEntries(from, limit uint64) []*db.TaskLog
	GetLogEntryCount() uint64
}
