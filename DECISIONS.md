# Phase 1 — Questions and Assumptions

Before finalising the scheduling model, these are the first questions I would ask the owner.

## Questions

Question: Can the tutor daily limit of 6 bookings ever be overridden?
--> If yes, then the system needs an explicit override flow and should record who overrode the rule and why.
--> Else, reject any booking that would become the tutor's 7th booking of the day.

Current assumption: No override is allowed.

---

Question: Do cancelled appointments count toward the tutor's 6-booking daily limit?
--> If yes, then cancellation frees the room and time conflict but does not reduce the tutor's daily booking count.
--> Else, cancelled appointments are excluded from the tutor's active daily booking count.

Current assumption: Cancelled appointments do not count toward the active daily limit.

---

Question: Should the current "two students with one tutor in the same room and same slot" arrangement remain supported?
--> If yes, then the data model cannot assume one appointment always has exactly one student, and the conflict rules need an explicit shared-slot exception.
--> Else, keep the normal one-to-one model and reject overlapping student/tutor/room bookings.

Current assumption: This is historical workaround behaviour and should not be supported as a normal scheduling rule.

---

Question: What exactly must be preserved when a booking changes after the 16:00 cut-off?
--> If only the fact that a change occurred matters, then a simple change flag/history record is enough.
--> If the previous schedule must be reconstructable, then preserve the previous tutor, room, time, and booking state.
--> If tutor acknowledgement is required, then notification/acknowledgement state must also be modelled.

Current assumption: Preserve enough history to show what changed, without implementing acknowledgement tracking.

---

Question: Are there exact opening and closing hours for each teaching day?
--> If yes, then appointment creation should reject bookings outside those hours.
--> Else, only enforce the known teaching days and avoid inventing time-of-day restrictions.

Current assumption: Exact daily opening/closing hours are not enforced because they are not specified.

---

Question: At exactly 4 hours before a lesson, is cancellation still free?
--> If yes, then the late-cancellation rule applies only when less than 4 hours remain.
--> Else, then the late-cancellation rule applies when 4 hours or less remain.

Current assumption: Exactly 4 hours before the lesson is still free cancellation.

# Working Assumptions

Until clarified, the implementation will follow these assumptions:

- Maximum 6 active bookings per tutor per day.
- No tutor-limit override.
- Cancelled appointments do not count toward the active daily tutor limit.
- One appointment represents one student, one tutor, and one room.
- The historical two-student shared-slot workaround is not supported.
- Post-16:00 changes preserve visible booking history.
- Exact opening and closing hours are not enforced.
- Exactly 4 hours before lesson start is still considered free cancellation.

# Phase 2 — Feature Selection

From the current workflow and business rules, I see the following features the scheduling system should provide.

## Features

### 1. Booking & Conflict Detection

Create a lesson by assigning a student, tutor, room, and time.

Before accepting the booking, ensure that:

- The student is not already booked elsewhere at that time.
- The tutor is not teaching another lesson at that time.
- The room is not occupied at that time.
- The tutor does not exceed the daily booking limit.
- The booking falls on a valid teaching day.

This feature owns the core scheduling correctness of the system.

---

### 2. Appointment Change / Rescheduling

Move an existing appointment to another:

- Time
- Tutor
- Room

The changed appointment must be checked against the same scheduling rules as a new booking.

Changes after the 16:00 schedule cut-off must remain identifiable rather than silently replacing what the tutor was previously told.

---

### 3. Cancellation

Cancel an existing appointment while preserving its history.

The system would need to:

- Mark the appointment as cancelled.
- Release its room/time.
- Release the tutor/student from scheduling conflicts.
- Determine whether the cancellation occurred within the 4-hour late-cancellation window.
- Preserve the information needed to distinguish normal and late cancellations.

---

### 4. No-show Tracking

Record that a student did not attend a scheduled lesson without treating the appointment as cancelled.

This keeps no-show and cancellation as separate outcomes.

The exact financial consequences require clarification before implementing more than the appointment state.

---

### 5. Tutor Schedule Communication

Provide tutors with their daily schedule and make later changes visible.

This could eventually include:

- Daily schedule delivery.
- Change notifications.
- Cancellation notifications.
- Visibility of changes made after the 16:00 cut-off.

The exact delivery mechanism is not sufficiently specified yet.

---

### 6. Schedule / Availability View

Provide staff with a view of the current schedule across:

- Students
- Tutors
- Rooms
- Time slots

This would make it easier to understand current availability before creating or changing appointments.

---

# Chosen Feature — Booking & Conflict Detection

I would build Booking & Conflict Detection first.

The current workflow depends on one receptionist manually coordinating students, tutors, rooms, and time slots in a spreadsheet. The most damaging failure described is creating a schedule that puts the same person or resource in two places at once.

Preventing those invalid bookings gives the system a reliable foundation before adding workflows around changing, cancelling, or communicating appointments.

It also fits the available implementation time because it can be delivered as one complete vertical slice:

Request
-> validate booking
-> detect conflicts
-> persist safely
-> return booking or clear rejection

Within this feature I would enforce:

- Student time conflicts.
- Tutor time conflicts.
- Room time conflicts.
- Tutor daily booking limit.
- Valid teaching day.
- Concurrent booking safety.

These are treated as business rules of one booking operation rather than separate features.

The implementation should be complete enough that two conflicting booking requests cannot both successfully enter the schedule.

---

## Why This Feature First

### Highest immediate value

A scheduling system is not trustworthy if it can create an impossible schedule.

Preventing conflicting bookings addresses that problem directly.

### Foundation for other features

Rescheduling is effectively:

existing booking
-> proposed new booking
-> run conflict detection
-> apply change

Availability views also depend on knowing which resources are already occupied.

Cancellation releases resources that the conflict detector considers occupied.

Therefore reliable booking rules become reusable foundations for several later features.

### Good fit for the available time

The feature has a small external surface but meaningful internal behaviour.

A minimal implementation only requires:

- Relevant data model.
- Booking endpoint.
- Conflict detection.
- Persistence.
- Concurrency protection.
- Focused tests.
- Seed data.

This is achievable without building unrelated platform functionality.

### Demonstrates correctness rather than breadth

The difficult part is not creating an appointment row.

The important part is guaranteeing that an accepted booking is actually valid, including when multiple booking requests happen concurrently.

I would rather finish this path correctly than partially implement booking, cancellation, notifications, and schedule views.

---

# What I Intentionally Leave Broken / Incomplete

Choosing Booking & Conflict Detection means several existing operational problems remain unsolved.

## Appointment changes

Existing appointments cannot yet be properly rescheduled.

The conflict engine can later be reused for this, but the change/history workflow is not part of this implementation.

## Post-16:00 change tracking

Because appointment changes are not implemented yet, the requirement to preserve changes made after the daily cut-off also remains incomplete.

The data model should avoid making this unnecessarily difficult to add later, but I would not build speculative history infrastructure now.

## Cancellation

Appointments cannot yet be cancelled through the implemented workflow.

The 4-hour late-cancellation behaviour therefore remains unimplemented.

## No-show

Attendance/no-show recording is not implemented.

## Tutor communication

Tutors are not automatically notified of bookings or schedule changes.

The existing communication problem therefore remains.

## Schedule UI

Staff do not receive a complete calendar-style scheduling interface.

The feature may expose enough data through the API to inspect bookings, but a full scheduling UI is not the focus.

---

# Trade-off

The result is intentionally not a complete replacement for the existing spreadsheet workflow.

What it does provide is one reliable building block:

Given a student, tutor, room, and time, the system can decide whether that lesson can safely be added to the schedule and persist that decision without allowing conflicting concurrent bookings.

I consider that more valuable within the available time than providing a wider workflow whose underlying schedule can still become inconsistent.