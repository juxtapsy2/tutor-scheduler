package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/brightpath/tutor-scheduler/internal/application"
	"github.com/brightpath/tutor-scheduler/internal/domain"
	"github.com/brightpath/tutor-scheduler/internal/repository"
)

type Handler struct {
	createAppointment *application.CreateAppointment
	listAppointments  *application.ListAppointments
}

func NewHandler(create *application.CreateAppointment, list *application.ListAppointments) *Handler {
	return &Handler{createAppointment: create, listAppointments: list}
}

func (h *Handler) CreateAppointment(c *gin.Context) {
	var input domain.CreateAppointmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, domain.BookingError{
			Code:    "BAD_REQUEST",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	appointment, err := h.createAppointment.Execute(c.Request.Context(), input)
	if err != nil {
		var be domain.BookingError
		if errors.As(err, &be) {
			switch be.Code {
			case "NOT_FOUND":
				c.JSON(http.StatusNotFound, be)
			case "BAD_REQUEST":
				c.JSON(http.StatusBadRequest, be)
			default:
				c.JSON(http.StatusConflict, be)
			}
			return
		}
		c.JSON(http.StatusInternalServerError, domain.BookingError{
			Code:    "INTERNAL",
			Message: "An unexpected error occurred",
		})
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

func (h *Handler) ListAppointments(c *gin.Context) {
	filters := repository.ListFilters{}

	if dateStr := c.Query("date"); dateStr != "" {
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, domain.BookingError{
				Code:    "BAD_REQUEST",
				Message: "Invalid date format, use YYYY-MM-DD",
			})
			return
		}
		filters.Date = &t
	}

	if sid := c.Query("studentId"); sid != "" {
		filters.StudentID = &sid
	}
	if tid := c.Query("tutorId"); tid != "" {
		filters.TutorID = &tid
	}
	if rid := c.Query("roomId"); rid != "" {
		filters.RoomID = &rid
	}

	appointments, err := h.listAppointments.Execute(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.BookingError{
			Code:    "INTERNAL",
			Message: "Failed to list appointments",
		})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

func RegisterRoutes(r *gin.Engine, h *Handler) {
	r.Use(corsMiddleware())

	api := r.Group("/api")
	api.POST("/appointments", h.CreateAppointment)
	api.GET("/appointments", h.ListAppointments)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
