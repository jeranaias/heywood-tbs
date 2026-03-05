# PIA Threshold Analysis

*The Heywood Initiative — Privacy Impact Assessment Threshold Checklist*

---

## Tool Information

| Field              | Entry |
|--------------------|-------|
| **Tool Name**      | |
| **Heywood Phase**  | Phase 1 / Phase 2 / Phase 3 / Phase 4 |
| **Developer**      | |
| **Date**           | |
| **Reviewer**       | |

---

## PII Determination

Answer each question below. If any answer is "Yes," a full Privacy Impact Assessment (PIA) is required per DoDI 5400.16.

| # | Question | Yes / No | Notes |
|---|----------|----------|-------|
| 1 | Does the tool collect, store, or process Personally Identifiable Information (PII)? | | *TBS PII includes: student names, EDIPIs, grades, class standing, peer evaluations, fitness scores, evaluation narratives* |
| 2 | Does the tool collect, store, or process Protected Health Information (PHI)? | | *PHI includes: injury records, medical readiness status, MHS GENESIS data. If yes, HIPAA review also required.* |
| 3 | Does the tool link PII across multiple data sources? | | *E.g., linking MCTIMS training records to student performance scores to HPB injury data* |
| 4 | Does the tool create new PII (e.g., user accounts, tracking IDs, BARS evaluation records)? | | |
| 5 | Is PII visible to users other than the data subject? | | *SPCs see their platoon's students. Staff sees all. Students see only themselves.* |
| 6 | Does the tool transmit PII outside TBS (e.g., to TECOM, MCTIMS, M&RA, DHA)? | | |
| 7 | Can data be aggregated/anonymized to avoid PII? | | *E.g., company averages instead of individual scores; duty positions instead of names* |

---

## TBS-Specific Data Categories

| Data Type | Contains PII? | Contains PHI? | Classification | Authorized Platform |
|-----------|:---:|:---:|---|---|
| Student names + EDIPIs | Yes | No | CUI | IL5 (GenAI.mil, Azure Gov, Dataverse) |
| Individual grades/scores | Yes (if linked to student) | No | CUI | IL5 |
| Aggregate company averages | No | No | Unclassified | Any MCEN platform |
| SPC evaluation narratives | Yes | No | CUI | IL5 |
| Peer evaluation data | Yes | No | CUI | IL5 |
| PFT/CFT scores (individual) | Yes | No | CUI | IL5 |
| Injury/medical data | Yes | Yes | CUI + PHI | IL5 + PIA + HIPAA |
| MHS GENESIS data | Yes | Yes | CUI + PHI | Requires DHA DSA + PIA |
| Training schedule (no names) | No | No | Unclassified | Any MCEN platform |
| AAR data (no student names) | No | No | CUI (if operational) | IL5 |
| Instructor qualifications | Yes (minimal) | No | CUI | IL5 |

---

## Data Classification

| Data Type | Classification Level | Handling Requirement |
|-----------|---------------------|---------------------|
| Input data | | |
| Output data | | |
| Stored data | | |
| Logged data | | |
| AI-processed data | | *Data sent to Azure OpenAI API — stays within IL5 boundary; not used for model training* |

---

## Threshold Determination

- [ ] No PII involved — PIA not required (aggregate/anonymized data only)
- [ ] Minimal PII (names + scores, no sensitive PII) — Threshold analysis sufficient; consult Privacy Officer
- [ ] PII involved with identifiers — Full PIA required (DD Form 2930)
- [ ] PHI involved — Full PIA + HIPAA review required (coordinate with DHA)

---

## Data Minimization Recommendations

*Before proceeding, document what PII can be eliminated or reduced:*

| Current Data Need | Can It Be Anonymized? | Recommended Approach |
|-------------------|-----------------------|---------------------|
| | | |
| | | |
| | | |

---

**Reviewer Signature:** _____________________ **Date:** _____________

**TBS Privacy Officer Concurrence:** _____________________ **Date:** _____________
