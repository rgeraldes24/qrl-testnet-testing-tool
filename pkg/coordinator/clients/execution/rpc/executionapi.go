package rpc

import (
	"context"
	"fmt"
	"math/big"
	"time"

	qrl "github.com/theQRL/go-zond"
	"github.com/theQRL/go-zond/common"
	"github.com/theQRL/go-zond/core/types"
	"github.com/theQRL/go-zond/qrlclient"
	"github.com/theQRL/go-zond/rpc"
)

type ExecutionClient struct {
	name             string
	endpoint         string
	headers          map[string]string
	rpcClient        *rpc.Client
	qrlClient        *qrlclient.Client
	concurrencyLimit int
	requestTimeout   time.Duration
	concurrencyChan  chan struct{}
}

// NewExecutionClient is used to create a new execution client
func NewExecutionClient(name, url string, headers map[string]string) (*ExecutionClient, error) {
	client := &ExecutionClient{
		name:             name,
		endpoint:         url,
		headers:          headers,
		concurrencyLimit: 50,
		requestTimeout:   30 * time.Second,
	}

	client.concurrencyChan = make(chan struct{}, client.concurrencyLimit)

	return client, nil
}

func (ec *ExecutionClient) Initialize(ctx context.Context) error {
	if ec.qrlClient != nil {
		return nil
	}

	rpcClient, err := rpc.DialContext(ctx, ec.endpoint)
	if err != nil {
		return err
	}

	for hKey, hVal := range ec.headers {
		rpcClient.SetHeader(hKey, hVal)
	}

	ec.rpcClient = rpcClient
	ec.qrlClient = qrlclient.NewClient(rpcClient)

	return nil
}

func (ec *ExecutionClient) enforceConcurrencyLimit(ctx context.Context) func() {
	select {
	case <-ctx.Done():
		return func() {}
	case ec.concurrencyChan <- struct{}{}:
		return func() {
			<-ec.concurrencyChan
		}
	}
}

func (ec *ExecutionClient) GetQRLClient() *qrlclient.Client {
	return ec.qrlClient
}

func (ec *ExecutionClient) GetClientVersion(ctx context.Context) (string, error) {
	var result string

	err := ec.rpcClient.CallContext(ctx, &result, "web3_clientVersion")

	return result, err
}

func (ec *ExecutionClient) GetChainSpec(ctx context.Context) (*ChainSpec, error) {
	chainID, err := ec.qrlClient.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	return &ChainSpec{
		ChainID: chainID.String(),
	}, nil
}

func (ec *ExecutionClient) GetNodeSyncing(ctx context.Context) (*SyncStatus, error) {
	status, err := ec.qrlClient.SyncProgress(ctx)
	if err != nil {
		return nil, err
	}

	if status == nil {
		// Not syncing
		ss := &SyncStatus{}
		ss.IsSyncing = false

		return ss, nil
	}

	return &SyncStatus{
		IsSyncing:     true,
		CurrentBlock:  status.CurrentBlock,
		HighestBlock:  status.HighestBlock,
		StartingBlock: status.StartingBlock,
	}, nil
}

func (ec *ExecutionClient) GetLatestBlock(ctx context.Context) (*types.Block, error) {
	reqCtx, reqCtxCancel := context.WithTimeout(ctx, ec.requestTimeout)
	defer reqCtxCancel()

	block, err := ec.qrlClient.BlockByNumber(reqCtx, nil)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (ec *ExecutionClient) GetBlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	reqCtx, reqCtxCancel := context.WithTimeout(ctx, ec.requestTimeout)
	defer reqCtxCancel()

	block, err := ec.qrlClient.BlockByHash(reqCtx, hash)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (ec *ExecutionClient) GetNonceAt(ctx context.Context, wallet common.Address, blockNumber *big.Int) (uint64, error) {
	closeFn := ec.enforceConcurrencyLimit(ctx)
	if closeFn == nil {
		return 0, fmt.Errorf("client busy")
	}

	defer closeFn()

	reqCtx, reqCtxCancel := context.WithTimeout(ctx, ec.requestTimeout)
	defer reqCtxCancel()

	return ec.qrlClient.NonceAt(reqCtx, wallet, blockNumber)
}

func (ec *ExecutionClient) GetBalanceAt(ctx context.Context, wallet common.Address, blockNumber *big.Int) (*big.Int, error) {
	closeFn := ec.enforceConcurrencyLimit(ctx)
	if closeFn == nil {
		return nil, fmt.Errorf("client busy")
	}

	defer closeFn()

	reqCtx, reqCtxCancel := context.WithTimeout(ctx, ec.requestTimeout)
	defer reqCtxCancel()

	return ec.qrlClient.BalanceAt(reqCtx, wallet, blockNumber)
}

func (ec *ExecutionClient) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	closeFn := ec.enforceConcurrencyLimit(ctx)
	if closeFn == nil {
		return nil, fmt.Errorf("client busy")
	}

	defer closeFn()

	reqCtx, reqCtxCancel := context.WithTimeout(ctx, ec.requestTimeout)
	defer reqCtxCancel()

	return ec.qrlClient.TransactionReceipt(reqCtx, txHash)
}

func (ec *ExecutionClient) GetBlockReceipts(ctx context.Context, blockHash common.Hash) ([]*types.Receipt, error) {
	closeFn := ec.enforceConcurrencyLimit(ctx)
	if closeFn == nil {
		return nil, fmt.Errorf("client busy")
	}

	defer closeFn()

	reqCtx, reqCtxCancel := context.WithTimeout(ctx, ec.requestTimeout)
	defer reqCtxCancel()

	return ec.qrlClient.BlockReceipts(reqCtx, rpc.BlockNumberOrHash{
		BlockHash: &blockHash,
	})
}

func (ec *ExecutionClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	closeFn := ec.enforceConcurrencyLimit(ctx)
	if closeFn == nil {
		return fmt.Errorf("client busy")
	}

	defer closeFn()

	reqCtx, reqCtxCancel := context.WithTimeout(ctx, ec.requestTimeout)
	defer reqCtxCancel()

	return ec.qrlClient.SendTransaction(reqCtx, tx)
}

func (ec *ExecutionClient) GetQRLCall(ctx context.Context, msg *qrl.CallMsg, blockNumber *big.Int) ([]byte, error) {
	closeFn := ec.enforceConcurrencyLimit(ctx)
	if closeFn == nil {
		return nil, fmt.Errorf("client busy")
	}

	defer closeFn()

	return ec.qrlClient.CallContract(ctx, *msg, blockNumber)
}
