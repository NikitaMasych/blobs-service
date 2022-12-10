package database

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"blobs/internal/data"
	"blobs/internal/types"
)

const (
	AssetsTable        = "assets"
	AssetsPKConstraint = "assets_pkey"
)

var (
	ErrAssetsConflict = errors.New("assets primary key conflict")
	AssetsColumns     = []string{"asset_code", "creator", "status"}
)

type AssetsQ struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
}

func NewAssetsQ(db *pgdb.DB) *AssetsQ {
	return &AssetsQ{
		db:       db,
		selector: squirrel.Select(AssetsColumns...).From(AssetsTable),
	}
}

func (p *AssetsQ) New() data.Assets {
	return NewAssetsQ(p.db.Clone())
}

func (p *AssetsQ) Create(asset data.Asset) error {
	query := squirrel.Insert(AssetsTable).
		Columns(AssetsColumns...).Values(asset.AssetCode, asset.Creator, asset.Status)

	if err := p.db.Exec(query); err != nil {
		cause := errors.Cause(err)
		if pqerr, ok := cause.(*pq.Error); ok {
			if pqerr.Constraint == AssetsPKConstraint {
				return ErrAssetsConflict
			}
		}
		return errors.Wrap(err, "failed to create the asset")
	}
	return nil
}

func (p *AssetsQ) UpdateStatus(status types.AssetStatus, assetCode types.AssetCode) error {
	stmt := squirrel.Update(AssetsTable).
		Set("status", status).
		Where(squirrel.Eq{"asset_code": assetCode})

	if err := p.db.Exec(stmt); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return errors.Wrap(err, "failed to update asset(s) status")
	}
	return nil
}

func (p *AssetsQ) FilterByStatus(status types.AssetStatus) data.Assets {
	p.selector = p.selector.Where(squirrel.Eq{"status": status})
	return p
}

func (p *AssetsQ) Select() ([]data.Asset, error) {
	assets := make([]data.Asset, 0)
	if err := p.db.Select(&assets, p.selector); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to select assets")
	}
	return assets, nil
}

func (p *AssetsQ) Get() (*data.Asset, error) {
	var asset data.Asset
	if err := p.db.Get(&asset, p.selector); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get the asset")
	}
	return &asset, nil
}
