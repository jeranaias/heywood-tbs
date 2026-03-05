# Tool Registry Entry

*The Heywood Initiative — Component Registry*

---

## Component Identification

| Field              | Entry |
|--------------------|-------|
| **Component Name** | Student Performance Dashboard |
| **Version**        | 1.0 |
| **Heywood Phase**  | Phase 1 |
| **Developer**      | SSgt Jesse C. Morgan |
| **Deployment Date**| Pending — estimated March–April 2026 |

---

## Description

| Field              | Entry |
|--------------------|-------|
| **Platform**       | Power BI |
| **Heywood Use Case** | #2 (Instructor Qual & Performance), #3 (Training Mgmt & Assessment), #6 (Human Performance Analysis) |

*Brief description (2-3 sentences):*

> Power BI dashboard providing a three-pillar view of student performance across Academics (32%), Military Skills (32%), and Leadership (36%), mirroring the TBS grading policy. Features four pages: Company Overview, Individual Student Detail (drill-through), At-Risk Students, and Phase Progression. Row-Level Security enforces three access tiers — Staff sees all students, SPCs see their company only, Students see their own individual record only.

---

## Users

| Field | Entry |
|-------|-------|
| **Roles Served** | Staff, SPC, Student |
| **Current User Count** | Estimated 6–10 Staff, 4–6 SPCs (pilot company), ~200 Students (pilot company) |
| **Companies Using** | Pilot: one company (TBD). Designed for expansion to all 8+ TBS companies. |

---

## Data Handling

| Field | Entry |
|-------|-------|
| **Contains PII?** | Yes — student names, EDIPIs, individual grades/scores, class standing, peer evaluation scores |
| **Contains PHI?** | No |
| **Data Classification** | CUI |
| **PIA Status** | Threshold complete — documented in `governance/pia/pia-student-performance.md`. Full PIA not required for Phase 1. Pending Privacy Officer concurrence. |
| **RBAC Implemented?** | Yes (design complete, implementation pending deployment) — Staff: no filter (all data); SPC: company-filtered via UserSecurity table + USERPRINCIPALNAME(); Student: EDIPI-filtered to own record only. Three RLS roles defined with DAX filter expressions. |

---

## Documentation and Support

| Field | Entry |
|-------|-------|
| **Documentation Location** | `dashboards/student-performance/spec.md` (Power BI specification with DAX measures, visual layout, RLS configuration, deployment notes) |
| **Current Maintainer** | SSgt Jesse C. Morgan |
| **Maintainer Contact** | jesse.morgan@usmc.mil |
| **Backup Maintainer** | TBD — to be identified from TBS organic staff during Phase 1 knowledge transfer |

---

## Status

| Field | Entry |
|-------|-------|
| **Current Status** | Under Revision (pre-deployment — specification complete, awaiting MCEN provisioning) |
| **Authorization** | Inherited (MCEN) — Power BI Service within existing MCEN M365 IL5 boundary |
| **Last Review** | March 2026 (initial specification and compliance review) |
| **Next Review** | Post-deployment review scheduled for 30 days after go-live |

---

## Dependencies

*List other Heywood components this depends on:*

| Dependency | Type | Notes |
|------------|------|-------|
| SharePoint Data Layer (StudentScores list) | Data source | Primary data source via Power BI SharePoint Online connector. Import mode with 4-hour scheduled refresh during duty hours. Dashboard cannot function without this list being provisioned and populated. |
| SharePoint Data Layer (UserSecurity table) | Auth | UserSecurity list or Excel table mapping Azure AD email to Role + Company + EDIPI. Required for RLS enforcement. Must be maintained as users change roles or companies. |
| Azure AD Security Groups | Auth | CAC authentication and workspace access control. Groups: TBS-Staff, TBS-SPC-[Company], TBS-Student-[Company]. Managed by MCEN Azure AD administrators. |

---

## Change Log

| Date | Version | Changed By | Description of Change |
|------|---------|------------|-----------------------|
| March 2026 | 1.0 | SSgt Jesse C. Morgan | Initial specification, DAX measures, visual layout, RLS design, compliance documentation |
| | | | |

---

*Register every deployed Heywood component. Update when status, maintainer, or version changes.*
