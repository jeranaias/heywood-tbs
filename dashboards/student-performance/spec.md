# Student Performance Dashboard — Power BI Specification

*Heywood Phase 1 | Data Source: SharePoint `StudentScores` list*

---

## Overview

Primary dashboard for tracking student performance across TBS's three-pillar grading system. Three role-specific views enforce need-to-know via Row-Level Security.

**Users:** Staff (full access), SPC (company-filtered), Student (individual record only)

---

## Data Model

### Source
- SharePoint Online list: `StudentScores`
- Connection: Power BI SharePoint Online connector (DirectQuery or Import)
- Refresh: Import mode with scheduled refresh every 4 hours during duty day

### Relationships
| From | To | Type | Key |
|------|----|------|-----|
| StudentScores.Company | DimCompany.Company | Many-to-One | Company |
| StudentScores.ClassNumber | DimClass.ClassNumber | Many-to-One | ClassNumber |
| StudentScores.CurrentPhase | DimPhase.Phase | Many-to-One | CurrentPhase |
| UserSecurity.UserEmail | Azure AD Email | RLS filter | UserEmail |

### Dimension Tables (created in Power BI)

**DimCompany**
```
Company = DATATABLE(
    "Company", STRING,
    "CompanyOrder", INTEGER,
    {
        {"Alpha", 1}, {"Bravo", 2}, {"Charlie", 3}, {"Delta", 4},
        {"Echo", 5}, {"Foxtrot", 6}, {"Golf", 7}, {"India", 8}, {"Mike", 9}
    }
)
```

**DimPhase**
```
Phase = DATATABLE(
    "Phase", STRING,
    "PhaseOrder", INTEGER,
    "PhaseShort", STRING,
    {
        {"Phase I - Individual Skills", 1, "Ph I"},
        {"Phase II - Squad", 2, "Ph II"},
        {"Phase III - Platoon", 3, "Ph III"},
        {"Phase IV - MAGTF", 4, "Ph IV"},
        {"Complete", 5, "Complete"}
    }
)
```

---

## DAX Measures

### Composite Calculations (mirror TBS grading policy)

```dax
// Academic Composite — average of completed exams and quiz average
Academic Composite =
VAR Exam1 = IF(ISBLANK([AcademicExam1]), BLANK(), [AcademicExam1])
VAR Exam2 = IF(ISBLANK([AcademicExam2]), BLANK(), [AcademicExam2])
VAR Exam3 = IF(ISBLANK([AcademicExam3]), BLANK(), [AcademicExam3])
VAR Exam4 = IF(ISBLANK([AcademicExam4]), BLANK(), [AcademicExam4])
VAR QuizAvg = IF(ISBLANK([AcademicQuizAvg]), BLANK(), [AcademicQuizAvg])
VAR AllScores = {Exam1, Exam2, Exam3, Exam4, QuizAvg}
VAR NonBlankScores = FILTER(AllScores, NOT(ISBLANK([Value])))
VAR ScoreCount = COUNTROWS(NonBlankScores)
RETURN
    IF(ScoreCount = 0, BLANK(), SUMX(NonBlankScores, [Value]) / ScoreCount)
```

```dax
// Military Skills Composite — weighted composite of all mil skills events
// PFT and CFT normalized to 0-100 scale (divide by 3)
MilSkills Composite =
VAR PFTNorm = DIVIDE([PFTScore], 3, BLANK())
VAR CFTNorm = DIVIDE([CFTScore], 3, BLANK())
VAR RifleScore = SWITCH([RifleQual], "Expert", 100, "Sharpshooter", 85, "Marksman", 70, "Unqualified", 0, BLANK())
VAR PistolScore = SWITCH([PistolQual], "Expert", 100, "Sharpshooter", 85, "Marksman", 70, "Unqualified", 0, BLANK())
VAR LandNavDayScore = SWITCH([LandNavDay], "Pass", 100, "Fail", 0, BLANK())
VAR LandNavNightScore = SWITCH([LandNavNight], "Pass", 100, "Fail", 0, BLANK())
VAR ObstacleScore = SWITCH([ObstacleCourse], "Pass", 100, "Fail", 0, BLANK())
VAR EnduranceScore = SWITCH([EnduranceCourse], "Pass", 100, "Fail", 0, BLANK())
VAR AllSkills = {PFTNorm, CFTNorm, RifleScore, PistolScore, LandNavDayScore, LandNavNightScore, [LandNavWritten], ObstacleScore, EnduranceScore}
VAR NonBlank = FILTER(AllSkills, NOT(ISBLANK([Value])))
VAR SkillCount = COUNTROWS(NonBlank)
RETURN
    IF(SkillCount = 0, BLANK(), SUMX(NonBlank, [Value]) / SkillCount)
```

```dax
// Leadership Composite — Week 12 (14% of total) + Week 22 (22% of total)
// Within each: SPC evaluation = 90%, Peer evaluation = 10%
Leadership Composite =
VAR Wk12SPC = [LeadershipWeek12]
VAR Wk12Peer = [PeerEvalWeek12]
VAR Wk22SPC = [LeadershipWeek22]
VAR Wk22Peer = [PeerEvalWeek22]
VAR Wk12Score =
    IF(
        NOT(ISBLANK(Wk12SPC)),
        IF(NOT(ISBLANK(Wk12Peer)),
            Wk12SPC * 0.9 + Wk12Peer * 0.1,
            Wk12SPC
        ),
        BLANK()
    )
VAR Wk22Score =
    IF(
        NOT(ISBLANK(Wk22SPC)),
        IF(NOT(ISBLANK(Wk22Peer)),
            Wk22SPC * 0.9 + Wk22Peer * 0.1,
            Wk22SPC
        ),
        BLANK()
    )
// Weight: Week 12 = 14/36 of leadership, Week 22 = 22/36 of leadership
RETURN
    SWITCH(
        TRUE(),
        NOT(ISBLANK(Wk12Score)) && NOT(ISBLANK(Wk22Score)),
            Wk12Score * (14/36) + Wk22Score * (22/36),
        NOT(ISBLANK(Wk12Score)),
            Wk12Score,
        NOT(ISBLANK(Wk22Score)),
            Wk22Score,
        BLANK()
    )
```

```dax
// Overall Composite — Academics 32% + MilSkills 32% + Leadership 36%
Overall Composite =
VAR Acad = [Academic Composite]
VAR Mil = [MilSkills Composite]
VAR Lead = [Leadership Composite]
VAR HasAcad = NOT(ISBLANK(Acad))
VAR HasMil = NOT(ISBLANK(Mil))
VAR HasLead = NOT(ISBLANK(Lead))
VAR TotalWeight = IF(HasAcad, 0.32, 0) + IF(HasMil, 0.32, 0) + IF(HasLead, 0.36, 0)
VAR WeightedSum = IF(HasAcad, Acad * 0.32, 0) + IF(HasMil, Mil * 0.32, 0) + IF(HasLead, Lead * 0.36, 0)
RETURN
    IF(TotalWeight = 0, BLANK(), WeightedSum / TotalWeight * 100 / 100)
```

### Aggregation Measures

```dax
// Company average composite
Company Avg Composite =
CALCULATE(
    AVERAGE(StudentScores[OverallComposite]),
    StudentScores[Status] = "Active"
)
```

```dax
// Count of active students
Active Student Count =
COUNTROWS(
    FILTER(StudentScores, StudentScores[Status] = "Active")
)
```

```dax
// At-risk student count
At Risk Count =
COUNTROWS(
    FILTER(
        StudentScores,
        StudentScores[AtRiskFlag] <> "None" && StudentScores[Status] = "Active"
    )
)
```

```dax
// At-risk percentage
At Risk Pct =
DIVIDE([At Risk Count], [Active Student Count], 0)
```

```dax
// Phase progression — students in each phase
Phase Count =
COUNTROWS(
    FILTER(StudentScores, StudentScores[Status] = "Active")
)
// Use with DimPhase slicer for per-phase breakdown
```

```dax
// Company rank calculation (for each student within their company)
Company Rank =
VAR CurrentStudent = StudentScores[StudentEDIPI]
VAR CurrentCompany = StudentScores[Company]
VAR CurrentScore = StudentScores[OverallComposite]
RETURN
    COUNTROWS(
        FILTER(
            StudentScores,
            StudentScores[Company] = CurrentCompany
            && StudentScores[Status] = "Active"
            && StudentScores[OverallComposite] > CurrentScore
        )
    ) + 1
```

```dax
// Class standing third
Class Standing Third =
VAR TotalActive = [Active Student Count]
VAR StudentRank = [Company Rank]
RETURN
    SWITCH(
        TRUE(),
        StudentRank <= ROUNDUP(TotalActive / 3, 0), "Top Third",
        StudentRank <= ROUNDUP(TotalActive * 2 / 3, 0), "Middle Third",
        "Bottom Third"
    )
```

### Trend Measures

```dax
// Academic trend (exam score progression)
Academic Trend Direction =
VAR E1 = [AcademicExam1]
VAR E2 = [AcademicExam2]
VAR E3 = [AcademicExam3]
VAR E4 = [AcademicExam4]
VAR Latest = COALESCE(E4, E3, E2, E1)
VAR Previous =
    SWITCH(
        TRUE(),
        NOT(ISBLANK(E4)), COALESCE(E3, E2, E1),
        NOT(ISBLANK(E3)), COALESCE(E2, E1),
        NOT(ISBLANK(E2)), E1,
        BLANK()
    )
RETURN
    IF(
        ISBLANK(Latest) || ISBLANK(Previous),
        "N/A",
        IF(Latest > Previous, "↑ Improving",
            IF(Latest < Previous, "↓ Declining", "→ Stable"))
    )
```

---

## Visual Layout

### Page 1: Company Overview (Staff + SPC default view)

| Position | Visual | Data | Size |
|----------|--------|------|------|
| Top-left | Card | Active Student Count | 1/4 width |
| Top-center-left | Card | Company Avg Composite (formatted 0.0) | 1/4 width |
| Top-center-right | Card | At Risk Count (red conditional) | 1/4 width |
| Top-right | Card | At Risk Pct (red if >15%) | 1/4 width |
| Middle-left | Stacked bar chart | Students by ClassStandingThird (Top/Mid/Bottom) grouped by Company | 1/2 width |
| Middle-right | Donut chart | Students by CurrentPhase | 1/2 width |
| Bottom | Table | LastName, FirstName, Rank, Platoon, AcademicComposite, MilSkillsComposite, LeadershipComposite, OverallComposite, CompanyRank, AtRiskFlag | Full width |

**Slicers:** Company (dropdown), ClassNumber, CurrentPhase, AtRiskFlag
**Conditional formatting on table:**
- OverallComposite < 75: red background
- AtRiskFlag != "None": red text
- CompanyRank = 1-3: bold gold

### Page 2: Individual Student Detail (Drill-through from Company Overview)

| Position | Visual | Data |
|----------|--------|------|
| Top | Card row | StudentName, Rank, Company, Platoon, SPC, CurrentPhase |
| Left column | Gauge | Academic Composite (target: 80, max: 100) |
| Center column | Gauge | MilSkills Composite (target: 80, max: 100) |
| Right column | Gauge | Leadership Composite (target: 80, max: 100) |
| Middle | Clustered bar | All 4 exam scores + quiz avg side by side |
| Middle-right | Table | All mil skills events with Pass/Fail/Score |
| Bottom-left | Line chart | Exam scores over time (E1→E4 trend) |
| Bottom-right | KPI card | Overall Composite with trend arrow |

**Drill-through from:** Page 1 table → click student row

### Page 3: At-Risk Students (Staff + SPC)

| Position | Visual | Data |
|----------|--------|------|
| Top | Card | Total at-risk count, breakdown by flag type |
| Middle | Table | LastName, FirstName, Company, Platoon, AtRiskFlag, OverallComposite, Academic/Mil/Leadership composites, SPCAssigned |
| Bottom | Clustered bar | At-risk count by Company (identify which companies need attention) |

**Conditional formatting:**
- "Multiple (<75%)": red row highlight
- "Declining Trend": orange row highlight

### Page 4: Phase Progression (Staff only)

| Position | Visual | Data |
|----------|--------|------|
| Top | Funnel chart | Student count by phase (should decrease from Phase I → Complete) |
| Middle | Matrix | Company (rows) × Phase (columns) → student count in each cell |
| Bottom | Table | Students with Status != Active (Medical Hold, Academic Hold, Dropped) |

---

## Row-Level Security Configuration

### RLS Roles

**Role: Staff**
```dax
// No filter — sees all data
// DAX filter: TRUE()
```

**Role: SPC**
```dax
// Filter to SPC's company
[Company] = LOOKUPVALUE(
    UserSecurity[Company],
    UserSecurity[UserEmail],
    USERPRINCIPALNAME()
)
```

**Role: Student**
```dax
// Filter to student's own EDIPI
[StudentEDIPI] = LOOKUPVALUE(
    UserSecurity[EDIPI],
    UserSecurity[UserEmail],
    USERPRINCIPALNAME()
)
```

### UserSecurity Table
Maintain a SharePoint list or Excel table mapping Azure AD email → Role + Company + EDIPI.

| UserEmail | Role | Company | EDIPI |
|-----------|------|---------|-------|
| john.doe@usmc.mil | Staff | N/A | N/A |
| jane.smith@usmc.mil | SPC | Alpha | N/A |
| lt.jones@usmc.mil | Student | Alpha | 1234567890 |

---

## Conditional Formatting Rules

| Column/Visual | Condition | Format |
|---------------|-----------|--------|
| OverallComposite | < 75 | Red background, white text |
| OverallComposite | 75-84.99 | Yellow background |
| OverallComposite | >= 85 | Green background |
| AtRiskFlag | != "None" | Red bold text |
| CompanyRank | 1-3 | Gold bold text |
| Academic/Mil/Leadership Composite | < 75 | Red text |
| Gauge visuals | < 75 zone | Red |
| Gauge visuals | 75-84.99 zone | Yellow |
| Gauge visuals | >= 85 zone | Green |
| Trend arrow | Declining | Red down arrow |
| Trend arrow | Improving | Green up arrow |
| Trend arrow | Stable | Gray right arrow |

---

## Deployment Notes

1. Connect Power BI Desktop to SharePoint Online list via Get Data → SharePoint Online List
2. Import mode recommended for Phase 1 (small data volume, better DAX performance)
3. Schedule refresh: every 4 hours during 0600-2200 EST
4. Publish to Power BI Service workspace with RLS roles configured
5. Test RLS: "View as Role" in Power BI Service for each role + test user
6. Share via Power BI app (not direct workspace access) for clean user experience
7. Mobile layout: configure for each page (SPCs will use tablets during field events)
