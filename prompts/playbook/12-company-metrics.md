# Prompt 12: Company Metrics Summary

**Use:** Generate a monthly command briefing slide for TBS leadership
**Platform:** GenAI.mil
**Complexity:** Intermediate
**Source:** Adapted from EDD Monthly Metrics Summary (Reporting #2)

---

## Prompt

```
Create a Power BI report page that summarizes monthly TBS company metrics for a command briefing.

Requirements:

1. Data sources: SharePoint lists for Students, Scores, Instructors, and Training Events.

2. Single-page layout designed for projection (16:9 aspect ratio):

   TOP BANNER: TBS logo placeholder | "Company [X] Monthly Metrics" | Reporting period dates

   LEFT COLUMN (3 gauges):
   - Average Leadership Grade (target: 80%+)
   - Average Military Skills Grade (target: 80%+)
   - Average Academics Grade (target: 80%+)

   CENTER: Trend chart showing all three pillar averages over the past 4 phases as lines. Include TBS-wide average as a dashed reference line.

   RIGHT COLUMN:
   - Students at risk (below 75% in any pillar) — count and list by billet position (no names)
   - Instructor qualification status (all current / X expiring)
   - Key events completed this period
   - Training days lost (weather, admin, other)

   BOTTOM ROW:
   - Top 3 sustains from AARs this period
   - Top 3 improves from AARs this period
   - "Prepared by [SPC NAME] | As of [DATE]"

3. Gauge colors:
   - Green: 80%+ (on target)
   - Yellow: 75-79% (caution)
   - Red: Below 75% (below standard)

4. Row-Level Security: Staff sees all companies. SPCs see only their own company's report.

Provide the DAX measures, gauge threshold parameters, and conditional formatting rules.
```
