CREATE TABLE `weetbot`.`subscribers`
(
	`broadcaster`  BIGINT NOT NULL,
	`subscriber`   BIGINT NOT NULL,
	`login`        VARCHAR(32) NOT NULL,
	`display_name` VARCHAR(32) NOT NULL,
	`type`         TINYINT(1) NOT NULL,
	`months`       TINYINT(3) NOT NULL,
	INDEX (`broadcaster`),
	INDEX (`subscriber`, `broadcaster`),
	INDEX (`login`),
	INDEX (`display_name`),
	INDEX (`type`),
	INDEX (`months`)
)
engine = innodb;