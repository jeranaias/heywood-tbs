# Prompt 05: AI Suitability Check (Frontier Check)

**Use:** Evaluate whether a TBS task is appropriate for AI assistance
**Platform:** GenAI.mil
**Complexity:** Beginner
**Source:** Adapted from EDD Frontier Check (#5)

---

## Prompt

```
I need to evaluate whether a task at The Basic School is suitable for AI assistance.

## TASK DESCRIPTION
[Describe the specific TBS task — e.g., "Writing counseling narratives for student evaluations," "Generating tactical scenarios for FEX planning," "Analyzing AAR trends across a training cycle"]

## TBS CONTEXT
This task is performed by: [SPC / S-3 Staff / Student / Instructor]
Frequency: [Daily / Weekly / Per event / Per training cycle]
Current method: [How it's done now — e.g., "SPC writes from memory after each event," "S-3 manually compiles from company reports"]
Stakes: [What happens if the output is wrong — e.g., "Student gets inaccurate counseling," "Training schedule conflict," "Incorrect grade calculation"]

## EVALUATION CRITERIA
For this task, assess:

1. **Inside the frontier?** Can AI reliably produce a correct or useful output for this type of task? What evidence supports this?
2. **Outside the frontier?** What aspects require SPC judgment, institutional knowledge of TBS culture, or real-world observation that AI cannot provide?
3. **Verification difficulty:** Can the SPC quickly tell if the AI output is right? Or does it require deep expertise to spot errors?
4. **Risk of undetected errors:** If the AI gets something wrong, how likely is it that the error goes unnoticed? What's the consequence for the student?
5. **Data sensitivity:** Does this task involve PII, PHI, or CUI? What platform restrictions apply?
6. **Recommended approach:**
   - **Delegate to AI** — Inside frontier, easy to verify, low stakes
   - **Centaur mode** — Human plans and reviews, AI executes the draft
   - **Cyborg mode** — Continuous human-AI collaboration, line by line
   - **Keep fully human** — Outside frontier, high stakes, or unverifiable output

Provide a clear recommendation with reasoning. If centaur or cyborg mode, specify what the human does vs. what the AI does.
```

---

## Notes
- Use this BEFORE starting any new AI-assisted workflow at TBS
- The "stakes" question is the most important — student evaluations and MOS assignments are career-defining
- When in doubt, default to centaur mode (human plans, AI drafts, human reviews)
