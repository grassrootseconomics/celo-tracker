package handler

import (
	"context"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/grassrootseconomics/celo-tracker/internal/cache"
	"github.com/grassrootseconomics/celo-tracker/internal/pub"
	"github.com/grassrootseconomics/celo-tracker/pkg/event"
	"github.com/grassrootseconomics/w3-celo"
)

type indexRemoveHandler struct {
	pub   pub.Pub
	cache cache.Cache
}

const indexRemoveEventName = "INDEX_REMOVE"

var (
	indexRemoveTopicHash = w3.H("0x24a12366c02e13fe4a9e03d86a8952e85bb74a456c16e4a18b6d8295700b74bb")
	indexRemoveEvent     = w3.MustNewEvent("AddressRemoved(address _token)")
	indexRemoveSig       = w3.MustNewFunc("remove(address)", "bool")
)

func NewIndexRemoveHandler(pub pub.Pub, cache cache.Cache) *indexRemoveHandler {
	return &indexRemoveHandler{
		pub:   pub,
		cache: cache,
	}
}

func (h *indexRemoveHandler) Name() string {
	return indexRemoveEventName
}

func (h *indexRemoveHandler) HandleLog(ctx context.Context, msg LogMessage) error {
	if msg.Log.Topics[0] == indexRemoveTopicHash {
		var address common.Address

		if err := indexRemoveEvent.DecodeArgs(msg.Log, &address); err != nil {
			return err
		}

		indexRemoveEvent := event.Event{
			Index:           msg.Log.Index,
			Block:           msg.Log.BlockNumber,
			ContractAddress: msg.Log.Address.Hex(),
			Success:         true,
			Timestamp:       msg.Timestamp,
			TxHash:          msg.Log.TxHash.Hex(),
			TxType:          indexRemoveEventName,
			Payload: map[string]any{
				"address": address.Hex(),
			},
		}

		if h.cache.IsWatchableIndex(address.Hex()) {
			h.cache.Remove(address.Hex())
		}

		return h.pub.Send(ctx, indexRemoveEvent)
	}

	return nil
}

func (h *indexRemoveHandler) HandleRevert(ctx context.Context, msg RevertMessage) error {
	if len(msg.InputData) < 8 {
		return nil
	}

	switch msg.InputData[:8] {
	case "29092d0e":
		var address common.Address

		if err := indexRemoveSig.DecodeArgs(w3.B(msg.InputData), &address); err != nil {
			return err
		}

		indexRemoveEvent := event.Event{
			Block:           msg.Block,
			ContractAddress: msg.ContractAddress,
			Success:         false,
			Timestamp:       msg.Timestamp,
			TxHash:          msg.TxHash,
			TxType:          indexRemoveEventName,
			Payload: map[string]any{
				"revertReason": msg.RevertReason,
				"address":      address.Hex(),
			},
		}

		return h.pub.Send(ctx, indexRemoveEvent)
	}

	return nil
}
