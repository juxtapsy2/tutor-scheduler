package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/brightpath/tutor-scheduler/internal/application"
	"github.com/brightpath/tutor-scheduler/internal/db"
	"github.com/brightpath/tutor-scheduler/internal/handler"
	"github.com/brightpath/tutor-scheduler/internal/repository"
)

func main() {
	ctx := context.Background()

	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer pool.Close()

	studentRepo := repository.NewStudentRepo(pool)
	tutorRepo := repository.NewTutorRepo(pool)
	roomRepo := repository.NewRoomRepo(pool)
	appointmentRepo := repository.NewAppointmentRepo(pool)

	createAppointment := &application.CreateAppointment{
		DB:           pool,
		Students:     studentRepo,
		Tutors:       tutorRepo,
		Rooms:        roomRepo,
		Appointments: appointmentRepo,
		IDGenerator:  func() string { return uuid.New().String() },
		Clock:        time.Now,
	}

	listAppointments := &application.ListAppointments{
		Appointments: appointmentRepo,
	}

	h := handler.NewHandler(createAppointment, listAppointments)

	r := gin.Default()
	handler.RegisterRoutes(r, h)

	fmt.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
