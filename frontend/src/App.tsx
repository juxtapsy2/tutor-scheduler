import { useState } from 'react';
import BookingForm from './features/booking/BookingForm';
import AppointmentList from './features/appointments/AppointmentList';

function App() {
  const [refreshKey, setRefreshKey] = useState(0);

  const handleBooked = () => {
    setRefreshKey(k => k + 1);
  };

  return (
    <div className="h-dvh bg-gray-100 overflow-y-auto">
      <div className="max-w-[960px] mx-auto px-4 py-8">
        <header className="mb-8">
          <h1 className="text-xl font-semibold">Bright Path Scheduler</h1>
        </header>
        <main>
          <BookingForm onBooked={handleBooked} />
          <AppointmentList refreshKey={refreshKey} />
        </main>
      </div>
    </div>
  );
}

export default App;
