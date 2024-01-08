create table if not exists users (
  id text not null primary key,
  email text not null unique,
  password text not null,
  signup_date timestamp not null
);

create table if not exists keys (
  id text not null primary key,
  user_id text not null references users(id),
  prefix text not null unique,
  key_hash text not null,
  creation_date text not null
);
