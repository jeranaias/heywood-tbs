# Compliance Checklist

*The Heywood Initiative — Pre-Deployment Compliance Verification*

---

## Component Information

| Field              | Entry |
|--------------------|-------|
| **Component Name** | |
| **Heywood Phase**  | Phase 1 / Phase 2 / Phase 3 / Phase 4 |
| **Developer**      | |
| **Date**           | |
| **Reviewer**       | |

---

## Compliance Items

Complete each item before the component proceeds to deployment.

### Data & Privacy

| # | Requirement | Status | Notes |
|---|-------------|--------|-------|
| 1 | Data classification completed (CUI, PII, PHI determination) | | |
| 2 | PIA threshold analysis completed (use `governance/pia/pia-threshold-analysis.md`) | | |
| 3 | TBS Privacy Officer consulted (if PII is involved) | | |
| 4 | SORN requirements reviewed (if PII is retrieved by identifier) | | |
| 5 | Data minimization applied (aggregate where possible, anonymize where required) | | |
| 6 | AI data handling verified (no PII/PHI in prompts unless on IL5 platform) | | |

### Authorization & Security

| # | Requirement | Status | Notes |
|---|-------------|--------|-------|
| 7 | ATO/IATT status verified for this component | | *Phase 1: inherited from MCEN. Phase 2+: IATT or ATO required.* |
| 8 | Platform authorization confirmed (Azure Gov IL5, GenAI.mil, MCEN M365) | | |
| 9 | Access controls implemented (CAC auth, RBAC for Staff/SPC/Student) | | |
| 10 | Row-Level Security configured and tested (each role sees only authorized data) | | |
| 11 | Audit logging enabled | | |
| 12 | Data encryption verified (at rest: CMK for IL5; in transit: TLS 1.2+) | | |

### Operational

| # | Requirement | Status | Notes |
|---|-------------|--------|-------|
| 13 | OPSEC review completed (no sensitive indicators in AI outputs) | | |
| 14 | Records management requirements identified (retention schedule) | | |
| 15 | Incident response procedures documented | | |
| 16 | Input validation implemented (prevent prompt injection, data corruption) | | |
| 17 | Component registered in Tool Registry (`governance/registry/`) | | |
| 18 | User documentation created | | |
| 19 | Backup/recovery plan in place (SharePoint versioning, Dataverse backup) | | |

---

## Approval

- [ ] All items complete — component may proceed to deployment
- [ ] Items outstanding — component must not deploy until resolved

**Outstanding items (if any):**
1.
2.

**Reviewer Signature:** _____________________ **Date:** _____________

**TBS S-6 / ISSM Concurrence (Phase 2+):** _____________________ **Date:** _____________
