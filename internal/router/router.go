package router

import (
	"context"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/core/types"
)

type (
	SuccessTx struct {
		Log       *types.Log
		Timestamp uint64
	}

	RevertTx struct {
		From            string
		InputData       string
		Block           uint64
		ContractAddress string
		Timestamp       uint64
		TxHash          string
	}

	Handler interface {
		SuccessTx(context.Context, SuccessTx) error
		RevertTx(context.Context, RevertTx) error
	}

	Router struct {
		logHandlers       map[common.Hash]Handler
		inputDataHandlers map[string]Handler
	}
)

func New() *Router {
	logHandlers := map[common.Hash]Handler{}

	inputDataNadlers := map[string]Handler{}

	return &Router{
		logHandlers:       logHandlers,
		inputDataHandlers: inputDataNadlers,
	}
}

func (r *Router) RouteSuccessTx(ctx context.Context, msg SuccessTx) error {
	handler, ok := r.logHandlers[msg.Log.Topics[0]]
	if ok {
		return handler.SuccessTx(ctx, msg)
	}

	return nil
}

func (r *Router) RouteRevertTx(ctx context.Context, msg RevertTx) error {
	handler, ok := r.inputDataHandlers[msg.InputData[:8]]
	if ok {
		return handler.RevertTx(ctx, msg)
	}

	return nil
}
