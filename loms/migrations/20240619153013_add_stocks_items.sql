-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stocks_items
(
  sku         INTEGER NOT NULL
    CONSTRAINT stocks_items_pk
      PRIMARY KEY,
  total_count INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stocks_items;
-- +goose StatementEnd
