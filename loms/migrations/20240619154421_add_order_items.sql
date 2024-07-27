-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_items
(
  order_id INTEGER           NOT NULL,
  item_id  INTEGER           NOT NULL,
  count    INTEGER DEFAULT 0 NOT NULL,
  status   TEXT              NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_items;
-- +goose StatementEnd
