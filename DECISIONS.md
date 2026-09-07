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