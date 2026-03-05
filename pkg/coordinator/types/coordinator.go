package types

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/theQRL/assertoor/pkg/coordinator/clients"
	"github.com/theQRL/assertoor/pkg/coordinator/db"
	"github.com/theQRL/assertoor/pkg/coordinator/logger"
	"github.com/theQRL/assertoor/pkg/coordinator/names"
	"github.com/theQRL/assertoor/pkg/coordinator/wallet"
)

type Coordinator interface {
	Logger() logrus.FieldLogger
	LogReader() logger.LogReader
	Database() *db.Database
	ClientPool() *clients.ClientPool
	WalletManager() *wallet.Manager
	ValidatorNames() *names.ValidatorNames
	GlobalVariables() Variables
	TestRegistry() TestRegistry

	GetTestByRunID(runID uint64) Test
	GetTestQueue() []Test
	GetTestHistory(testID string, firstRunID uint64, offset uint64, limit uint64) ([]Test, uint64)
	ScheduleTest(descriptor TestDescriptor, configOverrides map[string]any, allowDuplicate bool, skipQueue bool) (TestRunner, error)
	DeleteTestRun(runID uint64) error
}

type TestRegistry interface {
	AddLocalTest(testConfig *TestConfig) (TestDescriptor, error)
	AddExternalTest(ctx context.Context, extTestConfig *ExternalTestConfig) (TestDescriptor, error)
	DeleteTest(testID string) error
	GetTestDescriptors() []TestDescriptor
}
