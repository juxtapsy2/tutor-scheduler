import { useState } from 'react';
import BookingForm from './features/booking/BookingForm';
import AppointmentList from './features/appointments/AppointmentList';

function App() {
  const [refreshKey, setRefreshKey] = useState(0);

  const handleBooked = () => {
    setRefreshKey(k => k + 1);
  };

  return (
    <div className="h-dvh flex flex-col bg-gray-100 overflow-hidden">
      <header className="shrink-0 px-6 pt-4 pb-2">
        <h1 className="text-xl font-semibold">Bright Path Scheduler</h1>
      </header>
      <main className="flex-1 min-h-0 px-6 pb-4 flex flex-col gap-4 overflow-y-auto">
        <BookingForm onBooked={handleBooked} />
        <div className="flex-1 min-h-0">
          <AppointmentList refreshKey={refreshKey} />
        </div>
      </main>
    </div>
  );
}

export default App;
