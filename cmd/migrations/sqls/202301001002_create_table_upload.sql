-- +goose Up
-- +goose StatementBegin

CREATE TABLE `uploads` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `source_file_name` varchar(255) DEFAULT "",
  `destination_file_name` varchar(255) DEFAULT "",
  `content_type` varchar(255) DEFAULT "",
  `size` BIGINT DEFAULT 0,
  `user_id` BIGINT NOT NULL,
  `status` varchar(255) DEFAULT "",
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
-- +goose StatementEnd