# Prompt 09: Performance Tracking Data Model

**Use:** Design multi-table data model for the full Heywood data platform
**Platform:** GenAI.mil
**Complexity:** Intermediate
**Source:** Adapted from EDD Single List Design (#2)

---

## Prompt

```
Data platform: [SHAREPOINT_LIST / DATAVERSE_TABLE]

I need to design a multi-table data model for tracking student performance across The Basic School's 26-week training cycle.

Core entities and their purposes:
- **Students** — One record per student. Identified by EDIPI. Assigned to company, platoon, squad, SPC.
- **Events** — Training events from the POI. Categorized by phase (I-IV), type (academic, field, physical, evaluation), grading weight.
- **Scores** — Individual student scores per event. Links to Student and Event. Includes raw score, weighted score, pass/fail.
- **Evaluations** — SPC leadership evaluations (Week 12, Week 22). BARS dimensions if implemented. Narrative field.
- **Peer Reviews** — Peer rankings per evaluation point. Aggregated score per student.
- **Instructors** — SPC assignments, qualifications, certification dates.
- **AAR Entries** — Structured AAR data per event. Sustains, improves, action items, categories.

For each table, provide:
1. Recommended table name
2. Column definitions:

| Column Name | Type | Required | Description | PII? |
|-------------|------|----------|-------------|------|

3. Primary key and foreign key relationships
4. Calculated columns or rollup fields needed
5. Recommended views for each role (Staff / SPC / Student)
6. Index recommendations for performance

Also provide a text-based entity-relationship diagram showing how the tables connect.

Keep it practical — I can extend later, but the core model needs to be right from the start.
```
