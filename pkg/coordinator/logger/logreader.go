package logger

import "github.com/theQRL/assertoor/pkg/coordinator/db"

type LogReader interface {
	GetLogEntries(from, limit uint64) []*db.TaskLog
	GetLogEntryCount() uint64
}
