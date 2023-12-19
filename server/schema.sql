create table if not exists users (
  id uuid not null primary key,
  github_id text not null unique,
  email text not null unique,
  username text not null unique,
  github_avatar text not null,
  credits integer not null default 1000,
  joined_at timestamp not null
);

create table if not exists api_keys (
  id uuid not null primary key,
  owner uuid not null references users(id),
  prefix text not null,
  key_hash text not null,
  created timestamp not null
);
