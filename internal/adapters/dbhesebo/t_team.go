package dbhesebo

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"helper-sender-bot/internal/entity"
	"time"
)

func (r *Pgx) CreateTeam(ctx context.Context, team entity.Team) error {
	now := time.Now()
	sqlStr, args, err := r.sb.
		Insert("team").
		SetMap(sq.Eq{
			"team_name":     team.Name,
			"token":         team.Token,
			"team_lead_eid": team.LeadEID,
			"create_at":     now,
			"update_at":     now,
		}).
		ToSql()

	if err != nil {
		return fmt.Errorf("CreateTeam ToSql: %w", err)
	}
	_, err = r.Db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *Pgx) GetListTeam(ctx context.Context, limit, cursor int, search string) ([]string, int, error) {
	qb := r.sb.
		Select("team_name").
		From("team")

	if search != "" {
		qb = qb.Where(sq.ILike{"team_name": "%" + search + "%"})
	}

	qb = qb.
		OrderBy("create_at ASC").
		Limit(uint64(limit)).
		Offset(uint64(cursor))

	sqlStr, args, err := qb.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("GetListTeam ToSql: %w", err)
	}
	rows, err := r.Db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("GetListTeam Query: %w", err)
	}
	defer rows.Close()
	var teams []string
	for rows.Next() {
		var team string
		if err := rows.Scan(&team); err != nil {
			return nil, 0, fmt.Errorf("GetListTeam Scan: %w", err)
		}
		teams = append(teams, team)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("GetListTeam Err: %w", err)
	}

	nextCursor := -1
	if len(teams) == limit {
		nextCursor = cursor + 1 + limit
	}

	return teams, nextCursor, nil
}

func (r *Pgx) UpdateTeam(ctx context.Context, teamName string, token uuid.UUID, newTeam entity.Team) error {
	sqlStr, args, err := r.sb.
		Update("team").
		SetMap(sq.Eq{
			"update_at":     time.Now(),
			"token":         newTeam.Token,
			"team_lead_eid": newTeam.LeadEID,
		}).
		Where(sq.Eq{
			"token":     token,
			"team_name": teamName,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("GetTeam ToSql: %w", err)
	}
	_, err = r.Db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("GetTeam Exec: %w", err)
	}
	return nil
}

func (r *Pgx) DeleteTeam(ctx context.Context, teamName string, token uuid.UUID) error {
	sqlStr, args, err := r.sb.
		Delete("team").
		Where(sq.Eq{
			"team_name": teamName,
			"token":     token,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("DeleteTeam ToSql: %w", err)
	}
	_, err = r.Db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("DeleteTeam Exec: %w", err)
	}
	return nil
}

func (r *Pgx) OkTokenTeam(ctx context.Context, teamName string, token uuid.UUID) (bool, error) {
	sqlStr, args, err := r.sb.
		Select("COUNT(*)").
		From("team").
		Where(sq.Eq{
			"team_name": teamName,
			"token":     token,
		}).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("OkTokenTeam ToSql: %w", err)
	}

	var count int
	if err := r.Db.QueryRow(ctx, sqlStr, args...).Scan(&count); err != nil {
		return false, fmt.Errorf("OkTokenTeam QueryRow: %w", err)
	}
	return count > 0, nil
}
