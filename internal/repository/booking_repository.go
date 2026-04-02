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

func (r *BookingRepository) ListMyUpcoming(ctx context.Context, userID string) ([]domain.Booking, error) {
	query := `
		select b.id, b.slot_id, b.user_id, b.status, b.conference_link, b.created_at
		from bookings b
		join slots s on s.id = b.slot_id
		where b.user_id = $1
		  and s.start_at >= now()
		order by s.start_at asc
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		var booking domain.Booking
		if err := rows.Scan(
			&booking.ID,
			&booking.SlotID,
			&booking.UserID,
			&booking.Status,
			&booking.ConferenceLink,
			&booking.CreatedAt,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *BookingRepository) GetByID(ctx context.Context, bookingID string) (*domain.Booking, error) {
	query := `
		select id, slot_id, user_id, status, conference_link, created_at
		from bookings
		where id = $1
	`

	var booking domain.Booking
	err := r.db.QueryRow(ctx, query, bookingID).Scan(
		&booking.ID,
		&booking.SlotID,
		&booking.UserID,
		&booking.Status,
		&booking.ConferenceLink,
		&booking.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *BookingRepository) Cancel(ctx context.Context, bookingID string) (*domain.Booking, error) {
	query := `
		update bookings
		set status = 'cancelled'
		where id = $1
		returning id, slot_id, user_id, status, conference_link, created_at
	`

	var booking domain.Booking
	err := r.db.QueryRow(ctx, query, bookingID).Scan(
		&booking.ID,
		&booking.SlotID,
		&booking.UserID,
		&booking.Status,
		&booking.ConferenceLink,
		&booking.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}
