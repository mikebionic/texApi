CREATE TABLE admins (
    id SERIAL PRIMARY KEY,
    phone VARCHAR NOT NULL,
    password VARCHAR NOT NULL
);

CREATE TABLE languages (
    id SERIAL PRIMARY KEY,
    lang VARCHAR NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR NOT NULL,
    phone VARCHAR NOT NULL,
    address VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    is_verified BOOL NOT NULL DEFAULT false,
    subscriptions_id INTEGER NULL,
    notification_token VARCHAR NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    CONSTRAINT subscriptions_id
        FOREIGN KEY (subscriptions_id)
            REFERENCES subscriptions (id)
                ON DELETE SET NULL
);


CREATE TABLE about_us (
    id SERIAL PRIMARY KEY,
    is_active BOOL DEFAULT false
);

CREATE table about_us_translates (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    languages_id INTEGER NULL,
    about_us_id INTEGER NOT NULL,

    CONSTRAINT languages_id
        FOREIGN KEY (languages_id)
            REFERENCES languages (id)
                ON DELETE SET NULL,

    CONSTRAINT about_us_id
        FOREIGN KEY (about_us_id)
            REFERENCES about_us (id)
                ON DELETE CASCADE
);


CREATE TABLE services (
    id SERIAL PRIMARY KEY,
    image VARCHAR NULL,
    parent_id INTEGER NULL,

    CONSTRAINT parent_id
        FOREIGN KEY (parent_id)
            REFERENCES services (id)
                ON DELETE SET NULL
);

CREATE TABLE service_translates (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    languages_id INTEGER NULL,
    services_id INTEGER NOT NULL,

    CONSTRAINT languages_id
        FOREIGN KEY (languages_id)
            REFERENCES languages (id)
                ON DELETE SET NULL,

    CONSTRAINT services_id
        FOREIGN KEY (services_id)
            REFERENCES services (id)
                ON DELETE CASCADE
);

CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    start_at TIMESTAMP NOT NULL,
    end_at TIMESTAMP NOT NULL,
    days INTEGER NOT NULL,
    count INTEGER NOT NULL,
    price DECIMAL NOT NULL
);

CREATE TABLE subscription_translates (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    languages_id INTEGER NULL,
    subscriptions_id INTEGER NOT NULL,

    CONSTRAINT languages_id
        FOREIGN KEY (languages_id)
            REFERENCES languages (id)
                ON DELETE SET NULL,

    CONSTRAINT subscriptions_id
        FOREIGN KEY (subscriptions_id)
            REFERENCES subscriptions (id)
                ON DELETE CASCADE
);
