# Compliance Checklist

*The Heywood Initiative — Pre-Deployment Compliance Verification*

---

## Component Information

| Field              | Entry |
|--------------------|-------|
| **Component Name** | SharePoint Data Layer (6 lists: StudentScores, TrainingSchedule, Instructors, RequiredQualifications, QualificationRecords, EventFeedback) |
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
| 1 | Data classification completed (CUI, PII, PHI determination) | Complete | Four of six lists contain PII (StudentScores, Instructors, QualificationRecords, EventFeedback when student-attributed). Two lists contain no PII (TrainingSchedule, RequiredQualifications). All PII data classified as CUI. No PHI is collected in any Phase 1 list. Classification documented in PIA threshold analyses. |
| 2 | PIA threshold analysis completed (use `governance/pia/pia-threshold-analysis.md`) | Complete | Two PIA threshold analyses completed: (1) `pia-student-performance.md` covering StudentScores list and Student Performance Dashboard; (2) `pia-instructor-data.md` covering Instructors and QualificationRecords lists and Instructor Quals Dashboard. Both conclude threshold analysis is sufficient — full PIA not required for Phase 1. |
| 3 | TBS Privacy Officer consulted (if PII is involved) | Pending | PIA threshold analyses prepared and awaiting TBS Privacy Officer review and concurrence. Both analyses require Privacy Officer signature before PII data entry begins. Coordination initiated but concurrence not yet received. |
| 4 | SORN requirements reviewed (if PII is retrieved by identifier) | Complete | StudentScores records are retrievable by EDIPI (student self-service view via Power BI RLS). Instructor records are retrievable by EDIPI (staff queries). Applicable SORN: DMDC 02 DoD (Defense Manpower Data Center, DoD-wide personnel records) covers EDIPI-indexed records within DoD systems. No new SORN required — data is maintained within existing M365/SharePoint systems already covered by the M365 Government SORN. |
| 5 | Data minimization applied (aggregate where possible, anonymize where required) | Complete | Each list collects only fields necessary for its stated purpose. No SSNs, personal contact info for students, or medical data. Instructor contact info limited to duty phone and .mil email. Peer evaluation scores store only aggregated scores per student — individual evaluator identities are not recorded. EventFeedback supports anonymous submission. Documented in PIA data minimization tables. |
| 6 | AI data handling verified (no PII/PHI in prompts unless on IL5 platform) | N/A | Phase 1 SharePoint lists do not interact with any AI/LLM service. Data flows are: manual entry by SPCs/Staff into SharePoint, then SharePoint connector to Power BI. No Power Automate flows, no Azure OpenAI integration, no AI processing of list data in Phase 1. AI integration deferred to Phase 2+ with separate compliance review. |

### Authorization & Security

| # | Requirement | Status | Notes |
|---|-------------|--------|-------|
| 7 | ATO/IATT status verified for this component | Complete | Phase 1 operates entirely within the existing MCEN M365 authorization boundary. SharePoint Online is an authorized MCEN service — no separate ATO or IATT required. Confirmed per MCEN M365 IL5 authorization documentation. Phase 2+ components requiring custom code or external API connections will require IATT. |
| 8 | Platform authorization confirmed (Azure Gov IL5, GenAI.mil, MCEN M365) | Complete | SharePoint Online is authorized on MCEN at IL5. FedRAMP High + DoD IL5 Provisional Authorization (PA) covers SharePoint Online within Azure Government. TBS site collection will be created within the existing MCEN SharePoint tenant. |
| 9 | Access controls implemented (CAC auth, RBAC for Staff/SPC/Student) | In Progress | Design complete: SharePoint site permissions will use Azure AD security groups mapped to TBS roles (Staff, SPC, Student). CAC authentication inherited from MCEN Azure AD (all MCEN users authenticate via CAC). SharePoint site-level permissions: Staff = Full Control, SPC = Contribute (own company lists), Student = Read (EventFeedback = Contribute for write-only submission). Cannot be fully implemented until SharePoint site is provisioned on MCEN. |
| 10 | Row-Level Security configured and tested (each role sees only authorized data) | In Progress | Design complete: RLS is enforced at the Power BI layer, not at the SharePoint layer. SharePoint views will provide UI-level filtering (company-filtered views for SPCs) but true data security is via Power BI RLS roles (Staff = all data, SPC = company-filtered, Student = own EDIPI only). UserSecurity table design documented in dashboard specs. Cannot be configured and tested until Power BI is connected to live SharePoint data. |
| 11 | Audit logging enabled | In Progress | MCEN M365 Unified Audit Log is enabled tenant-wide by MCEN administrators. SharePoint audit events (item access, modification, deletion) are captured automatically for all SharePoint sites within the tenant. Verification that audit logging is active for the specific Heywood site will be confirmed post-provisioning. Retention period governed by MCEN tenant policy (minimum 90 days). |
| 12 | Data encryption verified (at rest: CMK for IL5; in transit: TLS 1.2+) | Complete | Inherited from MCEN M365 IL5 authorization. SharePoint Online on Azure Government provides: encryption at rest via Microsoft-managed keys (FIPS 140-2 validated), encryption in transit via TLS 1.2+. Customer-managed keys (CMK) available at tenant level if required by MCEN policy but not required for Phase 1 SharePoint lists. BitLocker volume-level encryption applied to all Azure Government storage. |

### Operational

| # | Requirement | Status | Notes |
|---|-------------|--------|-------|
| 13 | OPSEC review completed (no sensitive indicators in AI outputs) | N/A | Phase 1 SharePoint lists do not produce AI outputs. Lists contain training performance data (grades, qualifications, schedules) — no operational plans, movement data, or intelligence information. Training schedule data is administrative, not operational. OPSEC review will be required for Phase 2+ when AI processing generates outputs from this data. |
| 14 | Records management requirements identified (retention schedule) | In Progress | SharePoint list items are Marine Corps records subject to records management requirements. Applicable retention schedule: SECNAV M-5210.1, Chapter 4 (Training Records). Student performance records: retain for duration of training cycle + 2 years. Instructor qualification records: retain for duration of assignment + 2 years. SharePoint versioning and recycle bin provide short-term recovery. Long-term retention policy to be configured in M365 Compliance Center post-provisioning. |
| 15 | Incident response procedures documented | In Progress | MCEN provides enterprise incident response for M365 services (MCEN Help Desk + MCEN CSSP). For Heywood-specific incidents: (1) data integrity issues — SPCs report to S-3, developer restores from SharePoint version history; (2) unauthorized access — report to MCEN Help Desk + TBS S-6; (3) data loss — restore from SharePoint recycle bin (93-day retention) or site collection backup. Formal Heywood incident response SOP to be finalized before go-live. |
| 16 | Input validation implemented (prevent prompt injection, data corruption) | In Progress | SharePoint column validation rules designed: EDIPI = 10-digit number format, scores = 0-100 range, dates = valid date format, choice columns = predefined values only (e.g., Company = Alpha/Bravo/Charlie/etc.). No free-text fields that could be vectors for injection. Column-level validation rules will be configured during list creation. Cannot be tested until lists are provisioned. |
| 17 | Component registered in Tool Registry (`governance/registry/`) | Complete | Registered as `governance/registry/sharepoint-data-layer.md`. All six lists documented as a single registered component with data handling, dependencies, and maintenance information. |
| 18 | User documentation created | In Progress | SharePoint list schemas documented in `schemas/sharepoint/` with field definitions, relationships, and deployment notes. SPC data entry guide (how to enter student scores, what each field means) to be created before pilot deployment. Staff administration guide (how to manage list permissions, add/remove users) to be created before go-live. |
| 19 | Backup/recovery plan in place (SharePoint versioning, Dataverse backup) | In Progress | SharePoint Online provides: (1) version history — enabled on all lists, retains last 500 versions; (2) recycle bin — first-stage 93 days, second-stage site collection recycle bin; (3) Microsoft backup — 14-day recovery window for site collection restores via MCEN support. For Phase 1 data volumes (~200 student records, ~80 instructor records), these native capabilities are sufficient. Will verify versioning is enabled post-provisioning. |

---

## Approval

- [ ] All items complete — component may proceed to deployment
- [x] Items outstanding — component must not deploy until resolved

**Outstanding items (if any):**
1. Item 3: TBS Privacy Officer concurrence on PIA threshold analyses (required before PII data entry)
2. Item 9: RBAC configuration pending SharePoint site provisioning on MCEN
3. Item 10: RLS testing pending Power BI connection to live data
4. Item 11: Audit logging verification pending site provisioning
5. Item 14: M365 retention policy configuration pending site provisioning
6. Item 15: Formal incident response SOP to be finalized
7. Item 16: Column validation rules to be configured during list creation
8. Item 18: SPC data entry guide and Staff admin guide to be written
9. Item 19: Versioning verification pending site provisioning

**Note:** Items 2, 9, 10, 11, 14, 16, and 19 are blocked on SharePoint site provisioning — they are designed and documented but cannot be implemented and verified until the lists exist on MCEN. These are expected to be completed during the deployment sprint, not before.

**Reviewer Signature:** _____________________ **Date:** _____________

**TBS S-6 / ISSM Concurrence (Phase 2+):** _____________________ **Date:** _____________
