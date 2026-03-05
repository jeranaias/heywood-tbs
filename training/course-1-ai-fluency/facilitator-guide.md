# Course 1: AI Fluency for TBS Instructors

*Heywood Initiative — Adapted from EDD AI Fluency Fundamentals*
*Facilitator Guide*

---

## Course Information

| Field | Detail |
|-------|--------|
| **Duration** | 2 hours |
| **Audience** | TBS Staff Platoon Commanders, Assistant SPCs, Tactics/Weapons/PT Instructors |
| **Max Class Size** | 20 (one company's instructor cadre + staff) |
| **Prerequisites** | None — come as you are |
| **Platform** | GenAI.mil (primary), CamoGPT (backup) |
| **Classification** | UNCLASSIFIED — no live student data in exercises |

---

## Learning Objectives

By the end of this course, participants will:

1. Understand the six skills that separate sustained AI adopters from the 80% who quit after three weeks
2. Recognize AI's uneven capabilities — the "jagged frontier" — across TBS-specific tasks
3. Know when to trust AI output and when to verify (quality judgment)
4. Map AI into their actual SPC/instructor workflow (Centaur vs. Cyborg patterns)
5. Delegate work to AI using the three-variable equation
6. Follow TBS data handling rules for GenAI.mil (CUI authorized, PII anonymized, PHI never)

---

## Pre-Course Setup

### Facilitator Preparation
- [ ] Test GenAI.mil access on 2-3 classroom machines (CAC reader required)
- [ ] Print Red Pen Review exercise documents (3 per participant — see Appendix A)
- [ ] Print Quick Reference Card from SPC Quick-Start Guide (1 per participant)
- [ ] Print Workflow Mapping worksheet (1 per participant — see Appendix B)
- [ ] Have the TBS Prompt Playbook folder accessible (digital or printed index)
- [ ] Prepare backup: if GenAI.mil is down, all exercises work with the provided inline documents

### Room Setup
- Projector with MCEN workstation
- Participants at tables with MCEN laptops (ideal) or printed materials (fallback)
- Whiteboard or butcher paper for frontier mapping exercise

### Contingency
If GenAI.mil is unavailable during the session, all hands-on exercises include pre-generated AI outputs that participants can still evaluate. The learning objectives are about judgment and workflow, not tool operation.

---

## Module 1: Why 80% Quit (and How You Won't)
*0:00 – 0:15 (15 minutes)*

### Facilitator Notes

Open with this question to the room: *"How many of you have tried ChatGPT, GenAI.mil, or any AI tool at least once?"* (Most hands go up.) *"How many of you use it regularly — at least once a week for actual work?"* (Far fewer hands.)

**Key point:** Microsoft's internal research found that excitement about AI tools peaks after about three weeks — then 80% of users quit. Not because the tools are bad, but because nobody taught them how to use them well.

### Content

**Why people quit:**
- First attempts produce generic, unhelpful output → "this is useless"
- AI gives confident wrong answers → "I can't trust this"
- Doing it yourself seems faster → "it takes longer to check the AI than to just write it"

**Why that's fixable:**
- UK Government study: 20,000 civil servants across 12 departments
- With proper training: 80% kept using AI tools after 6 months
- Average time saved: 25 minutes per day
- Not from one big win — from dozens of small integrations across their workflow

**The management framing (this is the key insight):**
- AI is not Google. It's not a search engine.
- Think of it as a capable but inexperienced Lance Corporal who just checked into your section
- Smart, fast, eager to help — but doesn't know your unit, your standards, or your context
- Your job: break work into pieces, explain what good looks like, review the output, give feedback
- The six skills you'll learn today are management skills applied to AI

**TBS connection:**
- SPCs manage 40-50 students. You already know how to delegate, review, and develop.
- These same skills apply to AI. The Marines who treat AI like a team member (not a magic box) are the ones who keep using it.

### Transition
*"So what are these six skills? Let's walk through them with TBS examples."*

---

## Module 2: The Six Skills That Actually Matter
*0:15 – 0:40 (25 minutes)*

### Content

Walk through each skill with a TBS-specific example. For each skill, show the 101 behavior (what most people do) and the 201 behavior (what effective users do).

---

**Skill 1: Context Assembly**
*Knowing what information to provide, from which sources, and why*

| 101 Behavior | 201 Behavior |
|---|---|
| "Write me a counseling statement" | "Write a developmental counseling outline for a Phase I student. The student scored 71% on Academic Exam 1, is trending down from an 82% quiz average, and shows strong physical fitness (PFT 285). The counseling should address academic performance while acknowledging military skills strengths. Use a direct but developmental tone." |

*Facilitator: Show both prompts live on GenAI.mil. The difference in output quality is dramatic.*

**TBS examples of good context:**
- Phase of training, specific event name, learning objectives
- Student's trend (improving, declining, stable) — not their name
- What you want the output to look like (counseling format, SITREP format, 5-paragraph order)

---

**Skill 2: Quality Judgment**
*Knowing when to trust AI output and when to verify*

| 101 Behavior | 201 Behavior |
|---|---|
| Accepts AI counseling draft word-for-word | Checks: Are the developmental goals specific and measurable? Does the tone match what I'd actually say? Did it fabricate any references? |

**The verification hierarchy for TBS:**
- **High-stakes (always verify line by line):** Anything going in a student's record, fitness report language, admin separation recommendations, anything with a signature block
- **Medium-stakes (spot-check key claims):** SITREP content, training schedule coordination, AAR summaries
- **Low-stakes (quick review):** Brainstorming scenarios, first drafts of briefs, formatting help

*Facilitator: "If you wouldn't sign a document written by a brand-new Lieutenant without reading it, don't sign one written by AI without reading it."*

---

**Skill 3: Task Decomposition**
*Breaking work into AI-appropriate chunks*

| 101 Behavior | 201 Behavior |
|---|---|
| "Plan my company's entire Phase II field exercise" | Break it into: (1) generate scenario framework using METT-TC, (2) build timeline from scenario, (3) draft OPORD shell, (4) create evaluation criteria per event. Run each as a separate prompt. |

**TBS decomposition example — Weekly counseling prep for 6 students:**
1. Pull each student's scores from Power BI (human task — requires judgment on what matters)
2. For each student: paste anonymized scores into Prompt #2, get counseling outline (AI task)
3. Review and personalize each outline (human task — you know the student)
4. Conduct counseling (human task — always)

*Time saved: ~30 minutes per student × 6 = 3 hours/week*

---

**Skill 4: Iterative Refinement**
*Moving from 70% to 95% through structured passes*

| 101 Behavior | 201 Behavior |
|---|---|
| Accepts first draft or gives up | "Good start, but make the leadership section more specific — reference the patrol he led on Tuesday where he failed to issue a FRAGO when the situation changed." |

**Iteration example (live demo if GenAI.mil available):**
1. First prompt: Generate AAR summary → output is generic
2. Second prompt: "Add specific TBS context — this was a squad live-fire, Phase II, and the main issue was fire discipline during the assault" → output improves significantly
3. Third prompt: "Format this as three bullet sustains and three bullet improves, each under 20 words" → output is briefing-ready

*Facilitator: Run this live. Show participants that the third output is dramatically better than the first.*

---

**Skill 5: Workflow Integration**
*Embedding AI into how work actually gets done*

| 101 Behavior | 201 Behavior |
|---|---|
| "I'll try that AI thing when I have time" (never has time) | "Every Friday at 1500 I generate my SITREP using Prompt #13. It takes 15 minutes instead of 45." |

**The 25-minute finding:**
The UK study didn't find people saving 4 hours on one task. They found people saving 2-5 minutes on dozens of small tasks throughout the day. That adds up to 25 minutes.

**TBS workflow integration points:**
- Monday: Review at-risk students on Power BI, run counseling prep prompts
- After any training event: Run AAR analysis prompt within 24 hours
- Friday: Generate weekly SITREP
- Before FEX: Run scenario generator, equipment tracker
- End of phase: Company metrics summary for command brief

---

**Skill 6: Frontier Recognition**
*Knowing when you're outside AI's capability boundary*

| 101 Behavior | 201 Behavior |
|---|---|
| "AI is great at everything" or "AI is useless" | "AI writes a solid first-draft counseling outline, but it can't assess a student's leadership under stress. That's my job." |

**TBS frontier map (starter):**

| AI Handles Well | AI Handles Poorly | Moving Frontier |
|---|---|---|
| Formatting documents | Evaluating leadership in the field | Scenario generation quality |
| Summarizing AAR notes | Knowing a specific student's character | Predictive analytics on student risk |
| Generating scenario frameworks | Making pass/fail judgment calls | Adaptive study recommendations |
| Drafting counseling language | Understanding unit culture/dynamics | Cross-event trend analysis |
| Data visualization design | Physical assessment of skills | - |

**Critical warning:**
> *"AI will never replace an SPC's judgment about a student. It's a tool that handles the administrative work so you have more time for the work that requires your expertise — observing, evaluating, mentoring, and leading."*

---

### Research Note for Facilitator
BCG-Harvard study: AI helps novices perform 34% better by disseminating expert knowledge. But untrained experts who blindly trust AI actually perform 19% worse. The skills you just learned are what prevent that 19% decline.

### Transition
*"Now that you know the six skills, let's talk about how to decide which tasks to delegate."*

---

## Module 3: The Delegation Equation
*0:40 – 0:55 (15 minutes)*

### Content

**Mollick's Three-Variable Framework**

For any task, ask three questions:

```
Should I delegate this to AI?

1. Human Baseline Time → How long would this take me manually?
2. Probability of Success → How likely is AI to get it right (or close)?
3. AI Process Time → How long to prompt + evaluate + verify + fix?

Delegate when: (Human Time) > (AI Process Time ÷ Probability of Success)
```

### TBS Walk-Through Examples

**Example 1: Counseling outline for a struggling student**
- Human baseline: 45 minutes (pull data, organize thoughts, write structure)
- Probability of success: ~75% (good structure, needs personalization)
- AI process time: 10 min (prompt + review + personalize)
- **Verdict: Delegate.** Save ~35 minutes.

**Example 2: Checking a student's PFT score against the standard**
- Human baseline: 2 minutes (look at the order, compare the number)
- Probability of success: 95%
- AI process time: 3 minutes (open GenAI.mil, type prompt, read answer, verify anyway)
- **Verdict: Don't delegate.** Checking takes longer than just doing it.

**Example 3: Writing a weekly SITREP from scratch**
- Human baseline: 45 minutes
- Probability of success: 80% (format correct, content needs review)
- AI process time: 15 minutes (fill in template prompt, review, adjust)
- **Verdict: Delegate.** Save ~30 minutes every week.

**Example 4: Writing a 5-paragraph order for a Phase III patrol**
- Human baseline: 3-4 hours
- Probability of success: 60% (structure good, tactical details need significant review)
- AI process time: 1.5 hours (prompt, review, fix tactical errors, iterate)
- **Verdict: Delegate the structure**, but expect significant human revision. Still saves time.

**Example 5: Deciding whether to recommend a student for Academic Hold**
- Human baseline: N/A — this is a judgment call
- **Verdict: Never delegate.** AI can help you organize the data supporting your decision, but the decision itself is yours and your commander's.

### Key Principle
> *"Delegate the drafting, keep the deciding."*

### Transition
*"Let's take a 10-minute break. When we come back, you're going to put your quality judgment to the test."*

---

## BREAK
*0:55 – 1:05 (10 minutes)*

---

## Module 4: The Trust Problem — Quality Judgment
*1:05 – 1:30 (25 minutes)*

### The Red Pen Review Exercise

*This is the highest-value activity in the entire course.*

**Setup (3 minutes):**
Distribute three AI-generated TBS documents to each participant (see Appendix A). Each document contains deliberate errors that test different aspects of quality judgment.

**The three documents:**

**Document 1: Student Counseling Summary**
- AI-generated counseling outline for an anonymous Phase I student
- Mostly correct structure and tone
- **Hidden error:** References "MCO 1500.61" as the authority for TBS counseling requirements — this order number is fabricated
- **Hidden error:** States the student's academic composite is calculated as "average of all four exams" — omits quiz average from the formula

**Document 2: Weekly SITREP**
- AI-generated company SITREP for Alpha Company
- Professional format, accurate-looking data
- **Hidden error:** States "14 of 192 students are at-risk (7.3%)" — the math is actually 14/192 = 7.29%, which rounds correctly, but the original data said 14 of 198 students. The AI changed the denominator.
- **Hidden error:** Lists "Range 4 live fire" as completed on a date that was actually a federal holiday (Martin Luther King Day)

**Document 3: AAR Summary**
- AI-generated AAR synthesis from a land navigation practical
- Good structure with sustains/improves
- **Hidden error:** Attributes a quote to "the company commander's guidance" that was never stated — the AI fabricated a plausible-sounding directive
- **Hidden error:** Recommends "additional day land nav training per MCO 3500.27A" — this MCO exists but doesn't address TBS land navigation training requirements specifically

**Exercise (12 minutes):**
- Participants read all three documents with a red pen
- Mark anything they wouldn't sign, wouldn't forward, or don't trust
- Work individually first (5 min), then discuss with a partner (5 min), then group debrief (facilitator-led)

**Debrief (10 minutes):**

Reveal each error. Ask:
- *"Who caught the fabricated MCO number?"* (Usually <30% catch it)
- *"Who caught the denominator change in the SITREP?"* (Usually <20%)
- *"Who caught the fabricated commander quote?"* (Usually <40%)

**Key lessons from the exercise:**

1. **AI fabricates references with total confidence.** If you don't recognize an MCO number, look it up before you cite it. Every time.

2. **AI changes numbers subtly.** When you give it data, verify the numbers in the output match what you provided. Don't assume it copied correctly.

3. **AI fabricates plausible quotes and attributions.** If the output says "the commander directed..." make sure the commander actually said that.

4. **The errors that matter most are the ones that look right.** Nobody questions "MCO 1500.61" if they don't know it's fake. That's why quality judgment requires domain expertise — and why AI doesn't replace the expert.

### Verification Checklist for TBS Documents

After any AI-generated document, check:
- [ ] All order/regulation references — do they exist? Do they say what the AI claims?
- [ ] All numbers — do they match your source data? Is the math correct?
- [ ] All quotes and attributions — did that person actually say/direct that?
- [ ] All dates — do they fall on actual training days?
- [ ] Tone — would you actually say this to the student/commander?
- [ ] Completeness — did the AI leave out something important you provided?

### Transition
*"Now you know what to watch for. Let's talk about how to build AI into your actual weekly workflow."*

---

## Module 5: Your Workflow — Centaur, Cyborg, or Neither
*1:30 – 1:45 (15 minutes)*

### Content

**Two productive patterns:**

**Centaur Mode** — Clear handoffs between human and AI
```
You do Phase 1 (gather data, decide what matters)
    → Hand off to AI for Phase 2 (draft the document)
        → You do Phase 3 (review, edit, approve)
```
Best for: Counseling prep, SITREP generation, evaluation write-ups — anything with a signature block.

**Cyborg Mode** — Continuous back-and-forth
```
You start → AI contributes → You adjust → AI refines → You finalize
```
Best for: Scenario development, AAR analysis, brainstorming training improvements — creative/iterative work.

**TBS examples:**

| Task | Mode | Why |
|------|------|-----|
| Counseling outline | Centaur | High-stakes document — clear human review before delivery |
| FEX scenario building | Cyborg | Creative, iterative — you refine until the scenario is right |
| Weekly SITREP | Centaur | Standard format, clear handoff points |
| AAR trend analysis | Cyborg | Exploratory — you ask follow-up questions as patterns emerge |
| Peer eval criteria | Centaur | Must be standardized — draft once, review, finalize |
| Training schedule optimization | Cyborg | Multiple constraints to balance — iterate to find the best fit |

### Workflow Mapping Exercise (10 minutes)

*Distribute Appendix B worksheet*

**Instructions:**
1. Pick one recurring task from your SPC/instructor job (counseling, AARs, scheduling, reporting)
2. Break it into 3-5 subtasks
3. For each subtask, mark: **Human Only** / **AI Could Help** / **AI Should Do This**
4. Identify the pattern: Is this Centaur or Cyborg?
5. Estimate time savings per occurrence

**Share out:** 2-3 volunteers share their workflow map. Facilitator highlights patterns.

### Transition
*"Last module. Let's map the frontier for TBS — where AI helps and where it doesn't."*

---

## Module 6: Frontier Mapping and Your Assignment
*1:45 – 2:00 (15 minutes)*

### Frontier Mapping Exercise (10 minutes)

**On whiteboard or butcher paper, create three columns:**

| AI Handles Well at TBS | AI Handles Poorly at TBS | Moving Frontier |
|---|---|---|

**Seed with 3-5 examples (facilitator provides these):**

*AI Handles Well:*
- Formatting counseling outlines from bullet points
- Summarizing multiple AAR inputs into themes
- Generating first-draft SITREPs from raw data
- Creating scenario frameworks with METT-TC structure
- Calculating weighted grade composites

*AI Handles Poorly:*
- Evaluating a student's leadership under pressure
- Knowing whether a student is genuinely struggling or sandbagging
- Understanding platoon dynamics and interpersonal issues
- Making pass/fail calls on subjective performance
- Anything requiring physical observation (did the student actually clear that room correctly?)

**Ask participants to add to each column (5 minutes).** This becomes the company's working frontier map.

**Key insight:** The frontier moves. Tasks that AI handles poorly today may improve in 6 months. Tasks it handles well today will get even better. Your job is to keep updating this map based on experience.

### Your Assignment

1. **This week:** Pick ONE task from your workflow map (Module 5). Use GenAI.mil and the appropriate prompt from the playbook to accomplish it.

2. **Document what happened:** Did it save time? Was the output useful? What did you have to fix?

3. **Share results at next week's standup** — especially failures. Failure cases are the most valuable data for the frontier map.

4. **Data rules reminder:**
   - CUI is authorized on GenAI.mil
   - Anonymize student PII — use "Student A" or "the student," never names or EDIPIs
   - No PHI ever
   - No classified data ever
   - When in doubt: describe the data pattern, don't paste the data

### Closing

*"The prompts and dashboards are tools. The six skills you learned today are what make those tools useful. The goal isn't to become an AI expert — it's to get 25 minutes back every day for the work that actually requires your expertise: observing, evaluating, mentoring, and leading Marines."*

**Distribute:**
- SPC Quick-Start Guide (printed)
- Quick Reference Card (laminated if possible)
- Prompt Playbook index

---

## Appendix A: Red Pen Review Exercise Documents

### Document 1: Student Counseling Summary

> **DEVELOPMENTAL COUNSELING OUTLINE**
>
> **Student:** [Phase I Student, Alpha Company, 2nd Platoon]
> **Counseling Type:** Midpoint Academic Performance Review
> **Authority:** Per MCO 1500.61, TBS counseling requirements
>
> **Performance Summary:**
> The student's academic performance requires attention. Current academic composite stands at 73.5%, calculated as the average of all four exam scores (Exam 1: 78%, Exam 2: 71%, Exam 3: 68%, Exam 4: 77%). This places the student below the 75% threshold for satisfactory academic performance.
>
> **Strengths:**
> - Physical fitness scores are well above average (PFT: 285, CFT: 294)
> - Land navigation day practical: Pass
> - Consistent attendance and professional bearing
>
> **Areas for Improvement:**
> - Academic exam scores show a declining trend (78 → 71 → 68) before partial recovery (77)
> - Phase I written exam performance is the primary risk factor
> - Study habits and preparation methods should be evaluated
>
> **Developmental Plan:**
> 1. Establish a structured study schedule with accountability checks
> 2. Attend academic support sessions (Tuesday/Thursday 1800-1930)
> 3. Form a study group with peers who demonstrate strong academic performance
> 4. Review and correct all missed exam questions within 48 hours of score release
>
> **Follow-up:** Counseling review in 14 training days to assess improvement trajectory.

**Errors to find:**
1. "MCO 1500.61" — fabricated order number
2. "Average of all four exam scores" — TBS academic composite also includes quiz average; the formula is incomplete
3. Bonus catch: the student is said to be in Phase I but has all four exam scores — Exams 2-4 are Phase II-IV exams

---

### Document 2: Weekly SITREP

> **ALPHA COMPANY WEEKLY SITREP**
> **Period:** 13-17 January 2026
> **Classification:** CUI
>
> **1. PERSONNEL**
> - Assigned: 192 | Present: 186 | Medical Hold: 4 | Leave: 2
> - At-risk students: 14 of 192 (7.3%)
>   - Academic: 6 | MilSkills: 3 | Leadership: 2 | Multiple: 3
>
> **2. TRAINING COMPLETED**
> - 13 JAN: Range 4 live fire qualification (M16A4) — 184 fired, 2 refires scheduled
> - 14 JAN: Land Navigation written exam — avg score 81.3%
> - 15-16 JAN: Squad tactics practical evaluation
> - 17 JAN: Academic Exam 2 — scores pending
>
> **3. UPCOMING**
> - 20-21 JAN: Night land navigation practical
> - 22 JAN: Combat lifesaver certification
> - 23-24 JAN: MOUT training, Range 220A
>
> **4. ISSUES/CONCERNS**
> - 3 students flagged for declining trend across consecutive evaluations
> - Range 8 unavailable 20-24 JAN (WTBN priority) — shifted to Range 12
> - 1 instructor (SSgt Williams) RSO certification expires 28 JAN — renewal scheduled 22 JAN
>
> **5. COMMANDER'S ASSESSMENT**
> Company is tracking to standard. Academic at-risk count increased by 2 since last week; SPC counseling sessions scheduled for all 6 academic at-risk students NLT 21 JAN.

**Errors to find:**
1. 13 January 2026 is Martin Luther King Jr. Day (federal holiday) — no live fire would be scheduled
2. "14 of 192 (7.3%)" — earlier data said 198 students, not 192. The AI changed the denominator without flagging it.

---

### Document 3: AAR Summary

> **AAR SYNTHESIS: Land Navigation Day Practical**
> **Date:** 10 January 2026
> **Company:** Alpha | **Phase:** I
>
> **Event Summary:**
> 186 students executed day land navigation practical at Camp Barrett training area. Course consisted of 6 points over 5km in moderately wooded terrain. Time limit: 4 hours. Pass requirement: 4 of 6 points.
>
> **Results:** 158 Pass (84.9%) | 28 Fail (15.1%)
>
> **SUSTAINS:**
> 1. Route planning instruction in the classroom translated well to practical execution — students who attended the optional study session had a 94% pass rate
> 2. Safety plan was effective — no injuries, all students accounted for at all checkpoints
> 3. Per the company commander's guidance to "prioritize deliberate pace over speed," students who took longer generally scored better
>
> **IMPROVES:**
> 1. Terrain association skills remain weak — 60% of failures occurred on points requiring contour interpretation rather than direct azimuth
> 2. Protractor use: 12 students reported difficulty with their protractors, suggesting additional PMI is needed
> 3. Per MCO 3500.27A, additional day land nav training should be scheduled for students who failed to achieve the minimum standard
>
> **TRENDS:**
> This event's 84.9% pass rate is consistent with historical Alpha Company performance (82-87% across the last 4 classes). Failure patterns correlate with terrain association, not compass/pace count mechanics.
>
> **RECOMMENDATIONS:**
> 1. Add a 1-hour terrain association refresher between classroom and practical
> 2. Schedule 2-hour remedial land nav for the 28 students who failed, focusing on contour interpretation
> 3. Update the study guide to include more contour-based practice problems

**Errors to find:**
1. "Per the company commander's guidance to 'prioritize deliberate pace over speed'" — fabricated quote/attribution
2. "Per MCO 3500.27A" — this MCO exists (Training and Education) but doesn't specifically direct additional land nav training for TBS students who fail

---

## Appendix B: Workflow Mapping Worksheet

```
NAME: _________________________ ROLE: _________________________

TASK: _________________________________________________________
(Pick one recurring task from your job)

Step 1: Break it into subtasks

| # | Subtask | Human Only | AI Could Help | AI Should Do |
|---|---------|:---:|:---:|:---:|
| 1 | | | | |
| 2 | | | | |
| 3 | | | | |
| 4 | | | | |
| 5 | | | | |

Step 2: Identify the pattern
[ ] Centaur — Clear handoffs (human → AI → human)
[ ] Cyborg — Continuous back-and-forth
[ ] Human Only — AI doesn't help here

Step 3: Estimate impact
Current time per occurrence: _______ minutes
Estimated time with AI: _______ minutes
Frequency: _______ times per week/month
Estimated weekly savings: _______ minutes

Step 4: Which prompt(s) from the playbook would you use?
Prompt #___: _________________________________________________
Prompt #___: _________________________________________________

Step 5: What's your plan for this week?
I will use Prompt #___ for ______________________________________
by _____________ (day).
```

---

## Appendix C: Facilitator Certification

Per EDD SOP Section 8, facilitators for Course 1 must:

1. Have completed Course 1 themselves (as a participant)
2. Have used AI tools in their own workflow for at least 2 weeks
3. Have reviewed the full TBS Prompt Playbook
4. Be able to demonstrate live prompting on GenAI.mil
5. Understand TBS data handling rules (CUI/PII/PHI boundaries)

Recommended: Complete EDD Course 2 (Builder Orientation) before facilitating Course 1, though not required.

---

*This course is adapted from the Expert-Driven Development AI Fluency Fundamentals course. Original curriculum by SSgt Jesse C. Morgan.*
