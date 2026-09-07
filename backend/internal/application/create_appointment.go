package application

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/brightpath/tutor-scheduler/internal/domain"
	"github.com/brightpath/tutor-scheduler/internal/repository"
)

const (
	maxTutorDailyBookings = 6
	validDuration60       = 60 * time.Minute
	validDuration90       = 90 * time.Minute
)

type CreateAppointment struct {
	DB           *pgxpool.Pool
	Students     *repository.StudentRepo
	Tutors       *repository.TutorRepo
	Rooms        *repository.RoomRepo
	Appointments *repository.AppointmentRepo
	IDGenerator  func() string
	Clock        func() time.Time
}

func (s *CreateAppointment) Execute(ctx context.Context, input domain.CreateAppointmentInput) (*domain.Appointment, error) {
	if !input.StartAt.Before(input.EndAt) {
		return nil, domain.ErrInvalidTimeRange
	}

	if input.StartAt.Before(s.Clock()) {
		return nil, domain.ErrTimeInPast
	}

	if _, err := s.Students.FindByID(ctx, input.StudentID); err != nil {
		return nil, domain.ErrStudentNotFound
	}
	if _, err := s.Tutors.FindByID(ctx, input.TutorID); err != nil {
		return nil, domain.ErrTutorNotFound
	}
	if _, err := s.Rooms.FindByID(ctx, input.RoomID); err != nil {
		return nil, domain.ErrRoomNotFound
	}

	duration := input.EndAt.Sub(input.StartAt)
	if duration != validDuration60 && duration != validDuration90 {
		return nil, domain.ErrInvalidDuration
	}

	if input.StartAt.Weekday() == time.Monday {
		return nil, domain.ErrInvalidTeachingDay
	}

	tx, err := s.DB.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.Tutors.LockByID(ctx, tx, input.TutorID); err != nil {
		return nil, err
	}

	dailyCount, err := s.Appointments.CountTutorBookingsOnDay(ctx, tx, input.TutorID, input.StartAt)
	if err != nil {
		return nil, err
	}
	if dailyCount >= maxTutorDailyBookings {
		return nil, domain.ErrTutorDailyLimit
	}

	if existing, _ := s.Appointments.FindConflictsForStudent(ctx, tx, input.StudentID, input.StartAt, input.EndAt); existing != nil {
		return nil, domain.ErrStudentConflict
	}

	if existing, _ := s.Appointments.FindConflictsForTutor(ctx, tx, input.TutorID, input.StartAt, input.EndAt); existing != nil {
		return nil, domain.ErrTutorConflict
	}

	if existing, _ := s.Appointments.FindConflictsForRoom(ctx, tx, input.RoomID, input.StartAt, input.EndAt); existing != nil {
		return nil, domain.ErrRoomConflict
	}

	appointment := &domain.Appointment{
		ID:        s.IDGenerator(),
		StudentID: input.StudentID,
		TutorID:   input.TutorID,
		RoomID:    input.RoomID,
		StartAt:   input.StartAt,
		EndAt:     input.EndAt,
		Status:    domain.StatusBooked,
		CreatedAt: s.Clock(),
	}

	if err := s.Appointments.Create(ctx, tx, appointment); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return appointment, nil
}
