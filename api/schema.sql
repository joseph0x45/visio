create table users (
  id uuid not null primary key,
  github_id text not null unique,
  email text not null unique,
  username text not null unique,
  plan text not null default 'basic'
);