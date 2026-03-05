# Compliance Checklist

*The Heywood Initiative — Pre-Deployment Compliance Verification*

---

## Component Information

| Field              | Entry |
|--------------------|-------|
| **Component Name** | Power BI Dashboards (Student Performance Dashboard + Instructor Qualification Dashboard) |
| **Heywood Phase**  | Phase 1 |
| **Developer**      | SSgt Jesse C. Morgan |
| **Date**           | March 2026 |
| **Reviewer**       | Pending — TBS S-6 / ISSM |

---

## Compliance Items

Complete each item before the component proceeds to deployment.

### Data & Privacy

| # | Requirement | Status | Notes |
|---|-------------|--------|-------|
| 1 | Data classification completed (CUI, PII, PHI determination) | Complete | Both dashboards display CUI data containing PII (student names, EDIPIs, individual scores; instructor names, EDIPIs, qualification records). No PHI in either dashboard. Student Performance Dashboard: student PII visible to Staff/SPC/Student (role-filtered). Instructor Quals Dashboard: instructor PII visible to Staff only. Classification documented in corresponding PIA threshold analyses. |
| 2 | PIA threshold analysis completed (use `governance/pia/pia-threshold-analysis.md`) | Complete | Two PIA threshold analyses completed: (1) `pia-student-performance.md` covering Student Performance Dashboard; (2) `pia-instructor-data.md` covering Instructor Quals Dashboard. Both conclude threshold analysis is sufficient — full PIA not required for Phase 1. |
| 3 | TBS Privacy Officer consulted (if PII is involved) | Pending | PIA threshold analyses prepared and awaiting TBS Privacy Officer review and concurrence. Dashboards must not be published with live PII data until Privacy Officer concurrence is received. Testing with sample/synthetic data is permissible. |
| 4 | SORN requirements reviewed (if PII is retrieved by identifier) | Complete | Student self-service view retrieves records by EDIPI (via RLS filter matching USERPRINCIPALNAME to EDIPI). Covered under existing DoD SORN for personnel records maintained in authorized M365 systems. No new SORN required. Instructor dashboard does not provide self-service retrieval — staff queries only. |
| 5 | Data minimization applied (aggregate where possible, anonymize where required) | Complete | Dashboards provide aggregate views (company averages, phase distribution, qualification coverage rates) alongside individual records. Individual records are necessary for the dashboards' core functions (SPC counseling, qualification management). No extraneous PII fields displayed — dashboard visuals show only fields required for the specific analytical purpose. EDIPI used as filter key but not prominently displayed in visuals. |
| 6 | AI data handling verified (no PII/PHI in prompts unless on IL5 platform) | N/A | Phase 1 Power BI dashboards do not integrate with any AI/LLM service. All calculations are DAX measures executed within Power BI Service. No data is sent to Azure OpenAI, GenAI.mil, or any external AI endpoint. Smart narratives and Q&A features (Power BI AI visuals) are disabled for Phase 1. AI integration deferred to Phase 2+. |

### Authorization & Security

| # | Requirement | Status | Notes |
|---|-------------|--------|-------|
| 7 | ATO/IATT status verified for this component | Complete | Phase 1 operates entirely within the existing MCEN M365 authorization boundary. Power BI Service is an authorized MCEN service — no separate ATO or IATT required. Power BI reports connect to SharePoint Online (also within MCEN boundary) via native connector. No custom code, no external API calls, no data egress. |
| 8 | Platform authorization confirmed (Azure Gov IL5, GenAI.mil, MCEN M365) | Complete | Power BI Service is authorized on MCEN at IL5. FedRAMP High + DoD IL5 Provisional Authorization (PA) covers Power BI within Azure Government. Dashboards will be published to Power BI Service workspaces within the MCEN tenant. |
| 9 | Access controls implemented (CAC auth, RBAC for Staff/SPC/Student) | In Progress | Design complete. CAC authentication inherited from MCEN Azure AD — all Power BI Service access requires CAC login. RBAC implementation plan: (1) Two Power BI workspaces — staff workspace (both dashboards) and shared workspace (student performance only); (2) Power BI App for end-user distribution with role-based audience; (3) Azure AD security groups: TBS-Staff, TBS-SPC-[Company], TBS-Student-[Company]. Cannot be fully configured until workspaces are created in Power BI Service. |
| 10 | Row-Level Security configured and tested (each role sees only authorized data) | In Progress | RLS design complete and documented in dashboard specifications. Student Performance Dashboard: three RLS roles (Staff = no filter, SPC = company filter via UserSecurity table, Student = EDIPI filter via UserSecurity table). Instructor Quals Dashboard: two RLS roles (Staff = no filter, CompanyCommander = company filter). UserSecurity table schema defined. DAX filter expressions written and documented. Cannot be tested until dashboards are published to Power BI Service with live or representative data. "View as Role" testing plan documented. |
| 11 | Audit logging enabled | In Progress | Power BI Service activity logging is enabled tenant-wide by MCEN administrators. Power BI audit events include: report views, data exports, sharing changes, RLS role assignments, refresh history. Logs flow to M365 Unified Audit Log. Verification that audit events are being captured for Heywood workspaces will be confirmed post-deployment. |
| 12 | Data encryption verified (at rest: CMK for IL5; in transit: TLS 1.2+) | Complete | Inherited from MCEN Power BI Service IL5 authorization. Power BI on Azure Government provides: encryption at rest via Microsoft-managed keys (FIPS 140-2 validated) for datasets, reports, and cached data. Encryption in transit via TLS 1.2+ for all browser and API connections. Import mode datasets are encrypted at rest in Power BI Service storage. No additional encryption configuration required at the workspace level. |

### Operational

| # | Requirement | Status | Notes |
|---|-------------|--------|-------|
| 13 | OPSEC review completed (no sensitive indicators in AI outputs) | N/A | Phase 1 dashboards do not generate AI outputs. All visuals display structured data (scores, qualifications, schedules) entered by TBS staff. No natural language generation, no AI-generated narratives, no smart narratives enabled. Dashboard content is training performance data — no operational plans, movement data, or intelligence. |
| 14 | Records management requirements identified (retention schedule) | In Progress | Power BI reports and datasets are considered derived products — the authoritative records are the source SharePoint lists. Power BI Service retains published reports indefinitely until deleted. Dataset refresh history retained for 60 refreshes. Retention approach: Power BI reports are maintained as long as the underlying data exists; archived when a training class completes. Formal retention SOP to be documented before go-live. |
| 15 | Incident response procedures documented | In Progress | For Power BI-specific incidents: (1) dashboard unavailable — verify Power BI Service status via MCEN Service Health, contact MCEN Help Desk; (2) RLS bypass or unauthorized data access — immediately revoke workspace access, report to MCEN Help Desk + TBS S-6, review audit logs; (3) data refresh failure — check SharePoint connector credentials, verify scheduled refresh in Power BI Service; (4) incorrect data displayed — verify DAX measures, compare to source SharePoint data, correct and republish. Formal SOP to be finalized before go-live. |
| 16 | Input validation implemented (prevent prompt injection, data corruption) | N/A | Power BI dashboards are read-only visualization tools — they do not accept user input beyond slicer selections and drill-through navigation. No free-text input fields. No user-submitted data. All data originates from validated SharePoint lists (input validation addressed in SharePoint compliance checklist). Slicers use predefined values from dimension tables. |
| 17 | Component registered in Tool Registry (`governance/registry/`) | Complete | Two registry entries created: (1) `governance/registry/student-performance-dashboard.md` for Student Performance Dashboard; (2) `governance/registry/instructor-quals-dashboard.md` for Instructor Qualification Dashboard. |
| 18 | User documentation created | In Progress | Dashboard specifications documented in `dashboards/student-performance/spec.md` and `dashboards/instructor-quals/spec.md`. End-user guides needed: (1) SPC guide — how to navigate the Student Performance Dashboard, use slicers, drill through to student detail, interpret composites; (2) Staff guide — how to use both dashboards, manage RLS assignments, interpret alerts; (3) Student guide — how to access individual performance view. To be created before pilot deployment. |
| 19 | Backup/recovery plan in place (SharePoint versioning, Dataverse backup) | Complete | Power BI reports (.pbix files) maintained in version control alongside project files. If a published report is corrupted or accidentally deleted, it can be republished from the .pbix source file. Datasets can be re-imported from source SharePoint lists (no data loss since Power BI is the visualization layer, not the data store). Power BI Service workspace recovery available via MCEN support. Recovery time objective: same business day for republishing from .pbix backup. |

---

## Approval

- [ ] All items complete — component may proceed to deployment
- [x] Items outstanding — component must not deploy until resolved

**Outstanding items (if any):**
1. Item 3: TBS Privacy Officer concurrence on PIA threshold analyses (required before publishing with live PII data)
2. Item 9: Workspace creation and Azure AD security group configuration pending MCEN access
3. Item 10: RLS testing with "View as Role" pending dashboard publication to Power BI Service
4. Item 11: Audit logging verification for Heywood workspaces pending deployment
5. Item 14: Formal retention SOP to be documented
6. Item 15: Formal incident response SOP to be finalized
7. Item 18: End-user guides (SPC, Staff, Student) to be created

**Note:** Items 2, 9, 10, and 11 are blocked on Power BI Service workspace provisioning and SharePoint data availability. These are designed and documented but cannot be verified until the infrastructure exists. Items 14, 15, and 18 are documentation tasks that will be completed during the deployment sprint.

**Reviewer Signature:** _____________________ **Date:** _____________

**TBS S-6 / ISSM Concurrence (Phase 2+):** _____________________ **Date:** _____________
