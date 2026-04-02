package repository

import (
	"context"

	"room-booking/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomRepository struct {
	db *pgxpool.Pool
}

func NewRoomRepository(db *pgxpool.Pool) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error) {
	query := `
		insert into rooms (name, description, capacity)
		values ($1, $2, $3)
		returning id, name, description, capacity, created_at
	`

	var room domain.Room
	err := r.db.QueryRow(ctx, query, name, description, capacity).Scan(
		&room.ID,
		&room.Name,
		&room.Description,
		&room.Capacity,
		&room.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *RoomRepository) List(ctx context.Context) ([]domain.Room, error) {
	query := `
		select id, name, description, capacity, created_at
		from rooms
		order by created_at asc
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []domain.Room
	for rows.Next() {
		var room domain.Room
		if err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.Description,
			&room.Capacity,
			&room.CreatedAt,
		); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *RoomRepository) Exists(ctx context.Context, roomID string) (bool, error) {
	query := `select exists(select 1 from rooms where id = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, roomID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
