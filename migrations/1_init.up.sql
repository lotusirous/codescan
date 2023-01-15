CREATE TABLE IF NOT EXISTS repos (
    repo_id integer AUTO_INCREMENT PRIMARY KEY,
    http_url varchar(255),
    created integer,
    updated integer
);

CREATE TABLE IF NOT EXISTS scans (
    scan_id integer AUTO_INCREMENT KEY,
    repo_id integer,
    status varchar(255) NOT NULL,
    enqueued_at integer,
    started_at integer,
    finished_at integer
);

CREATE TABLE IF NOT EXISTS scan_results (
    scan_result_id integer  AUTO_INCREMENT PRIMARY KEY,
    scan_id integer,
    repo_id integer,

    -- repo detail
    commit varchar(255), -- latetst commit
    message text,
    findings JSON,
    CONSTRAINT uc_scan_repo UNIQUE (scan_id,repo_id)
);
