CREATE TABLE users 
(
  username VARCHAR(30) NOT NULL UNIQUE,
  password VARCHAR(30) NOT NULL,
  email VARCHAR(50) NOT NULL UNIQUE,
  logged_in BOOLEAN DEFAULT FALSE,
  PRIMARY KEY (username)
);

CREATE TABLE notes
(
 id UUID DEFAULT gen_random_uuid (),
 title TEXT NOT NULL,
 user VARCHAR (30) references users(username),
 text TEXT,
 created_at TIMESTAMP,
 updated_at TIMESTAMP,
 PRIMARY KEY (id)
);
