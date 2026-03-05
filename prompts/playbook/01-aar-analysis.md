# Prompt 01: AAR Analysis & Trend Extraction

**Use:** Post-event AAR processing — extract key findings from a verbal or written AAR
**Platform:** GenAI.mil
**Complexity:** Beginner
**Source:** Adapted from EDD Problem Discovery (#1)

---

## Prompt

```
I am a Staff Platoon Commander (SPC) at The Basic School, Quantico. I just completed an After Action Review for a training event.

EVENT: [EVENT_NAME — e.g., FEX 1, Squad Live Fire, Land Navigation Practical]
PHASE: [PHASE — Phase I Individual Skills / Phase II Squad / Phase III Platoon / Phase IV MAGTF]
UNIT SIZE: [NUMBER_OF_STUDENTS_INVOLVED]

AAR NOTES:
[PASTE_YOUR_AAR_NOTES_OR_SUMMARIZE_KEY_DISCUSSION_POINTS]

Analyze this AAR and provide:
1. **Sustains** — What went well that should be reinforced? (List 3-5)
2. **Improves** — What needs work? (List 3-5, ranked by impact on training objectives)
3. **Root Causes** — For each "improve," what is the likely root cause? (Student knowledge gap, planning failure, communication breakdown, resource limitation, etc.)
4. **Action Items** — Specific, assignable actions to address each "improve" (who should do what by when)
5. **Trend Check** — Based on your training knowledge, are any of these findings common across TBS classes? Flag any that suggest a systemic issue vs. a one-time event.

Format the output as a structured AAR summary I can file in the company training folder.
```

---

## Notes
- Do NOT paste student names or PII into the prompt — use billet/position (e.g., "Squad Leader, 2nd Squad")
- If comparing across multiple AARs, use Prompt 15 (AAR Synthesis & Trends) instead
