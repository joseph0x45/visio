create table if not exists users (
  id text not null primary key,
  email text not null unique,
  password text not null,
  signup_date timestamp not null
);
