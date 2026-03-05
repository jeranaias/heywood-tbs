# The Heywood Initiative: Technical Build Plan

**Classification:** UNCLASSIFIED // Distribution Unlimited
**Version:** 1.0
**Date:** March 2026
**Prepared by:** SSgt Jesse C. Morgan

---

## Executive Summary

The Heywood Initiative envisions a role-based AI agent system for The Basic School — a "Digital Staff Officer" for every Marine in the training pipeline. This document presents a 4-phase build plan that delivers Heywood incrementally using platforms already authorized inside MCEN, starting with zero additional authorization required.

**Key facts:**
- Total estimated cost: **$77K–$115K over 18 months**
- Phase 1 costs **$0–$1,200** and requires **no ATO**
- Every platform component is already IL5-authorized on Azure Government or MCEN
- Each phase delivers standalone value — TBS keeps everything built regardless of future funding
- The training methodology (Expert-Driven Development curriculum) is already built and ready to deploy

**For comparison:** West Point awarded Arch Systems $2.18M over 3 years for a comparable Power Platform/Dataverse modernization of their cadet management system. This plan delivers equivalent scope at 3–5% of that cost by using organic labor and existing infrastructure.

---

## What Heywood Actually Requires

The TBS AI Vision 2026 deck identifies 15 use cases and describes capabilities including "Multi-Agent Reinforcement Learning" and "Collective Intelligence." The technical reality is more straightforward — and that's good news, because it means this is achievable without a research program:

| Vision Deck Language | What's Actually Needed | Status |
|---------------------|----------------------|--------|
| Multi-Agent Reinforcement Learning | Retrieval-Augmented Generation (RAG) | Production-ready; available on Azure Gov IL5 |
| Collective Intelligence | Role-based AI agent with shared knowledge base | Power App + Azure OpenAI + Azure AI Search |
| AI Agent with RBAC | Power App with Azure Entra ID security groups | Standard Power Platform pattern |
| Multimodal (text, photo, audio) | GPT-4o on Azure OpenAI (supports text + image) | Available on Azure Gov today |
| Data Repository | Dataverse + Power BI | Available on MCEN today |
| Agentic Workflows | Power Automate flows + Azure OpenAI API calls | Standard Power Platform pattern |

The 3-layer architecture from the vision deck (Platform Model → RBAC → Decentralized Dashboards) maps directly to: **Azure OpenAI → Azure Entra ID → Power BI**.

---

## Platforms Available on MCEN (All IL5-Authorized)

| Platform | Authorization | Use In Heywood |
|----------|--------------|----------------|
| Azure OpenAI Service | FedRAMP High + DoD IL5 PA | LLM backend (GPT-4o, GPT-4.1) |
| Azure AI Search | IL5 on Azure Gov (with CMK) | Vector store for RAG pipeline |
| GenAI.mil | IL5, CUI-authorized, 1.1M users | Phase 1 AI interface (web-based) |
| Power Apps | Available on MCEN (M365) | Application front-end with RBAC |
| Power Automate | Available on MCEN (M365) | Workflow automation |
| Power BI | Available on MCEN (M365) | Dashboards with Row-Level Security |
| Dataverse | Available (premium license) | Unified student data store |
| SharePoint Online | Available on MCEN | Document storage, Phase 1–2 data |
| Azure Entra ID | In use (CAC auth) | RBAC enforcement |
| Azure Sentinel | Available on Azure Gov | Continuous monitoring for cATO |

---

## 4-Phase Build Plan

### Phase 1: Proof of Concept — Months 0–3

**Authorization:** None required. All tools within existing MCEN/M365 boundary.

**Cost:** $0–$1,200

**What gets delivered:**

1. **TBS-Adapted AI Training** (from Expert-Driven Development curriculum)
   - 2-hour AI Fluency workshop for one company's SPCs (4–6 instructors)
   - 30-minute Supervisor briefing for CO/XO
   - Builder Orientation for 3–5 designated TBS power users

2. **TBS Prompt Playbook** — 20 prompts adapted for TBS workflows
   - AAR preparation and synthesis
   - Counseling session prep
   - Scenario generation
   - Student performance analysis
   - PIA screening for student data
   - SharePoint schema design for TBS data
   - Dashboard specifications

3. **Power BI Student Performance Dashboard**
   - Three-pillar view: Leadership (36%), Military Skills (32%), Academics (32%)
   - Row-Level Security: Staff sees all, Instructor sees company, Student sees individual
   - Phase progression tracking across 26-week POI

4. **Power BI Instructor Qualification Dashboard**
   - Certification tracking with 30/60/90-day expiration alerts
   - Workload distribution across companies

5. **Governance Artifacts**
   - PIA threshold analysis for all student data
   - Compliance checklist for each deployed component
   - Tool registry for sustainment

**What TBS provides:** Access to one company's SPCs (4 hours), 30-minute CO/XO briefing slot, sample student data, M365 licenses

**Exit criteria for Phase 2:** 4+ SPCs using prompts weekly, at least one dashboard in active use, CO/XO approval to proceed

---

### Phase 2: Custom Agent with IATT — Months 3–6

**Authorization:** IATT (Interim Authority to Test) — 90–180 days, scoped to one company

**Cost:** $3,600–$4,650

**What gets delivered:**

1. **Heywood Power App** — Canvas app with three role views
   - Staff: Full visibility, aggregate analytics, schedule management
   - Instructor: Company-scoped student performance, counseling tools, AAR input
   - Student: Individual performance view, study recommendations
   - CAC authentication via Azure Entra ID

2. **Azure OpenAI + RAG Pipeline** — TBS-specific AI knowledge base
   - GPT-4o for reasoning/analysis
   - Azure AI Search index over: POI documents, doctrinal publications, TBS SOPs, AAR summaries
   - Responses grounded in TBS doctrine, not generic LLM knowledge

3. **BARS-Adapted Assessment Framework**
   - Structured behavioral indicators replacing subjective SPC rankings
   - Power App forms aligned to TBS grade pillars
   - AI-assisted narrative generation: SPC enters observations → structured counseling language

4. **AAR Intake and Synthesis System**
   - Structured AAR data entry forms
   - Automated theme extraction, categorization, and action item routing

5. **MCTIMS Coordination** — Formal data extract request submitted to TECOM MISSO

**What TBS provides:** Azure Gov subscription access, IATT application support (S-6), pilot company commitment, POI documents for RAG index, SME input on BARS behavioral indicators

---

### Phase 3: Scale and Data Integration — Months 6–12

**Authorization:** Full ATO or cATO (inherits ~60% of NIST 800-53 controls from Azure Gov)

**Cost:** $32,600–$49,800

**What gets delivered:**

1. **Dataverse Unified Student Store** — Single source of truth
   - Entities: Student, Event, Score, Instructor, BARS Evaluation, AAR Entry, HPB Record
   - Row-level security tied to Azure AD groups

2. **MCTIMS Integration** — Training record data flowing into Dataverse

3. **Power Automate Workflow Suite**
   - AAR processing: submit → AI extract themes → categorize → route action items
   - Counseling packet generation: aggregate all student data → AI narrative → SPC review
   - Instructor alerts: below-threshold students, expiring qualifications, overdue counseling
   - Automated grade processing with weighted calculations
   - Weekly command report generation

4. **Human Performance Data Integration** — PT scores, injury data, wellness indicators

5. **Predictive Analytics MVP** — Student risk flags (academic failure, injury risk, workload imbalance)
   - All predictions require human review — no automated decisions

6. **Cross-Company Expansion** — All 8 TBS companies (~1,600 students)

**What TBS provides:** Dataverse environment, TECOM MISSO data extracts, HPB data sharing, ATO coordination (S-6/ISSM), historical data for model training

---

### Phase 4: Full Heywood — Months 12–18

**Authorization:** ATO modification for expanded data scope

**Cost:** $41,200–$59,200

**What gets delivered:**

1. **MHS GENESIS Integration** (conditional on PIA approval) — Aggregate medical readiness data, staff access only
2. **MCTFS Personnel Data** — Prior MOS, TIS, assignments for predictive analysis
3. **Cross-Cycle Trend Analysis** — Which events predict final standing, earliest attrition indicators
4. **Scenario Generation Engine** — AI-generated tactical scenarios calibrated to phase, difficulty, and objectives
5. **Field-Deployable PWA** — Tablet interface for SPCs during field exercises (offline-capable)
6. **Knowledge Transfer Package** — Full documentation, 2–3 TBS-organic maintainers trained, Adaptation Guide for other schoolhouses

---

## Use Case Coverage by Phase

| # | Use Case | Ph1 | Ph2 | Ph3 | Ph4 |
|---|----------|:---:|:---:|:---:|:---:|
| 1 | Schedule Sync/Optimization | Prompts + Power BI | Power App | Automated alerts | Predictive |
| 2 | Instructor Qual & Performance | Power BI dashboard | Dataverse tracking | MCTIMS feed | Trend analysis |
| 3 | Training Mgmt & Assessment | GenAI prompts | BARS forms | Auto grading | Predictive |
| 4 | Curriculum Development | Adapted prompts | RAG over POIs | Gap analysis | AI scenario banks |
| 5 | Academic Support | Tutoring prompts | Study plans | Adaptive recs | Full AI tutor |
| 6 | Human Performance Analysis | Manual PT data | Collection forms | HPB integration | Injury prediction |
| 7 | Scenario Enhancement | GenAI prompts | RAG scenarios | Calibrated | Real-time engine |
| 8 | Data-Driven Feedback | Counseling prompts | Auto packets | Longitudinal | Predictive flags |
| 9 | AAR Synthesis | AAR analysis prompts | Structured intake | Cross-AAR trends | Auto routing |
| 10 | Logistics & Sustainment | Power BI readiness | Supply tracker app | Integrated data | Predictive |
| 11 | Forecast Maintenance/Supply | Manual tracking | Equipment app | Predictive alerts | Full optimization |
| 12 | Process Improvement | Process prompts | Mapping tools | Bottleneck ID | CI dashboard |
| 13 | Barracks & Facilities | Allocation visual | Booking app | Optimization | Predictive |
| 14 | Predictive Risk Assessment | Risk registers | Structured intake | Pattern analysis | Predictive models |
| 15 | Governance & Ethics | SOP + PIA template | RBAC + audit | Dashboard | Continuous |

---

## Cost Summary

| Phase | Duration | Cost | Cumulative |
|-------|----------|------|------------|
| Phase 1 | Months 0–3 | $0–$1,200 | $0–$1,200 |
| Phase 2 | Months 3–6 | $3,600–$4,650 | $3,600–$5,850 |
| Phase 3 | Months 6–12 | $32,600–$49,800 | $36,200–$55,650 |
| Phase 4 | Months 12–18 | $41,200–$59,200 | $77,400–$114,850 |

Primary cost drivers in Phase 3–4: Azure OpenAI API consumption, Power Platform premium licensing, ATO documentation support. Labor is organic (no contractor cost).

---

## Risk Summary

| Risk | Likelihood | Impact | Mitigation |
|------|:----------:|:------:|------------|
| MCTIMS access denied/delayed | High | High | System works without it; manual fallback |
| ATO takes longer than planned | Medium | High | Phase 1 maintains momentum; IATT bridges gap |
| MHS GENESIS PIA denied | High | Medium | Phase 4 designed to work without it |
| Azure Gov access blocked | Medium | High | Coordinate early; CamoGPT API as fallback |
| SPC resistance to BARS | Medium | Medium | Involve SPCs in design; AI narrative is the incentive |
| Jesse PCS/transfers | Medium | Critical | Continuous documentation; train replacements early |

---

## What I Provide vs What TBS Provides

| Jesse Morgan Provides | TBS Must Provide |
|----------------------|------------------|
| EDD training curriculum (5 courses, ready to deliver) | Access to instructors and leadership for training |
| 20 TBS-adapted prompt templates | SME review of prompt relevance/accuracy |
| Power BI dashboard design and development | M365 licenses, SharePoint permissions |
| Power App development (Phase 2+) | Azure Gov subscription, IATT coordination |
| Azure OpenAI RAG pipeline (Phase 2+) | POI documents and doctrine for knowledge base |
| BARS framework adaptation | SME input on behavioral indicators |
| All governance documentation | S-6/ISSM ATO coordination |
| Full documentation and knowledge transfer | 2–3 designated organic maintainers |

---

## Comparable Programs

| Program | Cost | Timeline | How Heywood Compares |
|---------|------|----------|---------------------|
| West Point Azimuth (cadet management modernization) | $2.18M | 3 years | Comparable scope at 3–5% cost |
| Army ATIS (training management, 2M users) | Enterprise | Multi-year | Heywood is schoolhouse-scoped, months not years |
| CGSC AI Wargaming | $0 | 1 week | Phase 1 follows same rapid-prototyping approach |
| GenAI4C SBIR (AI curriculum tools, 5 companies) | Phase I awards | TBD | Complementary — curriculum creation vs student assessment |

---

## Immediate Next Steps

1. **Brief CO** on this plan and the TBS AI Vision 2026 evaluation
2. **Identify TBS POC** willing to pilot with one company
3. **Deliver Phase 1 training** (AI Fluency + Supervisor Orientation) — requires only 2.5 hours of TBS time
4. **Deploy Phase 1 dashboards** on existing MCEN infrastructure
5. **Begin IATT coordination** with S-6 (long lead time — start early)
6. **Begin MCTIMS coordination** with TECOM MISSO (long lead time — start early)

---

**Point of Contact:** SSgt Jesse C. Morgan
