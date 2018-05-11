CREATE TABLE `weetbot`.`commands`
(
	`broadcaster` BIGINT NOT NULL,
	`active`      TINYINT(1) NOT NULL DEFAULT '1',
	`trigger`     VARCHAR(32) NOT NULL,
	`response`    VARCHAR(32) NOT NULL,
	`level`       TINYINT(3) NOT NULL DEFAULT '0',
	INDEX (`broadcaster`),
	INDEX (`trigger`, `broadcaster`),
	INDEX (`active`),
	INDEX (`trigger`),
	INDEX (`response`),
	INDEX (`level`)
)
engine = innodb;