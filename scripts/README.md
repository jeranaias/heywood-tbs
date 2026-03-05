# Heywood TBS - SharePoint Deployment Scripts

## Prerequisites

1. **PowerShell 5.1 or 7+** (pre-installed on MCEN workstations)
2. **PnP.PowerShell module** -- install once per machine:

```powershell
# If PSGallery is untrusted (common on MCEN), trust it first:
Set-PSRepository -Name PSGallery -InstallationPolicy Trusted

# Install the module for the current user (no admin required):
Install-Module PnP.PowerShell -Scope CurrentUser -Force
```

3. **CAC reader and valid CAC** inserted for MCEN authentication
4. **Site Collection Admin** or **Site Owner** permissions on the target SharePoint site

## Quick Start

### Deploy all 6 lists

```powershell
.\Deploy-SharePointLists.ps1 -SiteUrl "https://usmc.dps.mil/sites/TBS-Heywood"
```

A browser window will open for CAC/PIV authentication (WebLogin). Select your certificate and proceed.

### Preview what would be created (no changes)

```powershell
.\Deploy-SharePointLists.ps1 -SiteUrl "https://usmc.dps.mil/sites/TBS-Heywood" -WhatIf
```

### Deploy with sample data

```powershell
.\Deploy-SharePointLists.ps1 -SiteUrl "https://usmc.dps.mil/sites/TBS-Heywood" -ImportSampleData
```

Imports CSV files from `../schemas/sharepoint/sample-data/`. Each CSV must be named to match the list (e.g., `StudentScores.csv`).

### Remove all lists (cleanup for re-deployment)

```powershell
.\Deploy-SharePointLists.ps1 -SiteUrl "https://usmc.dps.mil/sites/TBS-Heywood" -RemoveAll
```

Prompts for confirmation before deleting. Use `-Confirm:$false` to skip the prompt (e.g., in automated pipelines).

### Verbose logging

Append `-Verbose` to any command for detailed output:

```powershell
.\Deploy-SharePointLists.ps1 -SiteUrl "https://usmc.dps.mil/sites/TBS-Heywood" -Verbose
```

## Lists Deployed

| # | List Name             | PII  | Description |
|---|----------------------|------|-------------|
| 1 | StudentScores         | Yes  | Individual student performance across TBS three-pillar grading |
| 2 | TrainingSchedule      | No   | Master training schedule with event-level detail |
| 3 | Instructors           | Yes  | Instructor roster with company assignments and workload |
| 4 | RequiredQualifications | No   | Master qualification reference table |
| 5 | QualificationRecords  | Yes  | Individual instructor qualification records with expiration tracking |
| 6 | EventFeedback         | Opt. | Post-event feedback (anonymous by default) |

## Idempotent Design

The script is safe to re-run:

- If a list already exists, it is skipped with a warning.
- If a column already exists on a list, it is skipped silently.
- If a view already exists, it is skipped.

To force a clean re-deployment, run with `-RemoveAll` first, then deploy again.

## MCEN-Specific Notes

### Tenant URL

MCEN SharePoint Online sites typically follow one of these patterns:

- `https://usmc.dps.mil/sites/<SiteName>`
- `https://usmc.sharepoint-mil.us/sites/<SiteName>`

Confirm the exact URL with your SharePoint admin. The `-SiteUrl` parameter validates that the URL starts with `https://`.

### Authentication

- **Default (CAC/PIV):** The script uses `Connect-PnPOnline -UseWebLogin`, which opens a browser window for certificate-based authentication. This is the standard method on MCEN workstations.
- **Service Account:** Pass a `PSCredential` object via `-Credential` for non-interactive use (e.g., scheduled tasks), though this is uncommon on MCEN due to CAC requirements.
- **MFA:** WebLogin handles MFA challenges automatically through the browser flow.

### Proxy / Network

If your MCEN workstation routes through a proxy, you may need to configure PowerShell:

```powershell
[System.Net.WebRequest]::DefaultWebProxy.Credentials = [System.Net.CredentialCache]::DefaultNetworkCredentials
```

### Calculated Fields

Calculated fields (AcademicComposite, MilSkillsComposite, LeadershipComposite, OverallComposite, DaysUntilExpiration, ExpirationStatus) are created with placeholder formulas (`=""`). The intended formula logic is documented in each schema JSON file under the `formula` key and in the field description. Update these formulas manually in SharePoint:

1. Go to **List Settings** > **Columns** > click the calculated column
2. Enter the SharePoint-compatible formula
3. Save

### PowerShell Execution Policy

If you get a "script cannot be loaded" error:

```powershell
Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned
```

## File Structure

```
heywood-tbs/
  schemas/sharepoint/
    student-scores.json          # Schema definition
    training-schedule.json
    instructors.json
    required-qualifications.json
    qualification-records.json
    event-feedback.json
    sample-data/                 # CSV files for -ImportSampleData
      StudentScores.csv
      TrainingSchedule.csv
      ...
  scripts/
    Deploy-SharePointLists.ps1   # This deployment script
    README.md                    # This file
```
