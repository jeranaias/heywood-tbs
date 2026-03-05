# PIA Threshold Analysis

*The Heywood Initiative — Privacy Impact Assessment Threshold Checklist*

---

## Tool Information

| Field              | Entry |
|--------------------|-------|
| **Tool Name**      | Instructor Data (SharePoint Instructors + QualificationRecords lists + Power BI Instructor Quals Dashboard) |
| **Heywood Phase**  | Phase 1 |
| **Developer**      | SSgt Jesse C. Morgan |
| **Date**           | March 2026 |
| **Reviewer**       | Pending — TBS Privacy Officer |

---

## PII Determination

Answer each question below. If any answer is "Yes," a full Privacy Impact Assessment (PIA) is required per DoDI 5400.16.

| # | Question | Yes / No | Notes |
|---|----------|----------|-------|
| 1 | Does the tool collect, store, or process Personally Identifiable Information (PII)? | Yes | Instructor names, EDIPIs, rank, MOS, contact information (duty phone, duty email), company assignment, role, PRD (Projected Rotation Date), qualification records with individually-identifiable certification dates and expiration dates. All stored across two SharePoint lists (Instructors, QualificationRecords) and surfaced via the Power BI Instructor Quals Dashboard. |
| 2 | Does the tool collect, store, or process Protected Health Information (PHI)? | No | No medical data, injury records, or health-related information is collected. Physical fitness data is not included in instructor tracking (only in student performance). No medical readiness status is stored. |
| 3 | Does the tool link PII across multiple data sources? | Yes (minimal, internal only) | QualificationRecords links to Instructors via InstructorEDIPI. Both lists are within the same SharePoint site and MCEN M365 boundary. No external data sources are linked — no MCTIMS, MCTFS, or personnel system integration in Phase 1. This is a standard relational join within a single application, not cross-system PII linkage. |
| 4 | Does the tool create new PII (e.g., user accounts, tracking IDs, BARS evaluation records)? | No | Instructor records reflect existing TBS staff data. EDIPIs are existing DoD identifiers. Qualification records document existing certifications. No new identifiers or PII elements are created by the system. Calculated fields (days until expiration, workload metrics) are derived measures, not new stored PII. |
| 5 | Is PII visible to users other than the data subject? | Yes | By design: Staff (S-3, Operations Officer, Company Commanders) see all instructor data including names, qualifications, workload, and PRD dates. This is a staff-only dashboard — instructors themselves do not have a self-service view in Phase 1. This mirrors the existing visibility that S-3 and company leadership already have over instructor assignment and qualification data via paper rosters and spreadsheets. Student role has no access to this dashboard. |
| 6 | Does the tool transmit PII outside TBS (e.g., to TECOM, MCTIMS, M&RA, DHA)? | No | All data remains within the MCEN M365 tenant boundary. SharePoint Online and Power BI Service are both within the same Azure Government IL5 environment. No data exports, automated flows, or API connections transmit instructor PII externally. No external sharing links are configured. No data is sent to TECOM, HQMC, or other external organizations. |
| 7 | Can data be aggregated/anonymized to avoid PII? | Partially | Company-level qualification coverage rates and aggregate workload statistics do not require PII. However, the core purpose — tracking which specific instructors hold which qualifications and when they expire — inherently requires individual identification. Anonymization would defeat the purpose: S-3 must know that "Capt Smith's RSO certification expires in 14 days," not that "an instructor in Bravo Company has an expiring cert." |

---

## TBS-Specific Data Categories

| Data Type | Contains PII? | Contains PHI? | Classification | Authorized Platform |
|-----------|:---:|:---:|---|---|
| Instructor names + EDIPIs | Yes | No | CUI | IL5 (MCEN M365 SharePoint) |
| Instructor contact info (duty phone, duty email) | Yes (minimal) | No | CUI | IL5 (MCEN M365 SharePoint) |
| Instructor rank, MOS, role | Yes (linked to named individual) | No | CUI | IL5 |
| PRD (Projected Rotation Date) | Yes (linked to named individual) | No | CUI | IL5 |
| Qualification records (certs, dates, expiration) | Yes (linked to instructor via EDIPI) | No | CUI | IL5 |
| Workload metrics (students assigned, events) | Yes (linked to instructor) | No | CUI | IL5 |
| Aggregate company qualification coverage | No | No | Unclassified | Any MCEN platform |
| RequiredQualifications master list | No | No | Unclassified | Any MCEN platform |

---

## Data Classification

| Data Type | Classification Level | Handling Requirement |
|-----------|---------------------|---------------------|
| Input data | CUI — instructor PII (names, EDIPIs, contact info, qualification records) | Stored in SharePoint Online within MCEN M365 IL5 boundary. Manual entry by authorized Staff (S-3 section) only. Instructors list is the master record; QualificationRecords references it via EDIPI. |
| Output data | CUI — Power BI dashboard displays individual instructor records, qualification status, workload | Staff-only access via Power BI Service. No SPC or Student access. Dashboard published to staff-only Power BI workspace. RLS configured for company commander views if needed. |
| Stored data | CUI — two SharePoint lists with PII fields | SharePoint Online encryption at rest (Microsoft-managed keys within Azure Gov IL5). RequiredQualifications list contains no PII (reference data only). |
| Logged data | CUI — SharePoint audit logs capture who accessed/modified records | M365 Unified Audit Log enabled by MCEN tenant policy. Logs contain user identity + action + timestamp. Retained per MCEN retention policy. |
| AI-processed data | N/A — Phase 1 dashboard does not use AI/LLM processing | No data is sent to Azure OpenAI or any AI service. All calculations are DAX measures within Power BI. AI-assisted qualification gap analysis is deferred to Phase 2+. |

---

## Threshold Determination

- [ ] No PII involved — PIA not required (aggregate/anonymized data only)
- [x] Minimal PII (names + scores, no sensitive PII) — Threshold analysis sufficient; consult Privacy Officer
- [ ] PII involved with identifiers — Full PIA required (DD Form 2930)
- [ ] PHI involved — Full PIA + HIPAA review required (coordinate with DHA)

**Rationale:** The system collects minimal PII limited to instructor names, EDIPIs, duty contact information, and professional qualification records. This PII is not sensitive PII as defined by DoDI 5400.11 — no SSNs, financial data, medical data, biometric data, or law enforcement records are involved. The contact information is duty-related only (duty phone, .mil email), not personal contact information. The data remains entirely within the MCEN M365 IL5 authorization boundary with no external transmission. Access is restricted to Staff roles only (S-3, Operations, Company Commanders) via Power BI workspace permissions and RLS. The audience for this data is narrower than the student performance dashboard — only staff leadership, not SPCs or students.

A full PIA (DD Form 2930) is not required for Phase 1 provided:

1. Data stays within MCEN M365 boundary (no external sharing)
2. Access is restricted to authorized Staff roles only
3. No personal contact information is stored (duty contact only)
4. No cross-system PII linkage to external personnel systems occurs
5. TBS Privacy Officer concurs with this threshold analysis

---

## Data Minimization Recommendations

*Before proceeding, document what PII can be eliminated or reduced:*

| Current Data Need | Can It Be Anonymized? | Recommended Approach |
|-------------------|-----------------------|---------------------|
| Instructor names (FirstName, LastName) | No — S-3 and company leadership must identify specific instructors for qualification management, workload balancing, and rotation planning | Retain. Staff-only visibility. Names are essential for actionable management decisions (e.g., "Send Capt Smith to RSO recertification"). |
| Instructor EDIPI | No — required as unique identifier to link Instructors list to QualificationRecords list and for future MCTIMS integration | Retain. EDIPI is the standard DoD identifier. Used as relational key, not prominently displayed in dashboard visuals. |
| Duty contact info (phone, email) | Partially — duty email could be derived from Azure AD, but duty phone is operationally needed for immediate contact about qualification gaps | Retain duty phone for operational necessity. Duty email retained for convenience but could be removed if Privacy Officer requires (derivable from name + Azure AD). |
| PRD (Projected Rotation Date) | No — essential for rotation planning and identifying qualification gaps that will arise when instructors PCS | Retain. This is the key data point for proactive qualification management. Without PRD, the "Quals at Risk from Rotation" analysis cannot function. |
| Individual qualification records | No — the core purpose requires tracking which specific instructor holds which certification and when it expires | Retain. Aggregate-only views ("Bravo Company has 3 RSO-qualified instructors") are useful but insufficient — S-3 must know exactly who is qualified to assign instructors to events. |

---

**Reviewer Signature:** _____________________ **Date:** _____________

**TBS Privacy Officer Concurrence:** _____________________ **Date:** _____________
