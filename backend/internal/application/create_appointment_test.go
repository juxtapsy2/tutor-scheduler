package application_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/brightpath/tutor-scheduler/internal/application"
	"github.com/brightpath/tutor-scheduler/internal/domain"
	"github.com/brightpath/tutor-scheduler/internal/repository"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://brightpath:brightpath@localhost:5432/brightpath?sslmode=disable"
	}

	var err error
	testPool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func cleanDB(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	_, err := testPool.Exec(ctx, `DELETE FROM appointments WHERE id NOT LIKE 'L%'`)
	if err != nil {
		t.Fatalf("failed to clean test appointments: %v", err)
	}
}

func newTestService() *application.CreateAppointment {
	return &application.CreateAppointment{
		DB:           testPool,
		Students:     repository.NewStudentRepo(testPool),
		Tutors:       repository.NewTutorRepo(testPool),
		Rooms:        repository.NewRoomRepo(testPool),
		Appointments: repository.NewAppointmentRepo(testPool),
		IDGenerator:  func() string { return "TEST-" + time.Now().Format("150405.000000000") },
		Clock:        time.Now,
	}
}

func TestValidBooking(t *testing.T) {
	cleanDB(t)
	svc := newTestService()
	ctx := context.Background()

	input := domain.CreateAppointmentInput{
		StudentID: "Le Minh Chau",
		TutorID:   "T1",
		RoomID:    "R1",
		StartAt:   time.Date(2027, 4, 1, 10, 0, 0, 0, time.Local),
		EndAt:     time.Date(2027, 4, 1, 11, 0, 0, 0, time.Local),
	}

	appointment, err := svc.Execute(ctx, input)
	if err != nil {
		t.Fatalf("expected valid booking, got error: %v", err)
	}
	if appointment.Status != domain.StatusBooked {
		t.Errorf("expected status BOOKED, got %s", appointment.Status)
	}
}

func TestInvalidTimeRange(t *testing.T) {
	cleanDB(t)
	svc := newTestService()
	ctx := context.Background()

	input := domain.CreateAppointmentInput{
		StudentID: "Le Minh Chau",
		TutorID:   "T1",
		RoomID:    "R1",
		StartAt:   time.Date(2027, 4, 1, 11, 0, 0, 0, time.Local),
		EndAt:     time.Date(2027, 4, 1, 10, 0, 0, 0, time.Local),
	}

	_, err := svc.Execute(ctx, input)
	if err == nil {
		t.Fatal("expected error for invalid time range")
	}
	var be domain.BookingError
	if !errors.As(err, &be) || be.Code != "BAD_REQUEST" {
		t.Errorf("expected BAD_REQUEST, got %v", err)
	}
}

func TestInvalidDuration45Minutes(t *testing.T) {
	cleanDB(t)
	svc := newTestService()
	ctx := context.Background()

	input := domain.CreateAppointmentInput{
		StudentID: "Le Minh Chau",
		TutorID:   "T1",
		RoomID:    "R1",
		StartAt:   time.Date(2027, 4, 1, 10, 0, 0, 0, time.Local),
		EndAt:     time.Date(2027, 4, 1, 10, 45, 0, 0, time.Local),
	}

	_, err := svc.Execute(ctx, input)
	if err == nil {
		t.Fatal("expected error for 45-minute duration")
	}
	var be domain.BookingError
	if !errors.As(err, &be) || be.Code != "BAD_REQUEST" {
		t.Errorf("expected BAD_REQUEST, got %v", err)
	}
}

func TestValidDuration60Minutes(t *testing.T) {
	cleanDB(t)
	svc := newTestService()
	ctx := context.Background()

	input := domain.CreateAppointmentInput{
		StudentID: "Le Minh Chau",
		TutorID:   "T1",
		RoomID:    "R1",
		StartAt:   time.Date(2027, 4, 1, 10, 0, 0, 0, time.Local),
		EndAt:     time.Date(2027, 4, 1, 11, 0, 0, 0, time.Local),
	}

	_, err := svc.Execute(ctx, input)
	if err != nil {
		t.Fatalf("expected 60-minute booking to succeed, got: %v", err)
	}
}

func TestValidDuration90Minutes(t *testing.T) {
	cleanDB(t)
	svc := newTestService()
	ctx := context.Background()

	input := domain.CreateAppointmentInput{
		StudentID: "Le Minh Chau",
		TutorID:   "T1",
		RoomID:    "R1",
		StartAt:   time.Date(2027, 4, 1, 10, 0, 0, 0, time.Local),
		EndAt:     time.Date(2027, 4, 1, 11, 30, 0, 0, time.Local),
	}

	_, err := svc.Execute(ctx, input)
	if err != nil {
		t.Fatalf("expected 90-minute booking to succeed, got: %v", err)
	}
}

func TestMondayRejected(t *testing.T) {
	cleanDB(t)
	svc := newTestService()
	ctx := context.Background()

	input := domain.CreateAppointmentInput{
		StudentID: "Le Minh Chau",
		TutorID:   "T1",
		RoomID:    "R1",
		StartAt:   time.Date(2027, 4, 5, 10, 0, 0, 0, time.Local), // Monday
		EndAt:     time.Date(2027, 4, 5, 11, 0, 0, 0, time.Local),
	}

	_, err := svc.Execute(ctx, input)
	if err == nil {
		t.Fatal("expected error for Monday booking")
	}
	var be domain.BookingError
	if !errors.As(err, &be) || be.Code != "BAD_REQUEST" {
		t.Errorf("expected BAD_REQUEST for Monday, got %v", err)
	}
}

func TestStudentNotFound(t *testing.T) {
	cleanDB(t)
	svc := newTestService()
	ctx := context.Background()

	input := domain.CreateAppointmentInput{
		StudentID: "NONEXISTENT",
		TutorID:   "T1",
		RoomID:    "R1",
		StartAt:   time.Date(2027, 4, 1, 10, 0, 0, 0, time.Local),
		EndAt:     time.Date(2027, 4, 1, 11, 0, 0, 0, time.Local),
	}

	_, err := svc.Execute(ctx, input)
	if err == nil {
		t.Fatal("expected error for nonexistent student")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected NOT_FOUND, got %v", err)
	}
}

func TestTutorNotFound(t *testing.T) {
	cleanDB(t)
	svc := newTestService()
	ctx := context.Background()

	input := domain.CreateAppointmentInput{
		StudentID: "Le Minh Chau",
		TutorID:   "T_NONEXISTENT",
		RoomID:    "R1",
		StartAt:   time.Date(2027, 4, 1, 10, 0, 0, 0, time.Local),
		EndAt:     time.Date(2027, 4, 1, 11, 0, 0, 0, time.Local),
	}

	_, err := svc.Execute(ctx, input)
	if err == nil {
		t.Fatal("expected error for nonexistent tutor")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected NOT_FOUND, got %v", err)
	}
}

func TestRoomNotFound(t *testing.T) {
	cleanDB(t)
	svc := newTestService()
	ctx := context.Background()

	input := domain.CreateAppointmentInput{
		StudentID: "Le Minh Chau",
		TutorID:   "T1",
		RoomID:    "R_NONEXISTENT",
		StartAt:   time.Date(2027, 4, 1, 10, 0, 0, 0, time.Local),
		EndAt:     time.Date(2027, 4, 1, 11, 0, 0, 0, time.Local),
	}

	_, err := svc.Execute(ctx, input)
	if err == nil {
		t.Fatal("expected error for nonexistent room")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected NOT_FOUND, got %v", err)
	}
}

func TestRoomConflict(t *testing.T) {
	cleanDB(t)
	svc := newTestService()
	ctx := context.Background()

	// First booking
	input1 := domain.CreateAppointmentInput{
		StudentID: "Le Minh Chau",
		TutorID:   "T1",
		RoomID:    "R1",
		StartAt:   time.Date(2027, 4, 1, 10, 0, 0, 0, time.Local),
		EndAt:     time.Date(2027, 4, 1, 11, 0, 0, 0, time.Local),
	}
	if _, err := svc.Execute(ctx, input1); err != nil {
		t.Fatalf("first booking should succeed: %v", err)
	}

	// Conflicting booking: same room, overlapping time, different student/tutor
	input2 := domain.CreateAppointmentInput{
		StudentID: "Tran Bao Long",
		TutorID:   "T2",
		RoomID:    "R1",
		StartAt:   time.Date(2027, 4, 1, 10, 30, 0, 0, time.Local),
		EndAt:     time.Date(2027, 4, 1, 11, 30, 0, 0, time.Local),
	}
	_, err := svc.Execute(ctx, input2)
	if err == nil {
		t.Fatal("expected room conflict error")
	}
	if !domain.IsConflict(err) {
		t.Errorf("expected conflict error, got %v", err)
	}
}

func TestTutorConflict(t *testing.T) {
	cleanDB(t)
	svc := newTestService()
	ctx := context.Background()

	// First booking
	input1 := domain.CreateAppointmentInput{
		StudentID: "Le Minh Chau",
		TutorID:   "T1",
		RoomID:    "R1",
		StartAt:   time.Date(2027, 4, 1, 10, 0, 0, 0, time.Local),
		EndAt:     time.Date(2027, 4, 1, 11, 0, 0, 0, time.Local),
	}
	if _, err := svc.Execute(ctx, input1); err != nil {
		t.Fatalf("first booking should succeed: %v", err)
	}

	// Conflicting booking: same tutor, overlapping time, different room
	input2 := domain.CreateAppointmentInput{
		StudentID: "Tran Bao Long",
		TutorID:   "T1",
		RoomID:    "R2",
		StartAt:   time.Date(2027, 4, 1, 10, 30, 0, 0, time.Local),
		EndAt:     time.Date(2027, 4, 1, 11, 30, 0, 0, time.Local),
	}
	_, err := svc.Execute(ctx, input2)
	if err == nil {
		t.Fatal("expected tutor conflict error")
	}
	if !domain.IsConflict(err) {
		t.Errorf("expected conflict error, got %v", err)
	}
}

func TestStudentConflict(t *testing.T) {
	cleanDB(t)
	svc := newTestService()
	ctx := context.Background()

	// First booking
	input1 := domain.CreateAppointmentInput{
		StudentID: "Le Minh Chau",
		TutorID:   "T1",
		RoomID:    "R1",
		StartAt:   time.Date(2027, 4, 1, 10, 0, 0, 0, time.Local),
		EndAt:     time.Date(2027, 4, 1, 11, 0, 0, 0, time.Local),
	}
	if _, err := svc.Execute(ctx, input1); err != nil {
		t.Fatalf("first booking should succeed: %v", err)
	}

	// Conflicting booking: same student, overlapping time, different tutor/room
	input2 := domain.CreateAppointmentInput{
		StudentID: "Le Minh Chau",
		TutorID:   "T2",
		RoomID:    "R2",
		StartAt:   time.Date(2027, 4, 1, 10, 30, 0, 0, time.Local),
		EndAt:     time.Date(2027, 4, 1, 11, 30, 0, 0, time.Local),
	}
	_, err := svc.Execute(ctx, input2)
	if err == nil {
		t.Fatal("expected student conflict error")
	}
	if !domain.IsConflict(err) {
		t.Errorf("expected conflict error, got %v", err)
	}
}

func TestAdjacentAppointmentsAllowed(t *testing.T) {
	cleanDB(t)
	svc := newTestService()
	ctx := context.Background()

	// First booking: 10:00 - 11:00
	input1 := domain.CreateAppointmentInput{
		StudentID: "Le Minh Chau",
		TutorID:   "T1",
		RoomID:    "R1",
		StartAt:   time.Date(2027, 4, 1, 10, 0, 0, 0, time.Local),
		EndAt:     time.Date(2027, 4, 1, 11, 0, 0, 0, time.Local),
	}
	if _, err := svc.Execute(ctx, input1); err != nil {
		t.Fatalf("first booking should succeed: %v", err)
	}

	// Adjacent booking: 11:00 - 12:00 (no overlap)
	input2 := domain.CreateAppointmentInput{
		StudentID: "Tran Bao Long",
		TutorID:   "T2",
		RoomID:    "R2",
		StartAt:   time.Date(2027, 4, 1, 11, 0, 0, 0, time.Local),
		EndAt:     time.Date(2027, 4, 1, 12, 0, 0, 0, time.Local),
	}
	if _, err := svc.Execute(ctx, input2); err != nil {
		t.Fatalf("adjacent booking should succeed, got: %v", err)
	}
}

func TestTutorDailyLimit(t *testing.T) {
	cleanDB(t)
	svc := newTestService()
	ctx := context.Background()

	// Book 6 slots for tutor T1 on the same day
	for i := 0; i < 6; i++ {
		start := time.Date(2027, 4, 2, 9+i, 0, 0, 0, time.Local)
		end := start.Add(60 * time.Minute)
		input := domain.CreateAppointmentInput{
			StudentID: "Le Minh Chau",
			TutorID:   "T1",
			RoomID:    "R1",
			StartAt:   start,
			EndAt:     end,
		}
		// Use different students to avoid student conflict
		switch i % 6 {
		case 0:
			input.StudentID = "Le Minh Chau"
		case 1:
			input.StudentID = "Tran Bao Long"
		case 2:
			input.StudentID = "Nguyen Thi Ha"
		case 3:
			input.StudentID = "Do Van Kien"
		case 4:
			input.StudentID = "Vu Ha My"
		case 5:
			input.StudentID = "Bui An Nhien"
		}
		if _, err := svc.Execute(ctx, input); err != nil {
			t.Fatalf("booking %d should succeed: %v", i+1, err)
		}
	}

	// 7th booking should fail
	input := domain.CreateAppointmentInput{
		StudentID: "Le Minh Chau",
		TutorID:   "T1",
		RoomID:    "R2",
		StartAt:   time.Date(2027, 4, 2, 16, 0, 0, 0, time.Local),
		EndAt:     time.Date(2027, 4, 2, 17, 0, 0, 0, time.Local),
	}
	_, err := svc.Execute(ctx, input)
	if err == nil {
		t.Fatal("expected tutor daily limit error")
	}
	if !domain.IsConflict(err) {
		t.Errorf("expected conflict error, got %v", err)
	}
}
