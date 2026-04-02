package repository

import (
	"context"
	"time"

	"room-booking/internal/domain"

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

func (r *SlotRepository) ListAvailable(ctx context.Context, roomID string, dayStart, dayEnd time.Time) ([]domain.Slot, error) {
	query := `
		select s.id, s.room_id, s.start_at, s.end_at
		from slots s
		left join bookings b
			on b.slot_id = s.id and b.status = 'active'
		where s.room_id = $1
		  and s.start_at >= $2
		  and s.start_at < $3
		  and b.id is null
		order by s.start_at asc
	`

	rows, err := r.db.Query(ctx, query, roomID, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []domain.Slot
	for rows.Next() {
		var slot domain.Slot
		if err := rows.Scan(&slot.ID, &slot.RoomID, &slot.Start, &slot.End); err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return slots, nil
}

func (r *SlotRepository) GetByID(ctx context.Context, slotID string) (*domain.Slot, error) {
	query := `
		select id, room_id, start_at, end_at
		from slots
		where id = $1
	`

	var slot domain.Slot
	err := r.db.QueryRow(ctx, query, slotID).Scan(
		&slot.ID,
		&slot.RoomID,
		&slot.Start,
		&slot.End,
	)
	if err != nil {
		return nil, err
	}

	return &slot, nil
}
