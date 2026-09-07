package application

import (
	"context"

	"github.com/brightpath/tutor-scheduler/internal/domain"
	"github.com/brightpath/tutor-scheduler/internal/repository"
)

type ListAppointments struct {
	Appointments *repository.AppointmentRepo
}

func (s *ListAppointments) Execute(ctx context.Context, filters repository.ListFilters) ([]domain.Appointment, error) {
	return s.Appointments.List(ctx, filters)
}
