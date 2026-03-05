# Prompt 04: Training Aid MVP Definition

**Use:** Scope a new TBS tool to its minimum viable version
**Platform:** GenAI.mil
**Complexity:** Beginner
**Source:** Adapted from EDD MVP Definition (#4)

---

## Prompt

```
I'm planning to build a training support tool for The Basic School.

Tool: [TOOL_DESCRIPTION_2_TO_3_SENTENCES]

My full vision includes these features:
[LIST_ALL_FEATURES_YOU_HAVE_THOUGHT_OF]

My constraints are:
- Development platform: Power Platform on MCEN (Power Apps, Power Automate, Power BI, SharePoint)
- Development time: [HOURS_AVAILABLE]
- Must work for: [NUMBER_OF_USERS — e.g., "one company of 200 students and 6 SPCs"]
- Must stay within MCEN authorization boundary (no external APIs in v1)
- Data classification: [UNCLASSIFIED / CUI — if CUI, must use IL5 platform]

For each feature in my list, analyze:
- Is this essential for a minimum viable product, or a "nice to have"?
- What's the simplest way to implement this on Power Platform?
- What could I cut and still have something useful?

Then define the Minimum Viable Product:
1. What is the ONE core problem this solves for TBS?
2. What are the 3-5 features that solve that core problem?
3. What can be cut entirely from v1?
4. What can be simplified (e.g., SharePoint list instead of Dataverse, manual entry instead of automated)?
5. What's the "walking skeleton" I could build in 4 hours to prove the concept works?
6. What's a recommended "Version 1" scope I could build in under 20 hours?

Be ruthless. A working simple tool beats a broken complex one. TBS runs on tight timelines.
```
