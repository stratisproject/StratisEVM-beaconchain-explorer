-- +goose Up
-- +goose StatementBegin
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_blocks_status1_blockroot ON blocks (blockroot) WHERE status = '1';
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_blocks_withdrawals_validator_blockroot_incl_slot ON blocks_withdrawals (validatorindex, block_root) INCLUDE (block_slot);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_blocks_proposer_slot_incl ON blocks (proposer, slot) INCLUDE (status, exec_block_number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX CONCURRENTLY IF EXISTS idx_blocks_status1_blockroot;
DROP INDEX CONCURRENTLY IF NOT EXISTS idx_blocks_withdrawals_validator_blockroot_incl_slot;
DROP INDEX CONCURRENTLY IF NOT EXISTS idx_blocks_proposer_slot_incl;
-- +goose StatementEnd
