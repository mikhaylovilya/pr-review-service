-- +goose Up
CREATE TABLE teams (
    id INT PRIMARY KEY,
    name VARCHAR(64) NOT NULL
);

CREATE TABLE pull_requests (
    id INT PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    author_id VARCHAR(64),
    status VARCHAR(64),
    created_at TIMESTAMP,
    merged_at TIMESTAMP
);

CREATE TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    team_name VARCHAR(64),
    is_active BOOLEAN,
    team_id INT,
    assigned_pr_id INT,
    CONSTRAINT fk_teams FOREIGN KEY (team_id) REFERENCES teams (id),
    CONSTRAINT fk_pull_requests FOREIGN KEY (
        assigned_pr_id
    ) REFERENCES pull_requests (id)
);

-- +goose Down
DROP TABLE user;
DROP TABLE team;
DROP TABLE pull_request;
