-- Enable pgcrypto for gen_random_uuid() if on older Postgres, 
-- though Postgres 13+ has it built-in.
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    product_id INT NOT NULL,
    status VARCHAR(20) NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Seed Admin User (Password is 'password')
INSERT INTO users (email, password) 
VALUES ('admin@example.com', '$2a$10$vI8aWBnW3fBr4ffgq7zh6.yY0K2AA.a.5h.r.s.F/j.q.w.l.z.y.');