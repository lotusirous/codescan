CREATE TABLE IF NOT EXISTS repos (
    repo_id integer AUTO_INCREMENT PRIMARY KEY,
    http_url varchar(255),
    name text,
    created integer,
    updated integer
);

CREATE TABLE IF NOT EXISTS scans (
    scan_id integer AUTO_INCREMENT KEY,
    repo_id integer,
    status varchar(255) NOT NULL CHECK (country IN ('Queued', 'In Progress', 'Success', 'Failure')),
    scan_error text, -- when the scan is failed.
    enqueued_at integer,
    started_at integer,
    finished_at integer
);

CREATE TABLE IF NOT EXISTS scan_results (
    scan_result_id integer  AUTO_INCREMENT PRIMARY KEY,
    scan_id integer,
    repo_id integer,
    created INTEGER,
    updated integer,
    -- repo detail
    commit varchar(255), -- The commit that produces the scan result.
    findings JSON,
    FOREIGN KEY (scan_id) REFERENCES scans (scan_id) ON DELETE CASCADE,
    CONSTRAINT uc_scan_repo UNIQUE (scan_id,repo_id)
);
