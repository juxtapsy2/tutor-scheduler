import { useState, useEffect } from 'react';
import { createAppointment, listStudents, listTutors, listRooms } from '../../shared/api';
import type { BookingError, Student, Tutor, Room } from '../../shared/types';

interface BookingFormProps {
  onBooked: () => void;
}

export default function BookingForm({ onBooked }: BookingFormProps) {
  const [students, setStudents] = useState<Student[]>([]);
  const [tutors, setTutors] = useState<Tutor[]>([]);
  const [rooms, setRooms] = useState<Room[]>([]);

  const [studentId, setStudentId] = useState('');
  const [tutorId, setTutorId] = useState('');
  const [roomId, setRoomId] = useState('');
  const [date, setDate] = useState('');
  const [startTime, setStartTime] = useState('');
  const [duration, setDuration] = useState<60 | 90>(60);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    Promise.all([listStudents(), listTutors(), listRooms()])
      .then(([s, t, r]) => {
        setStudents(s);
        setTutors(t);
        setRooms(r);
      })
      .catch(() => {});
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    const startAt = new Date(`${date}T${startTime}:00`);
    const endAt = new Date(startAt.getTime() + duration * 60 * 1000);

    try {
      await createAppointment({
        studentId,
        tutorId,
        roomId,
        startAt: startAt.toISOString(),
        endAt: endAt.toISOString(),
      });
      onBooked();
      setStudentId('');
      setTutorId('');
      setRoomId('');
      setDate('');
      setStartTime('');
      setDuration(60);
    } catch (err) {
      const be = err as BookingError;
      setError(`${be.error}: ${be.message}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="bg-white border border-gray-200 rounded-lg p-6 mb-6">
      <h2 className="text-lg font-semibold mb-4">New Booking</h2>

      <div className="flex gap-4 mb-4 flex-wrap">
        <label className="flex flex-col flex-1 min-w-[140px] text-sm font-medium">
          Student
          <select
            className="mt-1 px-3 py-2 border border-gray-300 rounded text-sm"
            value={studentId}
            onChange={e => setStudentId(e.target.value)}
            required
          >
            <option value="">Select student</option>
            {students.map(s => (
              <option key={s.id} value={s.id}>{s.name}</option>
            ))}
          </select>
        </label>
        <label className="flex flex-col flex-1 min-w-[140px] text-sm font-medium">
          Tutor
          <select
            className="mt-1 px-3 py-2 border border-gray-300 rounded text-sm"
            value={tutorId}
            onChange={e => setTutorId(e.target.value)}
            required
          >
            <option value="">Select tutor</option>
            {tutors.map(t => (
              <option key={t.id} value={t.id}>{t.name}</option>
            ))}
          </select>
        </label>
        <label className="flex flex-col flex-1 min-w-[140px] text-sm font-medium">
          Room
          <select
            className="mt-1 px-3 py-2 border border-gray-300 rounded text-sm"
            value={roomId}
            onChange={e => setRoomId(e.target.value)}
            required
          >
            <option value="">Select room</option>
            {rooms.map(r => (
              <option key={r.id} value={r.id}>{r.name} ({r.id})</option>
            ))}
          </select>
        </label>
      </div>

      <div className="flex gap-4 mb-4 flex-wrap">
        <label className="flex flex-col flex-1 min-w-[140px] text-sm font-medium">
          Date
          <input
            type="date"
            className="mt-1 px-3 py-2 border border-gray-300 rounded text-sm"
            value={date}
            onChange={e => setDate(e.target.value)}
            required
          />
        </label>
        <label className="flex flex-col flex-1 min-w-[140px] text-sm font-medium">
          Start Time
          <input
            type="time"
            className="mt-1 px-3 py-2 border border-gray-300 rounded text-sm"
            value={startTime}
            onChange={e => setStartTime(e.target.value)}
            required
          />
        </label>
        <label className="flex flex-col flex-1 min-w-[140px] text-sm font-medium">
          Duration
          <select
            className="mt-1 px-3 py-2 border border-gray-300 rounded text-sm"
            value={duration}
            onChange={e => setDuration(Number(e.target.value) as 60 | 90)}
          >
            <option value={60}>60 minutes</option>
            <option value={90}>90 minutes</option>
          </select>
        </label>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4 text-sm">
          {error}
        </div>
      )}

      <button
        type="submit"
        disabled={loading}
        className="bg-blue-600 text-white px-5 py-2 rounded text-sm hover:bg-blue-700 disabled:opacity-60 disabled:cursor-not-allowed"
      >
        {loading ? 'Booking...' : 'Book Lesson'}
      </button>
    </form>
  );
}
