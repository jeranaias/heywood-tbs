# SPC Quick-Start Guide: TBS Prompt Playbook on GenAI.mil

*Heywood Phase 1 — For Staff Platoon Commanders and Assistant SPCs*

---

## What This Is

20 ready-to-use AI prompts designed for your daily SPC workflow. Copy, paste, fill in the blanks, get useful output. Each prompt runs on GenAI.mil, which is already authorized on MCEN for CUI data.

---

## Getting Started (5 minutes)

### Step 1: Access GenAI.mil
1. Open your browser on any MCEN workstation
2. Navigate to **genai.mil**
3. Authenticate with your **CAC** (Common Access Card)
4. You'll see the GenAI chat interface

### Step 2: Start a New Conversation
- Click **New Chat** for each new task
- Do not continue old conversations for unrelated tasks — context from previous chats can confuse the output

### Step 3: Use a Prompt
1. Open the prompt file (from the playbook folder or the printed reference card)
2. **Copy** the entire prompt text
3. **Paste** it into GenAI.mil
4. **Replace** every `[BRACKETED_PLACEHOLDER]` with your specific information
5. Press Enter / Send
6. **Review the output critically** — you are the subject matter expert, the AI is a drafting tool

---

## Data Handling Rules

You must follow these rules every time:

| Data Type | Can I Put It in GenAI.mil? | Example |
|-----------|:-:|---------|
| CUI (controlled unclassified) | Yes | Training schedule dates, event names, POI references |
| PII (personally identifiable) | No — anonymize first | Use "Student A" or "the student" instead of names/EDIPIs |
| PHI (protected health info) | Never | No injury details, medical status, or BAS records |
| Classified | Never | Nothing from SIPR, nothing marked SECRET or above |
| Aggregate data | Yes | "3 of 48 students scored below 75% on Exam 1" |
| Doctrinal references | Yes | FM numbers, MCDP titles, TBS POI references |

**Rule of thumb:** Describe the *type* and *pattern* of data, not the actual records. Say "a student's four exam scores are 82, 71, 88, 79" not "2ndLt Smith's scores are..."

---

## Your Top 5 Prompts (Daily Use)

### 1. Counseling Prep (Prompt #2)
**When:** Before any counseling session with a student
**What it does:** Takes a student's scores and observations, generates a structured counseling outline with specific, actionable feedback
**Time saved:** ~30 min per counseling session

**How to use:**
1. Pull up the student's scores (from the Power BI dashboard or your records)
2. Open `02-counseling-prep.md`
3. Fill in: student rank (no name), scores across all three pillars, your observations, areas of concern
4. The AI generates: counseling structure, talking points, developmental goals, follow-up actions
5. **Review and edit** — add your own observations, remove anything that doesn't fit

### 2. AAR Analysis (Prompt #1)
**When:** After completing an AAR for any training event
**What it does:** Takes raw AAR notes and extracts themes, patterns, and actionable recommendations
**Time saved:** ~20 min per AAR

**How to use:**
1. Type up your AAR sustains/improves (or copy from notes)
2. Open `01-aar-analysis.md`
3. Paste the AAR content into the prompt
4. The AI identifies: recurring themes, root causes, priority recommendations, trends across events
5. Use the output to brief your company commander or update the event feedback tracker

### 3. Weekly SITREP (Prompt #13)
**When:** Weekly, before the command brief
**What it does:** Takes your raw data points and formats them into a standard SITREP format
**Time saved:** ~15 min per week

**How to use:**
1. Gather: training events completed, upcoming events, at-risk student count, personnel status, equipment issues
2. Open `13-sitrep-generator.md`
3. Fill in each section with your data points
4. The AI formats a clean SITREP ready for your commander
5. Review for accuracy, adjust tone, submit

### 4. Scenario Generator (Prompt #16)
**When:** Building FEX scenarios, developing training problems
**What it does:** Generates tactical scenarios calibrated to the current training phase using METT-TC format
**Time saved:** ~45 min per scenario

**How to use:**
1. Open `16-scenario-generator.md`
2. Specify: phase, skill level, terrain, learning objectives, threat type, OPFOR composition
3. The AI generates: situation, mission, terrain analysis, enemy forces, friendly forces, timeline
4. **Review thoroughly** — verify tactical realism against your experience
5. Adjust difficulty, add/remove complexity as needed

### 5. Peer Evaluation Structuring (Prompt #17)
**When:** Preparing peer evaluation forms or processing peer feedback
**What it does:** Creates structured peer evaluation criteria aligned to TBS competencies
**Time saved:** ~20 min per evaluation cycle

---

## Common Workflows

### Weekly SPC Workflow
| Day | Task | Prompt | Time |
|-----|------|--------|------|
| Monday | Review at-risk students, plan counseling | #2 Counseling Prep | 15 min/student |
| Tuesday-Thursday | Conduct training events | — | — |
| Friday AM | Process any AARs from the week | #1 AAR Analysis | 20 min/event |
| Friday PM | Generate weekly SITREP | #13 SITREP Generator | 15 min |

### Pre-FEX Planning
1. Use **#16 Scenario Generator** to draft scenarios (1-2 hours for a week's scenarios)
2. Use **#11 Training Schedule** to verify resource conflicts
3. Use **#18 Equipment Tracker** to check supply readiness

### End-of-Phase Student Review
1. Pull composite scores from Power BI dashboard
2. Run **#2 Counseling Prep** for each at-risk student
3. Use **#15 AAR Synthesis** to identify training event trends
4. Use **#12 Company Metrics** to prep the company commander's brief

---

## Tips for Better Output

1. **Be specific with placeholders.** "Phase I Land Nav Day Practical" is better than "a training event."

2. **Give context.** Tell the AI what phase, what skill level, what the student's trend looks like. More context = better output.

3. **Iterate.** If the first output isn't right, tell the AI what to change: "Make the counseling tone more developmental, less punitive" or "Add more detail on the defense portion."

4. **Don't trust blindly.** The AI doesn't know your students, your terrain, or your company's specific situation. Use it as a starting point, not a finished product.

5. **One task per conversation.** Start a new chat for each distinct task. Don't try to do counseling prep and scenario generation in the same conversation.

6. **Save good outputs.** When the AI produces something you'll reuse (a scenario format, a counseling template), save it locally for future reference.

---

## Troubleshooting

| Problem | Solution |
|---------|----------|
| GenAI.mil is slow or unresponsive | Try during off-peak hours (before 0800 or after 1700 EST). The platform has 1.1M+ users. |
| Output seems generic or unhelpful | Add more TBS-specific context. Mention the specific POI event, phase, and what you're actually trying to accomplish. |
| AI refuses to generate content | You may have included something it flags as sensitive. Remove any names, EDIPIs, or specific medical information and try again. |
| Output is too long | Add "Keep your response under 500 words" or "Provide a bullet-point summary only" to the end of the prompt. |
| AI "hallucinates" a doctrinal reference | Always verify any FM, MCO, or MCDP citation the AI produces. If you don't recognize it, look it up before using it. |

---

## Quick Reference Card

Print this section and keep it at your desk:

```
DAILY PROMPTS:
  #1  AAR Analysis .............. After any training event AAR
  #2  Counseling Prep ........... Before any student counseling
  #13 Weekly SITREP ............. Friday afternoon, before command brief
  #15 AAR Synthesis ............. End of week, identify trends

PLANNING PROMPTS:
  #16 Scenario Generator ........ FEX and field exercise planning
  #11 Training Schedule ......... Schedule visualization
  #18 Equipment Tracker ......... Pre-FEX supply readiness
  #19 Range/Classroom Booking ... Facility scheduling conflicts

EVALUATION PROMPTS:
  #17 Peer Eval Structuring ..... Before peer evaluation cycle
  #12 Company Metrics ........... Monthly command briefing prep
  #14 Instructor Quals .......... Quarterly qual review

DATA RULES:
  CUI → OK    |  PII → Anonymize    |  PHI → Never    |  Classified → Never
```

---

*Questions? Contact SSgt Morgan or refer to the full prompt playbook in the Heywood project repository.*
