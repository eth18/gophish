
-- +goose Up
ALTER TABLE `attachments` ADD COLUMN `tenant_id` INT AFTER `id`;
ALTER TABLE `campaigns` ADD COLUMN `tenant_id` INT AFTER `id`;
ALTER TABLE `email_requests` ADD COLUMN `tenant_id` INT AFTER `id`;
ALTER TABLE `events` ADD COLUMN `tenant_id` INT AFTER `id`;
ALTER TABLE `groups` ADD COLUMN `tenant_id` INT AFTER `id`;
ALTER TABLE `headers` ADD COLUMN `tenant_id` INT AFTER `id`;
ALTER TABLE `imap` ADD COLUMN `tenant_id` INT AFTER `user_id`;
ALTER TABLE `mail_logs` ADD COLUMN `tenant_id` INT AFTER `id`;
ALTER TABLE `pages` ADD COLUMN `tenant_id` INT AFTER `id`;
ALTER TABLE `permissions` ADD COLUMN `tenant_id` INT AFTER `id`;
ALTER TABLE `results` ADD COLUMN `tenant_id` INT AFTER `id`;
ALTER TABLE `roles` ADD COLUMN `tenant_id` INT AFTER `id`;
ALTER TABLE `smtp` ADD COLUMN `tenant_id` INT AFTER `id`;
ALTER TABLE `targets` ADD COLUMN `tenant_id` INT AFTER `id`;
ALTER TABLE `templates` ADD COLUMN `tenant_id` INT AFTER `id`;
ALTER TABLE `users` ADD COLUMN `tenant_id` INT AFTER `id`;
ALTER TABLE `webhooks` ADD COLUMN `tenant_id` INT AFTER `id`;

-- +goose Down
ALTER TABLE `attachments` DROP COLUMN `tenant_id`;
ALTER TABLE `campaigns` DROP COLUMN `tenant_id`;
ALTER TABLE `email_requests` DROP COLUMN `tenant_id`;
ALTER TABLE `events` DROP COLUMN `tenant_id`;
ALTER TABLE `groups` DROP COLUMN `tenant_id`;
ALTER TABLE `headers` DROP COLUMN `tenant_id`;
ALTER TABLE `imap` DROP COLUMN `tenant_id`;
ALTER TABLE `mail_logs` DROP COLUMN `tenant_id`;
ALTER TABLE `pages` DROP COLUMN `tenant_id`;
ALTER TABLE `permissions` DROP COLUMN `tenant_id`;
ALTER TABLE `results` DROP COLUMN `tenant_id`;
ALTER TABLE `roles` DROP COLUMN `tenant_id`;
ALTER TABLE `smtp` DROP COLUMN `tenant_id`;
ALTER TABLE `targets` DROP COLUMN `tenant_id`;
ALTER TABLE `templates` DROP COLUMN `tenant_id`;
ALTER TABLE `tenants` DROP COLUMN `tenant_id`;
ALTER TABLE `users` DROP COLUMN `tenant_id`;
ALTER TABLE `webhooks` DROP COLUMN `tenant_id`;
