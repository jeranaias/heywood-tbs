# Prompt 14: Instructor Qualification Tracker

**Use:** Design a Power BI report tracking SPC and instructor certifications
**Platform:** GenAI.mil
**Complexity:** Intermediate
**Source:** Adapted from EDD Training Compliance Tracker (Reporting #4)

---

## Prompt

```
Create a Power BI report that tracks instructor qualifications and certifications at The Basic School.

Requirements:

1. Data sources:
   - "Instructors" SharePoint list: Name, Rank, Company Assignment, Billet (SPC, Asst SPC, Guest), StartDate
   - "RequiredQualifications" list: QualName, Frequency (Annual, Semi-Annual, One-Time, Per-Assignment), AppliesTo (choice: All SPCs, Range Instructors, Guest Instructors)
   - "QualificationRecords" list: Instructor (lookup), Qualification (lookup), CompletionDate, ExpirationDate, CertifyingAuthority

   Example qualifications for TBS:
   - Range Safety Officer (RSO) certification
   - Combat Lifesaver (CLS) current
   - PME completion (appropriate for rank)
   - Annual training requirements (Cyber Awareness, SAPR, etc.)
   - Swim Qualification Instructor
   - MOUT Instructor Certification
   - Land Navigation Instructor

2. Report views:
   - **Compliance Matrix:** Rows = instructors, columns = required qualifications, cells = status icon
     - Green check: Current
     - Yellow warning: Expiring within 30 days
     - Red X: Expired
     - Gray dash: Not applicable to this instructor's billet
   - **Company Rollup:** Bar chart showing qualification compliance percentage by company
   - **Individual Detail:** Drill-through page showing one instructor's full qualification record with dates
   - **Expiration Forecast:** Line chart showing how many qualifications expire each month for the next 6 months — helps S-3 plan requalification events

3. Filters: Company, rank, specific qualification. KPI card showing overall TBS instructor compliance percentage.

4. Alert integration: Flag any instructor below 100% compliance to the S-3 weekly report.

Provide the data model relationships, DAX measures, and conditional formatting rules.
```
