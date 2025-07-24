-- +goose Up
DROP INDEX IF EXISTS ix_posts_git_mr_team_channel;
CREATE INDEX IF NOT EXISTS ix_posts_git_mr_team_channel_project_id ON posts_git_mr (team_name, channel_id, git_project_id);

-- +goose Down
CREATE INDEX IF NOT EXISTS ix_posts_git_mr_team_channel ON posts_git_mr (team_name, channel_id);
DROP INDEX IF EXISTS ix_posts_git_mr_team_channel_project_id;