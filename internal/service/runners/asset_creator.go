package runners

import (
	"blobs/internal/config"
	"blobs/internal/database"
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/connectors/submit"
	"gitlab.com/tokend/go/xdrbuild"

	"blobs/internal/data"
	"blobs/internal/service/runners/helpers"
	"blobs/internal/types"
)

type AssetCreator struct {
	log *logan.Entry
	cfg config.Config

	assets  chan data.Asset
	assetsQ data.Assets
}

func NewAssetCreator(cfg config.Config) *AssetCreator {
	return &AssetCreator{
		log:     cfg.Log(),
		cfg:     cfg,
		assets:  make(chan data.Asset),
		assetsQ: database.NewAssetsQ(cfg.DB()),
	}
}

func (c *AssetCreator) Run(ctx context.Context) {
	c.log.Info("asset creator started")
	go running.WithBackOff(ctx, c.log, "selector",
		c.selector, 20*time.Second, 30*time.Second, time.Minute)
	go running.WithBackOff(ctx, c.log, "receiver",
		c.receiver, 20*time.Second, 30*time.Second, time.Minute)
}

func (c *AssetCreator) selector(_ context.Context) error {
	c.log.Info("selecting assets to create")
	pendingAssets, err := c.assetsQ.
		New().
		FilterByStatus(types.PendingCreation).
		Select()
	if err != nil {
		return errors.Wrap(err, "failed to select assets to create")
	}

	for _, pending := range pendingAssets {
		c.log.
			WithFields(logan.F{"asset_code": pending.AssetCode}).
			Info("prepared for creating")

		c.assets <- pending
	}

	return nil
}

func (c *AssetCreator) receiver(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case pending := <-c.assets:
			c.log.Info("processing ", pending)
			if err := c.create(ctx, pending); err != nil {
				return errors.Wrap(err, "failed to create asset")
			}
		default:
			<-ticker.C
		}
	}
}

func (c *AssetCreator) create(ctx context.Context, pending data.Asset) error {
	envelope, err := c.composeEnvelope(pending)
	if err != nil {
		return errors.Wrap(err, "failed to compose envelope")
	}

	db := c.cfg.DB()
	err = db.Transaction(func() error {
		err = database.NewAssetsQ(db).UpdateStatus(types.Created, pending.AssetCode)
		if err != nil {
			return errors.Wrap(err, "failed to update pending creation status")
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

func (c *AssetCreator) composeEnvelope(pending data.Asset) (string, error) {
	tx := c.cfg.Builder().Transaction(c.cfg.Keys().Source)
	op, err := c.populateCreateAssetOp(pending)
	if err != nil {
		return "", errors.Wrap(err, "failed to populate create asset operation")
	}
	tx = tx.Op(op)
	tx.Sign(c.cfg.Keys().Signer)

	envelope, err := tx.Marshal()
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal tx")
	}

	return envelope, nil
}

func (c *AssetCreator) populateCreateAssetOp(pending data.Asset) (*xdrbuild.CreateAsset, error) {
	details, err := c.populateDetails(pending)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate details")
	}

	operation := &xdrbuild.CreateAsset{
		Code:                     string(pending.AssetCode),
		MaxIssuanceAmount:        helpers.MaxIssuanceAmount,
		PreIssuanceSigner:        c.cfg.Keys().Signer.Address(),
		InitialPreIssuanceAmount: helpers.MaxIssuanceAmount,
		TrailingDigitsCount:      helpers.Decimals,
		Policies:                 helpers.Policy,
		Type:                     helpers.OrdinaryAssetType,
		CreatorDetails:           details,
		AllTasks:                 &helpers.ZeroTasks,
	}
	return operation, nil
}

func (c *AssetCreator) populateDetails(pending data.Asset) (json.RawMessage, error) {
	assetDetails := helpers.AssetDetails{
		Name:  string(pending.AssetCode),
		Owner: string(pending.Creator),
	}
	bb, err := json.Marshal(assetDetails)
	if err != nil {
		return nil, err
	}

	return bb, nil
}
