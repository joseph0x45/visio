create table if not exists users (
  id text not null primary key,
  email text not null unique,
  password_hash text not null,
  signup_date text not null
);

create table if not exists keys (
  id text not null primary key,
  user_id text not null references users(id) on delete cascade,
  prefix text not null unique,
  key_hash text not null,
  creation_date text not null
);

create table if not exists faces (
  id text not null primary key,
  label text not null,
  user_id text not null references users(id) on delete cascade,
  descriptor text not null,
  unique (label, user_id)
);
