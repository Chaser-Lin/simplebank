CREATE TABLE `sessions` (
    `id` varchar(255) PRIMARY KEY,
    `username` varchar(255) NOT NULL,
    `refresh_token` varchar(512) NOT NULL,
    `user_agent` varchar(255) NOT NULL,
    `client_ip` varchar(255) NOT NULL,
    `is_blocked` boolean NOT NULL DEFAULT false,
    `expired_at` timestamp NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT now()
);

ALTER TABLE `sessions` ADD FOREIGN KEY (`username`) REFERENCES `users` (`username`);
