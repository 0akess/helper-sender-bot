package dbhesebo

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"helper-sender-bot/internal/entity"
	"time"
)

func (r *Pgx) CreateGitlabConfig(ctx context.Context, cfg entity.GitlabConfig) error {
	now := time.Now()
	sqlStr, args, err := r.sb.
		Insert("config_gitlab").
		SetMap(sq.Eq{
			"team_name":            cfg.Team,
			"git_url":              cfg.GitlabURL,
			"git_project_name":     cfg.ProjectName,
			"git_project_id":       cfg.ProjectID,
			"channel_id":           cfg.ChannelID,
			"reviewers":            cfg.Reviewers,
			"reviewers_count":      cfg.ReviewersCount,
			"ttl_review":           cfg.TTLReview,
			"qa_reviewers":         cfg.QAReviewers,
			"requires_qa_review":   cfg.RequiresQaReview,
			"push_qa_after_review": cfg.PushQaAfterReview,
			"create_at":            now,
			"update_at":            now,
		}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.Db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *Pgx) DeleteGitlabConfigByProjectID(ctx context.Context, gitProjectID int, gitURL, team string) error {
	sqlStr, args, err := r.sb.
		Delete("config_gitlab").
		Where(sq.Eq{
			"git_project_id": gitProjectID,
			"git_url":        gitURL,
			"team_name":      team,
		}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.Db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *Pgx) UpdateGitlabConfig(ctx context.Context, cfg entity.GitlabConfig, gitProjectID int, gitURL string) error {
	sqlStr, args, err := r.sb.
		Update("config_gitlab").
		SetMap(sq.Eq{
			"git_project_name":     cfg.ProjectName,
			"reviewers":            cfg.Reviewers,
			"reviewers_count":      cfg.ReviewersCount,
			"ttl_review":           cfg.TTLReview,
			"qa_reviewers":         cfg.QAReviewers,
			"requires_qa_review":   cfg.RequiresQaReview,
			"push_qa_after_review": cfg.PushQaAfterReview,
			"update_at":            time.Now(),
		}).
		Where(sq.Eq{
			"team_name":      cfg.Team,
			"git_url":        gitURL,
			"git_project_id": gitProjectID,
		}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.Db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *Pgx) GetAllGitlabConfigs(ctx context.Context) ([]entity.GitlabConfig, error) {
	return r.queryGitlabConfigs(ctx, r.sb.Select())
}

func (r *Pgx) GetGitlabConfigsByTeam(ctx context.Context, team string) ([]entity.GitlabConfig, error) {
	qb := r.sb.
		Select().
		Where(sq.Eq{"team_name": team})
	return r.queryGitlabConfigs(ctx, qb)
}

func (r *Pgx) GetGitlabConfig(ctx context.Context, projectID int, gitUrl string) (entity.GitlabConfig, error) {
	qb := r.sb.
		Select().
		Where(sq.Eq{
			"git_project_id": projectID,
			"git_url":        gitUrl,
		})
	cfgs, err := r.queryGitlabConfigs(ctx, qb)
	if err != nil {
		return entity.GitlabConfig{}, err
	}
	if len(cfgs) == 0 {
		return entity.GitlabConfig{}, fmt.Errorf("empty gitlab config")
	}
	return cfgs[0], nil
}

func (r *Pgx) queryGitlabConfigs(ctx context.Context, qb sq.SelectBuilder) ([]entity.GitlabConfig, error) {
	sqlStr, args, err := qb.
		Columns("team_name",
			"channel_id",
			"git_url",
			"git_project_name",
			"git_project_id",
			"reviewers",
			"reviewers_count",
			"ttl_review",
			"qa_reviewers",
			"requires_qa_review",
			"push_qa_after_review",
		).
		From("config_gitlab").
		OrderBy("create_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cfgs []entity.GitlabConfig
	for rows.Next() {
		var (
			cfg entity.GitlabConfig
		)

		if err := rows.Scan(
			&cfg.Team,
			&cfg.ChannelID,
			&cfg.GitlabURL,
			&cfg.ProjectName,
			&cfg.ProjectID,
			&cfg.Reviewers,
			&cfg.ReviewersCount,
			&cfg.TTLReview,
			&cfg.QAReviewers,
			&cfg.RequiresQaReview,
			&cfg.PushQaAfterReview,
		); err != nil {
			return nil, err
		}

		cfgs = append(cfgs, cfg)
	}
	return cfgs, rows.Err()
}
