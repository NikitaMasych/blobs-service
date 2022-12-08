package database

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"blobs/internal/data"
	"blobs/internal/types"
)

const pendingAssetsTable = "pending_assets"

var pendingAssetsColumns = []string{"asset_code", "tx_id", "creator", "status"}

type PendingAssetsQ struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
	limit    uint64
}

func NewPendingAssetsQ(db *pgdb.DB) *PendingAssetsQ {
	return &PendingAssetsQ{
		db:       db,
		selector: squirrel.Select(pendingAssetsColumns...).From(pendingAssetsTable),
		limit:    15,
	}
}

func (p *PendingAssetsQ) New() data.PendingAssets {
	return NewPendingAssetsQ(p.db.Clone())
}

func (p *PendingAssetsQ) Create(asset data.PendingAsset) error {
	query := squirrel.Insert(pendingAssetsTable).
		Columns(pendingAssetsColumns...).Values(asset.AssetCode, asset.TxId, asset.Creator, asset.Status)

	err := p.db.Exec(query)
	if err != nil {
		return errors.Wrap(err, "failed to execute insert query for pending_tokens table")
	}

	return nil
}

func (p *PendingAssetsQ) UpdateStatus(status types.PendingAssetStatus, txId string) error {
	stmt := squirrel.Update(pendingAssetsTable).
		Set("status", status).
		Where(squirrel.Eq{"tx_id": txId})

	err := p.db.Exec(stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return errors.Wrap(err, "failed to update status")
	}

	return nil
}

func (p *PendingAssetsQ) FilterByStatus(status types.PendingAssetStatus) data.PendingAssets {
	p.selector = p.selector.Where(squirrel.Eq{"status": status})
	return p
}

func (p *PendingAssetsQ) FilterByTxId(txId string) data.PendingAssets {
	p.selector = p.selector.Where(squirrel.Eq{"tx_id": txId})
	return p
}

func (p *PendingAssetsQ) Select() ([]data.PendingAsset, error) {
	result := make([]data.PendingAsset, 0, p.limit)
	err := p.db.Select(&result, p.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select pending tokens")
	}

	return result, nil
}

func (p *PendingAssetsQ) Get() (*data.PendingAsset, error) {
	var result data.PendingAsset
	err := p.db.Get(&result, p.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select pending tokens")
	}

	return &result, nil
}
