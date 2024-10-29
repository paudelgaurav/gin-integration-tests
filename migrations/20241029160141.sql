-- Create "projects" table
CREATE TABLE `projects` (
  `id` integer NULL,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `name` text NULL,
  `endpoint` text NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_projects_deleted_at" to table: "projects"
CREATE INDEX `idx_projects_deleted_at` ON `projects` (`deleted_at`);
