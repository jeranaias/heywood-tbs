# Prompt 10: Student Performance Dashboard

**Use:** Define Power BI dashboard requirements for TBS student performance
**Platform:** GenAI.mil
**Complexity:** Intermediate
**Source:** Adapted from EDD Dashboard Requirements (BI #1)

---

## Prompt

```
I'm building a Power BI dashboard for The Basic School student performance tracking.

Audience:
- TBS Staff (CO, XO, S-3): Need aggregate view across all 8 companies
- Staff Platoon Commanders (SPCs): Need company-level view with individual student drill-down
- Students: Need individual view of their own performance only

Data source: SharePoint list "StudentScores" (Phase 1) or Dataverse "Scores" table (Phase 3+)

TBS grading structure:
- Leadership: 36% of total grade (SPC evaluation 90% of leadership grade, Peer evaluation 10%)
- Military Skills: 32% of total grade (PFT, CFT, rifle/pistol, land nav, obstacle course, etc.)
- Academics: 32% of total grade (4 phased exams + quizzes)
- Minimum passing: 75% overall with 75% in each category
- Students ranked in thirds for MOS assignment

Help me define dashboard requirements:
1. What questions should each audience be able to answer?
2. What metrics matter most to each audience?
3. What time view makes sense? (Weekly during the 26-week cycle, by phase, cumulative)
4. What comparisons matter? (Student vs company average, company vs TBS average, current class vs historical)
5. What drill-down capability is needed?
6. What visual types best show three-pillar performance data?
7. What alert thresholds matter? (Student below 75%, student dropping more than 5% between evaluations)

Define the 3-5 MUST-HAVE report pages vs nice-to-haves. For each page, specify the visuals, filters, and DAX measures needed.

CRITICAL: Row-Level Security must enforce:
- Staff role: sees all companies
- SPC role: sees only their assigned company
- Student role: sees only their own record
```
