# Instructor Qualification Dashboard — Power BI Specification

*Heywood Phase 1 | Data Sources: SharePoint `Instructors`, `QualificationRecords`, `RequiredQualifications` lists*

---

## Overview

Tracks instructor qualification status, expiration alerts (30/60/90 day), and workload distribution across companies. Ensures TBS maintains readiness to execute training events with properly qualified personnel.

**Users:** Staff only (contains instructor PII and assignment data)

---

## Data Model

### Sources
- SharePoint Online lists: `Instructors`, `QualificationRecords`, `RequiredQualifications`
- Connection: Power BI SharePoint Online connector (Import mode)
- Refresh: Daily at 0500 EST

### Relationships

| From | To | Cardinality | Key |
|------|----|-------------|-----|
| QualificationRecords.InstructorEDIPI | Instructors.InstructorEDIPI | Many-to-One | InstructorEDIPI |
| QualificationRecords.QualCode | RequiredQualifications.QualCode | Many-to-One | QualCode |
| Instructors.CompanyAssigned | DimCompany.Company | Many-to-One | Company |

### Date Table (for time intelligence)

```dax
DateTable =
ADDCOLUMNS(
    CALENDAR(DATE(2025, 1, 1), DATE(2027, 12, 31)),
    "Year", YEAR([Date]),
    "Month", FORMAT([Date], "MMM YYYY"),
    "MonthNum", MONTH([Date]),
    "DaysFromToday", DATEDIFF(TODAY(), [Date], DAY)
)
```

---

## DAX Measures

### Expiration Tracking

```dax
// Days until qualification expires
Days Until Expiration =
DATEDIFF(TODAY(), QualificationRecords[ExpirationDate], DAY)
```

```dax
// Expiration status category
Expiration Status =
VAR DaysLeft = [Days Until Expiration]
RETURN
    SWITCH(
        TRUE(),
        DaysLeft < 0, "Expired",
        DaysLeft <= 30, "Critical (≤30 days)",
        DaysLeft <= 60, "Warning (≤60 days)",
        DaysLeft <= 90, "Caution (≤90 days)",
        "Current"
    )
```

```dax
// Count of expired qualifications
Expired Quals Count =
COUNTROWS(
    FILTER(
        QualificationRecords,
        QualificationRecords[ExpirationDate] < TODAY()
    )
)
```

```dax
// Count expiring within 30 days
Expiring 30 Days =
COUNTROWS(
    FILTER(
        QualificationRecords,
        QualificationRecords[ExpirationDate] >= TODAY()
        && QualificationRecords[ExpirationDate] <= TODAY() + 30
    )
)
```

```dax
// Count expiring within 60 days
Expiring 60 Days =
COUNTROWS(
    FILTER(
        QualificationRecords,
        QualificationRecords[ExpirationDate] >= TODAY()
        && QualificationRecords[ExpirationDate] <= TODAY() + 60
    )
)
```

```dax
// Count expiring within 90 days
Expiring 90 Days =
COUNTROWS(
    FILTER(
        QualificationRecords,
        QualificationRecords[ExpirationDate] >= TODAY()
        && QualificationRecords[ExpirationDate] <= TODAY() + 90
    )
)
```

### Coverage Analysis

```dax
// Total active instructors
Active Instructor Count =
COUNTROWS(
    FILTER(Instructors, Instructors[Status] = "Active")
)
```

```dax
// Qualification coverage rate — % of required quals currently held across all active instructors
Qual Coverage Rate =
VAR TotalRequired =
    SUMX(
        RequiredQualifications,
        RequiredQualifications[MinimumPerEvent]
    )
VAR TotalCurrent =
    COUNTROWS(
        FILTER(
            QualificationRecords,
            QualificationRecords[ExpirationDate] >= TODAY()
        )
    )
RETURN
    DIVIDE(TotalCurrent, TotalRequired, 0)
```

```dax
// Qualification gap — quals where current holders < minimum required
Qual Gaps =
VAR QualCode = RequiredQualifications[QualCode]
VAR MinRequired = RequiredQualifications[MinimumPerEvent]
VAR CurrentHolders =
    COUNTROWS(
        FILTER(
            QualificationRecords,
            QualificationRecords[QualCode] = QualCode
            && QualificationRecords[ExpirationDate] >= TODAY()
        )
    )
RETURN
    IF(CurrentHolders < MinRequired, MinRequired - CurrentHolders, 0)
```

```dax
// Instructors with no current qualifications (new or all expired)
Unqualified Instructors =
COUNTROWS(
    FILTER(
        Instructors,
        Instructors[Status] = "Active"
        && COUNTROWS(
            FILTER(
                QualificationRecords,
                QualificationRecords[InstructorEDIPI] = Instructors[InstructorEDIPI]
                && QualificationRecords[ExpirationDate] >= TODAY()
            )
        ) = 0
    )
)
```

### Workload Measures

```dax
// Average students per SPC
Avg Students Per SPC =
VAR SPCCount =
    COUNTROWS(
        FILTER(
            Instructors,
            Instructors[Role] = "Staff Platoon Commander"
            && Instructors[Status] = "Active"
        )
    )
RETURN
    DIVIDE(
        SUM(Instructors[StudentsAssigned]),
        SPCCount,
        0
    )
```

```dax
// Workload standard deviation — identifies imbalanced assignment
Workload StdDev =
VAR AvgLoad = [Avg Students Per SPC]
VAR SPCs =
    FILTER(
        Instructors,
        Instructors[Role] = "Staff Platoon Commander"
        && Instructors[Status] = "Active"
    )
VAR Variance =
    AVERAGEX(
        SPCs,
        (Instructors[StudentsAssigned] - AvgLoad) ^ 2
    )
RETURN
    SQRT(Variance)
```

```dax
// Events per instructor this month
Monthly Event Load =
CALCULATE(
    SUM(Instructors[EventsThisMonth]),
    Instructors[Status] = "Active"
)
```

### Rotation Planning

```dax
// Instructors departing within 90 days
Rotating Soon =
COUNTROWS(
    FILTER(
        Instructors,
        NOT(ISBLANK(Instructors[PRD]))
        && Instructors[PRD] <= TODAY() + 90
        && Instructors[Status] <> "Departed"
    )
)
```

```dax
// Quals at risk from rotating instructors — quals held only by soon-departing instructors
Quals At Risk From Rotation =
VAR RotatingEDIPIs =
    SELECTCOLUMNS(
        FILTER(
            Instructors,
            NOT(ISBLANK(Instructors[PRD]))
            && Instructors[PRD] <= TODAY() + 90
        ),
        "EDIPI", Instructors[InstructorEDIPI]
    )
RETURN
    COUNTROWS(
        FILTER(
            QualificationRecords,
            QualificationRecords[InstructorEDIPI] IN RotatingEDIPIs
            && QualificationRecords[ExpirationDate] >= TODAY()
        )
    )
```

---

## Visual Layout

### Page 1: Qualification Status Overview

| Position | Visual | Data | Size |
|----------|--------|------|------|
| Top-left | Card (red border) | Expired Quals Count | 1/4 width |
| Top-center-left | Card (orange border) | Expiring 30 Days | 1/4 width |
| Top-center-right | Card (yellow border) | Expiring 60 Days | 1/4 width |
| Top-right | Card (green border) | Active Instructor Count | 1/4 width |
| Middle-left | Stacked bar chart | Qualification status by category (Expired/Critical/Warning/Caution/Current) stacked per RequiredQualifications.Category | 2/3 width |
| Middle-right | Donut chart | Overall qual status distribution (% Current vs Warning vs Expired) | 1/3 width |
| Bottom | Table | InstructorName, QualName, DateEarned, ExpirationDate, Days Until Expiration, Expiration Status, RenewalStatus | Full width, sorted by DaysUntilExpiration ASC |

**Slicers:** CompanyAssigned, RequiredQualifications.Category, Expiration Status
**Conditional formatting:**
- Expired: red row
- Critical (≤30 days): red text
- Warning (≤60 days): orange text
- Caution (≤90 days): yellow text

### Page 2: Instructor Workload

| Position | Visual | Data |
|----------|--------|------|
| Top-left | Card | Avg Students Per SPC |
| Top-right | Card | Workload StdDev (red if > 5, meaning large imbalance) |
| Middle | Clustered bar chart | Students assigned per instructor, grouped by company, colored by role |
| Bottom-left | Bar chart | Events this month per instructor (top 20 busiest) |
| Bottom-right | Table | InstructorName, Company, Role, StudentsAssigned, EventsThisWeek, EventsThisMonth, CounselingsOverdue |

**Conditional formatting:**
- CounselingsOverdue > 0: red text
- StudentsAssigned > Avg + 1 StdDev: orange highlight (overloaded)

### Page 3: Coverage Gaps & Rotation Risk

| Position | Visual | Data |
|----------|--------|------|
| Top-left | Card | Qual Gaps (total shortfalls) |
| Top-right | Card | Rotating Soon (instructors PCSing within 90 days) |
| Middle-left | Table | Qualification gaps: QualName, Category, MinimumPerEvent, Current Holders, Shortfall | 1/2 width |
| Middle-right | Table | Rotating instructors: Name, Company, Role, PRD, Quals Held | 1/2 width |
| Bottom | Matrix | Company (rows) × Qual Category (columns) → coverage count in each cell, red where below minimum |

### Page 4: Individual Instructor Detail (Drill-through)

| Position | Visual | Data |
|----------|--------|------|
| Top | Card row | InstructorName, Rank, Role, Company, Platoon, Status, PRD |
| Middle | Table | All QualificationRecords for this instructor: QualName, DateEarned, ExpirationDate, Status, RenewalStatus |
| Bottom-left | Timeline visual | Qualification timeline (earned dates and expiration dates on a Gantt-style bar) |
| Bottom-right | Card | StudentsAssigned, EventsThisMonth, CounselingsOverdue |

**Drill-through from:** Page 1 table or Page 2 table → click instructor row

---

## Row-Level Security Configuration

This dashboard is **Staff only**. RLS is simpler:

**Role: Staff**
```dax
// Full access — no filter
TRUE()
```

**Role: CompanyCommander**
```dax
// If company commanders need access to their own company only
[CompanyAssigned] = LOOKUPVALUE(
    UserSecurity[Company],
    UserSecurity[UserEmail],
    USERPRINCIPALNAME()
)
```

---

## Conditional Formatting Rules

| Column/Visual | Condition | Format |
|---------------|-----------|--------|
| Days Until Expiration | < 0 | Red background, white text, bold |
| Days Until Expiration | 0-30 | Red text |
| Days Until Expiration | 31-60 | Orange text |
| Days Until Expiration | 61-90 | Yellow text |
| Days Until Expiration | > 90 | Green text |
| CounselingsOverdue | > 0 | Red bold |
| StudentsAssigned | > company avg + 1σ | Orange highlight |
| Qual Gaps shortfall | > 0 | Red background |
| Coverage matrix cells | < MinimumPerEvent | Red |
| Coverage matrix cells | = MinimumPerEvent | Yellow |
| Coverage matrix cells | > MinimumPerEvent | Green |

---

## Alert Configuration (Power BI Alerts)

Set up data-driven alerts in Power BI Service:

| Alert | Trigger | Recipient | Frequency |
|-------|---------|-----------|-----------|
| Qualification Expired | Expired Quals Count > 0 | S-3 Training Officer | Daily |
| Critical Expiration | Expiring 30 Days increases | Company Commander + S-3 | Daily |
| Workload Imbalance | Workload StdDev > 8 | S-3 Training Officer | Weekly |
| Counseling Overdue | Any CounselingsOverdue > 2 | Company Commander | Daily |

---

## Deployment Notes

1. Three SharePoint lists must be created and populated before dashboard connects
2. RequiredQualifications should be populated first (reference data)
3. Instructors list second (PII — requires PIA threshold completion)
4. QualificationRecords third (links to both)
5. Import mode with daily refresh at 0500 EST
6. Staff-only access — do not add to shared Power BI app with student view
7. Create separate Power BI workspace for staff dashboards
8. Test alerts with sample data before go-live
