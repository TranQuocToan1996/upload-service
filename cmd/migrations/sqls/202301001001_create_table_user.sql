-- +goose Up
-- +goose StatementBegin

CREATE TABLE `users` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_name` varchar(255) DEFAULT "",
  `password` varchar(255) DEFAULT "",
  `salt` varchar(255) DEFAULT "",
  `revoke_token_at` BIGINT DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
-- +goose StatementEnd