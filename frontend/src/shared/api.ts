import type { Appointment, CreateAppointmentInput } from './types';

const API_BASE = '/api';

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

  const res = await fetch(url);
  if (!res.ok) {
    throw new Error('Failed to list appointments');
  }
  return res.json();
}
