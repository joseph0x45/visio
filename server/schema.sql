create table users(
  id uuid not null primary key,
  email text not null unique,
  avatar_url text not null,
  plan text default 'basic'
);
