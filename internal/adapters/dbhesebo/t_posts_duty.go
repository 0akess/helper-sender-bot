package dbhesebo

import (
	"context"
	"helper-sender-bot/internal/entity"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (r *Pgx) GetListOpenPostDuty(ctx context.Context, channelID string) ([]entity.PostsInfoDuty, error) {
	sqlStr, args, err := r.sb.
		Select(
			"channel_id",
			"post_id",
			"create_at",
			"last_push_at",
			"in_progress",
		).
		From("posts_duty").
		Where(sq.Eq{"channel_id": channelID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.Db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entity.PostsInfoDuty
	for rows.Next() {
		var ti entity.PostsInfoDuty
		if err := rows.Scan(
			&ti.ChannelID,
			&ti.PostID,
			&ti.CreateAt,
			&ti.LastPushAt,
			&ti.InProgress,
		); err != nil {
			return nil, err
		}
		out = append(out, ti)
	}
	return out, rows.Err()
}

func (r *Pgx) CreatePostDuty(ctx context.Context, channelID, postID string, creatAt time.Time, inProgress bool) error {
	sqlStr, args, err := r.sb.
		Insert("posts_duty").
		SetMap(sq.Eq{
			"channel_id":   channelID,
			"post_id":      postID,
			"create_at":    creatAt,
			"last_push_at": time.Now(),
			"in_progress":  inProgress,
		}).
		Suffix(`
            ON CONFLICT (channel_id, post_id) DO NOTHING
        `).
		ToSql()

	if err != nil {
		return err
	}
	_, err = r.Db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *Pgx) DeletePostDutyOlderThan(ctx context.Context, channelID string, dateDif time.Time) error {
	sqlStr, args, err := r.sb.
		Delete("posts_duty").
		Where(sq.Eq{"channel_id": channelID}).
		Where(sq.Lt{"create_at": dateDif}).
		ToSql()

	if err != nil {
		return err
	}
	_, err = r.Db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *Pgx) DeletePostDuty(ctx context.Context, channelID, postID string) error {
	sqlStr, args, err := r.sb.
		Delete("posts_duty").
		Where(sq.Eq{
			"channel_id": channelID,
			"post_id":    postID,
		}).
		ToSql()

	if err != nil {
		return err
	}
	_, err = r.Db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *Pgx) MarkPostsDutyAsInProgress(ctx context.Context, ch, postID string) error {
	sqlStr, args, err := r.sb.
		Update("posts_duty").
		Set("in_progress", true).
		Where(sq.Eq{
			"channel_id": ch,
			"post_id":    postID,
		}).
		ToSql()

	if err != nil {
		return err
	}
	_, err = r.Db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *Pgx) UpdatePushAtPostDuty(ctx context.Context, ch, postID string) error {
	sqlStr, args, err := r.sb.
		Update("posts_duty").
		Set("last_push_at", time.Now()).
		Where(sq.Eq{
			"channel_id": ch,
			"post_id":    postID,
		}).
		ToSql()

	if err != nil {
		return err
	}
	_, err = r.Db.Exec(ctx, sqlStr, args...)
	return err
}
