create table if not exists users (
  id text not null primary key,
  email text not null unique,
  password text not null,
  signup_date timestamp not null
);

create table if not exists keys (
  owner text not null references users(id),
  prefix varchar(26) not null unique,
  key text not null unique primary key,
  key_creation_date timestamp not null
)