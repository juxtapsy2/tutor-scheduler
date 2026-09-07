# context.md

# Bright Path Scheduling System

## Overview

Bright Path Learning Centre is a tutoring centre in Da Nang.

The centre has:

- 12 tutors
- A little under 200 students/families
- 6 teaching rooms
- One-to-one lessons
- Lessons lasting either 60 or 90 minutes
- Teaching days from Tuesday through Sunday
- Monday is closed for room cleaning

Scheduling is currently managed through a shared spreadsheet and WhatsApp.

The system we are building replaces the spreadsheet-based scheduling workflow with a reliable scheduling system.

The primary concern is keeping the schedule consistent and preventing conflicting bookings.

---

## Scheduling Rules

### Opening days

The centre teaches Tuesday through Sunday.

Monday is closed.

Lessons generally run from mid-morning to mid-evening.

Exact opening and closing times are not specified and should not be invented unless required.

---

### Tutor daily capacity

A tutor may have at most 6 bookings in a single day.

The current receptionist sometimes exceeds this limit when necessary, but the intended system should enforce the limit.

Do not silently allow additional bookings beyond this limit unless an explicit override mechanism is later defined.

---

### Room availability

There are 6 rooms.

A room can contain only one lesson at a time.

Two lessons using the same room cannot overlap.

---

### Tutor availability

A tutor can only be in one room at a time.

Two overlapping lessons assigned to the same tutor are not allowed, even when different rooms are used.

---

### Student availability

Lessons are one-to-one.

A student cannot be booked into two places at the same time.

A student may have multiple lessons on the same day as long as their lesson times do not overlap.

---

### Cancellation

A family may cancel free of charge up to 4 hours before the lesson starts.

If a cancellation happens within 4 hours of the lesson:

- The lesson is charged in full.
- The tutor is still paid.

Cancellation frees the room and the scheduling slot.

A cancelled appointment should remain represented in the system rather than being deleted.

---

### No-show

A student who does not arrive is recorded as a no-show.

A no-show is different from a cancellation and should not be treated as a cancellation.

The source states that a no-show "weighs neither".

Do not infer additional financial behaviour from this phrase unless clarified.

---

### Schedule cut-off

Tomorrow's schedule is final at 16:00 today.

Changes after 16:00 are still allowed.

However, a post-cut-off change must remain visible as a change rather than silently overwriting what the tutor was previously told.

The system therefore needs to preserve enough information to understand that an appointment was changed after the previously communicated schedule.

The exact notification mechanism is not specified.

Do not invent WhatsApp, email, SMS, voice-call, or other notification infrastructure unless explicitly required later.

---

## Appointment

An appointment represents a scheduled lesson involving:

- One student
- One tutor
- One room
- One start time
- One end time
- One status

The system should clearly represent the appointment lifecycle.

A minimal status model may include:

- BOOKED
- CANCELLED
- COMPLETED
- NO_SHOW

Only introduce additional states when they solve a real requirement.

---

## Conflict Detection

Before creating or moving an appointment, the system must verify:

1. The selected room is available.
2. The selected tutor is available.
3. The selected student is available.
4. The tutor would not exceed 6 bookings on that day.

Conflict detection must compare actual time ranges rather than only matching start times.

Example:

Appointment A: 10:00 - 11:00
Appointment B: 10:30 - 11:30

Result: conflict.

Example:

Appointment A: 10:00 - 11:00
Appointment B: 11:00 - 12:00

Result: no overlap.

Use one consistent interval convention throughout the system.

---

## Appointment Changes

Moving or modifying an appointment is treated as another scheduling operation.

When an appointment changes:

- Validate the new time.
- Validate the new room.
- Validate the tutor.
- Validate the student.
- Recalculate the tutor's daily booking count.
- Ensure the resulting schedule remains valid.

A modification must not bypass normal scheduling rules.

If the change occurs after the 16:00 cut-off, the previous communicated state must not silently disappear.

The model should preserve enough information to understand what changed.

---

## Concurrency

Scheduling correctness must hold even when multiple users make requests simultaneously.

The following flow is not sufficient by itself:

check availability
-> resource appears available
-> insert booking

Two concurrent requests could both observe the same resource as available and both create bookings.

The persistence strategy must therefore prevent conflicting concurrent bookings from both succeeding.

An appropriate solution may use:

- Database constraints
- Transactions
- Locking
- Serializable operations
- Another equivalent concurrency-control mechanism

The exact strategy depends on the chosen database and implementation.

Do not claim concurrency safety when the implementation only performs a separate read followed by an unprotected insert.

---

## Historical Data

The supplied spreadsheet represents what actually happened during one particular week.

The data is anonymised and may contain situations that conflict with the intended scheduling rules.

Historical behaviour must not automatically be interpreted as a valid business rule.

In particular, the existing receptionist sometimes places two students with one tutor in the same room and same time slot as a deliberate half-price arrangement.

This historical behaviour conflicts with the stated one-to-one scheduling model.

Do not automatically support it as a normal booking rule.

Do not silently modify historical data merely because it violates intended rules.

When historical data and intended rules disagree:

1. Preserve the historical data.
2. Follow the explicit business rule for new scheduling behaviour.
3. Document the discrepancy if it affects implementation.

---

## Known Ambiguities

The following behaviour is not fully specified:

- Whether cancelled appointments count toward the tutor's daily booking limit.
- The exact financial meaning of the no-show rule.
- Whether the tutor daily limit can ever be formally overridden.
- The exact notification mechanism for booking changes.
- Exact daily opening and closing times.
- The exact audit/history model for post-cut-off changes.

Do not silently convert these ambiguities into requirements.

When implementation requires a decision:

1. Choose the smallest reasonable interpretation.
2. Keep the implementation easy to change.
3. Document the assumption.

---

## Product Scope

The primary product concern is reliable scheduling and conflict detection.

Focus on:

- Creating appointments
- Detecting conflicts
- Changing appointments
- Cancelling appointments
- Preserving meaningful appointment state
- Maintaining scheduling correctness under concurrent requests

Avoid expanding the system into unrelated tutoring-centre functionality unless required.

Out of scope unless explicitly introduced later:

- Authentication platforms
- Student portals
- Tutor portals
- Payments
- Billing systems
- Payroll
- Course management
- Analytics dashboards
- WhatsApp integration
- SMS integration
- Voice-call integration
- Email delivery infrastructure
- Calendar integrations
- Complex reminder systems
- Microservices
- Message brokers
- Kubernetes
- Event sourcing
- CQRS

---

## Product Goal

Build a small and reliable scheduling system that makes invalid schedules difficult or impossible to create.

The scheduling rules should be explicit, understandable, and consistently enforced.

Correctness is more important than feature breadth.

When requirements are unclear, do not invent behaviour.

Prefer simple decisions that can be explained and changed later.