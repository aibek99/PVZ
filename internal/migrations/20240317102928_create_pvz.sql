-- +goose Up
-- +goose StatementBegin
Create TABLE pvz(
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    address TEXT DEFAULT '',
    contact TEXT,
    deleted_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE box(
                    id BIGSERIAL PRIMARY KEY,
                    name TEXT UNIQUE NOT NULL,
                    cost NUMERIC(12,2) DEFAULT 0,
                    is_check BOOLEAN DEFAULT FALSE,
                    weight NUMERIC(12,2) DEFAULT 0,
                    deleted_at TIMESTAMPTZ,
                    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                    updated_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO box (name, cost, weight, is_check) VALUES
                                                   ('package', 5.0, 10.0, true),
                                                   ('box', 20.0, 30.0, true),
                                                   ('textile', 1.0, 0, false);
CREATE TABLE orders(
    order_id BIGINT UNIQUE PRIMARY KEY,
    weight NUMERIC(12,2) DEFAULT 0,
    box_id BIGINT,
    client_id BIGINT,
    returned_at TIMESTAMPTZ,
    accepted_at TIMESTAMPTZ,
    issued_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (box_id) REFERENCES box(id)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pvz;
DROP TABLE orders;
DROP TABLE box;
-- +goose StatementEnd

