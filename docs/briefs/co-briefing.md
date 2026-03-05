# The Heywood Initiative: CO/XO Briefing

**Classification:** UNCLASSIFIED // Distribution Unlimited
**Format:** 30-minute briefing (slide-by-slide content with speaker notes)
**Presenter:** SSgt Jesse C. Morgan
**Audience:** TBS Commanding Officer / Executive Officer
**Purpose:** Approve Phase 1 pilot deployment; establish TBS POC for Phase 2 coordination

---

## Slide 1: Title

**THE HEYWOOD INITIATIVE**
**Making the Vision Real**

SSgt Jesse C. Morgan
March 2026

*Phase 1 Proof of Concept — Ready to Deploy*

### Speaker Notes

Good morning, sir/ma'am. Thank you for the time. This briefing is 30 minutes. I'm going to show you what already exists, what it does, and what I need from TBS to put it in front of instructors. Everything I'm presenting today is built on platforms already authorized inside MCEN. Phase 1 requires no additional authorization and costs between zero and twelve hundred dollars.

---

## Slide 2: The Problem

**Time: 2 minutes**

- TBS evaluates approximately 1,600 students per year across 8 companies
- Three grading pillars — Academics (32%), Military Skills (32%), Leadership (36%) — tracked across a 26-week POI
- SPCs currently manage student data through a combination of spreadsheets, shared drives, and manual records
- No integrated view of individual student performance across all three pillars
- Instructor qualification tracking is manual and reactive — expiring certifications are caught late or missed
- The TBS AI Vision 2026 deck identified 15 use cases for AI integration — this briefing shows how to start delivering on them using platforms TBS already has access to

### Speaker Notes

The core problem is not a lack of data — TBS generates significant amounts of student performance data every cycle. The problem is that the data lives in disconnected locations, requires manual aggregation, and does not give SPCs or leadership a real-time picture of where students stand. An SPC preparing a counseling packet has to pull information from multiple sources and assemble it by hand. A company commander who wants to know which students are at risk of failing has to ask an SPC to build that picture manually. Instructor quals are the same story — we find out someone's range certification lapsed when we're trying to staff next week's live-fire.

The AI Vision 2026 deck lays out the right objectives. What I'm presenting is the technical plan to start achieving them — not as a research project, but as a series of deliverables that work on day one.

---

## Slide 3: What's Already Built

**Time: 3 minutes**

- **20 TBS-adapted AI prompts** for GenAI.mil
  - Counseling prep, AAR analysis, scenario generation, student performance analysis, PIA screening, and 15 more
  - Each prompt is written for the TBS context — references the 3-pillar system, POI phases, SPC workflows
- **SharePoint list schemas** for student data, training schedule, instructor qualifications, event feedback
  - Designed to feed directly into Power BI dashboards
- **Power BI dashboard specifications** with DAX measures for TBS 3-pillar grading
  - Student performance dashboard spec complete
  - Instructor qualification dashboard spec complete
- **Full governance package**
  - PIA threshold analysis
  - Compliance checklist for each deployed component
  - Tool registry entry for sustainment tracking
- **Deployment automation scripts**
- **Total cost to date: $0**
  - Built entirely on existing authorized platforms using organic labor

### Speaker Notes

I want to be clear about what "already built" means. These are not concepts or proposals. The prompt playbook is 20 finished prompts that an SPC can paste into GenAI.mil today and get useful output — counseling language, AAR analysis, scenario frameworks. The SharePoint schemas are ready to deploy as lists. The dashboard specs include the exact DAX calculations for weighted grading across the three pillars with row-level security already designed. The governance artifacts — PIA threshold, compliance checklist, tool registry — are written and ready for review.

None of this required procurement, contractor support, or additional authorization. It was built using platforms already available on MCEN.

---

## Slide 4: Phase 1 Deliverables

**Time: 5 minutes**

### 1. Power BI Student Performance Dashboard
- Three-pillar weighted view: Leadership 36%, Military Skills 32%, Academics 32%
- Company-level overview with student rankings
- At-risk student flagging: any student below 75% in any pillar is highlighted
- Individual student drill-through: click a name, see every grade, every event, trend over time
- Phase progression tracking across the 26-week POI
- Row-Level Security enforced:
  - Staff sees all companies
  - SPC sees own company only
  - Student sees own record only

### 2. Power BI Instructor Qualification Dashboard
- All instructor certifications in one view
- 30/60/90-day expiration alerts — color-coded
- Coverage gap analysis: can we staff next week's scheduled events?
- Workload distribution across SPCs
- Rotation risk flags: instructors PCSing who hold critical qualifications

### 3. GenAI.mil Prompt Playbook for SPCs
- 20 prompts covering the most time-consuming SPC administrative tasks
- Printed quick-reference cards for each prompt
- No training required beyond the initial workshop

### 4. AI Fluency Workshop
- 2-hour session adapted from Expert-Driven Development curriculum
- Hands-on with GenAI.mil using TBS-specific prompts
- Designed for TBS instructors — not a generic AI overview

### 5. Governance Artifacts
- PIA threshold analysis (ready for review/signature)
- Compliance checklist per component
- Tool registry for long-term sustainment

**All of the above operates within the existing MCEN M365 authorization boundary. No ATO required.**

### Speaker Notes

I want to walk through each deliverable because I want you to see that these are concrete, usable tools — not prototypes.

The student performance dashboard gives an SPC or company commander a single screen showing every student's standing across all three pillars. If a student drops below 75% in any pillar, the dashboard flags it automatically. You can click on any student and see their full history — every graded event, every score, trendline over time. Row-Level Security means an SPC only sees their company's data. A student, if given access, only sees their own record. Staff sees everything.

The instructor qualification dashboard solves the problem of finding out too late that someone's certification expired. It shows every qualification, color-coded by days remaining. Green is good, yellow is 60 days out, red is 30 days or less. It also shows coverage — if we have three events scheduled next week that require a specific qualification, and only two instructors hold it, that gap shows up automatically.

The prompt playbook is 20 finished prompts. An SPC pastes one into GenAI.mil, fills in the specifics, and gets structured output — counseling language, AAR analysis, scenario frameworks. These are not generic prompts. They reference the TBS grading system, the POI, and SPC-specific workflows.

The workshop is 2 hours, hands-on. Instructors leave with the ability to use every prompt in the playbook and understand when AI output needs human review.

None of this requires anything beyond what TBS already has on MCEN. Power BI, SharePoint, and GenAI.mil are all available today.

---

## Slide 5: How It Works — Student Performance

**Time: 3 minutes**

**[Dashboard mockup or screenshot would go here]**

### Data Flow
1. Student grades entered into SharePoint list (structured schema provided)
2. Power BI connects to SharePoint as data source
3. DAX measures calculate weighted pillar scores automatically
4. Dashboard refreshes on schedule or on demand

### Three-Pillar View
| Pillar | Weight | Components |
|--------|--------|------------|
| Academics | 32% | Exam scores, written assignments, course grades |
| Military Skills | 32% | Range qualifications, field exercise evaluations, practical applications |
| Leadership | 36% | Peer evaluations, SPC assessments, leadership billets, BARS indicators |

### Dashboard Pages
- **Company Overview** — All students ranked by composite score, color-coded by risk
- **Pillar Breakdown** — Performance distribution across each pillar
- **Individual Student** — Drill-through to full history and trend
- **At-Risk Report** — Students below threshold in any pillar, sortable by severity
- **Phase Progression** — Where each student stands relative to POI timeline

### Row-Level Security Model
| Role | Sees |
|------|------|
| TBS Staff (CO, XO, S-3) | All companies, all students |
| Company Commander | Own company |
| SPC | Own company |
| Student | Own record only |

Security enforced through Azure Entra ID groups — same groups already used for MCEN access.

### Speaker Notes

The data flow is straightforward. Student grades go into a SharePoint list using the schema I've already built. Power BI reads that list and calculates everything automatically — weighted scores, pillar breakdowns, risk flags. The dashboard refreshes on whatever schedule TBS wants, or an SPC can hit refresh manually.

Row-Level Security is the key feature here. This is not a shared spreadsheet where anyone can see everything. Security is enforced at the data layer through Azure Entra ID — the same identity system that controls CAC login to MCEN. An SPC logs into Power BI and only sees their company. A staff officer sees all companies. This is built into the dashboard design from the start, not added later.

The at-risk report is where this saves the most time. Instead of an SPC manually checking every student's grades across multiple sources, the dashboard automatically flags anyone below 75% in any pillar and ranks them by severity. That is the starting point for every counseling session and every company commander's weekly assessment.

---

## Slide 6: How It Works — Instructor Quals

**Time: 2 minutes**

**[Dashboard mockup or screenshot would go here]**

### Qualification Tracking
- Every instructor certification in a single SharePoint list
- Expiration dates drive automated alerting:
  - **Green:** 90+ days remaining
  - **Yellow:** 60 days — plan renewal
  - **Red:** 30 days — schedule now
  - **Black:** Expired — instructor cannot be assigned

### Coverage Analysis
- Cross-reference scheduled events against qualified instructor pool
- Identify gaps before they affect the training schedule
- Example: 3 live-fire events next week, only 2 RSO-qualified instructors available = gap flagged

### Workload Distribution
- How many events is each SPC assigned per week/month?
- Are some instructors carrying disproportionate load?
- Sortable by company, qualification type, or time period

### Rotation Risk
- Instructors with PCS orders who hold critical qualifications
- Lead time to train replacements
- Prioritized by qualification scarcity

### Speaker Notes

This dashboard solves a problem every company deals with — finding out an instructor's qualification expired when you're trying to build next week's schedule. The data lives in a SharePoint list. When a qualification is entered, the expiration date is required. Power BI calculates days remaining and color-codes automatically.

The coverage analysis is the most operationally useful view. It looks at scheduled events, checks which qualifications are required, and compares that against the pool of qualified instructors. If there is a gap, it shows up in red before the schedule goes out — not the day of execution.

Rotation risk is the long-term planning view. If an instructor with a rare qualification is PCSing in 60 days, leadership needs to know now so they can get someone else qualified. This dashboard surfaces that automatically.

---

## Slide 7: The Path Forward — 4 Phases

**Time: 5 minutes**

### Phase 1: Proof of Concept — Months 0-3
- **Cost:** $0–$1,200
- **Authorization:** None required (existing MCEN/M365)
- **Delivers:** Power BI dashboards, prompt playbook, AI training, governance docs
- **Scope:** One company pilot

### Phase 2: Custom AI Agent with IATT — Months 3-6
- **Cost:** $3,600–$4,650
- **Authorization:** IATT (Interim Authority to Test), 90–180 days, one company
- **Delivers:** Heywood Power App with role-based views, Azure OpenAI RAG pipeline over TBS doctrine, BARS-adapted assessment forms, AAR intake and synthesis system
- **New capability:** AI responses grounded in TBS doctrine, not generic LLM output

### Phase 3: Scale and Integration — Months 6-12
- **Cost:** $32,600–$49,800
- **Authorization:** Full ATO or cATO
- **Delivers:** Dataverse unified student store, MCTIMS integration, Power Automate workflows (automated counseling packets, grade processing, command reports), all 8 companies
- **New capability:** Automated workflows replace manual SPC administrative processes

### Phase 4: Full Heywood — Months 12-18
- **Cost:** $41,200–$59,200
- **Authorization:** ATO modification
- **Delivers:** Cross-cycle predictive analytics, AI scenario generation engine, field-deployable tablet app for SPCs during exercises, full knowledge transfer package
- **New capability:** Predictive student risk flags, offline-capable field tools

### Cost Comparison
| Program | Cost | Scope |
|---------|------|-------|
| **Heywood (18 months)** | **$77K–$115K** | TBS student management + AI integration |
| West Point Azimuth | $2.18M (3 years) | Comparable Power Platform modernization |
| Army ATIS | Enterprise scale | Multi-year, multi-site |

### Key Design Principle
Each phase delivers standalone value. If funding stops after any phase, TBS keeps everything built to that point. There is no phase that requires completion of later phases to be useful.

### Speaker Notes

I want to emphasize the design principle here: each phase stands on its own. If Phase 1 is all we ever do, TBS has working dashboards, a prompt playbook, and trained instructors — permanently. If we get through Phase 2 and funding or priorities change, TBS has a custom AI agent with doctrine-based responses and structured assessment tools. Nothing is throwaway.

The cost comparison to West Point is relevant because Azimuth is the closest comparable program — same mission (student management modernization), same platform (Power Platform on Azure Government), same scale. They contracted it out at $2.18M over 3 years. This plan achieves comparable scope at 3 to 5 percent of that cost because the labor is organic. I am building this. There is no contractor.

Phase 2 is where the AI capability becomes significant. GenAI.mil prompts in Phase 1 are useful but generic — the AI does not know TBS doctrine. In Phase 2, we build a RAG pipeline that indexes the POI, doctrinal publications, TBS SOPs, and AAR summaries. When an SPC asks the AI a question, it answers based on TBS-specific sources, not general internet knowledge.

Phase 3 is about scale and automation. Instead of one company, all eight. Instead of manual grade calculation, automated workflows. Instead of SPCs building counseling packets by hand, the system aggregates data and generates a draft that the SPC reviews and edits.

Phase 4 adds predictive capability and field tools. Which students are most likely to struggle in the next phase? Which instructors are approaching burnout based on workload data? A tablet app that works during field exercises without network connectivity.

---

## Slide 8: What I Need From TBS

**Time: 3 minutes**

### For Phase 1 (the ask today)

| Requirement | Details | Level of Effort |
|-------------|---------|----------------|
| Access to one company's SPC team | 4–6 instructors for a 4-hour training session | One afternoon |
| CO/XO time for this briefing | 30 minutes | Done |
| Sample student data | Anonymized or synthetic — enough to populate live dashboards | S-3 or company XO provides |
| M365 / Power BI Pro licenses | Already available on MCEN — just need assignment | S-6 action |
| TBS POC designation | One person to coordinate Phase 2 IT requirements | S-3 rep or company XO recommended |

### What Jesse Provides (no cost to TBS)

- All technical design, development, and deployment
- EDD training curriculum delivery (5 courses adapted for TBS)
- 20 TBS-specific prompt templates
- Power BI dashboard development
- All governance documentation (PIA, compliance, registry)
- Ongoing maintenance and iteration during pilot
- Full documentation for sustainment

### Speaker Notes

The ask is small relative to what gets delivered. I need one afternoon with one company's SPCs. I need sample data — it can be anonymized or completely synthetic, I just need realistic data to build live dashboards instead of static mockups. I need Power BI Pro licenses assigned, which are already part of MCEN's M365 licensing — S-6 just needs to assign them. And I need a single point of contact at TBS who I coordinate with for Phase 2 IT requirements, because IATT coordination and Azure Gov subscription access have lead times that need to start early.

Everything else — the build, the training, the governance documentation, the maintenance — I provide. There is no contractor cost. There is no procurement action. Phase 1 is my labor on platforms TBS already owns.

---

## Slide 9: Risk Mitigation

**Time: 3 minutes**

| Risk | Mitigation |
|------|------------|
| **PII exposure** | No real student PII until PIA threshold analysis is reviewed and approved by the appropriate authority. Phase 1 pilot uses anonymized or synthetic data. |
| **Unauthorized platform** | Phase 1 uses only platforms already authorized inside the MCEN M365 boundary — SharePoint, Power BI, GenAI.mil. No new systems. |
| **ATO timeline** | IATT/ATO coordination begins during Phase 1 in parallel. Phase 1 does not depend on it. If ATO takes longer than planned, Phase 1 tools remain operational. |
| **MCTIMS access denied or delayed** | Heywood functions without MCTIMS data. Manual data entry is the fallback — less efficient but fully operational. |
| **MHS GENESIS PIA denied** | Phase 4 is designed to work without medical data. It enhances the system but is not a dependency. |
| **SPC resistance to adoption** | SPCs are involved from the start — they see the tools, provide feedback, shape the workflow. The prompt playbook reduces their workload; it does not add to it. |
| **Jesse PCSes or transfers** | Replication guide written at every phase gate. 2–3 TBS-organic maintainers trained during Phase 3. All code, configurations, and documentation stored in TBS-owned repositories. |
| **Scope creep** | Each phase has defined exit criteria. Phase 2 does not begin without CO/XO approval at the Phase 1 review. |

### Speaker Notes

I want to address the risks directly because they are real and I want you to know I have accounted for them.

PII is the most important one. Phase 1 does not touch real student PII. We use anonymized or synthetic data until the PIA threshold analysis is reviewed and approved. That document is already written and ready for your review. When we move to real data, it will be after the appropriate privacy review is complete.

Platform authorization is straightforward — everything in Phase 1 is already inside MCEN. SharePoint, Power BI, GenAI.mil. There is nothing new to authorize.

The personnel risk is real. If I PCS, this project needs to survive without me. That is why documentation happens at every phase, not at the end. And in Phase 3, I specifically train 2 to 3 TBS personnel to maintain and extend the system. The goal is organic capability, not dependence on one person.

Scope creep is controlled by the phase gate structure. Each phase has defined deliverables and exit criteria. We brief results before moving to the next phase. If Phase 1 does not demonstrate value, we stop — and TBS still keeps the dashboards and prompts.

---

## Slide 10: The Ask

**Time: 2 minutes**

### Requesting CO/XO approval for the following:

1. **Approve Phase 1 pilot** with one company
   - Recommend Alpha or Bravo Company (whichever is earliest in their POI cycle)
   - Pilot runs 60 days

2. **Designate a TBS POC**
   - S-3 representative or company XO
   - Single point of contact for data access, scheduling, and Phase 2 IT coordination

3. **Schedule SPC training session**
   - 4 hours, within 2 weeks of approval
   - 4–6 instructors from the pilot company

4. **Next briefing: Phase 1 results at the 60-day mark**
   - Quantitative metrics: SPC time savings, dashboard adoption, prompt usage frequency
   - Qualitative feedback from pilot company SPCs
   - Recommendation to proceed, modify, or halt Phase 2

### Speaker Notes

The ask is four items. Approve the pilot with one company — I recommend whichever company is earliest in their POI cycle so we capture the most data during the pilot window. Designate one person as my POC so I am not chasing coordination across multiple offices. Schedule the SPC training within two weeks so we can start the 60-day clock. And agree to a 60-day results brief where I come back with hard numbers — how much time did SPCs save, how often are the dashboards being used, what does the feedback look like.

At that 60-day brief, I will present a recommendation to proceed to Phase 2, modify the approach, or stop. That decision point is built into the plan. You are not approving 18 months of work today. You are approving a 60-day proof of concept with a defined review point.

---

## Slide 11: Questions

**Time: 2 minutes**

### Anticipated Questions

**Q: What if the AI gives bad information to an SPC?**
A: All AI output in Phase 1 is advisory. An SPC uses GenAI.mil prompts the same way they would use any reference — the output is a starting point that requires human review and judgment. No AI output feeds directly into any official record without SPC review.

**Q: Who else is doing this in the Marine Corps?**
A: No other TBS-equivalent schoolhouse has deployed integrated AI tools at this level. CGSC (Army) ran a successful AI wargaming pilot. West Point contracted a comparable Power Platform modernization at $2.18M. We would be the first Marine Corps schoolhouse to build this organically.

**Q: What happens to student data after they graduate?**
A: Data retention follows existing TBS records management policy. Nothing in this system changes data retention requirements. Phase 3 would integrate with MCTIMS, which is the system of record.

**Q: Can other schoolhouses use this?**
A: Yes. Phase 4 includes a knowledge transfer package and adaptation guide specifically for replication at other TECOM schoolhouses. The architecture is not TBS-specific — it is a pattern that works anywhere the same platforms are available.

### Speaker Notes

I will take any questions. I have backup material on the Expert-Driven Development methodology if you want to understand the training approach in more detail.

---

## Slide 12: Backup — Expert-Driven Development Foundation

### What Is EDD?

Expert-Driven Development is the methodology behind everything I just briefed. It is a structured approach to building AI-integrated tools where the subject matter expert — in this case, an SPC or instructor — drives the development process using AI as a force multiplier.

### EDD Curriculum (Built and Ready)
- **5 courses** covering AI fluency through advanced tool building
- **51 prompt templates** across domains (analysis, creation, assessment, governance)
- **12-section SOP** for organizational AI deployment
- Full governance framework including PIA, compliance, and tool registry templates

### How It Applies to TBS
- The 20 TBS prompts in the playbook were built using EDD methodology
- The SharePoint schemas were designed using EDD data modeling prompts
- The dashboard specs were generated using EDD analytics prompts
- The governance artifacts follow EDD compliance templates

### Key Point
This is not a vendor product. It is not a contract deliverable. It is an organic DoD capability — a methodology that any qualified Marine can learn and apply. The tools I am deploying at TBS were built using this methodology, and the methodology itself is part of what gets transferred to TBS in the knowledge transfer package.

### Speaker Notes

EDD is worth understanding because it explains how one SSgt built everything I just showed you in a short period of time with zero budget. The methodology uses AI tools to accelerate every phase of development — requirements gathering, data modeling, dashboard design, governance documentation. The same methodology is what I teach SPCs in the workshop, so they can extend and adapt the tools themselves.

This is not about making TBS dependent on AI. It is about giving Marines a repeatable process for using AI to solve administrative problems so they can spend more time on the training mission.

---

## Appendix: Document References

| Document | Location | Status |
|----------|----------|--------|
| TBS AI Vision 2026 | TBS S-3 | Published |
| Heywood Build Plan v1 | This package | Complete |
| PIA Threshold Analysis | Governance package | Ready for review |
| Compliance Checklist | Governance package | Complete |
| Tool Registry Entry | Governance package | Complete |
| Prompt Playbook (20 prompts) | Playbook package | Complete |
| Student Dashboard Spec | Dashboard package | Complete |
| Instructor Quals Dashboard Spec | Dashboard package | Complete |
| SharePoint List Schemas | Schema package | Complete |
| EDD Curriculum | Training package | Complete |

---

**Classification:** UNCLASSIFIED // Distribution Unlimited
**Prepared by:** SSgt Jesse C. Morgan
**Date:** March 2026
**Version:** 1.0
