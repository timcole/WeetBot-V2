CREATE TABLE `weetbot`.`points`
(
	`broadcaster`  BIGINT NOT NULL,
	`chatter`      BIGINT NOT NULL,
	`login`        VARCHAR(32) NOT NULL,
	`display_name` VARCHAR(32) NOT NULL,
	`points`       INT NOT NULL,
	INDEX (`broadcaster`),
	INDEX (`chatter`, `broadcaster`),
	INDEX (`points`),
	INDEX (`display_name`),
	INDEX (`login`)
)
engine = innodb;