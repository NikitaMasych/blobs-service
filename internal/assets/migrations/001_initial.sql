-- noinspection SqlNoDataSourceInspectionForFile

-- +migrate Up

CREATE TABLE IF NOT EXISTS blobs(
    blobs_pkey SERIAL PRIMARY KEY,
    id TEXT UNIQUE,
    value TEXT,
    type INTEGER
);

CREATE TABLE IF NOT EXISTS assets(
    assets_pkey SERIAL PRIMARY KEY,
    asset_code TEXT,
    creator TEXT,
    status INTEGER
);

-- +migrate Down

DROP TABLE blobs CASCADE;
DROP TABLE assets CASCADE;