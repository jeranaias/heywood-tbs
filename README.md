# Heywood — TBS Digital Staff Officer

**Classification:** UNCLASSIFIED // Distribution Unlimited

A role-adaptive AI agent for The Basic School (TBS), Quantico. Heywood functions as a Digital Staff Officer — surfacing student performance data, managing tasks, syncing calendars and mail, and answering doctrine questions through natural conversation. Built on platforms authorized inside MCEN.

**Live Demo:** [heywood-tbs.nicefield-9a8db973.eastus.azurecontainerapps.io](https://heywood-tbs.nicefield-9a8db973.eastus.azurecontainerapps.io/)

## What It Does

| Capability | Description |
|-----------|-------------|
| **AI Chat** | Conversational access to student data, schedules, quals, and doctrine. Tool-use enabled — Heywood can look up students, check calendars, search the web, and create tasks from natural language. |
| **Role-Based Views** | XO sees battalion-wide analytics. Staff sees company dashboards. Students see their own record. All from the same app. |
| **Student Analytics** | 200-student roster with GPA, PFT/CFT, rifle qual, conduct, peer evals. At-risk identification with configurable thresholds. |
| **Task Management** | Heywood creates actionable tasks from conversation ("Draft counseling for Lt Smith" becomes a tracked task). Priority, status, assignment. |
| **Calendar & Mail** | Microsoft Outlook integration via Graph API. Personal calendar + master TBS schedule overlay. Mail summaries with unread count. Demo mode generates mock events. |
| **Instructor Quals** | Qualification tracking across 12 TBS-specific quals. Expiration monitoring, readiness percentages, gap identification. |
| **Settings & Config** | Web-based admin panel for data sources, AI provider, Outlook sync, and database connections. No CLI required. |

## Architecture

```
┌─────────────────────────────────────────────────┐
│                   React SPA                      │
│  Dashboard · Chat · Students · Calendar · Tasks  │
│  At-Risk · Quals · Schedule · Settings           │
├─────────────────────────────────────────────────┤
│              Go HTTP Server (stdlib)              │
│  REST API · Auth Middleware · Role Filtering      │
├──────────┬──────────┬───────────┬───────────────┤
│ DataStore│ AI/Chat  │ MS Graph  │ Auth Provider  │
│ Interface│ Service  │ Client    │ Interface      │
├──────────┼──────────┼───────────┼───────────────┤
│ JSON     │ OpenAI   │ Outlook   │ Demo (cookie)  │
│ SQLite   │ Azure    │ SharePoint│ CAC/PKI (x509) │
│ Postgres │ OpenAI   │ Teams     │                │
│ Excel    │          │           │                │
└──────────┴──────────┴───────────┴───────────────┘
```

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Frontend** | React 18, TypeScript, Tailwind CSS, Vite |
| **Backend** | Go 1.24, net/http (stdlib), FIPS 140-3 crypto |
| **AI** | OpenAI GPT-4o or Azure OpenAI (auto-detected from env) |
| **Database** | JSON (demo), SQLite (recommended), PostgreSQL (production) |
| **Microsoft 365** | Graph API — Outlook calendar/mail, SharePoint, Teams |
| **Auth** | Demo role picker or CAC/PKI (X.509 client certs) |
| **Deployment** | Docker multi-stage, Azure Container Apps, Azure Gov IL5 ready |
| **Search** | SearXNG (self-hosted, optional) |

## Project Structure

```
heywood-tbs/
├── app/
│   ├── cmd/server/          # Entry point
│   ├── internal/
│   │   ├── ai/              # Chat service, tool definitions, weather/news/traffic
│   │   ├── api/             # REST handlers, router, settings, graph endpoints
│   │   ├── auth/            # IdentityProvider interface, Demo + CAC providers
│   │   ├── calendar/        # CalendarProvider interface, Outlook + Mock impls
│   │   ├── data/            # DataStore interface, JSON/SQL/Excel/Hybrid stores
│   │   ├── middleware/      # Auth, CORS, security headers
│   │   ├── models/          # Shared types (Student, Task, CalendarEvent, etc.)
│   │   └── msgraph/         # Microsoft Graph client (OAuth2, calendar, mail, SP, Teams)
│   ├── web/
│   │   └── src/
│   │       ├── pages/       # 11 page components
│   │       ├── components/  # Layout, sidebar, charts
│   │       └── lib/         # API client, types, utilities
│   ├── data/                # JSON seed data, settings, user roster
│   ├── Dockerfile           # Multi-stage: node → go → alpine
│   └── Makefile
├── docs/                    # Build plan, phase details, briefs
├── governance/              # PIA, compliance checklists, tool registry
├── prompts/                 # 20 TBS-adapted prompt templates
├── training/                # EDD course materials
└── infrastructure/          # Azure deployment configs
```

## Data Sources

Heywood reads student/instructor data from configurable backends. All implement the same 27-method `DataStore` interface — application code never knows which backend is active.

| Source | Best For | Setup |
|--------|----------|-------|
| **JSON Files** | Demo, development | Default — no config needed |
| **SQLite** | Single-server production | Recommended — zero infrastructure, file-based |
| **PostgreSQL** | Multi-server / MCEN cloud | Connection string in settings |
| **Excel (.xlsx)** | Units transitioning from spreadsheets | Upload via admin page, auto-maps columns |
| **Hybrid** | Real units | Reference data from Excel/SP, mutable data in SQLite |

**Mutable data** (tasks, messages, notifications) is isolated per backend. In demo mode, mutable data is in-memory only and resets on restart.

## Microsoft 365 Integration

When configured with Graph API credentials, Heywood connects to:

- **Outlook Calendar** — Personal events + shared master calendar, role-filtered views
- **Outlook Mail** — Unread count badge, recent message summaries
- **SharePoint** — Site discovery, list browsing, document library file access
- **Microsoft Teams** — Team listing, channel browsing, shared file access

Supports commercial Azure, GCC High, and DoD national cloud endpoints. Uses client credentials flow (app-only, `Sites.Selected` permission scope).

Set environment variables: `GRAPH_TENANT_ID`, `GRAPH_CLIENT_ID`, `GRAPH_CLIENT_SECRET`, `GRAPH_CLOUD` (commercial/gcc-high/dod).

## Authentication

| Mode | How | Set By |
|------|-----|--------|
| **Demo** | Cookie-based role picker (XO / Staff / SPC / Student) | Default |
| **CAC/PKI** | X.509 client cert via `X-ARR-ClientCert` header → EDIPI extraction → role lookup from `user-roster.json` | `AUTH_MODE=cac` |

## Quick Start

```bash
# Clone and build
cd app
go mod download
cd web && npm ci && npm run build && cd ..
go build -o heywood ./cmd/server

# Run (demo mode — no config needed)
./heywood -dev -port 8080

# With AI (set one):
OPENAI_API_KEY=sk-... ./heywood -dev
# or
AZURE_OPENAI_ENDPOINT=https://... AZURE_OPENAI_KEY=... AZURE_OPENAI_DEPLOYMENT=gpt-4o ./heywood -dev
```

Open `http://localhost:8080`. Pick a role. Start talking to Heywood.

## Docker

```bash
cd app
docker build -t heywood-tbs .
docker run -p 8080:8080 \
  -e OPENAI_API_KEY=sk-... \
  heywood-tbs
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `OPENAI_API_KEY` | One AI provider | OpenAI API key |
| `AZURE_OPENAI_ENDPOINT` | One AI provider | Azure OpenAI endpoint URL |
| `AZURE_OPENAI_KEY` | One AI provider | Azure OpenAI API key |
| `AZURE_OPENAI_DEPLOYMENT` | With Azure | Deployment name (e.g., gpt-4o) |
| `AUTH_MODE` | No | `cac` for CAC/PKI auth, omit for demo |
| `GRAPH_TENANT_ID` | For M365 | Azure AD tenant ID |
| `GRAPH_CLIENT_ID` | For M365 | App registration client ID |
| `GRAPH_CLIENT_SECRET` | For M365 | App registration client secret |
| `GRAPH_CLOUD` | No | `commercial` (default), `gcc-high`, `dod` |
| `GRAPH_MASTER_CALENDAR_ID` | No | Shared calendar ID for TBS-wide events |
| `SEARXNG_URL` | No | SearXNG instance URL for web search |

## Codebase

- **~8,400 lines Go** across 8 packages
- **~4,400 lines TypeScript/React** across 11 pages
- **27-method DataStore interface** with 5 backend implementations
- **12 AI tools** for conversational data access
- **FIPS 140-3** native crypto (Go 1.24, no BoringCrypto/CGO)

## Foundation

Built on [Expert-Driven Development (EDD)](https://github.com/jeranaias/expertdrivendevelopment) — a 5-course AI training curriculum with 51 prompt templates, governance SOP, and reusable templates for DoD AI adoption.

## Authorization

- **Demo/Dev:** No ATO required (runs on any machine)
- **MCEN Deployment:** Azure Container Apps on Azure Gov, inherits FedRAMP High baseline. IATT for Azure OpenAI custom connector. Full ATO/cATO path documented.

## Data Handling

- **CUI:** Authorized on Azure OpenAI (IL5)
- **PII:** Minimized by design, PIA required at each phase gate
- **Classified:** Never authorized on any Heywood component

---

**Do not include classified, CUI, PII, or operationally sensitive information in this repository.**
