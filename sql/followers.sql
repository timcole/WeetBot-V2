CREATE TABLE `weetbot`.`followers`
(
	`broadcaster`  BIGINT NOT NULL,
	`follower`     BIGINT NOT NULL,
	`login`        VARCHAR(32) NOT NULL,
	`display_name` VARCHAR(32) NOT NULL,
	INDEX (`broadcaster`),
	INDEX (`follower`, `broadcaster`),
	INDEX (`display_name`),
	INDEX (`login`)
)
engine = innodb;