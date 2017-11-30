USE mysql;
UPDATE user SET host = '%' WHERE host = '1%';
FLUSH PRIVILEGES;

create database enchilada;

create table enchilada.hardware (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  host VARCHAR(128) NOT NULL,
  create_date timestamp default CURRENT_TIMESTAMP NOT NULL
);

create table enchilada.tags (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `host` VARCHAR(128) NOT NULL,
  `key` VARCHAR(128) NOT NULL,
  `value` VARCHAR(128) NOT NULL
);

create table enchilada.services (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `host` VARCHAR(128) NOT NULL,
  `name` VARCHAR(128) NOT NULL,
  `short_name` VARCHAR(128) NOT NULL,
  `version` VARCHAR(128) NOT NULL,
  `repo` VARCHAR(128) NOT NULL,
  `is_docker` INT NOT NULL,
  `install_date` timestamp default CURRENT_TIMESTAMP
);
