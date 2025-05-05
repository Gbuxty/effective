-- +goose Up
ALTER TABLE persons ALTER COLUMN age TYPE INTEGER USING (age::integer);

-- +goose Down
ALTER TABLE persons ALTER COLUMN age TYPE VARCHAR(255) USING (age::varchar);