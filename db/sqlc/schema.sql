CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
  username VARCHAR(30) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  email VARCHAR(50) NOT NULL UNIQUE,
  PRIMARY KEY (username)
);

CREATE TABLE IF NOT EXISTS notes (
 id UUID DEFAULT gen_random_uuid(),
 title TEXT NOT NULL,
 username VARCHAR(30) references users(username) ON DELETE CASCADE,
 text TEXT,
 created_at TIMESTAMP,
 updated_at TIMESTAMP,
 PRIMARY KEY (id)
);
