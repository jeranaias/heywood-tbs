# PIA Threshold Analysis

*The Heywood Initiative — Privacy Impact Assessment Threshold Checklist*

---

## Tool Information

| Field              | Entry |
|--------------------|-------|
| **Tool Name**      | Student Performance Data (SharePoint StudentScores list + Power BI Student Performance Dashboard) |
| **Heywood Phase**  | Phase 1 |
| **Developer**      | SSgt Jesse C. Morgan |
| **Date**           | March 2026 |
| **Reviewer**       | Pending — TBS Privacy Officer |

---

## PII Determination

Answer each question below. If any answer is "Yes," a full Privacy Impact Assessment (PIA) is required per DoDI 5400.16.

| # | Question | Yes / No | Notes |
|---|----------|----------|-------|
| 1 | Does the tool collect, store, or process Personally Identifiable Information (PII)? | Yes | Student names, EDIPIs, individual grades/scores, class standing, peer evaluation scores, SPC evaluation scores. All stored in SharePoint StudentScores list and surfaced via Power BI dashboard. |
| 2 | Does the tool collect, store, or process Protected Health Information (PHI)? | No | No injury records, medical readiness data, or MHS GENESIS data. PFT/CFT scores are fitness performance data, not PHI. No medical hold reasons are stored — only the status "Medical Hold" without clinical detail. |
| 3 | Does the tool link PII across multiple data sources? | No | Phase 1 uses a single SharePoint list (StudentScores) as the sole data source. No cross-linking to MCTIMS, MCTFS, HPB, or other external systems. Power BI dimension tables (DimCompany, DimPhase) contain no PII. The UserSecurity table maps Azure AD emails to roles but does not link student PII across systems. |
| 4 | Does the tool create new PII (e.g., user accounts, tracking IDs, BARS evaluation records)? | No | StudentScores records are manually entered by SPCs from existing TBS grading records. No new identifiers are generated — EDIPI is the existing DoD identifier. Composite scores are calculated measures in Power BI (DAX), not new stored data. Phase 1 does not include BARS evaluation records (Phase 2). |
| 5 | Is PII visible to users other than the data subject? | Yes | By design: SPCs see their platoon/company students' names, scores, and standings. Staff sees all students across all companies. Students see only their own individual record. Row-Level Security (RLS) in Power BI enforces these boundaries. This mirrors the existing paper-based and Excel-based visibility that SPCs and Staff already have — the dashboard does not expand the audience for student PII beyond current practice. |
| 6 | Does the tool transmit PII outside TBS (e.g., to TECOM, MCTIMS, M&RA, DHA)? | No | All data remains within the MCEN M365 tenant boundary. SharePoint Online and Power BI Service are both within the same Azure Government IL5 environment. No data exports, API connections, or automated flows transmit PII externally. No external sharing links are configured. |
| 7 | Can data be aggregated/anonymized to avoid PII? | Partially | Company-level aggregates (averages, at-risk counts, phase distribution) do not require PII. However, the core purpose — SPC tracking of individual student performance for counseling and grading — requires student-level PII (names, EDIPIs, individual scores). Data minimization is applied: only performance-relevant fields are collected, no SSNs, no medical data, no financial data. |

---

## TBS-Specific Data Categories

| Data Type | Contains PII? | Contains PHI? | Classification | Authorized Platform |
|-----------|:---:|:---:|---|---|
| Student names + EDIPIs | Yes | No | CUI | IL5 (MCEN M365 SharePoint) |
| Individual grades/scores | Yes (linked to student via EDIPI) | No | CUI | IL5 (MCEN M365 SharePoint + Power BI) |
| Aggregate company averages | No | No | Unclassified | Any MCEN platform |
| SPC evaluation scores (Week 12, Week 22) | Yes (linked to student) | No | CUI | IL5 |
| Peer evaluation scores (Week 12, Week 22) | Yes (linked to student) | No | CUI | IL5 |
| PFT/CFT scores (individual) | Yes (linked to student) | No | CUI | IL5 |
| At-risk flag and status | Yes (linked to student) | No | CUI | IL5 |
| Training schedule (no names) | No | No | Unclassified | Any MCEN platform |

---

## Data Classification

| Data Type | Classification Level | Handling Requirement |
|-----------|---------------------|---------------------|
| Input data | CUI — student PII (names, EDIPIs, individual scores) | Stored in SharePoint Online within MCEN M365 IL5 boundary. Manual entry by authorized SPCs and Staff only. No bulk upload from external sources in Phase 1. |
| Output data | CUI — Power BI dashboard displays individual student records | RLS-filtered views only. No data export enabled for SPC/Student roles. Staff export restricted to authorized personnel. Dashboard shared via Power BI App, not direct workspace access. |
| Stored data | CUI — SharePoint list items with PII fields | SharePoint Online encryption at rest (Microsoft-managed keys within Azure Gov IL5). Item-level permissions not required — RLS enforced at Power BI layer. SharePoint column-level encryption not available but data is within IL5 boundary. |
| Logged data | CUI — SharePoint audit logs capture who accessed/modified records | M365 Unified Audit Log enabled by MCEN tenant policy. Logs retained per MCEN retention policy (minimum 90 days). Logs contain user identity + action + timestamp. |
| AI-processed data | N/A — Phase 1 dashboard does not use AI/LLM processing | No data is sent to Azure OpenAI or any AI service. All calculations are DAX measures within Power BI. AI processing of student data is deferred to Phase 2+ with separate PIA. |

---

## Threshold Determination

- [ ] No PII involved — PIA not required (aggregate/anonymized data only)
- [x] Minimal PII (names + scores, no sensitive PII) — Threshold analysis sufficient; consult Privacy Officer
- [ ] PII involved with identifiers — Full PIA required (DD Form 2930)
- [ ] PHI involved — Full PIA + HIPAA review required (coordinate with DHA)

**Rationale:** The system collects minimal PII (student names, EDIPIs, and academic/military/leadership performance scores). This PII is not sensitive PII as defined by DoDI 5400.11 (no SSNs, financial data, medical data, biometric data, or law enforcement records). The data remains entirely within the MCEN M365 IL5 authorization boundary with no external transmission. Access is controlled via RBAC (Azure AD security groups) and Row-Level Security (Power BI). The visibility model mirrors existing TBS practices — SPCs already maintain this data in Excel spreadsheets and paper records. A full PIA (DD Form 2930) is not required for Phase 1 provided:

1. Data stays within MCEN M365 boundary (no external sharing)
2. No cross-system PII linkage occurs
3. RBAC and RLS are properly configured and tested before data entry
4. TBS Privacy Officer concurs with this threshold analysis

---

## Data Minimization Recommendations

*Before proceeding, document what PII can be eliminated or reduced:*

| Current Data Need | Can It Be Anonymized? | Recommended Approach |
|-------------------|-----------------------|---------------------|
| Student names (FirstName, LastName) | No — SPCs must identify students for counseling, grading, and daily management | Retain. RLS ensures students cannot see other students' names/scores. Names visible only to assigned SPC and Staff. |
| Student EDIPI | No — required as unique identifier to link records and enforce RLS for student self-service view | Retain. EDIPI is the standard DoD identifier; using an alternative would create a new PII element. Not displayed in dashboard visuals (used as filter key only). |
| Individual scores and composites | No — individual performance tracking is the core purpose of the tool | Retain. Aggregate-only views would eliminate the tool's utility for SPC counseling and student self-assessment. Company averages are provided as additional context, not as a replacement. |
| Peer evaluation scores | Partially — individual peer scores are linked to the evaluated student, not the evaluator | Evaluator identities are not stored. Only the aggregated peer evaluation score per student is recorded. This is the minimum PII footprint for peer evaluation data. |

---

**Reviewer Signature:** _____________________ **Date:** _____________

**TBS Privacy Officer Concurrence:** _____________________ **Date:** _____________
