package dbhesebo

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"helper-sender-bot/internal/entity"
	"time"
)

func (r *Pgx) CreatePostGitMR(ctx context.Context, p entity.PostGitMR) error {
	now := time.Now()

	sqlStr, args, err := r.sb.
		Insert("posts_git_mr").
		SetMap(sq.Eq{
			"team_name":      p.TeamName,
			"channel_id":     p.ChannelID,
			"git_url":        p.GitURL,
			"git_project_id": p.GitProjectID,
			"git_mr_id":      p.GitMRID,
			"post_id":        p.PostID,
			"ttl_review":     p.TTLReview,
			"reviewers":      p.Reviewers,
			"create_at":      now,
		}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.Db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *Pgx) ExistsPostGitMR(ctx context.Context, gitURL string, projectID, mrID int) (bool, error) {
	sqlStr, args, err := r.sb.
		Select("1").
		From("posts_git_mr").
		Where(sq.Eq{
			"git_url":        gitURL,
			"git_project_id": projectID,
			"git_mr_id":      mrID,
		}).
		Limit(1).
		ToSql()
	if err != nil {
		return false, err
	}

	var has int
	err = r.Db.QueryRow(ctx, sqlStr, args...).Scan(&has)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return false, err
}

func (r *Pgx) GetPostGitMR(ctx context.Context, gitURL string, projectID, mrID int) (entity.PostGitMR, error) {
	sqlStr, args, err := r.sb.
		Select(
			"team_name",
			"channel_id",
			"git_url",
			"git_project_id",
			"git_mr_id",
			"post_id",
			"ttl_review",
			"reviewers",
			"create_at",
			"pushed_review",
		).
		From("posts_git_mr").
		Where(sq.Eq{
			"git_url":        gitURL,
			"git_project_id": projectID,
			"git_mr_id":      mrID,
		}).
		Limit(1).
		ToSql()
	if err != nil {
		return entity.PostGitMR{}, err
	}

	var p entity.PostGitMR
	err = r.Db.QueryRow(ctx, sqlStr, args...).Scan(
		&p.TeamName,
		&p.ChannelID,
		&p.GitURL,
		&p.GitProjectID,
		&p.GitMRID,
		&p.PostID,
		&p.TTLReview,
		&p.Reviewers,
		&p.CreateAt,
		&p.PushedReview,
	)
	return p, err
}

func (r *Pgx) GetListPostGitMR(
	ctx context.Context,
	team, channel string,
	gitProjectID int,
) ([]entity.PostGitMR, error) {
	sqlStr, args, err := r.sb.
		Select(
			"team_name",
			"channel_id",
			"git_url",
			"git_project_id",
			"git_mr_id",
			"post_id",
			"ttl_review",
			"reviewers",
			"create_at",
			"pushed_review",
		).
		From("posts_git_mr").
		Where(sq.Eq{
			"team_name":      team,
			"channel_id":     channel,
			"git_project_id": gitProjectID,
		}).
		OrderBy("create_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entity.PostGitMR
	for rows.Next() {
		var p entity.PostGitMR
		if err := rows.Scan(
			&p.TeamName,
			&p.ChannelID,
			&p.GitURL,
			&p.GitProjectID,
			&p.GitMRID,
			&p.PostID,
			&p.TTLReview,
			&p.Reviewers,
			&p.CreateAt,
			&p.PushedReview,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *Pgx) UpdatePostGitMRPushed(ctx context.Context, gitURL string, projectID, mrID int) error {
	sqlStr, args, err := r.sb.
		Update("posts_git_mr").
		Set("pushed_review", true).
		Where(sq.Eq{
			"git_url":        gitURL,
			"git_project_id": projectID,
			"git_mr_id":      mrID,
		}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.Db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *Pgx) DeletePostGitMR(ctx context.Context, gitURL string, projectID, mrID int) error {
	sqlStr, args, err := r.sb.
		Delete("posts_git_mr").
		Where(sq.Eq{
			"git_url":        gitURL,
			"git_project_id": projectID,
			"git_mr_id":      mrID,
		}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.Db.Exec(ctx, sqlStr, args...)
	return err
}
