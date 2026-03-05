# Prompt 11: Training Schedule Board

**Use:** Visualize the TBS 26-week POI and weekly training calendar
**Platform:** GenAI.mil
**Complexity:** Intermediate
**Source:** Adapted from EDD Training Schedule Board (BI #4)

---

## Prompt

```
Build a Power BI report that displays the TBS weekly training schedule from a SharePoint list.

Requirements:

1. Data source: SharePoint list "TrainingSchedule" with columns:
   - EventName (text)
   - Date (date)
   - StartTime (text HH:MM)
   - EndTime (text HH:MM)
   - Location (choice: [CLASSROOM, RANGE, FIELD, PT_AREA, POOL, OTHER])
   - Instructor (person — SPC or guest instructor)
   - Company (choice: Alpha through Golf, India, Mike, All)
   - Phase (choice: Phase I, Phase II, Phase III, Phase IV)
   - EventType (choice: Classroom, Field Exercise, Physical Training, Evaluation, Admin, Liberty)
   - POIReference (text — lesson number from Program of Instruction)
   - GradedEvent (yes/no)

2. Report pages:
   - **Weekly Calendar View:** Matrix with days as columns, time blocks as rows, events as colored cells (color by EventType). Show event name and location in each cell.
   - **Phase Overview:** Gantt-style view showing the 26-week training cycle divided into 4 phases with major events (FEX 1, FEX 2, Eight Day War, AMFEX) highlighted.
   - **Instructor Workload:** Bar chart of events per instructor for selected time period. Flag any instructor with >15 events/week.
   - **Facility Usage:** Heat map showing which locations are most utilized by day/time.

3. Filters: Date range slicer (default to current week), company filter, phase filter, event type filter.

4. Navigation: Click an event in any view to see full details (time, location, instructor, POI reference, graded status).

Provide the DAX measures needed and visual configuration for each page.
```
