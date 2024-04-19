CREATE TABLE users (
  id       BIGSERIAL PRIMARY KEY,
  name     text      NOT NULL UNIQUE,
  password text      NOT NULL,
  email    text      NOT NULL,
  is_admin boolean   NOT NULL DEFAULT false
);
