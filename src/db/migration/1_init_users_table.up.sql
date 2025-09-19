CREATE TABLE IF NOT EXISTS users (
        username VARCHAR(1000) NOT NULL,
        password VARCHAR(1000) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
      );