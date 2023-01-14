CREATE TABLE IF NOT EXISTS repos (
    repo_id integer AUTO_INCREMENT PRIMARY KEY,
    commit VARCHAR(255),
    http_url varchar(255),
    name varchar(255),
    created integer,
    updated integer
);

CREATE TABLE IF NOT EXISTS scans (
    scan_id serial PRIMARY KEY,
    repo_id integer REFERENCES repos (repo_id),
    status varchar(255) NOT NULL,
    enqueued_at integer,
    started_at integer,
    finished_at integer
);

