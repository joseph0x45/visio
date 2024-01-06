create table if not exists users (
  id text not null primary key,
  email text not null unique,
  password text not null,
  signup_date timestamp not null
);

create table if not exists keys (
  key_owner text not null references users(id),
  prefix text not null unique primary key,
  key_hash text not null unique,
  key_creation_date timestamp not null
)