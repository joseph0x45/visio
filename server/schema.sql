create table if not exists users (
  id text not null primary key,
  github_id text not null unique,
  email text not null,
  username text not null unique,
  avatar text not null,
  signup_date timestamp not null
);

create table if not exists api_keys (
  id text not null primary key,
  owner text not null references users(id),
  prefix text not null,
  key_hash text not null,
  created timestamp not null
);
