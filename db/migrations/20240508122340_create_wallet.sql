-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets (
                                     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                     user_id UUID NOT NULL,
                                     Balance BIGINT NOT NULL,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     transaction_id UUID ,
                                     FOREIGN KEY (user_id) REFERENCES users(id),
                                     FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE wallets;
-- +goose StatementEnd
