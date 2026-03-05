# Tool Registry Entry

*The Heywood Initiative — Component Registry*

---

## Component Identification

| Field              | Entry |
|--------------------|-------|
| **Component Name** | Instructor Qualification Dashboard |
| **Version**        | 1.0 |
| **Heywood Phase**  | Phase 1 |
| **Developer**      | SSgt Jesse C. Morgan |
| **Deployment Date**| Pending — estimated March–April 2026 |

---

## Description

| Field              | Entry |
|--------------------|-------|
| **Platform**       | Power BI |
| **Heywood Use Case** | #2 (Instructor Qual & Performance) |

*Brief description (2-3 sentences):*

> Power BI dashboard tracking instructor qualification status across TBS with 30/60/90-day expiration alerts, workload distribution analysis, and rotation risk planning. Features four pages: Qualification Status Overview, Instructor Workload, Coverage Gaps & Rotation Risk, and Individual Instructor Detail (drill-through). Staff-only access — ensures S-3 and company leadership can proactively manage qualification currency and identify coverage gaps before they impact training execution.

---

## Users

| Field | Entry |
|-------|-------|
| **Roles Served** | Staff |
| **Current User Count** | Estimated 6–10 Staff (S-3 section, Operations Officer, Company Commanders) |
| **Companies Using** | All companies — instructor qualification management is a school-wide staff function, not company-scoped |

---

## Data Handling

| Field | Entry |
|-------|-------|
| **Contains PII?** | Yes — instructor names, EDIPIs, duty contact information, rank, MOS, PRD, individual qualification records |
| **Contains PHI?** | No |
| **Data Classification** | CUI |
| **PIA Status** | Threshold complete — documented in `governance/pia/pia-instructor-data.md`. Full PIA not required for Phase 1. Pending Privacy Officer concurrence. |
| **RBAC Implemented?** | Yes (design complete, implementation pending deployment) — Staff: no filter (all instructor data across all companies); CompanyCommander: optional company-filtered role for commanders who need only their own company's instructors. Published to staff-only Power BI workspace — SPC and Student roles have no access. |

---

## Documentation and Support

| Field | Entry |
|-------|-------|
| **Documentation Location** | `dashboards/instructor-quals/spec.md` (Power BI specification with DAX measures, visual layout, RLS configuration, alert configuration, deployment notes) |
| **Current Maintainer** | SSgt Jesse C. Morgan |
| **Maintainer Contact** | jesse.morgan@usmc.mil |
| **Backup Maintainer** | TBD — to be identified from TBS S-3 section during Phase 1 knowledge transfer |

---

## Status

| Field | Entry |
|-------|-------|
| **Current Status** | Under Revision (pre-deployment — specification complete, awaiting MCEN provisioning and SharePoint data layer) |
| **Authorization** | Inherited (MCEN) — Power BI Service within existing MCEN M365 IL5 boundary |
| **Last Review** | March 2026 (initial specification and compliance review) |
| **Next Review** | Post-deployment review scheduled for 30 days after go-live |

---

## Dependencies

*List other Heywood components this depends on:*

| Dependency | Type | Notes |
|------------|------|-------|
| SharePoint Data Layer (Instructors list) | Data source | Instructor roster with PII, company assignments, workload metrics. Must be provisioned and populated before dashboard connects. |
| SharePoint Data Layer (QualificationRecords list) | Data source | Individual qualification records with earned/expiration dates. Links to Instructors via InstructorEDIPI. |
| SharePoint Data Layer (RequiredQualifications list) | Data source | Reference table defining required qualifications, categories, and minimum-per-event thresholds. No PII — populate first as reference data. |
| SharePoint Data Layer (UserSecurity table) | Auth | Maps Azure AD email to Role + Company for RLS enforcement. Required for CompanyCommander role filtering. |
| Azure AD Security Groups | Auth | CAC authentication and staff-only workspace access control. Group: TBS-Staff. Managed by MCEN Azure AD administrators. |

---

## Change Log

| Date | Version | Changed By | Description of Change |
|------|---------|------------|-----------------------|
| March 2026 | 1.0 | SSgt Jesse C. Morgan | Initial specification, DAX measures (expiration tracking, coverage analysis, workload, rotation planning), visual layout, RLS design, alert configuration, compliance documentation |
| | | | |

---

*Register every deployed Heywood component. Update when status, maintainer, or version changes.*
