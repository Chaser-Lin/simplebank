CREATE TABLE `accounts` (
    `id` bigint PRIMARY KEY AUTO_INCREMENT,
    `owner` varchar(255) NOT NULL,
    `balance` double NOT NULL,
    `currency` varchar(255) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT now()
);

CREATE TABLE `entries` (
   `id` bigint PRIMARY KEY AUTO_INCREMENT,
   `account_id` bigint NOT NULL,
   `amount` double NOT NULL COMMENT 'can be positive or negative',
   `created_at` timestamp NOT NULL DEFAULT now()
);

CREATE TABLE `transfers` (
     `id` bigint PRIMARY KEY AUTO_INCREMENT,
     `from_account_id` bigint NOT NULL,
     `to_account_id` bigint NOT NULL,
     `amount` double NOT NULL COMMENT 'must be positive',
     `created_at` timestamp NOT NULL DEFAULT now()
);

CREATE INDEX `accounts_index_0` ON `accounts` (`owner`);

CREATE INDEX `entries_index_1` ON `entries` (`account_id`);

CREATE INDEX `transfers_index_2` ON `transfers` (`from_account_id`);

CREATE INDEX `transfers_index_3` ON `transfers` (`from_account_id`, `to_account_id`);

-- ALTER TABLE `entries` ADD FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`);
--
-- ALTER TABLE `transfers` ADD FOREIGN KEY (`from_account_id`) REFERENCES `accounts` (`id`);
--
-- ALTER TABLE `transfers` ADD FOREIGN KEY (`to_account_id`) REFERENCES `accounts` (`id`);
