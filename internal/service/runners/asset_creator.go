package runners

import (
	"blobs/internal/config"
	"blobs/internal/database"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/connectors/submit"
	"gitlab.com/tokend/go/xdrbuild"
	"time"

	"blobs/internal/data"
	"blobs/internal/service/runners/helpers"
	"blobs/internal/types"
)

type AssetCreator struct {
	log *logan.Entry
	cfg config.Config

	pendingAssets chan data.PendingAsset
	pendingQ      data.PendingAssets
}

func NewAssetCreator(cfg config.Config) *AssetCreator {
	return &AssetCreator{
		log:           cfg.Log(),
		cfg:           cfg,
		pendingAssets: make(chan data.PendingAsset),
		pendingQ:      database.NewPendingAssetsQ(cfg.DB()),
	}
}

func (c *AssetCreator) Run(ctx context.Context) {
	c.log.Info("asset creator started")
	go running.WithBackOff(ctx, c.log, "tx_checker",
		c.selector, 20*time.Second, 30*time.Second, time.Minute)
	go running.WithBackOff(ctx, c.log, "asset_creator",
		c.receiver, 20*time.Second, 30*time.Second, time.Minute)
}

func (c *AssetCreator) selector(_ context.Context) error {
	c.log.Info("selecting pending assets")
	pendingAssets, err := c.pendingQ.
		New().
		FilterByStatus(types.Pending).
		Select()
	if err != nil {
		return errors.Wrap(err, "failed to select pending assets")
	}

	for _, pending := range pendingAssets {
		c.log.
			WithFields(logan.F{"tx_id": pending.TxId, "asset_code": pending.AssetCode}).
			Info("prepared for creating")

		c.pendingAssets <- pending
	}

	return nil
}

func (c *AssetCreator) receiver(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil
		case pending := <-c.pendingAssets:

			c.log.Info("processing ", pending)

			err := c.create(ctx, pending)
			if err != nil {
				return errors.Wrap(err, "failed to create asset")
			}
		default:
			<-ticker.C
		}
	}
}

func (c *AssetCreator) create(ctx context.Context, pending data.PendingAsset) error {
	tx := c.cfg.Builder().Transaction(c.cfg.Keys().Source)

	op, err := c.populateCreateAssetOp(pending)
	if err != nil {
		return errors.Wrap(err, "failed to populate create asset operation")
	}

	tx = tx.Op(op)

	tx.Sign(c.cfg.Keys().Signer)

	envelope, err := tx.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal tx")
	}

	db := c.cfg.DB().Clone()
	q := database.NewPendingAssetsQ(db)

	err = db.Transaction(func() error {
		err = q.UpdateStatus(types.Approved, pending.TxId)
		if err != nil {
			return errors.Wrap(err, "failed to update pending status")
		}

		_, err = c.cfg.Submit().Submit(ctx, envelope, true, true)
		if err != nil {
			if txFailed, ok := err.(submit.TxFailure); ok {
				if len(txFailed.OperationResultCodes) != 0 &&
					txFailed.OperationResultCodes[0] == "op_reference_duplication" {
					return errors.Wrap(err, "op_reference_duplication")
				}
				return txFailed
			}
			return err
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "db tx failed")
	}

	c.log.WithFields(logan.F{"asset": pending.AssetCode}).Info("creating finished")

	return nil
}

func (c *AssetCreator) populateCreateAssetOp(pending data.PendingAsset) (*xdrbuild.CreateAsset, error) {
	details, err := c.populateDetails(pending)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate details")
	}

	operation := &xdrbuild.CreateAsset{
		Code:           pending.AssetCode,
		Policies:       helpers.Policy,
		Type:           helpers.OrdinaryAssetType,
		CreatorDetails: details,
		AllTasks:       &helpers.ZeroTasks,
	}

	return operation, nil
}

func (c *AssetCreator) populateDetails(pending data.PendingAsset) (json.RawMessage, error) {
	assetDetails := helpers.AssetDetails{
		Name:          pending.AssetCode,
		ContractOwner: pending.Creator,
	}
	bb, err := json.Marshal(assetDetails)
	if err != nil {
		return nil, err
	}

	return bb, nil
}
