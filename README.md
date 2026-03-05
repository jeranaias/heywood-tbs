# The Heywood Initiative — TBS AI Agent System

**Classification:** UNCLASSIFIED // Distribution Unlimited

An AI agent system for The Basic School (TBS), Quantico, functioning as a "Digital Staff Officer" for every Marine in the training pipeline. Built incrementally on platforms already authorized inside MCEN.

## Overview

Heywood delivers role-based AI capabilities (Staff / Instructor / Student) across 15 use cases identified in the TBS AI Vision 2026, from AAR synthesis and counseling preparation to predictive performance analytics and scenario generation.

## 4-Phase Build Approach

| Phase | Timeline | Cost | Key Deliverables |
|-------|----------|------|------------------|
| **Phase 1: Proof of Concept** | Months 0-3 | $0-$1,200 | GenAI.mil prompts, Power BI dashboards, EDD training |
| **Phase 2: Custom Agent** | Months 3-6 | $3,600-$4,650 | Power App with RBAC, Azure OpenAI RAG pipeline, BARS framework |
| **Phase 3: Scale** | Months 6-12 | $32,600-$49,800 | Dataverse, MCTIMS integration, predictive analytics, all 8 companies |
| **Phase 4: Full Heywood** | Months 12-18 | $41,200-$59,200 | MHS GENESIS, cross-cycle trends, scenario engine, field-deployable PWA |

**Total: $77K-$115K over 18 months** (3-5% of comparable programs)

Each phase delivers standalone value. If funding stops at any phase, TBS keeps everything built.

## Technology Stack

- **LLM:** Azure OpenAI on Azure Government (IL5) — GPT-4o, GPT-4.1
- **RAG:** Azure AI Search (vector store + hybrid search)
- **App:** Power Apps (Canvas + Model-Driven)
- **Automation:** Power Automate
- **Dashboards:** Power BI with Row-Level Security
- **Data:** Dataverse (Phase 3+), SharePoint (Phase 1-2)
- **Auth:** Azure Entra ID (CAC authentication, RBAC)
- **Monitoring:** Azure Sentinel (cATO compliance)

## Project Structure

```
heywood-tbs/
├── docs/
│   ├── plan/              # Full build plan and phase details
│   └── briefs/            # Leadership briefing materials
├── prompts/
│   ├── playbook/          # 20 TBS-adapted prompt templates (Phase 1)
│   └── custom/            # Custom prompts for TBS-specific use cases
├── governance/
│   ├── pia/               # Privacy Impact Assessment artifacts
│   ├── compliance/        # Compliance checklists per component
│   └── registry/          # Tool registry entries
├── training/
│   ├── courses/           # TBS-adapted EDD course materials
│   └── materials/         # Supplementary training materials
├── dashboards/
│   ├── student-performance/  # Power BI dashboard specs
│   └── instructor-quals/     # Instructor qualification dashboard specs
└── schemas/
    ├── sharepoint/        # SharePoint list schemas (Phase 1-2)
    └── dataverse/         # Dataverse entity models (Phase 3+)
```

## Foundation

Built on the [Expert-Driven Development (EDD)](https://github.com/jeranaias/expertdrivendevelopment) curriculum:
- 5-course training progression (AI Fluency → Builder → Platform → Advanced → Supervisor)
- 51 interactive prompt templates across 10 categories
- Governance SOP with 12 sections and 4-layer framework
- 7 reusable governance templates

## Authorization

- **Phase 1:** No ATO required (GenAI.mil + Power BI + SharePoint within MCEN boundary)
- **Phase 2:** IATT for Azure OpenAI + Power App custom connector
- **Phase 3-4:** Full ATO or cATO (inherits Azure Gov FedRAMP High baseline)

## Data Handling

- **CUI:** Authorized on GenAI.mil and Azure OpenAI (IL5)
- **PII:** Requires PIA at each phase gate; minimized by design
- **PHI:** MHS GENESIS integration only in Phase 4, conditional on PIA approval
- **Classified:** Never authorized on any Heywood component

---

**Reminder:** Do not include any classified, CUI, PII, or operationally sensitive information in this repository. All content must be UNCLASSIFIED.
