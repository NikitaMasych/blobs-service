-- noinspection SqlNoDataSourceInspectionForFile

-- +migrate Down

DROP TABLE blobs CASCADE;
DROP TABLE assets CASCADE;