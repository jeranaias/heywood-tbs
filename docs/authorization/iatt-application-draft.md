# Interim Authority to Test (IATT) Application — Draft

*The Heywood Initiative — Phase 2*

---

## 1. System Identification

| Field | Entry |
|-------|-------|
| **System Name** | Heywood TBS AI-Assisted Training Management System |
| **Acronym** | Heywood |
| **Version** | 2.0 (Phase 2 — Custom Agent with IATT) |
| **System Owner** | The Basic School, Training and Education Command |
| **System Administrator** | SSgt Jesse C. Morgan |
| **ISSM** | [TBS S-6 ISSM — name TBD] |
| **Authorizing Official** | MCEN ISSM + TBS S-6 + TECOM CIO |
| **Requested Duration** | 180 days (renewable) |
| **Requested Start Date** | [TBD — upon completion of Phase 1 exit criteria] |

---

## 2. System Description

### 2.1 Purpose

Heywood is an AI-assisted training management system for The Basic School that integrates student performance tracking, instructor qualification management, AI-powered document generation, and doctrine retrieval into a single Role-Based Access Controlled (RBAC) application.

Phase 2 extends the Phase 1 proof of concept (Power BI dashboards + SharePoint lists + GenAI.mil prompts) by adding:
- A custom Power App with three role-specific views (Staff, SPC, Student)
- Azure OpenAI integration for AI-assisted counseling narrative generation and document drafting
- Azure AI Search for Retrieval-Augmented Generation (RAG) over TBS doctrine and POI documents
- Power Automate workflows for AAR processing and counseling packet generation

### 2.2 Scope of IATT

| Aspect | Boundary |
|--------|----------|
| **Users** | One TBS company (~200 students, 4-6 SPCs, company staff) |
| **Total user accounts** | ~220 (10 staff/instructor, ~210 student) |
| **Data** | TBS training data (CUI), student performance scores, instructor qualifications, anonymized AAR content |
| **Environment** | Non-production pilot; test data and one company's live data only |
| **Networks** | MCEN (NIPRNet) only — no SIPR, no external connections |
| **Physical location** | The Basic School, MCB Quantico, VA |

### 2.3 System Components

| Component | Platform | Purpose | IL Level |
|-----------|----------|---------|----------|
| Heywood Power App | Power Platform (MCEN) | User interface with RBAC views | IL5 (inherited) |
| Azure OpenAI Service | Azure Government | GPT-4o reasoning, text-embedding-ada-002 for embeddings | IL5 |
| Azure AI Search | Azure Government | Vector store for RAG over TBS doctrine | IL5 |
| Azure App Service | Azure Government | Custom connector between Power App and Azure OpenAI | IL5 |
| Azure Key Vault | Azure Government | API keys, connection strings, CMK storage | IL5 |
| SharePoint Online | MCEN M365 | Data storage (6 lists from Phase 1) | IL5 (inherited) |
| Power Automate | Power Platform (MCEN) | Workflow automation (AAR processing, counseling packets) | IL5 (inherited) |
| Power BI Service | MCEN M365 | Dashboards (carried forward from Phase 1) | IL5 (inherited) |
| Azure Entra ID | MCEN | Authentication (CAC) and authorization (RBAC groups) | IL5 (inherited) |

---

## 3. Authorization Boundary

### 3.1 Boundary Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    MCEN AUTHORIZATION BOUNDARY                   │
│                                                                  │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────────┐   │
│  │  Azure Entra  │    │  SharePoint  │    │   Power BI       │   │
│  │  ID (CAC)     │    │  Online      │    │   Service        │   │
│  │  [Inherited]  │    │  [Inherited]  │    │   [Inherited]    │   │
│  └──────┬───────┘    └──────┬───────┘    └────────┬─────────┘   │
│         │                   │                      │             │
│  ┌──────┴───────────────────┴──────────────────────┴─────────┐  │
│  │                    Power Platform (MCEN)                    │  │
│  │  ┌──────────────────┐    ┌──────────────────────────────┐  │  │
│  │  │  Heywood Power   │    │  Power Automate Flows        │  │  │
│  │  │  App (Canvas)    │    │  - AAR Processing            │  │  │
│  │  │  - Staff View    │    │  - Counseling Packet Gen     │  │  │
│  │  │  - SPC View      │    │  - Instructor Alerts         │  │  │
│  │  │  - Student View  │    └──────────────────────────────┘  │  │
│  │  └────────┬─────────┘                                      │  │
│  └───────────┼────────────────────────────────────────────────┘  │
│              │ Custom Connector (HTTPS/TLS 1.2+)                 │
│  ┌───────────┼────────────────────────────────────────────────┐  │
│  │           ▼          AZURE GOVERNMENT (IL5)                │  │
│  │  ┌──────────────┐   ┌──────────────┐   ┌──────────────┐   │  │
│  │  │  Azure App   │   │  Azure       │   │  Azure AI    │   │  │
│  │  │  Service     │   │  OpenAI      │   │  Search      │   │  │
│  │  │  (Connector) │──▶│  Service     │   │  (RAG Index) │   │  │
│  │  │              │   │  GPT-4o      │   │  S1 + CMK    │   │  │
│  │  │              │   │  Embeddings  │   │              │   │  │
│  │  └──────────────┘   └──────────────┘   └──────────────┘   │  │
│  │           │                                                │  │
│  │  ┌────────┴─────┐   ┌──────────────┐                      │  │
│  │  │  Azure Key   │   │  Azure       │                      │  │
│  │  │  Vault       │   │  Monitor     │                      │  │
│  │  │  (Secrets)   │   │  (Logging)   │                      │  │
│  │  └──────────────┘   └──────────────┘                      │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  NO EXTERNAL CONNECTIONS — all traffic stays within MCEN/Azure  │
│  Gov boundary. No internet egress. No SIPR connectivity.         │
└─────────────────────────────────────────────────────────────────┘
```

### 3.2 Data Flow

```
1. User authenticates via CAC → Azure Entra ID validates → assigns role (Staff/SPC/Student)
2. Power App loads role-specific view → queries SharePoint for data (filtered by RBAC)
3. User requests AI assistance (e.g., counseling narrative):
   a. Power App sends request to Custom Connector (App Service)
   b. App Service forwards to Azure OpenAI with system prompt + user context
   c. Azure OpenAI generates response
   d. Response returns through App Service → Power App
   e. No PII is sent to Azure OpenAI — context is anonymized at the App Service layer
4. RAG queries:
   a. User query → App Service → text-embedding-ada-002 (vectorize query)
   b. Vector → Azure AI Search (retrieve relevant doctrine chunks)
   c. Chunks + query → GPT-4o (generate grounded answer)
   d. Response returns to user with source citations
5. Power Automate flows:
   a. AAR submission → extract themes via Azure OpenAI → categorize → update SharePoint
   b. Counseling packet request → aggregate student data → generate narrative → SPC review queue
```

### 3.3 External Interfaces

| Interface | Direction | Data | Protocol | Classification |
|-----------|-----------|------|----------|----------------|
| Azure Entra ID → Power App | Inbound | Auth tokens (CAC) | SAML 2.0 / OAuth 2.0 | CUI |
| Power App → SharePoint | Bidirectional | Training data | HTTPS/TLS 1.2+ | CUI |
| Power App → App Service | Outbound | AI queries (anonymized) | HTTPS/TLS 1.2+ | CUI |
| App Service → Azure OpenAI | Outbound | Prompts + context | HTTPS/TLS 1.2+ | CUI |
| App Service → Azure AI Search | Outbound | Search queries | HTTPS/TLS 1.2+ | CUI |
| Power BI → SharePoint | Inbound | Data refresh | HTTPS/TLS 1.2+ | CUI |

**No external interfaces exist.** All traffic remains within the MCEN/Azure Government boundary.

---

## 4. Security Controls

### 4.1 Inherited Controls (from Azure Gov FedRAMP High + MCEN)

Azure Government maintains FedRAMP High authorization, which satisfies approximately 60% of NIST 800-53 Rev 5 controls. The following control families are fully or substantially inherited:

| Control Family | Inherited From | Coverage |
|----------------|---------------|----------|
| AC (Access Control) | Azure Entra ID + MCEN | Partial — system-specific RBAC is Heywood responsibility |
| AU (Audit) | Azure Monitor + M365 Audit | Full |
| AT (Awareness & Training) | EDD Course 1-5 curriculum | System-specific |
| CA (Assessment & Authorization) | This IATT document | System-specific |
| CM (Configuration Management) | Azure Gov + MCEN | Partial — Heywood configs are system-specific |
| CP (Contingency Planning) | Azure Gov SLA + SharePoint versioning | Partial |
| IA (Identification & Auth) | Azure Entra ID + CAC | Full |
| IR (Incident Response) | MCEN SOC + TBS S-6 | Partial — system-specific procedures below |
| MA (Maintenance) | Azure Gov | Full |
| MP (Media Protection) | Azure Gov encryption | Full |
| PE (Physical & Environmental) | Azure Gov data centers | Full |
| PL (Planning) | This document | System-specific |
| PM (Program Management) | Heywood governance artifacts | System-specific |
| PS (Personnel Security) | MCEN CAC enrollment | Full |
| RA (Risk Assessment) | This document, Section 6 | System-specific |
| SA (System & Services Acquisition) | Azure Gov | Full |
| SC (System & Communications) | Azure Gov TLS + encryption | Partial |
| SI (System & Information Integrity) | Azure Defender + MCEN | Partial |

### 4.2 System-Specific Controls

**AC-2: Account Management**
- All accounts are Azure AD accounts provisioned via MCEN standard process
- No local accounts exist
- Account removal follows MCEN standard offboarding (automatic on PCS/EAS)
- RBAC groups: Heywood-Staff, Heywood-SPC, Heywood-Student
- Group membership managed by TBS S-6

**AC-3: Access Enforcement**
- Three-tier RBAC enforced at every layer:

| Role | Power App View | SharePoint Access | Power BI RLS | AI Features |
|------|---------------|-------------------|-------------|-------------|
| Staff | Full admin | All lists, all rows | All companies | Full (counseling, RAG, analytics) |
| SPC | Company view | Filtered to company | Company only | Counseling prep, AAR, scenarios |
| Student | Individual view | Own record only | Own scores only | Study plan generation only |

**AC-6: Least Privilege**
- Students cannot see other students' data
- SPCs cannot see other companies' data
- Only Staff can access instructor qualification data
- Azure OpenAI endpoint accessible only via App Service managed identity — no direct user access

**AU-2/AU-3: Audit Events / Content**
- Azure Monitor captures all API calls to Azure OpenAI, AI Search, App Service
- M365 Unified Audit Log captures all SharePoint and Power App activity
- Power BI audit log captures all dashboard access and data refreshes
- Logs retained per MCEN standard retention (90 days online, 1 year archive)
- Audit events include: user identity (CAC-linked), timestamp, action, resource, outcome

**IA-2: Identification and Authentication**
- CAC authentication required for all access (no username/password option)
- Multi-factor: CAC (something you have) + PIN (something you know)
- No service accounts — App Service uses Azure Managed Identity

**SC-8: Transmission Confidentiality**
- All data in transit encrypted with TLS 1.2+
- No HTTP endpoints exist — HTTPS only
- Certificate management via Azure-managed certificates

**SC-12/SC-13: Cryptographic Key Management / Protection**
- Azure AI Search index encrypted with Customer-Managed Key (CMK) via Azure Key Vault
- Required for IL5 compliance
- Key rotation: annual, automated via Key Vault policy
- Azure OpenAI data encryption: platform-managed keys (Microsoft-managed, FedRAMP High compliant)
- SharePoint: platform-managed encryption at rest (BitLocker + per-file encryption)

**SC-28: Protection of Information at Rest**
- SharePoint: encrypted at rest (Microsoft-managed keys)
- Azure AI Search: CMK encryption via Key Vault
- Azure OpenAI: no persistent storage of prompts or responses (stateless API)
- Key Vault: HSM-backed key storage

**SI-4: Information System Monitoring**
- Azure Monitor alerts configured for:
  - Unusual API call volume (>200% of baseline)
  - Failed authentication attempts (>5 in 10 minutes)
  - Data export attempts
  - Azure OpenAI token usage anomalies
- Alert routing: TBS S-6 + System Administrator

---

## 5. Data Handling

### 5.1 Data Classification

| Data Category | Classification | PII | PHI | Handling |
|---------------|---------------|:---:|:---:|----------|
| Student names, EDIPIs | CUI | Yes | No | Stored in SharePoint (IL5), never sent to Azure OpenAI |
| Student scores (individual) | CUI | Yes | No | Stored in SharePoint, aggregated in Power BI, anonymized before AI processing |
| Student scores (aggregate) | CUI | No | No | Used freely in AI queries and dashboards |
| Instructor names, EDIPIs | CUI | Yes | No | Stored in SharePoint (IL5), staff-only access |
| Qualification records | CUI | Yes | No | Linked to instructor EDIPI, staff-only access |
| Training schedule | CUI | No | No | No PII, shared across roles |
| Event feedback | CUI | Optional | No | Anonymous by default for student submissions |
| Doctrine documents (POI, FMs) | CUI | No | No | Indexed in Azure AI Search for RAG retrieval |
| AAR content | CUI | No | No | Anonymized before storage and AI processing |
| AI-generated narratives | CUI | No | No | Generated from anonymized context, reviewed by SPC before use |

### 5.2 PII Protection in AI Pipeline

**The PII anonymization layer** is the critical security control for this system:

```
User enters data in Power App (may include PII)
         ↓
App Service anonymization layer:
  - Strips names → replaces with "the student" / "Student A"
  - Strips EDIPIs → removes entirely
  - Strips platoon/company identifiers when context allows
  - Retains: scores, trends, event types, phase, doctrinal context
         ↓
Azure OpenAI receives anonymized context only
         ↓
AI response returns (no PII possible — none was provided)
         ↓
SPC reviews and re-associates with specific student manually
```

### 5.3 PIA Status

| Component | PIA Status | Document |
|-----------|-----------|----------|
| SharePoint lists (Phase 1) | Threshold analysis complete — full PIA not required | `governance/pia/pia-student-performance.md`, `pia-instructor-data.md` |
| Power App + AI features (Phase 2) | Threshold analysis required before IATT approval | To be completed |
| Azure AI Search (doctrine index) | Not required — no PII in indexed documents | N/A |

---

## 6. Risk Assessment

### 6.1 System-Specific Risks

| # | Risk | Likelihood | Impact | Mitigation |
|---|------|:---:|:---:|---|
| 1 | PII leaks to Azure OpenAI despite anonymization layer | Low | High | Anonymization at App Service layer (server-side, not client-side). Input validation rejects any 10-digit number pattern (EDIPI format). Automated testing of anonymization filter before deployment. |
| 2 | AI generates inaccurate counseling language that SPC uses without review | Medium | Medium | All AI outputs require explicit SPC review and approval before use. Disclaimer watermark on all AI-generated text. Training (EDD Course 1) emphasizes quality judgment. |
| 3 | RAG retrieves outdated doctrine | Low | Medium | Doctrine index versioned and dated. AI Search index rebuild triggered when source documents are updated. Responses include source document title and date. |
| 4 | Unauthorized access to student data beyond role | Low | High | Three-tier RBAC enforced at SharePoint (list permissions), Power App (conditional views), and Power BI (RLS). Tested before deployment with each role type. |
| 5 | Azure OpenAI service disruption | Medium | Low | System degrades gracefully — Power App continues to function without AI features. SharePoint data access and Power BI dashboards are not affected. |
| 6 | Excessive Azure consumption costs | Low | Medium | Azure budgets and spending alerts configured. Token-level monitoring on Azure OpenAI. Monthly cost review. |

### 6.2 POA&M Items

| # | Finding | Severity | Milestone | Target Date | Status |
|---|---------|----------|-----------|-------------|--------|
| 1 | PIA threshold for Phase 2 AI components | Medium | Complete PIA threshold analysis for Power App + Azure OpenAI | IATT application - 30 days | Pending |
| 2 | Anonymization filter testing | High | Complete automated test suite for PII anonymization layer | IATT application - 14 days | Pending |
| 3 | RBAC validation testing | Medium | Test all three roles with test accounts, document results | IATT approval + 7 days | Pending |
| 4 | Incident response SOP | Medium | Document Heywood-specific IR procedures | IATT approval + 14 days | Pending |
| 5 | Audit log review SOP | Low | Establish weekly audit log review process | IATT approval + 30 days | Pending |
| 6 | Continuity of operations plan | Low | Document system recovery procedures | IATT approval + 30 days | Pending |

---

## 7. Test Plan

### 7.1 Security Testing

| Test | Method | Pass Criteria |
|------|--------|---------------|
| RBAC enforcement | Log in as each role, attempt to access data outside role boundary | Zero unauthorized data access |
| PII anonymization | Send 50 test prompts containing synthetic PII through the pipeline | Zero PII reaches Azure OpenAI endpoint (verified via Azure Monitor logs) |
| CAC authentication | Attempt access without CAC, with expired CAC, with valid CAC | Only valid CAC succeeds |
| TLS enforcement | Attempt HTTP connection to all endpoints | All connections refuse non-HTTPS |
| CMK encryption | Verify AI Search index encryption via Azure portal | CMK key ID matches Key Vault key |
| Audit logging | Perform 20 representative actions, verify all appear in audit log | 100% action capture with required fields |

### 7.2 Functional Testing

| Test | Method | Pass Criteria |
|------|--------|---------------|
| Power App — Staff view | Staff user accesses app, views all companies | All data visible, all features functional |
| Power App — SPC view | SPC user accesses app, views company data | Only assigned company visible, AI features work |
| Power App — Student view | Student accesses app | Only own record visible, limited features |
| AI counseling narrative | SPC requests counseling draft via Power App | Relevant, well-structured output within 10 seconds |
| RAG doctrine retrieval | Query TBS-specific doctrine question | Relevant chunks retrieved, grounded response, source cited |
| Power Automate — AAR flow | Submit test AAR via Power App | Themes extracted, categorized, stored in SharePoint |
| Graceful degradation | Disable Azure OpenAI endpoint | Power App continues to function, AI features show "unavailable" message |

---

## 8. Operational Procedures

### 8.1 Monitoring

| Activity | Frequency | Responsible |
|----------|-----------|-------------|
| Azure Monitor alert review | Daily (automated alerts) | System Administrator |
| Audit log review | Weekly | System Administrator + S-6 |
| Azure OpenAI usage/cost review | Weekly | System Administrator |
| RBAC group membership review | Monthly | TBS S-6 |
| CMK key rotation verification | Annually | System Administrator |

### 8.2 Incident Response

For any suspected security incident involving Heywood:

1. **Contain:** Disable the affected component (Power App, custom connector, or Azure OpenAI endpoint)
2. **Notify:** TBS S-6 ISSM within 1 hour, MCEN SOC per standing procedures
3. **Preserve:** Capture Azure Monitor logs and M365 audit logs for the incident timeframe
4. **Assess:** Determine if PII was exposed, if unauthorized access occurred, scope of impact
5. **Remediate:** Fix root cause, re-test, document
6. **Report:** Incident report to Authorizing Official within 72 hours

### 8.3 Contingency

| Scenario | Response |
|----------|----------|
| Azure OpenAI unavailable | System continues without AI features; manual workflows resume |
| SharePoint unavailable | MCEN M365 SLA applies; Power BI shows cached data |
| Power App unavailable | Users access Power BI directly for read-only data |
| Complete system outage | Revert to manual processes used before Heywood; no data loss (SharePoint versioning) |
| Personnel change (System Admin PCSes) | Backup administrator trained, documentation in Heywood repository |

---

## 9. Estimated Cost

| Resource | Monthly Cost | 6-Month IATT Total |
|----------|-------------|---------------------|
| Azure OpenAI (GPT-4o) | $200-500 | $1,200-3,000 |
| Azure OpenAI (embeddings) | $50-100 | $300-600 |
| Azure AI Search (S1) | $250 | $1,500 |
| Azure App Service (B1) | $55 | $330 |
| Azure Key Vault | $10 | $60 |
| Azure Monitor | $25 | $150 |
| Azure Gov overhead | $50 | $300 |
| Power Platform Premium (10 users) | $400 | $2,400 |
| **Total** | **$1,040-1,390** | **$6,240-8,340** |

---

## 10. Signatures

| Role | Name | Signature | Date |
|------|------|-----------|------|
| System Owner | [TBS CO/XO] | | |
| System Administrator | SSgt Jesse C. Morgan | | |
| Information System Security Manager | [TBS S-6 ISSM] | | |
| Authorizing Official | [MCEN ISSM / TECOM CIO] | | |

---

## 11. Attachments

| # | Document | Location |
|---|----------|----------|
| A | Architecture Brief — Heywood Build Plan v1 | `docs/briefs/heywood_build_plan_v1.md` |
| B | PIA Threshold — Student Performance Data | `governance/pia/pia-student-performance.md` |
| C | PIA Threshold — Instructor Data | `governance/pia/pia-instructor-data.md` |
| D | Compliance Checklist — SharePoint | `governance/compliance/compliance-sharepoint.md` |
| E | Compliance Checklist — Power BI | `governance/compliance/compliance-powerbi.md` |
| F | Tool Registry — All Phase 1 Components | `governance/registry/` |
| G | SharePoint List Schemas | `schemas/sharepoint/` |
| H | Power BI Dashboard Specifications | `dashboards/` |

---

*This is a draft application. Final submission requires TBS S-6 ISSM review, legal review, and Authorizing Official coordination. Format and content requirements may vary by command — adapt to local IATT template if one exists.*
