CREATE TABLE `sessions`(
  `id` integer PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime,
  `updated_at` datetime,
  `deleted_at` datetime,
  `session_type` text,
  `exit_code` integer,
  `error` text
);
CREATE TABLE sqlite_sequence(name,seq);
CREATE INDEX `idx_sessions_deleted_at` ON `sessions`(`deleted_at`);
CREATE TABLE `languages`(
  `id` integer PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime,
  `updated_at` datetime,
  `deleted_at` datetime,
  `name` text NOT NULL,
  CONSTRAINT `uni_languages_name` UNIQUE(`name`)
);
CREATE INDEX `idx_languages_deleted_at` ON `languages`(`deleted_at`);
CREATE TABLE `language_data`(
  `id` integer PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime,
  `updated_at` datetime,
  `deleted_at` datetime,
  `name` text NOT NULL,
  `compiled` numeric,
  `static_type` numeric,
  `object_oriented` numeric,
  `functional` numeric,
  `style` text,
  `language_id` integer,
  CONSTRAINT `fk_language_data_language` FOREIGN KEY(`language_id`) REFERENCES `languages`(`id`),
  CONSTRAINT `uni_language_data_name` UNIQUE(`name`)
);
CREATE INDEX `idx_language_data_language_id` ON `language_data`(`language_id`);
CREATE INDEX `idx_language_data_deleted_at` ON `language_data`(`deleted_at`);
CREATE TABLE `runtimes`(
  `id` integer PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime,
  `updated_at` datetime,
  `deleted_at` datetime,
  `filename` text NOT NULL,
  `language_id` integer NOT NULL,
  CONSTRAINT `fk_languages_runtimes` FOREIGN KEY(`language_id`) REFERENCES `languages`(`id`)
);
CREATE UNIQUE INDEX `idx_lang_file` ON `runtimes`(`filename`,`language_id`);
CREATE INDEX `idx_runtimes_deleted_at` ON `runtimes`(`deleted_at`);
CREATE TABLE `runtime_data`(
  `id` integer PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime,
  `updated_at` datetime,
  `deleted_at` datetime,
  `github_url` text NOT NULL,
  `runtime_id` integer NOT NULL,
  CONSTRAINT `fk_runtime_data_runtime` FOREIGN KEY(`runtime_id`) REFERENCES `runtimes`(`id`),
  CONSTRAINT `uni_runtime_data_github_url` UNIQUE(`github_url`)
);
CREATE INDEX `idx_runtime_data_deleted_at` ON `runtime_data`(`deleted_at`);
CREATE TABLE `github_stars`(
  `id` integer PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime,
  `updated_at` datetime,
  `deleted_at` datetime,
  `star_count` integer NOT NULL,
  `runtime_data_id` integer NOT NULL,
  CONSTRAINT `fk_runtime_data_stars` FOREIGN KEY(`runtime_data_id`) REFERENCES `runtime_data`(`id`)
);
CREATE INDEX `idx_github_stars_deleted_at` ON `github_stars`(`deleted_at`);
CREATE TABLE `docker_images`(
  `id` integer PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime,
  `updated_at` datetime,
  `deleted_at` datetime,
  `tag` text NOT NULL,
  `runtime_id` integer NOT NULL,
  CONSTRAINT `fk_runtimes_docker_image` FOREIGN KEY(`runtime_id`) REFERENCES `runtimes`(`id`),
  CONSTRAINT `uni_docker_images_tag` UNIQUE(`tag`),
  CONSTRAINT `uni_docker_images_runtime_id` UNIQUE(`runtime_id`)
);
CREATE INDEX `idx_docker_images_deleted_at` ON `docker_images`(`deleted_at`);
CREATE TABLE `image_sizes`(
  `id` integer PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime,
  `updated_at` datetime,
  `deleted_at` datetime,
  `size` text NOT NULL,
  `image_id` integer NOT NULL,
  CONSTRAINT `fk_docker_images_image_size` FOREIGN KEY(`image_id`) REFERENCES `docker_images`(`id`),
  CONSTRAINT `uni_image_sizes_image_id` UNIQUE(`image_id`)
);
CREATE INDEX `idx_image_sizes_deleted_at` ON `image_sizes`(`deleted_at`);
CREATE TABLE `container_runs`(
  `id` integer PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime,
  `updated_at` datetime,
  `deleted_at` datetime,
  `tag` text NOT NULL,
  `iterations` integer NOT NULL,
  `sample_size` integer NOT NULL,
  `percent` real NOT NULL,
  `seconds` real NOT NULL,
  `iterations_per_second` integer NOT NULL,
  `image_id` integer NOT NULL,
  CONSTRAINT `fk_docker_images_container_runs` FOREIGN KEY(`image_id`) REFERENCES `docker_images`(`id`)
);
CREATE INDEX `idx_container_runs_deleted_at` ON `container_runs`(`deleted_at`);
