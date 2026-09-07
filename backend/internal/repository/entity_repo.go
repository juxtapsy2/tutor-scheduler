package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/brightpath/tutor-scheduler/internal/domain"
)

type StudentRepo struct {
	db *pgxpool.Pool
}

func NewStudentRepo(db *pgxpool.Pool) *StudentRepo {
	return &StudentRepo{db: db}
}

func (r *StudentRepo) FindByID(ctx context.Context, id string) (*domain.Student, error) {
	var s domain.Student
	err := r.db.QueryRow(ctx, `SELECT id, name FROM students WHERE id = $1`, id).
		Scan(&s.ID, &s.Name)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

type TutorRepo struct {
	db *pgxpool.Pool
}

func NewTutorRepo(db *pgxpool.Pool) *TutorRepo {
	return &TutorRepo{db: db}
}

func (r *TutorRepo) FindByID(ctx context.Context, id string) (*domain.Tutor, error) {
	var t domain.Tutor
	err := r.db.QueryRow(ctx, `SELECT id, name FROM tutors WHERE id = $1`, id).
		Scan(&t.ID, &t.Name)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TutorRepo) LockByID(ctx context.Context, tx pgx.Tx, id string) error {
	_, err := tx.Exec(ctx, `SELECT id FROM tutors WHERE id = $1 FOR SHARE`, id)
	return err
}

type RoomRepo struct {
	db *pgxpool.Pool
}

func NewRoomRepo(db *pgxpool.Pool) *RoomRepo {
	return &RoomRepo{db: db}
}

func (r *RoomRepo) FindByID(ctx context.Context, id string) (*domain.Room, error) {
	var rm domain.Room
	err := r.db.QueryRow(ctx, `SELECT id, name FROM rooms WHERE id = $1`, id).
		Scan(&rm.ID, &rm.Name)
	if err != nil {
		return nil, err
	}
	return &rm, nil
}

func (r *RoomRepo) ListAll(ctx context.Context) ([]domain.Room, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name FROM rooms ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []domain.Room
	for rows.Next() {
		var rm domain.Room
		if err := rows.Scan(&rm.ID, &rm.Name); err != nil {
			return nil, err
		}
		rooms = append(rooms, rm)
	}
	return rooms, rows.Err()
}
