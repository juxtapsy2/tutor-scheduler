package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/brightpath/tutor-scheduler/internal/domain"
)

type AppointmentRepo struct {
	db *pgxpool.Pool
}

func NewAppointmentRepo(db *pgxpool.Pool) *AppointmentRepo {
	return &AppointmentRepo{db: db}
}

func (r *AppointmentRepo) Create(ctx context.Context, tx pgx.Tx, a *domain.Appointment) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO appointments (id, student_id, tutor_id, room_id, start_at, end_at, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		a.ID, a.StudentID, a.TutorID, a.RoomID, a.StartAt, a.EndAt, string(a.Status), a.CreatedAt,
	)
	return err
}

func (r *AppointmentRepo) FindConflictsForStudent(ctx context.Context, tx pgx.Tx, studentID string, start, end time.Time) (*domain.Appointment, error) {
	var a domain.Appointment
	err := tx.QueryRow(ctx,
		`SELECT id, student_id, tutor_id, room_id, start_at, end_at, status, created_at
		 FROM appointments
		 WHERE student_id = $1 AND status = 'BOOKED'
		   AND start_at < $3 AND end_at > $2
		 LIMIT 1`,
		studentID, start, end,
	).Scan(&a.ID, &a.StudentID, &a.TutorID, &a.RoomID, &a.StartAt, &a.EndAt, &a.Status, &a.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AppointmentRepo) FindConflictsForTutor(ctx context.Context, tx pgx.Tx, tutorID string, start, end time.Time) (*domain.Appointment, error) {
	var a domain.Appointment
	err := tx.QueryRow(ctx,
		`SELECT id, student_id, tutor_id, room_id, start_at, end_at, status, created_at
		 FROM appointments
		 WHERE tutor_id = $1 AND status = 'BOOKED'
		   AND start_at < $3 AND end_at > $2
		 LIMIT 1`,
		tutorID, start, end,
	).Scan(&a.ID, &a.StudentID, &a.TutorID, &a.RoomID, &a.StartAt, &a.EndAt, &a.Status, &a.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AppointmentRepo) FindConflictsForRoom(ctx context.Context, tx pgx.Tx, roomID string, start, end time.Time) (*domain.Appointment, error) {
	var a domain.Appointment
	err := tx.QueryRow(ctx,
		`SELECT id, student_id, tutor_id, room_id, start_at, end_at, status, created_at
		 FROM appointments
		 WHERE room_id = $1 AND status = 'BOOKED'
		   AND start_at < $3 AND end_at > $2
		 LIMIT 1`,
		roomID, start, end,
	).Scan(&a.ID, &a.StudentID, &a.TutorID, &a.RoomID, &a.StartAt, &a.EndAt, &a.Status, &a.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AppointmentRepo) CountTutorBookingsOnDay(ctx context.Context, tx pgx.Tx, tutorID string, day time.Time) (int, error) {
	dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	dayEnd := dayStart.Add(24 * time.Hour)
	var count int
	err := tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM appointments
		 WHERE tutor_id = $1 AND status = 'BOOKED'
		   AND start_at >= $2 AND start_at < $3`,
		tutorID, dayStart, dayEnd,
	).Scan(&count)
	return count, err
}

type ListFilters struct {
	Date      *time.Time
	StudentID *string
	TutorID   *string
	RoomID    *string
}

func (r *AppointmentRepo) List(ctx context.Context, filters ListFilters) ([]domain.Appointment, error) {
	query := `SELECT id, student_id, tutor_id, room_id, start_at, end_at, status, created_at FROM appointments WHERE 1=1`
	args := []any{}
	argIdx := 1

	if filters.Date != nil {
		dayStart := time.Date(filters.Date.Year(), filters.Date.Month(), filters.Date.Day(), 0, 0, 0, 0, filters.Date.Location())
		dayEnd := dayStart.Add(24 * time.Hour)
		query += fmt.Sprintf(" AND start_at >= $%d AND start_at < $%d", argIdx, argIdx+1)
		args = append(args, dayStart, dayEnd)
		argIdx += 2
	}
	if filters.StudentID != nil {
		query += fmt.Sprintf(" AND student_id = $%d", argIdx)
		args = append(args, *filters.StudentID)
		argIdx++
	}
	if filters.TutorID != nil {
		query += fmt.Sprintf(" AND tutor_id = $%d", argIdx)
		args = append(args, *filters.TutorID)
		argIdx++
	}
	if filters.RoomID != nil {
		query += fmt.Sprintf(" AND room_id = $%d", argIdx)
		args = append(args, *filters.RoomID)
		argIdx++
	}

	query += " ORDER BY start_at ASC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []domain.Appointment
	for rows.Next() {
		var a domain.Appointment
		if err := rows.Scan(&a.ID, &a.StudentID, &a.TutorID, &a.RoomID, &a.StartAt, &a.EndAt, &a.Status, &a.CreatedAt); err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}
	return appointments, rows.Err()
}
