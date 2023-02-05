CREATE TABLE IF NOT EXISTS users (
  username VARCHAR(30) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  email VARCHAR(50) NOT NULL UNIQUE,
  PRIMARY KEY (username)
);

CREATE TABLE IF NOT EXISTS notes (
 id UUID,
 title TEXT NOT NULL,
 username VARCHAR(30) references users(username) ON DELETE CASCADE NOT NULL,
 text TEXT,
 created_at TIMESTAMP NOT NULL,
 updated_at TIMESTAMP NOT NULL,
 PRIMARY KEY (id)
);
