export interface Appointment {
  id: string;
  studentId: string;
  tutorId: string;
  roomId: string;
  startAt: string;
  endAt: string;
  status: string;
  createdAt: string;
}

export interface BookingError {
  error: string;
  message: string;
}

export interface CreateAppointmentInput {
  studentId: string;
  tutorId: string;
  roomId: string;
  startAt: string;
  endAt: string;
}
