CREATE TABLE `users` (
     `username` varchar(255) PRIMARY KEY,
     `hashed_password` varchar(255) NOT NULL,
     `full_name` varchar(255) NOT NULL,
     `email` varchar(255) NOT NULL UNIQUE,
     `password_changed_at` timestamp NOT NULL DEFAULT "1970-01-01 01:01:01",
     `created_at` timestamp NOT NULL DEFAULT now()
);

ALTER TABLE `accounts` ADD FOREIGN KEY (`owner`) REFERENCES `users` (`username`);

CREATE UNIQUE INDEX `accounts_index_1` ON `accounts` (`owner`, `currency`);
-- ALTER TABLE `accounts` ADD CONSTRAINT `owner_currency_key` UNIQUE (`owner`, `currency`);
-- 第二条效果与第一条相同
