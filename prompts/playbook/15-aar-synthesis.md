# Prompt 15: AAR Synthesis & Cross-Event Trends

**Use:** Analyze patterns across multiple AARs to identify systemic training issues
**Platform:** GenAI.mil
**Complexity:** Intermediate
**Source:** Custom for TBS (no direct EDD analog)

---

## Prompt

```
I am TBS Staff (S-3) or a Staff Platoon Commander analyzing After Action Reviews across multiple training events.

ANALYSIS SCOPE:
- Company: [COMPANY_LETTER or "All Companies"]
- Phase: [PHASE or "All Phases"]
- Time period: [DATE_RANGE]
- Number of AARs being analyzed: [NUMBER]

AAR SUMMARIES:
[FOR EACH AAR, PROVIDE:
Event name | Date | Phase | Key sustains (2-3) | Key improves (2-3) | Action items taken

Example:
- Squad Attack FEX | 15 Jan | Phase II | Good fire discipline, strong patrol orders | Slow casualty evacuation, weak radio procedures | Added radio drill to pre-FEX prep
- Platoon Defense | 22 Jan | Phase III | Effective use of terrain, strong security plan | Poor integration of supporting arms, late logistics coordination | Added fire support planning block
]

Analyze these AARs and provide:

1. **Recurring Sustains** — What consistently goes well? (Ranked by frequency)
2. **Recurring Improves** — What keeps showing up as a problem? (Ranked by frequency and severity)
3. **Trend Analysis** — Are improves getting better or worse over time? Are Phase I issues resolving by Phase III?
4. **Root Cause Patterns** — Group the recurring improves by root cause:
   - Knowledge gap (students don't know the doctrine)
   - Practice gap (students know it but can't execute under pressure)
   - Planning gap (orders/preparation insufficient)
   - Communication gap (information not flowing between elements)
   - Resource/equipment issue (training limitation, not student failure)
5. **Curriculum Recommendations** — Based on the trends, what training adjustments would address the root causes? Be specific (add a class, add a rehearsal, adjust the POI sequence, etc.)
6. **Comparison to Known TBS Patterns** — If you have knowledge of common TBS training challenges, note whether these trends are typical or unusual.

Format as a one-page executive summary suitable for a TBS staff meeting, followed by a detailed appendix.
```

---

## Notes
- This prompt works best with 5+ AARs — fewer than that and you're finding anecdotes, not trends
- Do NOT include student names — use event names, dates, and unit designations only
- The curriculum recommendations are AI suggestions that require SME validation before implementation
