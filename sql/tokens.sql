CREATE TABLE `weetbot`.`tokens`
(
	`id`            BIGINT NOT NULL,
	`access_token`  VARCHAR(32) NULL,
	`refresh_token` VARCHAR(64) NULL,
	PRIMARY KEY (`id`),
	INDEX (`access_token`),
	INDEX (`refresh_token`)
)
engine = innodb;