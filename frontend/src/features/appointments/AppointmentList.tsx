import { useState, useEffect } from 'react';
import { listAppointments } from '../../shared/api';
import type { Appointment } from '../../shared/types';

interface AppointmentListProps {
  refreshKey: number;
}

const PAGE_SIZE = 10;

export default function AppointmentList({ refreshKey }: AppointmentListProps) {
  const [appointments, setAppointments] = useState<Appointment[]>([]);
  const [date, setDate] = useState('');
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);

  const fetchAppointments = async () => {
    setLoading(true);
    try {
      const filters = date ? { date } : undefined;
      const data = await listAppointments(filters);
      setAppointments(data.reverse());
      setPage(1);
    } catch {
      setAppointments([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAppointments();
  }, [date, refreshKey]);

  const totalPages = Math.max(1, Math.ceil(appointments.length / PAGE_SIZE));
  const pageData = appointments.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE);

  return (
    <div className="bg-white border border-gray-200 rounded-lg p-4 h-full flex flex-col">
      <div className="flex items-end gap-4 mb-3 shrink-0">
        <h2 className="text-lg font-semibold">Appointments</h2>
        <label className="flex flex-col text-sm font-medium ml-auto">
          <input
            type="date"
            className="px-3 py-1.5 border border-gray-300 rounded text-sm"
            value={date}
            onChange={e => setDate(e.target.value)}
          />
        </label>
        <button
          onClick={() => setDate('')}
          className="px-3 py-1.5 border border-gray-300 rounded text-sm text-gray-700 hover:bg-gray-50"
        >
          Clear
        </button>
      </div>

      {loading ? (
        <p className="text-gray-500 text-sm">Loading...</p>
      ) : (
        <div className="flex-1 min-h-0 flex flex-col">
          <div className="flex-1 min-h-0 overflow-y-auto border border-gray-200 rounded">
            <table className="w-full text-sm">
              <thead className="sticky top-0 bg-gray-50">
                <tr className="border-b border-gray-200">
                  <th className="text-left py-1.5 px-2 font-semibold text-xs">ID</th>
                  <th className="text-left py-1.5 px-2 font-semibold text-xs">Student</th>
                  <th className="text-left py-1.5 px-2 font-semibold text-xs">Tutor</th>
                  <th className="text-left py-1.5 px-2 font-semibold text-xs">Room</th>
                  <th className="text-left py-1.5 px-2 font-semibold text-xs">Start</th>
                  <th className="text-left py-1.5 px-2 font-semibold text-xs">End</th>
                  <th className="text-left py-1.5 px-2 font-semibold text-xs">Status</th>
                </tr>
              </thead>
              <tbody>
                {pageData.length === 0 ? (
                  <tr>
                    <td colSpan={7} className="py-8 text-center text-gray-400 text-sm">
                      No appointments found.
                    </td>
                  </tr>
                ) : (
                  pageData.map(a => (
                    <tr
                      key={a.id}
                      className={`border-b border-gray-100 ${
                        a.status === 'CANCELLED'
                          ? 'text-gray-400 line-through'
                          : a.status === 'NO_SHOW'
                            ? 'text-red-600'
                            : a.status === 'COMPLETED'
                              ? 'text-gray-400'
                              : ''
                      }`}
                    >
                      <td className="py-1.5 px-2">{a.id}</td>
                      <td className="py-1.5 px-2">{a.studentId}</td>
                      <td className="py-1.5 px-2">{a.tutorId}</td>
                      <td className="py-1.5 px-2">{a.roomId}</td>
                      <td className="py-1.5 px-2">{new Date(a.startAt).toLocaleString()}</td>
                      <td className="py-1.5 px-2">{new Date(a.endAt).toLocaleString()}</td>
                      <td className="py-1.5 px-2">{a.status}</td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>

          <div className="flex items-center justify-between mt-2 text-xs text-gray-600 shrink-0">
            <span>
              {appointments.length} appointment{appointments.length !== 1 ? 's' : ''}
              {' '}&middot;{' '}Page {page} of {totalPages}
            </span>
            <div className="flex gap-2">
              <button
                onClick={() => setPage(p => Math.max(1, p - 1))}
                disabled={page <= 1}
                className="px-2 py-0.5 border border-gray-300 rounded text-xs disabled:opacity-40 disabled:cursor-not-allowed hover:bg-gray-50"
              >
                Prev
              </button>
              <button
                onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                disabled={page >= totalPages}
                className="px-2 py-0.5 border border-gray-300 rounded text-xs disabled:opacity-40 disabled:cursor-not-allowed hover:bg-gray-50"
              >
                Next
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
