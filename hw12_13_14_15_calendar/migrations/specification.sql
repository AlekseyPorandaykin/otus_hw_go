CREATE TABLE IF NOT EXISTS users (
    id INT PRIMARY KEY,
    user_name VARCHAR ( 50 )  NOT NULL
);
alter table users
    owner to calendar;

CREATE TABLE IF NOT EXISTS events  (
    id VARCHAR ( 50 ) PRIMARY KEY,
    title TEXT  NOT NULL,
    description TEXT NOT NULL,
    datetime_from TIMESTAMP NOT NULL,
    datetime_to TIMESTAMP NOT NULL,
    created_by INTEGER NOT NULL,
    start_notify TIMESTAMP NULL,
    CONSTRAINT fk_user
        FOREIGN KEY(created_by)
            REFERENCES users(id)
    );
alter table events
    owner to calendar;

INSERT INTO users(id, user_name) VALUES (1, 'test_user')
    ON CONFLICT (id) DO NOTHING;