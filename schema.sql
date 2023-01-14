CREATE TABLE repositories (
  repo_id integer,
  commit VARCHAR(255),
  http_url varchar(255),
  name varchar(255),
  created integer,
  updated integer
);

CREATE TABLE scans (
  scan_id serial PRIMARY KEY,
  repo_id integer REFERENCES repositories (id) NOT NULL,
  status varchar(255) NOT NULL,
  enqueued_at timestamp NOT NULL,
  started_at timestamp,
  finished_at timestamp
);

CREATE TABLE scan_results (
  result_id serial PRIMARY KEY,
  scan_id integer REFERENCES scans (id) NOT NULL,
  rule_id varchar(255) NOT NULL,
  description varchar(255) NOT NULL,
  severity varchar(255) NOT NULL,
  location varchar(255) NOT NULL,
  created_at timestamp NOT NULL
);

