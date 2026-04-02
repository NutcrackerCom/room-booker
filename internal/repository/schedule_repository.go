package repository

import (
	"context"

	"room-booking/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleRepository struct {
	db *pgxpool.Pool
}

func NewScheduleRepository(db *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) ExistsByRoomID(ctx context.Context, roomID string) (bool, error) {
	query := `select exists(select 1 from schedules where room_id = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, roomID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *ScheduleRepository) Create(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error) {
	query := `
		insert into schedules (room_id, days_of_week, start_time, end_time)
		values ($1, $2, $3, $4)
		returning id, room_id, days_of_week, start_time::text, end_time::text
	`

	var schedule domain.Schedule
	err := r.db.QueryRow(ctx, query, roomID, daysOfWeek, startTime, endTime).Scan(
		&schedule.ID,
		&schedule.RoomID,
		&schedule.DaysOfWeek,
		&schedule.StartTime,
		&schedule.EndTime,
	)
	if err != nil {
		return nil, err
	}

	return &schedule, nil
}
