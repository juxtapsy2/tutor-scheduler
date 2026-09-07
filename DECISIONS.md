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

# Phase 3 — Booking & Conflict Detection Design

The feature selected for implementation is Booking & Conflict Detection.

The responsibility of this feature is:

Given a student, tutor, room, and lesson time, either create a valid booking or reject it with a clear reason.

A successful booking must never leave the schedule in a state that violates the rules defined for students, tutors, rooms, or tutor daily capacity.

The intended flow is:

Booking request
-> validate input
-> validate business rules
-> persist atomically
-> return created appointment

If any scheduling rule fails:

Booking request
-> detect conflict
-> reject booking
-> persist nothing


## Data Model

The booking feature only needs four core entities:

- Student
- Tutor
- Room
- Appointment

The model should remain intentionally small.

### Student

Student
- id
- name

Only identity is required by this feature.

Other student information is outside the scheduling responsibility.


### Tutor

Tutor
- id
- name

Tutor availability is derived from existing appointments.

I would not create a separate tutor availability table unless the business later introduces explicit working schedules, leave, or blocked periods.


### Room

Room
- id
- name / number

There are six rooms.

Availability is derived from appointments using the room.


### Appointment

Appointment
- id
- student_id
- tutor_id
- room_id
- start_at
- end_at
- status
- created_at

Relationships:

Student 1 ---- * Appointment
Tutor   1 ---- * Appointment
Room    1 ---- * Appointment

Each appointment represents one lesson involving exactly:

- one student
- one tutor
- one room

This follows the intended one-to-one lesson model.

The existing historical two-student shared-slot workaround is not modelled as valid future scheduling behaviour.


## Appointment Status

Although cancellation is not part of the feature being implemented, historical seed data may contain appointments that are not currently active.

For that reason, appointment status belongs in the model even though this phase only creates new BOOKED appointments.

A minimal representation is:

BOOKED
CANCELLED
COMPLETED
NO_SHOW

New appointments created through this feature start as:

BOOKED

Conflict detection should only consider appointments whose state still occupies scheduling resources.

For the current implementation:

BOOKED -> blocks scheduling resources
CANCELLED -> does not block scheduling resources
COMPLETED -> historical; does not affect future scheduling
NO_SHOW -> historical; does not affect future scheduling

The API implemented in this phase does not expose cancellation or attendance-state transitions.

The status field exists because the scheduling model must correctly understand imported historical data, not because those additional workflows are being implemented now.


## Time Representation

Appointments should use complete timestamps rather than separate string date/time fields.

Preferred representation:

start_at
end_at

The invariant is:

start_at < end_at

Time intervals are treated as half-open:

[start_at, end_at)

Therefore:

10:00 - 11:00
11:00 - 12:00

do not overlap.

But:

10:00 - 11:00
10:30 - 11:30

do overlap.

The overlap condition is:

existing.start_at < requested.end_at
AND
existing.end_at > requested.start_at

This same definition must be used for room, tutor, and student conflict detection.


## Lesson Duration

The centre uses 60-minute and 90-minute lessons.

Therefore a new appointment must have a duration of either:

60 minutes
or
90 minutes

Any other duration should be rejected.

This can be calculated from:

end_at - start_at

rather than trusting a separate duration value supplied by the client.


## Teaching Day Rule

The centre operates Tuesday through Sunday.

Monday is closed.

Therefore:

If appointment date is Monday
-> reject the booking.

Else
-> continue validation.

Exact daily opening and closing times are not enforced because they are not specified.


## Student Conflict Rule

A student cannot be booked into two places at the same time.

For a requested appointment:

Find active appointments where:

student_id = requested student
AND time overlaps requested time

If one exists
-> reject with STUDENT_CONFLICT.

Else
-> continue.

A student may have multiple lessons on the same day if they do not overlap.


## Tutor Conflict Rule

A tutor can only teach one lesson at a time.

For a requested appointment:

Find active appointments where:

tutor_id = requested tutor
AND time overlaps requested time

If one exists
-> reject with TUTOR_CONFLICT.

Else
-> continue.

Different rooms do not make overlapping tutor bookings valid.


## Room Conflict Rule

A room can contain only one lesson at a time.

For a requested appointment:

Find active appointments where:

room_id = requested room
AND time overlaps requested time

If one exists
-> reject with ROOM_CONFLICT.

Else
-> continue.

The identity of the tutor does not matter for this rule.


## Tutor Daily Capacity Rule

A tutor may have at most six bookings in one day.

Before inserting a booking:

Count active bookings for:

requested tutor
on requested local calendar day

If count >= 6
-> reject with TUTOR_DAILY_LIMIT.

Else
-> continue.

Current assumption:

CANCELLED appointments do not count toward the six active bookings.

The seventh booking must not be silently accepted.


## Validation Order

Validation should fail as early as practical.

Suggested order:

1. Required fields exist.
2. Referenced student exists.
3. Referenced tutor exists.
4. Referenced room exists.
5. start_at < end_at.
6. Duration is 60 or 90 minutes.
7. Appointment is on a valid teaching day.
8. Tutor daily capacity is available.
9. Student has no overlapping booking.
10. Tutor has no overlapping booking.
11. Room has no overlapping booking.
12. Persist appointment.

The exact conflict-check order is not a business requirement.

It only affects which error is returned when a request violates multiple rules at the same time.

The database remains responsible for ensuring that concurrency cannot bypass these checks.


## Concurrency Problem

Application-level validation alone is not enough.

Consider:

Request A:
checks Room 2 at 10:00
-> available

Request B:
checks Room 2 at 10:00
-> available

Request A:
creates booking

Request B:
creates booking

Without database-level concurrency protection, both requests can succeed.

That would violate the central purpose of this feature.


## Concurrency Strategy

Recommended database: PostgreSQL.

PostgreSQL is useful here because overlapping time ranges can be protected using exclusion constraints.

The database can enforce that active appointments do not overlap for the same:

- room
- tutor
- student

Conceptually:

same room + overlapping time range -> forbidden
same tutor + overlapping time range -> forbidden
same student + overlapping time range -> forbidden

This protection occurs during the database write itself.

Therefore two concurrent requests cannot both commit conflicting appointments even if both application requests initially observed the schedule as available.

Application conflict checks still exist because they allow us to return understandable errors before attempting the insert.

The database constraint is the final correctness boundary.


## Database-Enforced Rules

The database should enforce structural invariants wherever practical.

Database responsibilities:

- Appointment primary key.
- Student foreign key.
- Tutor foreign key.
- Room foreign key.
- Required appointment fields.
- Valid appointment status.
- start_at must be before end_at.
- Prevent overlapping active appointments for the same room.
- Prevent overlapping active appointments for the same tutor.
- Prevent overlapping active appointments for the same student.

The overlap constraints are particularly important because they protect against concurrent requests.


## Application-Enforced Rules

The application layer should enforce:

- Lesson duration must be 60 or 90 minutes.
- Monday bookings are not allowed.
- Tutor maximum six active bookings per day.
- Friendly conflict error responses.
- Request validation.
- Translation of database constraint failures into domain errors.

These rules are easier to express and explain in business logic than as complex database constraints.

The application should not assume its own checks are sufficient for concurrency-sensitive overlap rules.


## Tutor Daily Limit Concurrency

The six-booking daily limit also has a concurrency problem.

Example:

Tutor currently has 5 bookings.

Request A:
counts 5
-> allowed

Request B:
counts 5
-> allowed

Both insert.

Tutor now has 7 bookings.

Therefore the daily-limit check must participate in the same transaction/concurrency strategy as booking creation.

A simple approach is to lock the tutor record while creating an appointment.

Booking transaction:

BEGIN

lock tutor row

count tutor bookings for requested day

if count >= 6
    reject

perform remaining checks / insert appointment

COMMIT

Booking requests for different tutors can still proceed concurrently.

Booking requests for the same tutor are serialised around the capacity decision.

This is simpler than building a separate capacity-counter system.


## Transaction Boundary

Creating a booking is one atomic operation.

Conceptually:

BEGIN

validate tutor daily capacity

insert appointment

COMMIT

If validation fails:

ROLLBACK

If an exclusion constraint detects a concurrent conflict:

ROLLBACK

No partial booking should remain in the database.


## Why Both Application Checks and Database Constraints?

Application checks provide good user-facing behaviour.

Example:

ROOM_CONFLICT

is more useful than exposing a raw database constraint error.

Database constraints provide correctness.

The application can be wrong, stale, or racing with another request.

Therefore:

Application layer
-> explains the rule

Database
-> guarantees the invariant

For concurrency-sensitive scheduling rules, I prefer having both.


## Booking API

The implemented API surface can remain very small.

Primary endpoint:

POST /appointments

Request:

{
  "studentId": "...",
  "tutorId": "...",
  "roomId": "...",
  "startAt": "...",
  "endAt": "..."
}

Successful response:

201 Created

{
  "id": "...",
  "studentId": "...",
  "tutorId": "...",
  "roomId": "...",
  "startAt": "...",
  "endAt": "...",
  "status": "BOOKED"
}


## Expected Booking Errors

400 Bad Request

Used for invalid input such as:

- invalid time range
- unsupported lesson duration
- booking on Monday


404 Not Found

Used when:

- student does not exist
- tutor does not exist
- room does not exist


409 Conflict

Used when the requested booking is valid in structure but cannot coexist with the current schedule.

Possible domain errors:

STUDENT_CONFLICT
TUTOR_CONFLICT
ROOM_CONFLICT
TUTOR_DAILY_LIMIT


## Conflict Response

A conflict response should be explicit enough for the caller to understand what must change.

Example:

{
  "error": "ROOM_CONFLICT",
  "message": "Room is already booked during the requested time."
}

Do not expose raw SQL or database errors to the caller.


## Read API

A small read endpoint is useful for verification:

GET /appointments

Optional filters may include:

date
studentId
tutorId
roomId

Only add filters required to inspect or demonstrate the scheduler.

A full calendar/query API is not part of this feature.


## Booking Service

The central application operation can conceptually be:

createAppointment(input)

Responsibilities:

1. Validate request.
2. Load required entities.
3. Validate teaching day.
4. Validate duration.
5. Start transaction.
6. Protect tutor daily-capacity calculation.
7. Check tutor daily count.
8. Check scheduling conflicts.
9. Insert appointment.
10. Commit.
11. Return appointment.

Database conflict errors caused by concurrent writes should be translated back into the appropriate domain conflict response.


## Suggested Internal Structure

Keep the implementation small.

api/
  appointment_handler

application/
  create_appointment

domain/
  appointment
  scheduling_rules

repository/
  appointment_repository
  student_repository
  tutor_repository
  room_repository

The names are illustrative.

Do not create layers or interfaces that provide no concrete value.


## Indexing

Queries will frequently search appointments by:

student + time
tutor + time
room + time
tutor + day

Indexes should support these access patterns.

Do not add speculative indexes unrelated to the booking flow.


## Seed Data

The supplied week's spreadsheet should be loaded as historical seed data.

Important distinction:

Seed data represents what happened.
Business rules represent what should be allowed going forward.

Therefore historical rows should not be silently removed because they violate current rules.

If conflicting historical rows make database invariants impossible to apply directly, the import strategy must make that discrepancy explicit rather than rewriting history.

Possible approaches include:

- Importing historical records separately from enforceable active scheduling records.
- Marking historical invalid records appropriately if the source provides enough information.
- Documenting records that cannot satisfy the new scheduling invariants.

The exact choice should be based on the actual supplied data rather than invented ahead of inspection.


## Tests

The feature should be considered complete only when the important business cases are tested.


### Valid booking

Given:

- available student
- available tutor
- available room
- valid teaching day
- tutor below daily limit
- 60 or 90 minute duration

When:
a booking is created

Then:
the appointment is persisted as BOOKED.


### Room conflict

Existing:

Room 1
10:00 - 11:00

Requested:

Room 1
10:30 - 11:30

Then:
ROOM_CONFLICT.


### Tutor conflict

Existing:

Tutor A
10:00 - 11:00
Room 1

Requested:

Tutor A
10:30 - 11:30
Room 2

Then:
TUTOR_CONFLICT.


### Student conflict

Existing:

Student A
10:00 - 11:00

Requested:

Student A
10:30 - 11:30

Then:
STUDENT_CONFLICT.


### Adjacent appointments

Existing:

10:00 - 11:00

Requested:

11:00 - 12:00

Then:
booking is allowed.


### Tutor capacity

Given:
Tutor A has 5 active bookings that day

When:
booking number 6 is created

Then:
booking succeeds.

Given:
Tutor A has 6 active bookings that day

When:
another booking is requested

Then:
TUTOR_DAILY_LIMIT.


### Invalid teaching day

Requested appointment:
Monday

Then:
booking is rejected.


### Invalid duration

Requested duration:
45 minutes

Then:
booking is rejected.

Requested duration:
60 minutes

Then:
booking may proceed.

Requested duration:
90 minutes

Then:
booking may proceed.


### Concurrent room booking

Two requests attempt to book:

same room
overlapping time

concurrently.

Expected:

one succeeds
one fails

Never:
both succeed.


### Concurrent tutor booking

Two requests attempt overlapping bookings for the same tutor concurrently.

Expected:

one succeeds
one fails.


### Concurrent student booking

Two requests attempt overlapping bookings for the same student concurrently.

Expected:

one succeeds
one fails.


### Concurrent tutor daily limit

Tutor currently has 5 bookings.

Two requests attempt to create booking number 6 concurrently.

Expected:

one succeeds
one fails with TUTOR_DAILY_LIMIT.

Tutor must never end with 7 bookings.


## Explicitly Not Implemented

This phase does not implement:

- Rescheduling.
- Cancellation workflow.
- Late-cancellation charging.
- No-show workflow.
- Tutor notification delivery.
- Post-16:00 change history.
- Tutor acknowledgement.
- Calendar UI.
- Booking-limit overrides.
- Shared two-student lesson support.

These are intentionally left outside the selected feature rather than partially implemented.


## Design Principle

The main invariant for this implementation is:

If POST /appointments returns success, the resulting booking is compatible with every scheduling rule implemented by this feature.

The API should never return success first and depend on someone manually fixing the schedule later.

That guarantee is more important than implementing a larger number of scheduling workflows.