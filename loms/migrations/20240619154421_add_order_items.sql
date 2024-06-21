-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_items
(
  order_id INTEGER           NOT NULL
    CONSTRAINT order_items_orders_id_fk
      REFERENCES orders,
  item_id  INTEGER           NOT NULL
    CONSTRAINT order_items_stocks_items_sku_fk
      REFERENCES stocks_items,
  count    INTEGER DEFAULT 0 NOT NULL,
  status   TEXT              NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_items;
-- +goose StatementEnd
