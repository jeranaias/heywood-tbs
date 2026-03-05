# Prompt 08: Student Data Schema Design

**Use:** Design SharePoint list or Dataverse table structure for TBS student data
**Platform:** GenAI.mil
**Complexity:** Intermediate
**Source:** Adapted from EDD Requirements to Data Structure (#1)

---

## Prompt

```
I'm building a student performance tracking system for The Basic School on MCEN.

The tool needs to track/store:
[LIST_EVERYTHING — e.g.,
- Student identification (by EDIPI, not name, for data minimization)
- Company and platoon assignment
- SPC assignment
- Academic scores (4 phased exams, quizzes)
- Military skills scores (PFT, CFT, rifle/pistol qual, land nav, obstacle course)
- Leadership evaluation scores (Week 12 and Week 22 SPC evaluations)
- Peer evaluation scores
- BARS evaluation data (if implemented)
- Phase completion status (Phase I through Phase IV)
- Overall class standing (rank order within company)
]

Users need to:
[LIST_WHAT_EACH_ROLE_DOES — e.g.,
- Staff: View all companies, run aggregate reports, export data
- SPCs: Enter scores for their platoon, view company-level comparisons, generate counseling prep
- Students: View their own scores and standing only
]

Data platform: [SHAREPOINT_LIST / DATAVERSE_TABLE — SharePoint for Phase 1-2, Dataverse for Phase 3+]

Design the data structure:
1. What lists/tables do I need? (Consider: Students, Events, Scores, Evaluations, Instructors)
2. For each list/table, what columns are required?
3. For each column: data type (text, number, choice, person, date, lookup, calculated)
4. What relationships exist between lists/tables?
5. What columns should be required vs optional?
6. What default values make sense?
7. How does Row-Level Security work? (Staff sees all, SPC sees company, Student sees self)

Show the structure in table format. Flag any PII columns.
```
