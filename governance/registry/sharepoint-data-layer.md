# Tool Registry Entry

*The Heywood Initiative — Component Registry*

---

## Component Identification

| Field              | Entry |
|--------------------|-------|
| **Component Name** | SharePoint Data Layer (6 lists) |
| **Version**        | 1.0 |
| **Heywood Phase**  | Phase 1 |
| **Developer**      | SSgt Jesse C. Morgan |
| **Deployment Date**| Pending — estimated March 2026 |

---

## Description

| Field              | Entry |
|--------------------|-------|
| **Platform**       | SharePoint Online |
| **Heywood Use Case** | #2 (Instructor Qual & Performance), #3 (Training Mgmt & Assessment), #6 (Human Performance Analysis), #9 (AAR Synthesis) |

*Brief description (2-3 sentences):*

> Six SharePoint Online lists serving as the Phase 1 data layer for all Heywood dashboards and future Power App integration. Lists include: StudentScores (individual student performance across 3-pillar grading), TrainingSchedule (master training calendar), Instructors (instructor roster and workload), RequiredQualifications (reference table of TBS qualification requirements), QualificationRecords (individual instructor certification tracking), and EventFeedback (post-event structured feedback for AAR synthesis). Designed for migration to Dataverse in Phase 3 with minimal schema changes.

---

## Users

| Field | Entry |
|-------|-------|
| **Roles Served** | Staff, SPC, Student |
| **Current User Count** | Estimated 6–10 Staff (direct list access + Power BI), 4–6 SPCs (data entry + Power BI), ~200 Students (read-only via Power BI + EventFeedback write) |
| **Companies Using** | Pilot: one company (TBD) for StudentScores. All companies for Instructors, RequiredQualifications, QualificationRecords, TrainingSchedule. |

---

## Data Handling

| Field | Entry |
|-------|-------|
| **Contains PII?** | Yes — four of six lists contain PII. StudentScores: student names, EDIPIs, individual scores. Instructors: instructor names, EDIPIs, contact info. QualificationRecords: linked to instructor EDIPI. EventFeedback: optional student attribution. Two lists contain no PII: TrainingSchedule, RequiredQualifications. |
| **Contains PHI?** | No |
| **Data Classification** | CUI (for PII-containing lists); Unclassified (for TrainingSchedule, RequiredQualifications) |
| **PIA Status** | Threshold complete — two PIA threshold analyses cover all six lists: `governance/pia/pia-student-performance.md` (StudentScores) and `governance/pia/pia-instructor-data.md` (Instructors, QualificationRecords). TrainingSchedule, RequiredQualifications, and EventFeedback (when anonymous) do not require PIA. Pending Privacy Officer concurrence for PII-containing lists. |
| **RBAC Implemented?** | Yes (design complete, implementation pending deployment) — SharePoint site permissions via Azure AD security groups: Staff = Full Control, SPC = Contribute (company-scoped views), Student = Read (EventFeedback = Contribute for write-only submission). True row-level data security enforced at Power BI layer, not SharePoint. SharePoint views provide UI-level company filtering for SPCs. |

---

## Documentation and Support

| Field | Entry |
|-------|-------|
| **Documentation Location** | `schemas/sharepoint/README.md` (overview, relationships, RLS model, deployment notes); `schemas/sharepoint/*.json` (individual list schemas with field definitions and validation rules) |
| **Current Maintainer** | SSgt Jesse C. Morgan |
| **Maintainer Contact** | jesse.morgan@usmc.mil |
| **Backup Maintainer** | TBD — to be identified from TBS S-3 section or S-6 during Phase 1 knowledge transfer |

---

## Status

| Field | Entry |
|-------|-------|
| **Current Status** | Under Revision (pre-deployment — schemas complete, awaiting MCEN SharePoint site provisioning) |
| **Authorization** | Inherited (MCEN) — SharePoint Online within existing MCEN M365 IL5 boundary |
| **Last Review** | March 2026 (initial schema design, compliance review, PIA threshold analyses) |
| **Next Review** | Post-deployment review scheduled for 30 days after go-live; schema review before Phase 2 Dataverse migration |

---

## Dependencies

*List other Heywood components this depends on:*

| Dependency | Type | Notes |
|------------|------|-------|
| MCEN M365 Tenant | Auth | SharePoint site collection must be provisioned within the MCEN tenant. Requires site creation permissions from MCEN SharePoint administrator or TBS S-6. |
| Azure AD Security Groups | Auth | CAC authentication and RBAC enforcement. Groups: TBS-Staff, TBS-SPC-[Company], TBS-Student-[Company]. Must be created before SharePoint permissions can be configured. |
| Power BI Dashboards | Other | SharePoint lists are the data source for both Power BI dashboards. Lists must be provisioned and populated before dashboards can connect. Circular dependency — dashboards depend on lists, but lists are designed to feed dashboards. |

---

## Change Log

| Date | Version | Changed By | Description of Change |
|------|---------|------------|-----------------------|
| March 2026 | 1.0 | SSgt Jesse C. Morgan | Initial schema design for all 6 lists: StudentScores, TrainingSchedule, Instructors, RequiredQualifications, QualificationRecords, EventFeedback. Field definitions, validation rules, relationships, and deployment documentation. |
| | | | |

---

### Individual List Summary

| # | List Name | Records (est.) | PII | Primary Users | Power BI Dashboard |
|---|-----------|:--------------:|:---:|---------------|-------------------|
| 1 | StudentScores | ~200/company | Yes | SPC (data entry), Staff (read) | Student Performance Dashboard |
| 2 | TrainingSchedule | ~300/class | No | Staff (data entry), SPC/Student (read) | Future — Training Schedule Board |
| 3 | Instructors | ~80 | Yes | Staff/S-3 (data entry + read) | Instructor Qualification Dashboard |
| 4 | RequiredQualifications | ~30–50 | No | Staff/S-3 (reference data) | Instructor Qualification Dashboard |
| 5 | QualificationRecords | ~400+ | Yes | Staff/S-3 (data entry + read) | Instructor Qualification Dashboard |
| 6 | EventFeedback | Growing | Optional | Student (write), SPC/Staff (read) | Future — AAR Synthesis |

---

*Register every deployed Heywood component. Update when status, maintainer, or version changes.*
