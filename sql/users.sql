CREATE TABLE `weetbot`.`users`
(
	`id`           BIGINT NOT NULL,
	`login`        VARCHAR(32) NOT NULL,
	`display_name` VARCHAR(32) NOT NULL,
	`followers`    INT NOT NULL,
	`views`        INT NOT NULL,
	`avatar`       VARCHAR(124) NULL,
	`banner`       VARCHAR(124) NULL,
	`type`         TINYINT(1) NOT NULL,
	`join`         TINYINT(1) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`),
	INDEX (`display_name`),
	INDEX (`followers`),
	INDEX (`views`),
	INDEX (`join`),
	INDEX (`type`),
	UNIQUE (`login`)
)
engine = innodb;

INSERT INTO `weetbot`.`users` (`id`, `login`, `display_name`, `followers`, `views`, `type`, `join`) VALUES ('51684790', 'modesttim', 'ModestTim', 0, 0, 1, 1);