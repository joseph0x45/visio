create table users (
  id uuid not null primary key,
  github_id text not null unique,
  username text not null unique,
  email text not null unique,
  avatar text not null,
  plan text default 'basic'
);

create table keys (
  id uuid not null primary key,
  prefix text not null unique,
  owner uuid not null references users(id),
  key_hash text not null
);

create table faces (
  id uuid not null primary key,
  created_by uuid not null references users(id),
  descriptor text not null,
  created_at timestamp not null default now(),
  last_updated timestamp not null default now()
);
