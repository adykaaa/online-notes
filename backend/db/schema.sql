CREATE TABLE users 
(
  id UUID DEFAULT gen_random_uuid (),
  email VARCHAR(50) NOT NULL UNIQUE,
  username VARCHAR(30) NOT NULL UNIQUE,
  password VARCHAR(30) NOT NULL,
  logged_in BOOLEAN DEFAULT FALSE,
  PRIMARY KEY (id)
);

CREATE TABLE notes
(
 id UUID DEFAULT gen_random_uuid (),
 title TEXT NOT NULL,
 user UUID references users(id),
 text TEXT,
 created_at TIMESTAMP,
 updated_at TIMESTAMP,
 PRIMARY KEY (id)
);
