ALTER TABLE users
    ADD COLUMN h3_index BIGINT,
    ADD COLUMN longitude DECIMAL(9, 6),
    ADD COLUMN latitude DECIMAL(9, 6),
    ADD INDEX idx_h3_index (h3_index);
