# Prompt 19: Range/Classroom Booking System

**Use:** Design a Power App for scheduling TBS ranges, classrooms, and training areas
**Platform:** GenAI.mil
**Complexity:** Intermediate
**Source:** Adapted from EDD Room Booking System (Admin Tools #3)

---

## Prompt

```
Create a Power Apps canvas app for booking TBS training facilities. Requirements:

1. Data source: SharePoint list "FacilityBookings" with columns:
   - Facility (choice: [LIST_YOUR_FACILITIES — e.g.,
     Range 4, Range 8, Range 10, KD Range,
     Classroom A, Classroom B, Classroom C, Auditorium,
     Sand Table Room, MOUT Town, Confidence Course,
     Obstacle Course, LZ Vulture, Training Area 1, Training Area 2,
     PT Field, Pool])
   - FacilityType (choice: Range, Classroom, Field Training Area, Physical Training, Other)
   - BookedBy (person)
   - Company (choice: Alpha through Golf, India, Staff, TECOM)
   - Date (date)
   - StartTime (text HH:MM)
   - EndTime (text HH:MM)
   - Purpose (text — e.g., "Squad Live Fire," "Land Nav Written Test," "TDG Session")
   - POIReference (text — lesson number if applicable)
   - Recurring (yes/no)
   - SafetyOfficerRequired (yes/no — auto-set to yes for all Range types)
   - ApprovalStatus (choice: Pending, Approved, Denied)

2. Features:
   - **Daily Calendar View:** Show all bookings by facility for selected date. Color by FacilityType.
   - **Conflict Detection:** Prevent double-booking. Alert if overlapping times on same facility.
   - **Booking Form:** Validation — end time after start time, date must be today or future, Range bookings auto-flag SafetyOfficerRequired.
   - **My Bookings:** Filtered gallery for current user to see and cancel their own reservations.
   - **Weekly Overview:** Matrix showing all facilities (rows) by days (columns) with booking density.
   - **Approval Workflow:** Range bookings require S-3 approval via Power Automate. Classrooms are auto-approved.

3. Color coding: Green = available, Red = booked, Yellow = pending approval.

Generate the app structure, then the conflict-detection formula first.
```
