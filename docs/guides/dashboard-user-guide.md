# Dashboard User Guide: Power BI Dashboards

*Heywood Phase 1 — For Staff, SPCs, and Students*

---

## Accessing the Dashboards

1. Open your browser on any MCEN workstation
2. Navigate to **app.powerbigov.us** (Power BI for Government)
3. Authenticate with your **CAC**
4. Open the **Heywood TBS** app from your Apps list
5. You will see only the dashboards authorized for your role

### What You'll See by Role

| Role | Student Performance Dashboard | Instructor Quals Dashboard |
|------|:---:|:---:|
| **Staff** (S-3, Ops) | All companies, all students | Full access |
| **SPC** | Your company only | No access |
| **Student** | Your own record only | No access |

---

## Dashboard 1: Student Performance

### Understanding the TBS Grading System

TBS evaluates students across three pillars with these weights:

```
┌──────────────────────────────────────────────────────┐
│                  OVERALL COMPOSITE                    │
│                                                      │
│  ┌──────────┐  ┌──────────────┐  ┌────────────────┐ │
│  │ACADEMICS │  │ MIL SKILLS   │  │  LEADERSHIP    │ │
│  │   32%    │  │    32%       │  │     36%        │ │
│  │          │  │              │  │                │ │
│  │ 4 Exams  │  │ PFT    CFT  │  │ Week 12 (14%)  │ │
│  │ Quiz Avg │  │ Rifle  Land │  │ Week 22 (22%)  │ │
│  │          │  │ Pistol Nav  │  │                │ │
│  │          │  │ Obstacle    │  │ SPC Eval: 90%  │ │
│  │          │  │ Endurance   │  │ Peer Eval: 10% │ │
│  └──────────┘  └──────────────┘  └────────────────┘ │
└──────────────────────────────────────────────────────┘
```

**A score below 75% in any pillar triggers an at-risk flag.**

---

### Page 1: Company Overview

*Default view for Staff and SPCs*

**Top row — Key metrics:**
- **Active Students:** Total students currently in training
- **Company Avg Composite:** Average overall score (target: above 80)
- **At Risk Count:** Students flagged below 75% in any pillar (red = action needed)
- **At Risk %:** Percentage of company at risk (red if above 15%)

**Middle section:**
- **Class Standing Chart:** How many students are in the top, middle, and bottom third — by company (Staff view) or your company only (SPC view)
- **Phase Distribution:** Donut chart showing how many students are in each training phase

**Bottom table:** Student roster with all composite scores and rankings
- Click any column header to sort
- Red-highlighted rows = at-risk students
- Gold text on rank 1-3 = company top performers
- Click any student row to **drill through** to their individual detail (Page 2)

**Slicers (filters) on the right:**
- Company: filter to one company (Staff) or pre-filtered (SPC)
- Class Number: switch between classes
- Phase: show only students in a specific phase
- At Risk Flag: filter to show only at-risk students

---

### Page 2: Individual Student Detail

*Drill-through page — click a student from Page 1*

**Top row:** Student info (rank, company, platoon, SPC, current phase)

**Three gauges:**
| Gauge | What It Shows | Color Zones |
|-------|---------------|-------------|
| Academic Composite | Average of exams + quiz avg | Green: 85+ / Yellow: 75-84 / Red: below 75 |
| MilSkills Composite | Weighted average of PFT, CFT, weapons, land nav, courses | Same zones |
| Leadership Composite | SPC eval (90%) + peer eval (10%), weighted by week | Same zones |

**Exam breakdown:** Bar chart showing all four exam scores and quiz average side by side. Look for trends — are scores improving or declining?

**Military skills table:** Shows pass/fail status and scores for every mil skills event.

**Trend line:** Exam scores plotted over time (Exam 1 → 4). An arrow indicates:
- **Green up arrow:** Improving
- **Red down arrow:** Declining (may warrant counseling)
- **Gray right arrow:** Stable

**Overall composite KPI:** Large number at bottom right with trend direction.

---

### Page 3: At-Risk Students

*Staff and SPCs — students needing intervention*

**What "at-risk" means:**

| Flag | Meaning | Action |
|------|---------|--------|
| Academic (<75%) | Below 75% in academic composite | Schedule remedial study, review exam performance |
| MilSkills (<75%) | Below 75% in military skills | Identify specific weak events, schedule extra training |
| Leadership (<75%) | Below 75% in leadership composite | SPC counseling, peer mentoring |
| Multiple (<75%) | Below 75% in more than one pillar | Priority case — SPC + company commander awareness |
| Declining Trend | Scores dropping across consecutive evaluations | Early intervention before threshold is crossed |

**The table shows:** Each at-risk student with their flag type, overall composite, all three pillar scores, and assigned SPC.

**Bar chart at bottom:** At-risk count by company — helps Staff identify which companies need additional support.

**When to take action:**
- **Any student flagged:** SPC should schedule a counseling session within 5 training days
- **Multiple flags:** Notify company commander, document intervention plan
- **Declining trend (no flag yet):** Preventive counseling — address before the student crosses the threshold

---

### Page 4: Phase Progression (Staff Only)

**Funnel chart:** Shows student count decreasing from Phase I → Complete. Large drops between phases may indicate systematic issues.

**Matrix:** Company × Phase grid. Each cell shows the count of students in that phase for that company. Helps identify if one company is behind schedule.

**Status table:** Students not in "Active" status:
- Medical Hold (Mike Co): student transferred to Mike Company for medical recovery
- Academic Hold: student paused for remedial academics
- Dropped: student removed from training

---

## Dashboard 2: Instructor Qualification Tracking

*Staff only*

### Page 1: Qualification Status Overview

**Top row — Alert cards:**
| Card | What It Means | Color |
|------|---------------|-------|
| Expired | Quals past their expiration date — **immediate action** | Red border |
| Expiring ≤30 days | Critical — schedule renewal now | Orange border |
| Expiring ≤60 days | Warning — begin renewal coordination | Yellow border |
| Active Instructors | Total instructors on roster | Green border |

**Stacked bar chart:** Shows qualification status breakdown by category (Range Safety, Weapons, Tactics, etc.). Tall red/orange bars = categories with coverage problems.

**Bottom table:** Every qualification record sorted by soonest expiration. Color-coded:
- **Red row:** Expired
- **Red text:** ≤30 days
- **Orange text:** ≤60 days
- **Yellow text:** ≤90 days
- **Green text:** Current (>90 days)

---

### Page 2: Instructor Workload

**Key metric:** Average students per SPC and workload standard deviation.
- StdDev > 5 means an imbalance — some SPCs have significantly more students than others.

**Bar chart:** Students assigned per instructor, grouped by company. Look for outliers — an SPC with 60 students when the average is 45 needs load balancing.

**Event load chart:** Events per instructor this month — top 20 busiest. Are some instructors overcommitted?

**Table columns to watch:**
- **CounselingsOverdue > 0:** Red flag — students aren't getting required counseling sessions
- **EventsThisMonth very high:** Instructor may be spread too thin

---

### Page 3: Coverage Gaps & Rotation Risk

**Qualification gaps table:** Lists qualifications where the number of current holders is below the minimum required for training events. Red rows = you can't execute that event until someone gets qualified.

**Rotating instructors table:** Instructors with PRD (Projected Rotation Date) within 90 days. Shows what quals they hold — if they're the only holder of a critical qual, you need a replacement trained before they depart.

**Coverage matrix:** Company × Qual Category grid.
- **Red cell:** Below minimum required
- **Yellow cell:** At minimum (no redundancy)
- **Green cell:** Above minimum (healthy)

---

### Page 4: Individual Instructor Detail (Drill-Through)

Click any instructor from Pages 1-3 to see their full record:
- All qualifications held with earn/expiration dates
- Timeline visual showing when each qual was earned and when it expires
- Current student load and event assignments

---

## Taking Action on Dashboard Data

### For SPCs

| You See | You Do |
|---------|--------|
| Student in red (below 75%) | Open Prompt #2 (Counseling Prep), schedule session within 5 days |
| Student trend declining | Counseling session, identify root cause, document intervention |
| Multiple students struggling on same event | Use Prompt #15 (AAR Synthesis) to identify patterns, report to company commander |
| High counseling overdue count for your platoon | Block time this week, prioritize students closest to threshold |

### For Staff

| You See | You Do |
|---------|--------|
| Company with high at-risk % | Meet with company commander, review SPC workload |
| Qual expiring in <30 days | Coordinate renewal course with issuing authority |
| Coverage gap on critical qual | Identify candidates, schedule qualification course |
| Instructor overloaded | Rebalance student assignments across SPCs |
| Company behind in phase progression | Review training schedule for conflicts, coordinate with S-3 |

---

## Refreshing Data

- **Student Performance:** Data refreshes every **4 hours** during the duty day (0600-2200)
- **Instructor Quals:** Data refreshes **daily at 0500**
- If you need an immediate refresh, contact the Heywood administrator
- Data comes from SharePoint lists — updates to SharePoint appear in dashboards at the next scheduled refresh

---

## Mobile Access

Both dashboards have mobile-optimized layouts for tablets. SPCs can access during field events:

1. Open the Power BI mobile app (available on MCEN-managed devices)
2. Sign in with CAC
3. Navigate to the Heywood TBS app
4. The mobile layout shows key metrics and the student table in a vertical format

---

*Questions or issues? Contact SSgt Morgan or your Heywood POC.*
