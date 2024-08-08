package processor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/grassrootseconomics/celo-tracker/internal/chain"
	"github.com/grassrootseconomics/celo-tracker/internal/router"
)

type (
	ProcessorOpts struct {
		Chain  chain.Chain
		Router router.Router
		Logg   *slog.Logger
	}

	Processor struct {
		chain  chain.Chain
		router router.Router
		logg   *slog.Logger
	}
)

func NewProcessor(o ProcessorOpts) *Processor {
	return &Processor{
		chain: o.Chain,
		logg:  o.Logg,
	}
}

func (p *Processor) ProcessBlock(ctx context.Context, blockNumber uint64) error {
	block, err := p.chain.GetBlock(ctx, blockNumber)
	if err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("block %d error: %v", blockNumber, err)
	}

	receipts, err := p.chain.GetReceipts(ctx, block.Number())
	if err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("receipts fetch error: block %d: %v", blockNumber, err)
	}

	for _, receipt := range receipts {
		if receipt.Status > 0 {
			for _, log := range receipt.Logs {
				if err := p.router.RouteSuccessTx(
					ctx,
					router.SuccessTx{
						Log:       log,
						Timestamp: block.Time(),
					},
				); err != nil && !errors.Is(err, context.Canceled) {
					return err
				}
			}
		} else {
			tx, err := p.chain.GetTransaction(ctx, receipt.TxHash)
			if err != nil && !errors.Is(err, context.Canceled) {
				return fmt.Errorf("get transaction error: tx %s: %v", receipt.TxHash.Hex(), err)
			}

			if tx.To() == nil {
				return nil
			}

			from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
			if err != nil {
				return fmt.Errorf("transaction decode error: tx %s: %v", receipt.TxHash.Hex(), err)
			}

			if err := p.router.RouteRevertTx(
				ctx,
				router.RevertTx{
					From:            from.Hex(),
					InputData:       common.Bytes2Hex(tx.Data()),
					Block:           blockNumber,
					ContractAddress: tx.To().Hex(),
					Timestamp:       block.Time(),
					TxHash:          receipt.TxHash.Hex(),
				},
			); err != nil && !errors.Is(err, context.Canceled) {
				return err
			}
		}
	}

	return nil
}
