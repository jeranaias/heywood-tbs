# Tool Registry Entry

*The Heywood Initiative — Component Registry*

---

## Component Identification

| Field              | Entry |
|--------------------|-------|
| **Component Name** | TBS Prompt Playbook |
| **Version**        | 1.0 |
| **Heywood Phase**  | Phase 1 |
| **Developer**      | SSgt Jesse C. Morgan |
| **Deployment Date**| March 2026 |

---

## Description

| Field              | Entry |
|--------------------|-------|
| **Platform**       | GenAI.mil |
| **Heywood Use Case** | #1 (Schedule Sync/Optimization), #3 (Training Mgmt & Assessment), #4 (Curriculum Development), #5 (Academic Support), #7 (Scenario Enhancement), #8 (Data-Driven Feedback), #9 (AAR Synthesis), #10 (Logistics & Sustainment), #12 (Process Improvement), #13 (Barracks & Facilities) |

*Brief description (2-3 sentences):*

> Collection of 20 structured prompts adapted from the Expert-Driven Development curriculum for TBS-specific workflows, designed for use on GenAI.mil (IL5, CUI-authorized). Prompts cover AAR analysis, counseling preparation, scenario generation, data schema design, dashboard specification, and governance screening. Each prompt includes bracketed placeholders for TBS-specific context, enabling SPCs and Staff to leverage AI for routine analytical and administrative tasks without AI expertise.

---

## Users

| Field | Entry |
|-------|-------|
| **Roles Served** | Staff, SPC |
| **Current User Count** | Estimated 4–6 SPCs (pilot company) + 2–4 Staff (S-3 section, Operations) |
| **Companies Using** | Pilot: one company (TBD). Prompts are company-agnostic — any TBS user with GenAI.mil access can use them. |

---

## Data Handling

| Field | Entry |
|-------|-------|
| **Contains PII?** | No — prompts themselves contain no PII. Users are instructed to anonymize data before pasting into GenAI.mil: use duty positions and roles instead of names, describe data types and field names instead of pasting actual records. PII is authorized on GenAI.mil (IL5) but playbook guidance encourages anonymization as a best practice. |
| **Contains PHI?** | No — prompts do not reference or request medical/health data. Playbook README explicitly states: "PHI is never authorized on GenAI.mil." |
| **Data Classification** | Unclassified (prompt templates themselves). CUI is authorized on GenAI.mil if users paste CUI data into prompts — this is permissible per GenAI.mil usage policy. |
| **PIA Status** | Not required — prompts contain no PII. The GenAI.mil platform has its own PIA and authorization. Users are responsible for data handling within GenAI.mil per their own PIA obligations. |
| **RBAC Implemented?** | N/A — prompts are static markdown files distributed to authorized users. GenAI.mil access is controlled by GenAI.mil's own CAC authentication (1.1M+ authorized DoD users). No Heywood-specific access control needed — prompts are unclassified reference material. |

---

## Documentation and Support

| Field | Entry |
|-------|-------|
| **Documentation Location** | `prompts/playbook/README.md` (index, usage instructions, data handling rules); `prompts/playbook/01-aar-analysis.md` through `prompts/playbook/20-feedback-form.md` (individual prompt files) |
| **Current Maintainer** | SSgt Jesse C. Morgan |
| **Maintainer Contact** | jesse.morgan@usmc.mil |
| **Backup Maintainer** | TBD — any TBS user trained via EDD AI Fluency workshop can maintain and adapt prompts |

---

## Status

| Field | Entry |
|-------|-------|
| **Current Status** | Active (prompts are complete and ready for use — no infrastructure deployment required) |
| **Authorization** | Inherited (MCEN) — GenAI.mil is an authorized IL5 platform with 1.1M+ DoD users. No separate authorization required for prompt usage. |
| **Last Review** | March 2026 (initial prompt development and adaptation for TBS workflows) |
| **Next Review** | June 2026 (post-pilot review — assess which prompts are most used, refine based on SPC feedback, add new prompts for emerging use cases) |

---

## Dependencies

*List other Heywood components this depends on:*

| Dependency | Type | Notes |
|------------|------|-------|
| GenAI.mil | Other | Platform dependency — prompts are designed for GenAI.mil's interface and capabilities (IL5, CUI-authorized, GPT-4-class models). Prompts would work on other LLM platforms but data handling rules may differ. |
| EDD AI Fluency Training | Other | Users should complete the 2-hour AI Fluency workshop before using prompts independently. Training covers prompt engineering fundamentals, critical evaluation of AI outputs, and TBS-specific data handling rules. Not a technical dependency — prompts function without training, but effectiveness depends on user understanding. |

---

## Change Log

| Date | Version | Changed By | Description of Change |
|------|---------|------------|-----------------------|
| March 2026 | 1.0 | SSgt Jesse C. Morgan | Initial playbook: 20 prompts adapted from EDD curriculum for TBS workflows. Categories: AAR analysis, counseling prep, user stories, MVP definition, AI suitability, PIA screening, data classification, schema design, data modeling, dashboards, training schedule, company metrics, SITREP generation, instructor quals, AAR synthesis, scenario generation, peer evaluation, equipment tracking, range booking, event feedback. |
| | | | |

---

### Prompt Index

| # | Prompt | File | Primary Use Case |
|---|--------|------|-----------------|
| 1 | AAR Analysis & Trend Extraction | `01-aar-analysis.md` | #9 (AAR Synthesis) |
| 2 | Counseling Session Prep | `02-counseling-prep.md` | #8 (Data-Driven Feedback) |
| 3 | Student Performance User Stories | `03-user-stories.md` | #3 (Training Mgmt) |
| 4 | Training Aid MVP | `04-mvp-definition.md` | #4 (Curriculum Development) |
| 5 | AI Suitability Check | `05-frontier-check.md` | #12 (Process Improvement) |
| 6 | PIA Screening for Student Data | `06-pia-threshold.md` | #3 (Training Mgmt) |
| 7 | TBS Data Classification | `07-data-classification.md` | #3 (Training Mgmt) |
| 8 | Student Data Schema | `08-data-schema.md` | #3 (Training Mgmt) |
| 9 | Performance Tracking Data Model | `09-data-model.md` | #3 (Training Mgmt) |
| 10 | Student Performance Dashboard | `10-student-dashboard.md` | #3 (Training Mgmt) |
| 11 | Training Schedule Board | `11-training-schedule.md` | #1 (Schedule Sync) |
| 12 | Company Metrics Summary | `12-company-metrics.md` | #12 (Process Improvement) |
| 13 | Weekly SITREP Generator | `13-sitrep-generator.md` | #10 (Logistics & Sustainment) |
| 14 | Instructor Qualification Tracker | `14-instructor-quals.md` | #3 (Training Mgmt) |
| 15 | AAR Synthesis & Trends | `15-aar-synthesis.md` | #9 (AAR Synthesis) |
| 16 | Tactical Scenario Generator | `16-scenario-generator.md` | #7 (Scenario Enhancement) |
| 17 | Peer Evaluation Structuring | `17-peer-eval.md` | #8 (Data-Driven Feedback) |
| 18 | Equipment/Supply Tracker | `18-equipment-tracker.md` | #10 (Logistics & Sustainment) |
| 19 | Range/Classroom Booking | `19-range-booking.md` | #13 (Barracks & Facilities) |
| 20 | Post-Event Feedback Collection | `20-feedback-form.md` | #5 (Academic Support) |

---

*Register every deployed Heywood component. Update when status, maintainer, or version changes.*
