package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SlotRepository struct {
	db *pgxpool.Pool
}

func NewSlotRepository(db *pgxpool.Pool) *SlotRepository {
	return &SlotRepository{db: db}
}

func (r *SlotRepository) Create(ctx context.Context, roomID string, startAt, endAt time.Time) error {
	query := `
		insert into slots (room_id, start_at, end_at)
		values ($1, $2, $3)
		on conflict (room_id, start_at) do nothing
	`
	_, err := r.db.Exec(ctx, query, roomID, startAt, endAt)
	return err
}

func (r *SlotRepository) CountByRoomID(ctx context.Context, roomID string) (int, error) {
	query := `select count(*) from slots where room_id = $1`

	var count int
	err := r.db.QueryRow(ctx, query, roomID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

