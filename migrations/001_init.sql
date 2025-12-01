-- 001_init.sql
CREATE TABLE IF NOT EXISTS subscriptions (
                                             id SERIAL PRIMARY KEY,
                                             client_id VARCHAR(255) NOT NULL UNIQUE,
    topics TEXT[] DEFAULT '{}',
    tags TEXT[] DEFAULT '{}',
    created_at TIMESTAMP DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS messages (
                                        id SERIAL PRIMARY KEY,
                                        title TEXT NOT NULL,
                                        body TEXT NOT NULL,
                                        topic VARCHAR(255) NOT NULL,
    tags TEXT[] DEFAULT '{}',
    created_at TIMESTAMP DEFAULT now()
    );
