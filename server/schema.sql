create table users(
  id uuid not null primary key,
  github_id text not null unique,
  username text not null unique,
  email text not null unique,
  avatar text not null,
  plan text default 'basic'
);
