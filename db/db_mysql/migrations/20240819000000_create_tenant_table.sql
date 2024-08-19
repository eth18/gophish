
-- +goose Up
CREATE TABLE IF NOT EXISTS `tenants` (
    tenant_id INT AUTO_INCREMENT PRIMARY KEY,
    guid VARCHAR(36) NOT NULL UNIQUE,
    tenant_name VARCHAR(255) NOT NULL,
    tenant_identifier VARCHAR(255) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE `tenant`;
