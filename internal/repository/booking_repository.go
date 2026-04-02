package repository

import (
	"context"

	"room-booking/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(ctx context.Context, slotID, userID string, conferenceLink *string) (*domain.Booking, error) {
	query := `
		insert into bookings (slot_id, user_id, status, conference_link)
		values ($1, $2, 'active', $3)
		returning id, slot_id, user_id, status, conference_link, created_at
	`

	var booking domain.Booking
	err := r.db.QueryRow(ctx, query, slotID, userID, conferenceLink).Scan(
		&booking.ID,
		&booking.SlotID,
		&booking.UserID,
		&booking.Status,
		&booking.ConferenceLink,
		&booking.CreatedAt,
	)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, domain.ErrSlotBooked
		}
		return nil, err
	}

	return &booking, nil
}
