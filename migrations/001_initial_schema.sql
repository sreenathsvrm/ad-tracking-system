CREATE TABLE ads (
    id          VARCHAR(36) PRIMARY KEY,
    image_url   TEXT NOT NULL,
    target_url  TEXT NOT NULL
);

CREATE TABLE clicks (
    id            SERIAL PRIMARY KEY,
    ad_id         VARCHAR(36) NOT NULL,
    timestamp     TIMESTAMP NOT NULL,
    ip            TEXT NOT NULL,
    playback_time INT NOT NULL
);