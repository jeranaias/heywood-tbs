# Heywood — Digital Staff Officer for The Basic School

> *"Good Morning, Sir. Here's your brief for March 6, 2026."*

Every morning at TBS, staff officers spend 30-60 minutes piecing together the day — pulling student data from spreadsheets, checking the training schedule, reviewing at-risk lists, scanning email for updates. Heywood does it in 3 seconds.

Heywood is an AI-powered Digital Staff Officer that gives every Marine at The Basic School — from the XO to individual students — instant, role-appropriate access to the data they need. No spreadsheet hunting. No email chains. Just ask.

**[Live Demo](https://heywood-tbs.nicefield-9a8db973.eastus.azurecontainerapps.io/)** — try it now, no login required

---

## One Question, Full Situational Awareness

Ask Heywood for a morning brief and it pulls from every data source simultaneously — student performance, today's schedule, qualification gaps, weather, news, and your calendar — then delivers it in seconds, tailored to your role.

![XO Morning Brief](screenshots/chat-morning-brief.png)

*The XO gets the full picture: schedule with travel advisories, weather and uniform call, at-risk students ranked by severity, qualification gaps by type, and proactive recommendations — all from a single question.*

---

## Role-Adaptive Views

Same app, different experience. Heywood automatically adjusts what you see based on who you are.

| Role | What They See |
|------|--------------|
| **Executive Officer** | Battalion-wide analytics, all companies, master calendar, full at-risk visibility |
| **Staff Officer** | TBS-wide data access, training schedule, instructor qualifications |
| **SPC (Company)** | Company-scoped student roster, company training events |
| **Student** | Personal record only — grades, schedule, upcoming quals |

---

## What's Inside

### Student Performance Dashboard
200-student roster with composite scoring across Academics, Military Skills, and Leadership. Class rank, phase tracking, trend indicators, and status flags (Active, At Risk, Academic Hold, Medical Hold).

![Dashboard](screenshots/dashboard.png)

### At-Risk Identification
Automatic flagging when any pillar drops below 75 or trends negative. 97 students flagged in the demo dataset — sortable, searchable, with direct drill-down to individual records.

![At-Risk Students](screenshots/at-risk.png)

### AI Chat with Tool Use
Heywood doesn't just generate text — it queries live data. Ask "How is 2ndLt Perez doing?" and it pulls the actual student record. Ask "What's on my schedule?" and it checks the calendar API. 12 tools available:

- `lookup_student` / `search_students` — find and analyze student data
- `get_at_risk` / `get_student_stats` — performance analytics
- `get_schedule` / `lookup_calendar` — training and personal schedule
- `get_qual_records` / `get_qual_stats` — instructor qualification tracking
- `web_search` — doctrine references, regulations, current info
- `create_task` — turn conversation into tracked action items

### Student Roster
Full searchable, sortable, paginated roster. Filter by phase, at-risk status, or search by name/ID. Click any student for detailed performance breakdown.

![Students](screenshots/students.png)

### Admin Settings
Web-based configuration — no CLI, no config files. Data source selection, AI provider status, Outlook integration toggle, database connections. Guided setup wizard with progress tracking.

![Settings](screenshots/settings.png)

### Additional Pages
- **Instructor Quals** — 12 TBS-specific qualifications tracked per instructor. Expiration monitoring, coverage gap analysis, readiness percentages.
- **Training Schedule** — Full TBS schedule with event types, locations, lead instructors, graded/ungraded status.
- **My Calendar** — Personal + master calendar overlay. Demo mode generates realistic mock events.
- **Task Inbox** — AI-generated tasks from conversation, with priority, status, and assignment tracking.

---

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                    React SPA (Vite)                   │
│   11 Pages · Tailwind CSS · Role-Adaptive Routing    │
├─────────────────────────────────────────────────────┤
│               Go HTTP Server (net/http)               │
│   35+ REST Endpoints · Auth Middleware · FIPS 140-3   │
├───────────┬───────────┬────────────┬────────────────┤
│ DataStore │  AI Chat  │  MS Graph  │  Auth Provider  │
│ Interface │  Service  │  Client    │  Interface      │
├───────────┼───────────┼────────────┼────────────────┤
│ JSON      │ OpenAI    │ Calendar   │ Demo (cookie)   │
│ SQLite    │ Azure     │ Mail       │ CAC/PKI (x509)  │
│ PostgreSQL│ OpenAI    │ SharePoint │                  │
│ Excel     │           │ Teams      │                  │
└───────────┴───────────┴────────────┴────────────────┘
```

**Frontend:** React 18, TypeScript, Tailwind CSS, Vite
**Backend:** Go 1.24, stdlib net/http, FIPS 140-3 native crypto (no BoringCrypto/CGO)
**AI:** OpenAI GPT-4o or Azure OpenAI — auto-detected from environment
**Database:** JSON (demo) / SQLite (recommended) / PostgreSQL (production)
**Microsoft 365:** Graph API — Outlook, SharePoint, Teams (commercial + GCC High + DoD clouds)
**Auth:** Demo role picker or CAC/PKI via X.509 client certificates
**Deployment:** Multi-stage Docker, Azure Container Apps, Azure Gov IL5 ready

---

## Data Sources

Units track data differently. Heywood adapts to them — not the other way around.

| Source | Best For | How It Works |
|--------|----------|-------------|
| **JSON** | Demo, dev | Default — ships with 200 synthetic students |
| **SQLite** | Single-server | Zero config, file-based, recommended for most units |
| **PostgreSQL** | MCEN cloud | Connection string, production-grade with pooling |
| **Excel (.xlsx)** | Transitioning units | Upload via admin page, auto-maps column headers |
| **SharePoint Lists** | Units already on SP | Connect via Graph API, maps lists to data types |

All backends implement the same 27-method `DataStore` interface. Application code never knows which one is active. **Hybrid mode** lets units keep reference data in Excel while mutable data (tasks, messages) lives in SQLite.

---

## Microsoft 365 Integration

One set of Graph API credentials unlocks:

- **Outlook Calendar** — personal + shared master calendar, merged view
- **Outlook Mail** — unread count, recent message summaries
- **SharePoint** — site discovery, list browsing, document libraries, file access
- **Teams** — team listing, channels, shared files

Supports **commercial Azure**, **GCC High**, and **DoD** national cloud endpoints. Client credentials flow with `Sites.Selected` permission scope (passes MCEN security review).

---

## Authentication

| Mode | Mechanism | Use Case |
|------|-----------|----------|
| **Demo** | Cookie-based role picker | Evaluation, training, development |
| **CAC/PKI** | X.509 cert → EDIPI → role lookup | Production on MCEN |

Set `AUTH_MODE=cac` and provide a `user-roster.json` mapping EDIPIs to roles. Works behind Azure App Service with client certificate forwarding.

---

## Quick Start

```bash
cd app
go mod download && cd web && npm ci && npm run build && cd ..
go build -o heywood ./cmd/server
OPENAI_API_KEY=sk-... ./heywood -dev -port 8080
```

Open `http://localhost:8080`. Pick a role. Ask for a morning brief.

### Docker

```bash
cd app
docker build -t heywood-tbs .
docker run -p 8080:8080 -e OPENAI_API_KEY=sk-... heywood-tbs
```

---

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `OPENAI_API_KEY` | One AI provider | OpenAI API key |
| `AZURE_OPENAI_ENDPOINT` | One AI provider | Azure OpenAI endpoint |
| `AZURE_OPENAI_KEY` | With Azure | Azure OpenAI key |
| `AZURE_OPENAI_DEPLOYMENT` | With Azure | Model deployment name |
| `AUTH_MODE` | No | `cac` for CAC/PKI, omit for demo |
| `GRAPH_TENANT_ID` | For M365 | Azure AD tenant ID |
| `GRAPH_CLIENT_ID` | For M365 | App registration client ID |
| `GRAPH_CLIENT_SECRET` | For M365 | Client secret |
| `GRAPH_CLOUD` | No | `commercial` / `gcc-high` / `dod` |
| `GRAPH_MASTER_CALENDAR_ID` | No | Shared calendar for TBS-wide events |
| `SEARXNG_URL` | No | SearXNG URL for web search |

---

## By the Numbers

- **~8,400 lines Go** across 8 packages
- **~4,400 lines TypeScript/React** across 11 pages
- **35+ REST API endpoints**
- **27-method DataStore interface** with 5 backend implementations
- **12 AI tools** for conversational data access
- **4 Microsoft Graph integrations** (Calendar, Mail, SharePoint, Teams)
- **3 cloud endpoints** supported (Commercial, GCC High, DoD)
- **2 auth modes** (Demo + CAC/PKI)
- **FIPS 140-3** compliant crypto, zero CGO dependencies
- **Single Docker image** for commercial + Azure Gov

---

## Foundation

Built on [Expert-Driven Development (EDD)](https://github.com/jeranaias/expertdrivendevelopment) — a 5-course AI training curriculum with 51 prompt templates and governance SOP for responsible DoD AI adoption.

---

**Classification:** UNCLASSIFIED // Distribution Unlimited

Do not include classified, CUI, PII, or operationally sensitive information in this repository. All student data is synthetic.
