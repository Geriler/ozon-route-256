-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders
(
  id      serial,
  user_id INTEGER NOT NULL,
  status TEXT DEFAULT 'new' NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE SEQUENCE order_id_manual_seq INCREMENT 10 START 10;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SEQUENCE order_id_manual_seq;

DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
