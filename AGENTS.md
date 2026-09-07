# Engineering Guidelines

## General Principles

Build the simplest solution that correctly satisfies context.md.

Prefer:

- Clear code
- Small modules
- Explicit behaviour
- Simple data models
- Predictable error handling
- Testable business logic
- Minimal dependencies

Avoid premature abstraction.

Do not introduce architecture, infrastructure, or patterns without a concrete need.

Do not overengineer for hypothetical future requirements.

---

## Source of Truth

context.md defines product behaviour and business rules.

AGENTS.md defines engineering conventions.

Do not duplicate business rules unnecessarily across the codebase.

When code behaviour conflicts with context.md, context.md takes precedence unless the requirement has been explicitly changed.

Do not invent missing requirements.

---

## Business Logic

Business rules must have a clear owner.

Avoid duplicating the same scheduling rule across handlers, repositories, frontend components, and helper functions.

The backend/application layer is authoritative for scheduling behaviour.

Frontend validation may improve user experience but must not be relied on for business correctness.

Scheduling logic should be independently testable where practical.

---

## Architecture

Prefer a simple layered flow:

HTTP/API
-> application/use-case logic
-> scheduling/domain rules
-> persistence/database

Keep boundaries useful rather than ceremonial.

Do not create interfaces, factories, services, repositories, or abstractions merely because a design pattern suggests them.

Introduce an abstraction only when it:

- Separates meaningful responsibilities
- Makes testing materially easier
- Protects an important boundary
- Removes real duplication

---

## HTTP Handlers

Keep HTTP handlers/controllers thin.

They should primarily:

1. Parse the request.
2. Validate basic input.
3. Call the relevant application operation.
4. Convert the result into an HTTP response.

Do not place substantial scheduling logic directly inside HTTP handlers.

---

## Business Errors

Expected scheduling failures should be explicit.

Examples may include:

- ROOM_CONFLICT
- TUTOR_CONFLICT
- STUDENT_CONFLICT
- TUTOR_DAILY_LIMIT
- INVALID_APPOINTMENT_TIME
- APPOINTMENT_NOT_FOUND
- INVALID_APPOINTMENT_STATE

Use names appropriate to the implementation.

Do not return a generic internal error when the system knows that a business rule rejected the operation.

Separate expected business failures from unexpected system failures.

---

## Database

Use the database as the final authority for data integrity where practical.

Prefer database constraints for structural invariants such as:

- Primary keys
- Foreign keys
- Required values
- Unique constraints
- Valid relationships

Use application logic for domain rules that require contextual calculations.

Be careful with rules involving multiple rows or time ranges.

---

## Concurrency

Scheduling operations must be safe under concurrent requests.

Never assume this sequence is safe:

1. Query for conflicts.
2. Find none.
3. Insert appointment.

Another request may insert a conflicting appointment between steps 2 and 3.

Use an appropriate strategy such as:

- Transaction isolation
- Row locking
- Advisory locking
- Database exclusion constraints
- Serializable transactions
- Another database-supported concurrency mechanism

Choose the simplest mechanism appropriate for the selected database.

Concurrency protection should live near the persistence boundary rather than depending on frontend behaviour.

---

## Transactions

Operations that must succeed or fail together should run atomically.

Examples include:

- Validating and creating a booking
- Moving an appointment
- Updating state while preserving associated history

Do not leave appointments in partially updated states when an operation fails.

Keep transactions as small as practical.

---

## Data Model

Model actual business concepts directly.

Prefer explicit relationships over encoded or implicit behaviour.

Avoid introducing fields without a current purpose.

Avoid generic abstractions such as:

- entity_type
- entity_data
- metadata blobs
- generic resource tables

when explicit columns or relationships better represent the domain.

Use JSON fields only when the structure is genuinely flexible.

---

## Appointment History

Avoid destructive updates when they would erase meaningful scheduling history.

When appointment changes need to remain visible, preserve enough information to understand what changed.

Do not introduce a large event-sourcing system just to maintain simple booking history.

Prefer the smallest model that satisfies the requirement.

---

## Time Handling

Use one consistent time representation across the application.

Prefer storing complete timestamps when practical.

Clearly define the overlap rule.

For half-open intervals:

[start, end)

the following appointments do not overlap:

10:00 - 11:00
11:00 - 12:00

Centralise overlap logic rather than recreating it in multiple places.

Avoid string-based time comparisons.

---

## Clock Dependency

Business logic that depends on "now" should not directly depend on the operating system clock throughout the codebase.

Prefer an injectable or centralised clock abstraction when needed.

Keep it simple.

Do not build an elaborate time framework.

Tests involving cut-off behaviour should be deterministic.

---

## Configuration

Business constants must have a clear source.

Examples include:

- Maximum tutor bookings per day
- Cancellation window
- Schedule cut-off time

Do not duplicate magic values across files.

Do not introduce a configuration framework solely for a few stable domain constants.

Use the simplest maintainable representation appropriate to the project.

---

## Validation

Separate basic input validation from business validation.

Basic validation includes:

- Missing required fields
- Invalid identifiers
- Malformed dates
- End time before start time

Business validation includes:

- Room availability
- Tutor availability
- Student availability
- Tutor daily capacity
- Appointment state transitions

Do not rely solely on frontend validation.

---

## Repository / Persistence Layer

Keep persistence logic focused on storing and retrieving data.

Avoid embedding HTTP concerns in repository code.

Avoid moving domain decisions into SQL unless the database itself is intentionally enforcing the invariant.

Complex queries are acceptable when they make scheduling correctness stronger or more atomic.

Use parameterised queries.

Never build SQL by concatenating untrusted input.

---

## Dependencies

Prefer existing dependencies and the standard library.

Before adding a dependency, consider whether the problem can be solved clearly without it.

Add dependencies only when they provide meaningful value.

Avoid bringing in large frameworks for small utilities.

Do not add infrastructure dependencies for hypothetical future scale.

---

## Naming

Prefer names that describe business meaning.

Good examples:

- createAppointment
- findConflictingAppointments
- cancelAppointment
- tutorDailyBookingCount

Avoid vague names such as:

- processData
- handleThing
- executeLogic
- manager
- helper
- utils2

Use consistent terminology from context.md.

Prefer "appointment" or "lesson" consistently rather than alternating between many synonyms.

---

## Functions

Keep functions focused.

A function should usually perform one understandable operation.

Avoid:

- Extremely large functions
- Excessive nesting
- Long parameter lists
- Boolean parameters whose meaning is unclear
- Hidden side effects

Prefer early returns when they improve readability.

Do not split trivial logic into dozens of tiny functions solely to reduce line count.

---

## Comments

Write comments to explain why, not what.

Avoid comments that merely repeat the code.

Good comments explain:

- Non-obvious business decisions
- Concurrency behaviour
- Database constraints
- Intentional trade-offs
- Workarounds

Delete outdated comments when behaviour changes.

---

## Error Handling

Handle errors explicitly.

Do not ignore errors.

Do not silently recover from failed scheduling operations.

Keep internal errors separate from user-facing messages.

Avoid exposing raw database errors through the API.

Wrap errors only when doing so adds useful context.

---

## Logging

Log useful operational information.

Avoid excessive logging.

Never log sensitive personal information unnecessarily.

Do not log full request payloads by default if they contain student or tutor information.

Logs should help answer:

- What operation failed?
- Why did it fail?
- Which system component failed?

---

## Testing

Prioritise business behaviour.

Important areas include:

- Room conflicts
- Tutor conflicts
- Student conflicts
- Tutor daily capacity
- Appointment creation
- Appointment movement
- Cancellation
- No-show state
- Post-cut-off changes
- Relevant state transitions
- Concurrency behaviour

Prefer testing observable behaviour over implementation details.

Tests should remain readable.

A test should make clear:

- Given
- When
- Then

Do not heavily mock domain logic.

Use real database integration tests when database behaviour is essential to correctness.

---

## Frontend

If a frontend is present:

- Use functional components.
- Prefer TypeScript.
- Keep components focused.
- Keep API calls separate from presentation where practical.
- Keep authoritative business rules on the backend.
- Display scheduling conflicts clearly.
- Avoid duplicating complex validation logic from the backend.
- Avoid excessive global state.

Do not spend disproportionate effort on visual polish before core scheduling behaviour is reliable.

---

## API Design

Keep the API small and predictable.

Prefer resource-oriented routes.

Possible operations may resemble:

POST /appointments
GET /appointments
GET /appointments/:id
PATCH /appointments/:id
POST /appointments/:id/cancel

Exact routes depend on the implementation.

Use appropriate HTTP status codes.

Keep request and response structures explicit.

Avoid generic endpoints such as:

POST /execute
POST /action
POST /process

when the domain operation can be named clearly.

---

## Security

Treat all external input as untrusted.

Use parameterised database queries.

Validate identifiers and input formats.

Do not commit secrets.

Use environment variables for credentials and environment-specific settings.

Avoid unnecessarily collecting or exposing personal data.

---

## File Structure

Keep the repository structure simple.

A typical backend may look like:

src/
  api/
  application/
  domain/
  repository/

tests/

Exact folder names are not mandatory.

Do not create deep directory trees unless the project genuinely requires them.

Prefer discoverability over architectural ceremony.

---

## Refactoring

Do not rewrite working code merely for stylistic preference.

Before refactoring:

1. Understand existing behaviour.
2. Identify the concrete problem.
3. Keep the change focused.
4. Preserve behaviour unless behaviour is intentionally changing.
5. Run relevant tests.

Avoid mixing major refactors with feature changes.

---

## Code Changes

Before modifying the codebase:

1. Read the relevant existing code.
2. Read context.md when business behaviour is involved.
3. Identify the smallest useful change.
4. Implement it.
5. Run relevant tests.
6. Fix failures caused by the change.

Do not perform unrelated cleanup during a focused task unless necessary.

---

## Documentation

Document non-obvious decisions.

Keep documentation concise and useful.

Do not document assumptions as confirmed requirements.

When an ambiguity requires an implementation decision, state that it is an assumption.

Avoid maintaining duplicate documentation describing the same rule in multiple places.

---

## Performance

Prioritise correctness first.

Do not optimise without evidence of a performance problem.

Avoid obvious inefficiencies such as unnecessary repeated database queries inside loops.

Use indexes when query patterns clearly require them.

Do not introduce caches until there is a demonstrated need.

---

## Scope Control

Do not expand the system beyond current requirements.

Avoid building:

- Microservices
- Event buses
- Message brokers
- Distributed caches
- Kubernetes configuration
- CQRS
- Event sourcing
- Generic workflow engines
- Plugin systems
- Complex role systems

unless a concrete requirement makes them necessary.

A smaller correct implementation is preferable to a larger speculative architecture.

---

## Definition of Good Code

Good code in this project is:

- Correct
- Simple
- Explicit
- Readable
- Testable
- Maintainable
- Safe under realistic concurrent use

Optimise for those properties before abstraction, extensibility, cleverness, or visual polish.