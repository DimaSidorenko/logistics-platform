-- +goose Up
-- +goose StatementBegin
CREATE TABLE order_items (
     order_id BIGINT NOT NULL,
     sku_id BIGINT NOT NULL,
     items_count BIGINT NOT NULL CHECK (items_count > 0),
     PRIMARY KEY (order_id, sku_id),
     CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES orders (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table order_items;
-- +goose StatementEnd
