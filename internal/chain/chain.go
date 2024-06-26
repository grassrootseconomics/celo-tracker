package chain

import (
	"context"
	"math/big"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/grassrootseconomics/celoutils/v3"
)

type Chain interface {
	GetBlocks(context.Context, []uint64) ([]types.Block, error)
	GetBlock(context.Context, uint64) (*types.Block, error)
	GetLatestBlock(context.Context) (uint64, error)
	GetTransaction(context.Context, common.Hash) (*types.Transaction, error)
	GetReceipts(context.Context, *types.Block) ([]types.Receipt, error)
	GetRevertReason(context.Context, common.Hash, *big.Int) (string, error)
	Provider() *celoutils.Provider
	IsArchiveNode() bool
}
