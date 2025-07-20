package dbhesebo

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"helper-sender-bot/internal/entity"
	"time"
)

func (r *Pgx) CreateCfgDuty(ctx context.Context, c entity.Chat, team string) error {
	now := time.Now()
	sqlStr, args, err := r.sb.
		Insert("config_duty").
		SetMap(sq.Eq{
			"team_name":          team,
			"channel_id":         c.ChannelID,
			"duty_ttl":           c.DutyTtlInMinute,
			"start_emoji":        c.EmojiStart,
			"done_emoji":         c.EmojiDone,
			"workday_start_hour": c.WorkdayStart,
			"workday_end_hour":   c.WorkdayEnd,
			"duty_repeat_ttl":    c.DutyRepeatTtlInMinute,
			"create_at":          now,
			"update_at":          now,
		}).
		ToSql()

	if err != nil {
		return fmt.Errorf("CreateCfgDuty ToSql: %w", err)
	}
	_, err = r.Db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *Pgx) UpdateCfgDuty(ctx context.Context, team, channel string, upd entity.Chat) error {
	sqlStr, args, err := r.sb.
		Update("config_duty").
		SetMap(sq.Eq{
			"duty_ttl":           upd.DutyTtlInMinute,
			"start_emoji":        upd.EmojiStart,
			"done_emoji":         upd.EmojiDone,
			"workday_start_hour": upd.WorkdayStart,
			"workday_end_hour":   upd.WorkdayEnd,
			"duty_repeat_ttl":    upd.DutyRepeatTtlInMinute,
			"update_at":          time.Now(),
		}).
		Where(sq.Eq{
			"team_name":  team,
			"channel_id": channel,
		}).
		ToSql()

	if err != nil {
		return fmt.Errorf("UpdateCfgDuty ToSql: %w", err)
	}
	if len(args) <= 2 {
		return nil
	}
	_, err = r.Db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *Pgx) DeleteCfgDuty(ctx context.Context, team, channel string) error {
	sqlStr, args, err := r.sb.
		Delete("config_duty").
		Where(sq.Eq{
			"team_name":  team,
			"channel_id": channel,
		}).
		ToSql()

	if err != nil {
		return fmt.Errorf("DeleteCfgDuty ToSql: %w", err)
	}
	_, err = r.Db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *Pgx) GetListCfgDuty(ctx context.Context) ([]entity.Chat, error) {
	qb := r.sb.Select()
	return r.fetchCfgDuty(ctx, qb)
}

func (r *Pgx) GetListCfgDutyByTeam(ctx context.Context, team string) ([]entity.Chat, error) {
	qb := r.sb.
		Select().
		Where(sq.Eq{"team_name": team})
	return r.fetchCfgDuty(ctx, qb)
}

func (r *Pgx) fetchCfgDuty(ctx context.Context, qb sq.SelectBuilder) ([]entity.Chat, error) {
	qb = qb.Columns(
		"team_name",
		"channel_id",
		"duty_ttl",
		"start_emoji",
		"done_emoji",
		"workday_start_hour",
		"workday_end_hour",
		"duty_repeat_ttl",
	)

	sqlStr, args, err := qb.From("config_duty").ToSql()
	if err != nil {
		return nil, fmt.Errorf("fetchChats ToSql: %w", err)
	}

	rows, err := r.Db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entity.Chat
	for rows.Next() {
		var c entity.Chat
		if err := rows.Scan(
			&c.Team,
			&c.ChannelID,
			&c.DutyTtlInMinute,
			&c.EmojiStart,
			&c.EmojiDone,
			&c.WorkdayStart,
			&c.WorkdayEnd,
			&c.DutyRepeatTtlInMinute,
		); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
