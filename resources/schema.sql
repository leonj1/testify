USE mysql;
UPDATE user SET host = '%' WHERE host = '1%';
FLUSH PRIVILEGES;

create database testify;

create table testify.check (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(128) NOT NULL,
  `confession_name` VARCHAR(128) NOT NULL,
  `status` VARCHAR(32) NOT NULL,
  `journal_date` timestamp default CURRENT_TIMESTAMP NOT NULL
);

create table testify.journal (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `by` VARCHAR(128) NOT NULL,
  `confession_name` VARCHAR(128) NOT NULL,
  `status` VARCHAR(32) NOT NULL,
  `journal_date` timestamp default CURRENT_TIMESTAMP NOT NULL
);

create table testify.confession (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(128) NOT NULL,
  `entity_type` VARCHAR(128) NOT NULL,
  `last_update_date` timestamp default CURRENT_TIMESTAMP NOT NULL,
  `last_update_by` VARCHAR(128) NOT NULL,
  `last_update_status` VARCHAR(32) NOT NULL
);
