CREATE EXTENSION IF NOT EXISTS btree_gist;

DO $$ BEGIN
    CREATE TYPE appointment_status AS ENUM ('BOOKED', 'CANCELLED', 'COMPLETED', 'NO_SHOW');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS students (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tutors (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS rooms (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS appointments (
    id TEXT PRIMARY KEY,
    student_id TEXT NOT NULL REFERENCES students(id),
    tutor_id TEXT NOT NULL REFERENCES tutors(id),
    room_id TEXT NOT NULL REFERENCES rooms(id),
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    status appointment_status NOT NULL DEFAULT 'BOOKED',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT valid_time_range CHECK (start_at < end_at)
);

ALTER TABLE appointments
    ADD CONSTRAINT no_room_overlap
    EXCLUDE USING gist (
        room_id WITH =,
        tstzrange(start_at, end_at, '[)') WITH &&
    ) WHERE (status = 'BOOKED');

ALTER TABLE appointments
    ADD CONSTRAINT no_tutor_overlap
    EXCLUDE USING gist (
        tutor_id WITH =,
        tstzrange(start_at, end_at, '[)') WITH &&
    ) WHERE (status = 'BOOKED');

ALTER TABLE appointments
    ADD CONSTRAINT no_student_overlap
    EXCLUDE USING gist (
        student_id WITH =,
        tstzrange(start_at, end_at, '[)') WITH &&
    ) WHERE (status = 'BOOKED');

CREATE INDEX IF NOT EXISTS idx_appointments_student_time ON appointments (student_id, start_at, end_at);
CREATE INDEX IF NOT EXISTS idx_appointments_tutor_time ON appointments (tutor_id, start_at, end_at);
CREATE INDEX IF NOT EXISTS idx_appointments_room_time ON appointments (room_id, start_at, end_at);
CREATE INDEX IF NOT EXISTS idx_appointments_tutor_day ON appointments (tutor_id, start_at);
