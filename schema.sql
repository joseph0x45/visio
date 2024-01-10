create table if not exists users (
  id text not null primary key,
  email text not null unique,
  password_hash text not null,
  signup_date text not null
);

create table if not exists keys (
  id text not null primary key,
  user_id text not null references users(id),
  prefix text not null,
  key_hash text not null,
  creation_date text not null,
  unique (user_id, prefix)
);
