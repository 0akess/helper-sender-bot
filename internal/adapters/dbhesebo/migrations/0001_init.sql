-- +goose Up
CREATE TABLE IF NOT EXISTS team
(
    team_name     TEXT PRIMARY KEY,
    token         UUID        NOT NULL,
    team_lead_eid TEXT        NOT NULL,
    create_at     TIMESTAMPTZ not null,
    update_at     TIMESTAMPTZ not null
);

CREATE TABLE IF NOT EXISTS config_duty
(
    channel_id         TEXT PRIMARY KEY,
    team_name          TEXT        NOT NULL,
    duty_ttl           INT         NOT NULL,
    duty_repeat_ttl    INT         NOT NULL,
    start_emoji        TEXT        NOT NULL,
    done_emoji         TEXT        NOT NULL,
    workday_start_hour SMALLINT    NOT NULL,
    workday_end_hour   SMALLINT    NOT NULL,
    create_at          TIMESTAMPTZ not null,
    update_at          TIMESTAMPTZ not null,

    CONSTRAINT fk_chat_team
        FOREIGN KEY (team_name)
            REFERENCES team (team_name)
            ON UPDATE CASCADE
            ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS ix_config_duty_by_team ON config_duty (team_name);

CREATE TABLE IF NOT EXISTS posts_duty
(
    channel_id   TEXT        NOT NULL
        CONSTRAINT fk_posts_duty
            REFERENCES config_duty (channel_id)
            ON UPDATE CASCADE
            ON DELETE CASCADE,
    post_id      TEXT        NOT NULL,
    create_at    TIMESTAMPTZ NOT NULL,
    last_push_at TIMESTAMPTZ NOT NULL,
    in_progress  BOOLEAN     NOT NULL,
    PRIMARY KEY (channel_id, post_id)
);
CREATE INDEX IF NOT EXISTS idx_posts_duty_create_at ON posts_duty (create_at);

CREATE TABLE IF NOT EXISTS config_gitlab
(
    team_name            TEXT        NOT NULL,
    channel_id           TEXT        NOT NULL,

    git_url              TEXT        NOT NULL,
    git_project_id       INT         NOT NULL,
    git_project_name     TEXT        NOT NULL,

    reviewers            TEXT[]      NOT NULL,
    reviewers_count      SMALLINT    NOT NULL,
    ttl_review           JSONB,
    qa_reviewers         TEXT,
    requires_qa_review   BOOL,
    push_qa_after_review BOOL,

    create_at            TIMESTAMPTZ NOT NULL,
    update_at            TIMESTAMPTZ NOT NULL,

    CONSTRAINT uq_team_channel_project UNIQUE (team_name, channel_id, git_url, git_project_id),

    CONSTRAINT uq_project_global UNIQUE (git_url, git_project_id),

    CONSTRAINT fk_gitlab_team
        FOREIGN KEY (team_name)
            REFERENCES team (team_name)
            ON UPDATE CASCADE
            ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS ix_config_gitlab_team_channel ON config_gitlab (team_name, channel_id);
CREATE INDEX IF NOT EXISTS ix_config_gitlab_project_global ON config_gitlab (git_url, git_project_id);

CREATE TABLE IF NOT EXISTS posts_git_mr
(
    team_name      TEXT        NOT NULL,
    channel_id     TEXT        NOT NULL,
    git_url        TEXT        NOT NULL,
    git_project_id INT         NOT NULL,
    git_mr_id      INT         NOT NULL,
    create_at      TIMESTAMPTZ NOT NULL,
    pushed_review  BOOL DEFAULT FALSE,
    ttl_review     JSONB       NOT NULL,
    reviewers      TEXT        NOT NULL,
    post_id        TEXT        NOT NULL,
    PRIMARY KEY (git_url, git_project_id, git_mr_id)
);
CREATE INDEX IF NOT EXISTS ix_posts_git_mr_update_at ON posts_git_mr (create_at);
CREATE INDEX IF NOT EXISTS ix_posts_git_mr_team_channel ON posts_git_mr (team_name, channel_id);

-- +goose Down
drop table if exists config_gitlab;
drop table if exists posts_duty;
drop table if exists team;
drop table if exists config_duty;
drop table if exists posts_git_mr;