package domain

import "time"

type AppointmentStatus string

const (
	StatusBooked    AppointmentStatus = "BOOKED"
	StatusCancelled AppointmentStatus = "CANCELLED"
	StatusCompleted AppointmentStatus = "COMPLETED"
	StatusNoShow    AppointmentStatus = "NO_SHOW"
)

type Student struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Tutor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Appointment struct {
	ID        string            `json:"id"`
	StudentID string            `json:"studentId"`
	TutorID   string            `json:"tutorId"`
	RoomID    string            `json:"roomId"`
	StartAt   time.Time         `json:"startAt"`
	EndAt     time.Time         `json:"endAt"`
	Status    AppointmentStatus `json:"status"`
	CreatedAt time.Time         `json:"createdAt"`
}

type CreateAppointmentInput struct {
	StudentID string    `json:"studentId" binding:"required"`
	TutorID   string    `json:"tutorId" binding:"required"`
	RoomID    string    `json:"roomId" binding:"required"`
	StartAt   time.Time `json:"startAt" binding:"required"`
	EndAt     time.Time `json:"endAt" binding:"required"`
}

type BookingError struct {
	Code    string `json:"error"`
	Message string `json:"message"`
}

func (e BookingError) Error() string {
	return e.Code + ": " + e.Message
}

var (
	ErrStudentNotFound    = BookingError{"NOT_FOUND", "Student not found"}
	ErrTutorNotFound      = BookingError{"NOT_FOUND", "Tutor not found"}
	ErrRoomNotFound       = BookingError{"NOT_FOUND", "Room not found"}
	ErrInvalidTimeRange   = BookingError{"BAD_REQUEST", "start_at must be before end_at"}
	ErrInvalidDuration    = BookingError{"BAD_REQUEST", "Lesson duration must be 60 or 90 minutes"}
	ErrInvalidTeachingDay = BookingError{"BAD_REQUEST", "Lessons cannot be scheduled on Monday"}
	ErrTimeInPast         = BookingError{"BAD_REQUEST", "Cannot book an appointment in the past"}
	ErrStudentConflict    = BookingError{"STUDENT_CONFLICT", "Student already has a booking during the requested time"}
	ErrTutorConflict      = BookingError{"TUTOR_CONFLICT", "Tutor is already teaching during the requested time"}
	ErrRoomConflict       = BookingError{"ROOM_CONFLICT", "Room is already booked during the requested time"}
	ErrTutorDailyLimit    = BookingError{"TUTOR_DAILY_LIMIT", "Tutor already has 6 bookings on this day"}
)

func IsNotFound(err error) bool {
	if be, ok := err.(BookingError); ok {
		return be.Code == "NOT_FOUND"
	}
	return false
}

func IsConflict(err error) bool {
	if be, ok := err.(BookingError); ok {
		switch be.Code {
		case "STUDENT_CONFLICT", "TUTOR_CONFLICT", "ROOM_CONFLICT", "TUTOR_DAILY_LIMIT":
			return true
		}
	}
	return false
}
