create database if not exists usersdb;
use usersdb;
create table if not exists users (
  id varchar(36) NOT NULL PRIMARY KEY,
  email text NOT NULL,
  first_name text NOT NULL,
  last_name text NOT NULL,
  role int
);
