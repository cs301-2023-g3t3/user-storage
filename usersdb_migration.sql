create database if not exists usersdb;
use usersdb;
create table if not exists users (
  id varchar(36) NOT NULL PRIMARY KEY,
  email varchar(255) NOT NULL,
  first_name varchar(255) NOT NULL,
  last_name varchar(255) NOT NULL,
  role varchar(255)
);
