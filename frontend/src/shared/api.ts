import type { Appointment, CreateAppointmentInput, Student, Tutor, Room } from './types';

const API_BASE = '/api';

async function fetchJSON<T>(url: string): Promise<T> {
  const res = await fetch(url);
  if (!res.ok) throw new Error('Request failed');
  return res.json();
}

export async function createAppointment(input: CreateAppointmentInput): Promise<Appointment> {
  const res = await fetch(`${API_BASE}/appointments`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  });

  const data = await res.json();

  if (!res.ok) {
    throw data;
  }

  return data;
}

export async function listAppointments(filters?: {
  date?: string;
  studentId?: string;
  tutorId?: string;
  roomId?: string;
}): Promise<Appointment[]> {
  const params = new URLSearchParams();
  if (filters?.date) params.set('date', filters.date);
  if (filters?.studentId) params.set('studentId', filters.studentId);
  if (filters?.tutorId) params.set('tutorId', filters.tutorId);
  if (filters?.roomId) params.set('roomId', filters.roomId);

  const query = params.toString();
  const url = query ? `${API_BASE}/appointments?${query}` : `${API_BASE}/appointments`;

  return fetchJSON<Appointment[]>(url);
}

export async function listStudents(): Promise<Student[]> {
  return fetchJSON<Student[]>(`${API_BASE}/students`);
}

export async function listTutors(): Promise<Tutor[]> {
  return fetchJSON<Tutor[]>(`${API_BASE}/tutors`);
}

export async function listRooms(): Promise<Room[]> {
  return fetchJSON<Room[]>(`${API_BASE}/rooms`);
}
