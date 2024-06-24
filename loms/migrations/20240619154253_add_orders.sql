-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders
(
  id      serial
    CONSTRAINT orders_pk
      PRIMARY KEY,
  user_id INTEGER NOT NULL,
  status TEXT DEFAULT 'new' NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
