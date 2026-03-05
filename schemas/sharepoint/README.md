# SharePoint List Schemas — Heywood Phase 1

These JSON schemas define the SharePoint Online lists that serve as the Phase 1 data layer. All lists deploy within the existing MCEN M365 authorization boundary — no ATO required.

## Lists

| File | List Name | PII | Records | Purpose |
|------|-----------|:---:|---------|---------|
| `student-scores.json` | StudentScores | Yes | ~200/company | Individual student performance across 3-pillar grading system |
| `training-schedule.json` | TrainingSchedule | No | ~300/class | Master training calendar with event-level detail |
| `instructors.json` | Instructors | Yes | ~80 | Instructor roster, company assignments, workload metrics |
| `required-qualifications.json` | RequiredQualifications | No | ~30-50 | Master list of qualifications required for TBS events |
| `qualification-records.json` | QualificationRecords | Yes | ~400+ | Individual instructor qual records with expiration tracking |
| `event-feedback.json` | EventFeedback | Optional | Growing | Post-event structured feedback — feeds AAR synthesis |

## Relationships

```
RequiredQualifications (reference)
        ↓ QualCode
QualificationRecords ←→ Instructors (via InstructorEDIPI)
                              ↓ LeadInstructor / SupportInstructors
TrainingSchedule ←→ EventFeedback (via EventCode)
        ↓ CompanyAssigned / ClassNumber
StudentScores (via Company + ClassNumber)
```

## Row-Level Security Model

| Role | StudentScores | TrainingSchedule | Instructors | QualRecords | EventFeedback |
|------|:---:|:---:|:---:|:---:|:---:|
| **Staff** | All | All | All | All | All |
| **SPC** | Own company | Own company + All | Own company | Own company | Own company |
| **Student** | Own record only | Own company (read) | No access | No access | Write-only |

## Deployment Notes

1. Create lists in SharePoint Online via browser or PnP PowerShell
2. These schemas are design specifications — not direct import files
3. Calculated columns require SharePoint column formulas or Power BI DAX measures
4. Person columns use SharePoint People Picker tied to Azure AD
5. Row-Level Security is enforced at the Power BI layer (SharePoint views provide UI filtering but not true security)
6. PII lists require PIA threshold analysis completion before data entry begins
