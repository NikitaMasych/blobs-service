-- noinspection SqlNoDataSourceInspectionForFile

-- +migrate Up
CREATE TABLE IF NOT EXISTS blobs(
    blobs_pkey SERIAL PRIMARY KEY,
    id VARCHAR(100),
    value VARCHAR(100),
    type INTEGER
);
-- +migrate Down

DROP TABLE blobs CASCADE;