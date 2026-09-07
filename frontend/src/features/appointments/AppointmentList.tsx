import { useState, useEffect } from 'react';
import { listAppointments } from '../../shared/api';
import type { Appointment } from '../../shared/types';

interface AppointmentListProps {
  refreshKey: number;
}

export default function AppointmentList({ refreshKey }: AppointmentListProps) {
  const [appointments, setAppointments] = useState<Appointment[]>([]);
  const [date, setDate] = useState('');
  const [loading, setLoading] = useState(false);

  const fetchAppointments = async () => {
    setLoading(true);
    try {
      const filters = date ? { date } : undefined;
      const data = await listAppointments(filters);
      setAppointments(data);
    } catch {
      setAppointments([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAppointments();
  }, [date, refreshKey]);

  return (
    <div className="bg-white border border-gray-200 rounded-lg p-6">
      <h2 className="text-lg font-semibold mb-4">Appointments</h2>

      <div className="flex gap-4 items-end mb-4">
        <label className="flex flex-col text-sm font-medium">
          Filter by date
          <input
            type="date"
            className="mt-1 px-3 py-2 border border-gray-300 rounded text-sm"
            value={date}
            onChange={e => setDate(e.target.value)}
          />
        </label>
        <button
          onClick={() => setDate('')}
          className="px-4 py-2 border border-gray-300 rounded text-sm text-gray-700 hover:bg-gray-50"
        >
          Clear
        </button>
      </div>

      {loading ? (
        <p className="text-gray-500 text-sm">Loading...</p>
      ) : appointments.length === 0 ? (
        <p className="text-gray-500 text-sm">No appointments found.</p>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-gray-200">
                <th className="text-left py-2 px-3 font-semibold bg-gray-50">ID</th>
                <th className="text-left py-2 px-3 font-semibold bg-gray-50">Student</th>
                <th className="text-left py-2 px-3 font-semibold bg-gray-50">Tutor</th>
                <th className="text-left py-2 px-3 font-semibold bg-gray-50">Room</th>
                <th className="text-left py-2 px-3 font-semibold bg-gray-50">Start</th>
                <th className="text-left py-2 px-3 font-semibold bg-gray-50">End</th>
                <th className="text-left py-2 px-3 font-semibold bg-gray-50">Status</th>
              </tr>
            </thead>
            <tbody>
              {appointments.map(a => (
                <tr
                  key={a.id}
                  className={`border-b border-gray-200 ${
                    a.status === 'CANCELLED'
                      ? 'text-gray-400 line-through'
                      : a.status === 'NO_SHOW'
                        ? 'text-red-600'
                        : a.status === 'COMPLETED'
                          ? 'text-gray-400'
                          : ''
                  }`}
                >
                  <td className="py-2 px-3">{a.id}</td>
                  <td className="py-2 px-3">{a.studentId}</td>
                  <td className="py-2 px-3">{a.tutorId}</td>
                  <td className="py-2 px-3">{a.roomId}</td>
                  <td className="py-2 px-3">{new Date(a.startAt).toLocaleString()}</td>
                  <td className="py-2 px-3">{new Date(a.endAt).toLocaleString()}</td>
                  <td className="py-2 px-3">{a.status}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
