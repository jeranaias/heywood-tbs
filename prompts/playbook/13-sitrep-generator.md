# Prompt 13: Weekly SITREP Generator

**Use:** Generate a formatted weekly status report for the TBS training chain
**Platform:** GenAI.mil
**Complexity:** Beginner
**Source:** Adapted from EDD Status Report Generator (Reporting #1)

---

## Prompt

```
Help me draft a weekly SITREP for my TBS company. I'll provide the raw data, you format it.

REPORTING PERIOD: [START_DATE] to [END_DATE]
COMPANY: [COMPANY_LETTER]
PHASE: [CURRENT_PHASE — Phase I / II / III / IV]
WEEK: [WEEK_NUMBER] of 26

PERSONNEL STATUS:
- Present for duty: [NUMBER]
- Leave/liberty: [NUMBER]
- Medical hold (Mike Company): [NUMBER]
- TAD: [NUMBER]
- Total assigned: [NUMBER]

TRAINING COMPLETED THIS WEEK:
[LIST_EVENTS_COMPLETED — e.g.,
- Land Navigation Written Test (Phase I, graded)
- Squad Attack Live Fire (Phase II, graded)
- PFT (graded)
- 3x PT sessions (non-graded)
]

TRAINING PLANNED NEXT WEEK:
[LIST_UPCOMING_EVENTS]

SIGNIFICANT OBSERVATIONS:
[ANY_NOTABLE_EVENTS — e.g., "3 students fell below 75% in academics after Exam 2," "Zero safety incidents during live fire," "Guest instructor from 2d MarDiv for urban ops block"]

ISSUES/CONCERNS:
[ANY_PROBLEMS — e.g., "Range 8 unavailable Tuesday due to maintenance," "2 SPCs TDY, coverage plan in place"]

Format this as a clean SITREP with:
1. DTG header
2. Unit identification
3. Personnel summary (one-line)
4. Training summary (completed and planned, in bullets)
5. Commander's assessment (2-3 sentences synthesizing the week — draft this based on my inputs, I'll edit)
6. Issues requiring higher attention (if any)

Keep it to one page. Use standard military formatting.
```
