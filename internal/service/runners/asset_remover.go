package runners

import (
	"blobs/internal/config"
	"blobs/internal/database"
	"context"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/connectors/submit"
	"gitlab.com/tokend/go/xdrbuild"

	"blobs/internal/data"
	"blobs/internal/types"
)

type AssetRemover struct {
	log *logan.Entry
	cfg config.Config

	assets  chan data.Asset
	assetsQ data.Assets
}

func NewAssetRemover(cfg config.Config) *AssetRemover {
	return &AssetRemover{
		log:     cfg.Log(),
		cfg:     cfg,
		assets:  make(chan data.Asset),
		assetsQ: database.NewAssetsQ(cfg.DB()),
	}
}

func (r *AssetRemover) Run(ctx context.Context) {
	r.log.Info("asset remover started")
	go running.WithBackOff(ctx, r.log, "selector",
		r.selector, 20*time.Second, 30*time.Second, time.Minute)
	go running.WithBackOff(ctx, r.log, "receiver",
		r.receiver, 20*time.Second, 30*time.Second, time.Minute)
}

func (r *AssetRemover) selector(_ context.Context) error {
	r.log.Info("selecting assets to remove")
	pendingAssets, err := r.assetsQ.
		New().
		FilterByStatus(types.PendingRemoval).
		Select()
	if err != nil {
		return errors.Wrap(err, "failed to select assets to remove")
	}

	for _, pending := range pendingAssets {
		r.log.
			WithFields(logan.F{"asset_code": pending.AssetCode}).
			Info("prepared for removal")

		r.assets <- pending
	}

	return nil
}

func (r *AssetRemover) receiver(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case pending := <-r.assets:
			r.log.Info("processing ", pending)
			if err := r.remove(ctx, pending); err != nil {
				return errors.Wrap(err, "failed to remove asset")
			}
		default:
			<-ticker.C
		}
	}
}

func (r *AssetRemover) remove(ctx context.Context, pending data.Asset) error {
	envelope, err := r.composeEnvelope(pending)
	if err != nil {
		return errors.Wrap(err, "failed to compose envelope")
	}

	db := r.cfg.DB()
	err = db.Transaction(func() error {
		err = database.NewAssetsQ(db).UpdateStatus(types.Removed, pending.AssetCode)
		if err != nil {
			return errors.Wrap(err, "failed to update pending removal status")
		}

		_, err = r.cfg.Submit().Submit(ctx, envelope, true, true)
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

	r.log.WithFields(logan.F{"asset": pending.AssetCode}).Info("removal finished")

	return nil
}

func (r *AssetRemover) composeEnvelope(pending data.Asset) (string, error) {
	tx := r.cfg.Builder().Transaction(r.cfg.Keys().Source)
	op := &xdrbuild.RemoveAsset{
		Code: string(pending.AssetCode),
	}
	tx = tx.Op(op)
	tx.Sign(r.cfg.Keys().Signer)

	envelope, err := tx.Marshal()
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal tx")
	}
	return envelope, nil
}
