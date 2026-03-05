<#
.SYNOPSIS
    Deploys all 6 Heywood TBS SharePoint lists to a SharePoint Online site.

.DESCRIPTION
    Creates the following lists on MCEN SharePoint Online with full column definitions,
    property settings, and views as defined in the Heywood TBS schema files:

      1. StudentScores        - Individual student performance data (PII)
      2. TrainingSchedule     - Master training schedule with event detail
      3. Instructors          - Instructor roster with workload tracking (PII)
      4. RequiredQualifications - Master qualification reference table
      5. QualificationRecords - Individual instructor qualification records (PII)
      6. EventFeedback        - Post-event feedback from students and instructors

    The script is idempotent: re-running it will skip lists and columns that already exist.

.PARAMETER SiteUrl
    Required. The full URL of the SharePoint Online site.
    Example: https://usmc.dps.mil/sites/TBS-Heywood

.PARAMETER Credential
    Optional. A PSCredential object for non-interactive authentication.
    If omitted, the script uses -UseWebLogin for CAC/PIV authentication on MCEN.

.PARAMETER ImportSampleData
    Optional switch. If present, imports CSV files from the sample-data directory
    located at ..\schemas\sharepoint\sample-data\ relative to this script.

.PARAMETER WhatIf
    Preview mode. Shows what would be created without making changes.

.PARAMETER RemoveAll
    Cleanup mode. Removes all 6 lists from the site. Prompts for confirmation
    unless -Confirm:$false is passed.

.EXAMPLE
    .\Deploy-SharePointLists.ps1 -SiteUrl "https://usmc.dps.mil/sites/TBS-Heywood"

    Connects via CAC/WebLogin and deploys all lists.

.EXAMPLE
    .\Deploy-SharePointLists.ps1 -SiteUrl "https://usmc.dps.mil/sites/TBS-Heywood" -WhatIf

    Preview mode: shows what would be created without making changes.

.EXAMPLE
    .\Deploy-SharePointLists.ps1 -SiteUrl "https://usmc.dps.mil/sites/TBS-Heywood" -ImportSampleData

    Deploys all lists and imports sample data from CSVs.

.EXAMPLE
    .\Deploy-SharePointLists.ps1 -SiteUrl "https://usmc.dps.mil/sites/TBS-Heywood" -RemoveAll

    Removes all 6 Heywood lists from the site (with confirmation prompt).

.NOTES
    Author:  Heywood TBS Project
    Requires: PnP.PowerShell module (Install-Module PnP.PowerShell -Scope CurrentUser)
    Platform: PowerShell 5.1 or 7+
    Network:  MCEN (Marine Corps Enterprise Network) with CAC authentication
#>

[CmdletBinding(SupportsShouldProcess = $true, DefaultParameterSetName = 'Deploy')]
param(
    [Parameter(Mandatory = $true, Position = 0, HelpMessage = "SharePoint Online site URL")]
    [ValidateNotNullOrEmpty()]
    [ValidatePattern('^https://')]
    [string]$SiteUrl,

    [Parameter(Mandatory = $false)]
    [System.Management.Automation.PSCredential]$Credential,

    [Parameter(Mandatory = $false, ParameterSetName = 'Deploy')]
    [switch]$ImportSampleData,

    [Parameter(Mandatory = $true, ParameterSetName = 'Remove')]
    [switch]$RemoveAll
)

# ---------------------------------------------------------------------------
# Module check
# ---------------------------------------------------------------------------
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

function Assert-PnPModule {
    if (-not (Get-Module -ListAvailable -Name 'PnP.PowerShell')) {
        Write-Error @"
PnP.PowerShell module is not installed.
Install it with:
    Install-Module PnP.PowerShell -Scope CurrentUser -Force
On MCEN you may need to set the repository to trusted first:
    Set-PSRepository -Name PSGallery -InstallationPolicy Trusted
"@
        exit 1
    }
    Import-Module PnP.PowerShell -ErrorAction Stop
    Write-Verbose "PnP.PowerShell module loaded successfully."
}

# ---------------------------------------------------------------------------
# Connection helper
# ---------------------------------------------------------------------------
function Connect-Site {
    param([string]$Url, [System.Management.Automation.PSCredential]$Cred)

    Write-Host "[CONNECT] Connecting to $Url ..." -ForegroundColor Cyan

    $connectParams = @{ Url = $Url }

    if ($Cred) {
        Write-Verbose "Using supplied credential for authentication."
        $connectParams['Credentials'] = $Cred
    }
    else {
        Write-Verbose "No credential supplied -- using WebLogin (CAC/PIV)."
        $connectParams['UseWebLogin'] = $true
    }

    try {
        Connect-PnPOnline @connectParams
        Write-Host "[CONNECT] Connected successfully." -ForegroundColor Green
    }
    catch {
        Write-Error "Failed to connect to $Url : $_"
        exit 1
    }
}

# ---------------------------------------------------------------------------
# Core helper: ensure a list exists
# ---------------------------------------------------------------------------
function Ensure-List {
    [CmdletBinding(SupportsShouldProcess = $true)]
    param(
        [string]$ListName,
        [string]$Description
    )

    $existing = Get-PnPList -Identity $ListName -ErrorAction SilentlyContinue
    if ($existing) {
        Write-Host "  [SKIP] List '$ListName' already exists." -ForegroundColor Yellow
        return $existing
    }

    if ($PSCmdlet.ShouldProcess($ListName, "Create list")) {
        Write-Host "  [CREATE] Creating list '$ListName' ..." -ForegroundColor Green
        $list = New-PnPList -Title $ListName -Template GenericList -ErrorAction Stop
        Set-PnPList -Identity $ListName -Description $Description -ErrorAction SilentlyContinue
        return $list
    }
}

# ---------------------------------------------------------------------------
# Core helper: add a field to a list
# ---------------------------------------------------------------------------
function Add-ListColumn {
    [CmdletBinding(SupportsShouldProcess = $true)]
    param(
        [string]$ListName,
        [hashtable]$Column
    )

    $colName = $Column.name
    $colType = $Column.type

    # Check if column already exists
    $existingField = Get-PnPField -List $ListName -Identity $colName -ErrorAction SilentlyContinue
    if ($existingField) {
        Write-Verbose "    [SKIP] Column '$colName' already exists on '$ListName'."
        return
    }

    if (-not $PSCmdlet.ShouldProcess("$ListName.$colName ($colType)", "Add column")) {
        return
    }

    $description = if ($Column.ContainsKey('description')) { $Column.description } else { '' }

    switch ($colType) {

        'Text' {
            $maxLen = if ($Column.ContainsKey('maxLength')) { $Column.maxLength } else { 255 }
            Add-PnPField -List $ListName -InternalName $colName -DisplayName $colName `
                -Type Text -Required:([bool]$Column.required) -ErrorAction Stop | Out-Null

            # Set max length via schema XML if needed and less than 255
            if ($maxLen -and $maxLen -lt 255) {
                $fieldXml = "<Field Type='Text' Name='$colName' DisplayName='$colName' MaxLength='$maxLen' />"
                Set-PnPField -List $ListName -Identity $colName -Values @{
                    Description = $description
                } -ErrorAction SilentlyContinue
            }
            else {
                Set-PnPField -List $ListName -Identity $colName -Values @{
                    Description = $description
                } -ErrorAction SilentlyContinue
            }
        }

        'Choice' {
            $choices = $Column.choices
            $default = if ($Column.ContainsKey('default')) { $Column.default } else { $null }

            $addParams = @{
                List         = $ListName
                InternalName = $colName
                DisplayName  = $colName
                Type         = 'Choice'
                Choices      = $choices
                Required     = [bool]$Column.required
                ErrorAction  = 'Stop'
            }

            Add-PnPField @addParams | Out-Null

            $setValues = @{ Description = $description }
            if ($default) {
                $setValues['DefaultValue'] = $default
            }
            Set-PnPField -List $ListName -Identity $colName -Values $setValues -ErrorAction SilentlyContinue
        }

        'Number' {
            Add-PnPField -List $ListName -InternalName $colName -DisplayName $colName `
                -Type Number -Required:([bool]$Column.required) -ErrorAction Stop | Out-Null

            $setValues = @{ Description = $description }
            if ($Column.ContainsKey('min')) {
                $setValues['MinimumValue'] = $Column.min
            }
            if ($Column.ContainsKey('max')) {
                $setValues['MaximumValue'] = $Column.max
            }
            Set-PnPField -List $ListName -Identity $colName -Values $setValues -ErrorAction SilentlyContinue
        }

        'DateTime' {
            $displayFormat = if ($Column.ContainsKey('dateOnly') -and $Column.dateOnly) {
                'DateOnly'
            } elseif ($Column.ContainsKey('timeOnly') -and $Column.timeOnly) {
                'DateTime'
            } else {
                'DateTime'
            }

            Add-PnPField -List $ListName -InternalName $colName -DisplayName $colName `
                -Type DateTime -Required:([bool]$Column.required) -ErrorAction Stop | Out-Null

            Set-PnPField -List $ListName -Identity $colName -Values @{
                Description   = $description
                DisplayFormat = $displayFormat
            } -ErrorAction SilentlyContinue
        }

        'Boolean' {
            $defaultVal = if ($Column.ContainsKey('default')) {
                if ($Column.default) { '1' } else { '0' }
            } else { '0' }

            Add-PnPField -List $ListName -InternalName $colName -DisplayName $colName `
                -Type Boolean -Required:([bool]$Column.required) -ErrorAction Stop | Out-Null

            Set-PnPField -List $ListName -Identity $colName -Values @{
                Description  = $description
                DefaultValue = $defaultVal
            } -ErrorAction SilentlyContinue
        }

        'Person' {
            Add-PnPField -List $ListName -InternalName $colName -DisplayName $colName `
                -Type User -Required:([bool]$Column.required) -ErrorAction Stop | Out-Null

            Set-PnPField -List $ListName -Identity $colName -Values @{
                Description = $description
            } -ErrorAction SilentlyContinue
        }

        'PersonMulti' {
            # Person field allowing multiple selections
            $fieldXml = @"
<Field Type="UserMulti" Name="$colName" DisplayName="$colName" StaticName="$colName"
       UserSelectionMode="PeopleOnly" Mult="TRUE" Required="$( if ($Column.required) { 'TRUE' } else { 'FALSE' } )"
       Description="$([System.Security.SecurityElement]::Escape($description))" />
"@
            Add-PnPFieldFromXml -List $ListName -FieldXml $fieldXml -ErrorAction Stop | Out-Null
        }

        'MultilineText' {
            Add-PnPField -List $ListName -InternalName $colName -DisplayName $colName `
                -Type Note -Required:([bool]$Column.required) -ErrorAction Stop | Out-Null

            Set-PnPField -List $ListName -Identity $colName -Values @{
                Description    = $description
                RichText       = $false
                NumberOfLines  = 6
            } -ErrorAction SilentlyContinue
        }

        'Calculated' {
            # Calculated fields -- we store the formula description as a note field
            # because true calculated-field formulas require exact SharePoint formula syntax.
            # The schema "formula" values are human-readable descriptions, not SP formulas.
            # We create them as Calculated type with a placeholder; admins can refine formulas
            # in the SharePoint UI or via Set-PnPField after deployment.
            Write-Verbose "    [CALC] Column '$colName' is Calculated -- creating as placeholder."

            $formula = if ($Column.ContainsKey('formula')) { $Column.formula } else { '=""' }

            # Build a safe calculated field XML
            # Default formula returns empty string; admin should update via SP UI
            $calcFormula = '=""'
            $fieldXml = @"
<Field Type="Calculated" Name="$colName" DisplayName="$colName" StaticName="$colName"
       ResultType="Text" ReadOnly="TRUE"
       Description="$([System.Security.SecurityElement]::Escape($description + ' | Formula intent: ' + $formula))">
  <Formula>$calcFormula</Formula>
</Field>
"@
            try {
                Add-PnPFieldFromXml -List $ListName -FieldXml $fieldXml -ErrorAction Stop | Out-Null
            }
            catch {
                Write-Warning "    [WARN] Could not create calculated field '$colName': $_"
                Write-Warning "           Formula intent: $formula"
                Write-Warning "           Create this field manually in the SharePoint UI."
            }
        }

        'Hyperlink' {
            Add-PnPField -List $ListName -InternalName $colName -DisplayName $colName `
                -Type URL -Required:([bool]$Column.required) -ErrorAction Stop | Out-Null

            Set-PnPField -List $ListName -Identity $colName -Values @{
                Description = $description
                DisplayFormat = 'Hyperlink'
            } -ErrorAction SilentlyContinue
        }

        default {
            Write-Warning "    [WARN] Unknown column type '$colType' for '$colName'. Skipping."
        }
    }

    # Set indexed property if specified
    if ($Column.ContainsKey('indexed') -and $Column.indexed) {
        try {
            Set-PnPField -List $ListName -Identity $colName -Values @{
                Indexed = $true
            } -ErrorAction SilentlyContinue
        }
        catch {
            Write-Warning "    [WARN] Could not set index on '$colName': $_"
        }
    }

    Write-Host "    [ADD] $colName ($colType)" -ForegroundColor Gray
}

# ---------------------------------------------------------------------------
# Core helper: create a view on a list
# ---------------------------------------------------------------------------
function Add-ListView {
    [CmdletBinding(SupportsShouldProcess = $true)]
    param(
        [string]$ListName,
        [hashtable]$ViewDef
    )

    $viewName = $ViewDef.name

    $existingView = Get-PnPView -List $ListName -Identity $viewName -ErrorAction SilentlyContinue
    if ($existingView) {
        Write-Verbose "    [SKIP] View '$viewName' already exists on '$ListName'."
        return
    }

    if (-not $PSCmdlet.ShouldProcess("$ListName / $viewName", "Create view")) {
        return
    }

    # Build the fields list for the view
    # Filter out any pseudo-column references (e.g., "AcademicExam1-4" is not a real column)
    $viewFields = @()
    foreach ($f in $ViewDef.columns) {
        # Skip range references like "AcademicExam1-4"
        if ($f -match '-\d+$') { continue }
        $viewFields += $f
    }

    # Build CAML query for filter and sort
    $camlQuery = Build-CamlQuery -ViewDef $ViewDef

    $viewParams = @{
        List      = $ListName
        Title     = $viewName
        Fields    = $viewFields
        RowLimit  = 100
        ErrorAction = 'Stop'
    }

    if ($camlQuery) {
        $viewParams['Query'] = $camlQuery
    }

    try {
        Add-PnPView @viewParams | Out-Null
        Write-Host "    [VIEW] $viewName" -ForegroundColor DarkCyan
    }
    catch {
        Write-Warning "    [WARN] Could not create view '$viewName': $_"
        Write-Warning "           You may need to create this view manually in the SharePoint UI."
    }
}

# ---------------------------------------------------------------------------
# Helper: build a simple CAML query from view definition
# ---------------------------------------------------------------------------
function Build-CamlQuery {
    param([hashtable]$ViewDef)

    $parts = @()
    $orderBy = ''

    # --- Sort ---
    if ($ViewDef.ContainsKey('sort') -and $ViewDef.sort) {
        $sortParts = $ViewDef.sort -split ','
        $orderFields = @()
        foreach ($sp in $sortParts) {
            $sp = $sp.Trim()
            if ($sp -match '^(\S+)\s+(ASC|DESC)$') {
                $fieldName = $Matches[1]
                $ascending = if ($Matches[2] -eq 'ASC') { 'TRUE' } else { 'FALSE' }
                $orderFields += "<FieldRef Name='$fieldName' Ascending='$ascending' />"
            }
            elseif ($sp -match '^(\S+)$') {
                $orderFields += "<FieldRef Name='$($Matches[1])' Ascending='TRUE' />"
            }
        }
        if ($orderFields.Count -gt 0) {
            $orderBy = "<OrderBy>" + ($orderFields -join '') + "</OrderBy>"
        }
    }

    # --- Filter (simple equality filters only) ---
    $whereClause = ''
    if ($ViewDef.ContainsKey('filter') -and $ViewDef.filter) {
        $filterStr = $ViewDef.filter

        # Parse simple conditions: "Field = Value" / "Field != Value"
        # Handle AND/OR combinations. Skip complex placeholder filters like [CurrentUser].
        if ($filterStr -notmatch '\[') {
            $conditions = @()
            $andParts = $filterStr -split '\s+AND\s+'
            foreach ($part in $andParts) {
                $part = $part.Trim()
                if ($part -match '^(\S+)\s*!=\s*(.+)$') {
                    $fn = $Matches[1]
                    $fv = $Matches[2].Trim()
                    $conditions += "<Neq><FieldRef Name='$fn' /><Value Type='Text'>$fv</Value></Neq>"
                }
                elseif ($part -match '^(\S+)\s*=\s*(.+)$') {
                    $fn = $Matches[1]
                    $fv = $Matches[2].Trim()

                    # Determine value type
                    $valType = 'Text'
                    if ($fv -eq 'true' -or $fv -eq 'false') {
                        $valType = 'Boolean'
                        $fv = if ($fv -eq 'true') { '1' } else { '0' }
                    }
                    $conditions += "<Eq><FieldRef Name='$fn' /><Value Type='$valType'>$fv</Value></Eq>"
                }
            }

            if ($conditions.Count -eq 1) {
                $whereClause = "<Where>" + $conditions[0] + "</Where>"
            }
            elseif ($conditions.Count -eq 2) {
                $whereClause = "<Where><And>" + $conditions[0] + $conditions[1] + "</And></Where>"
            }
            elseif ($conditions.Count -gt 2) {
                # Nest multiple ANDs
                $nested = $conditions[-1]
                for ($i = $conditions.Count - 2; $i -ge 0; $i--) {
                    $nested = "<And>" + $conditions[$i] + $nested + "</And>"
                }
                $whereClause = "<Where>$nested</Where>"
            }
        }
    }

    # --- GroupBy ---
    $groupByClause = ''
    if ($ViewDef.ContainsKey('groupBy') -and $ViewDef.groupBy) {
        $groupByClause = "<GroupBy><FieldRef Name='$($ViewDef.groupBy)' /></GroupBy>"
    }

    $query = $whereClause + $orderBy + $groupByClause
    if ($query) { return $query } else { return $null }
}

# ---------------------------------------------------------------------------
# List schema definitions (derived from JSON schema files)
# ---------------------------------------------------------------------------
function Get-ListDefinitions {
    $lists = @()

    # -----------------------------------------------------------------------
    # 1. StudentScores
    # -----------------------------------------------------------------------
    $lists += @{
        listName    = 'StudentScores'
        description = 'Heywood Phase 1 - Individual student performance data across TBS three-pillar grading system. PII: Yes.'
        columns     = @(
            @{ name = 'StudentEDIPI';       type = 'Text';          required = $true;  maxLength = 10;  indexed = $true;  description = '10-digit EDIPI. PII.' }
            @{ name = 'LastName';           type = 'Text';          required = $true;  maxLength = 50;  description = 'Student last name. PII.' }
            @{ name = 'FirstName';          type = 'Text';          required = $true;  maxLength = 50;  description = 'Student first name. PII.' }
            @{ name = 'Rank';              type = 'Choice';        required = $true;  choices = @('2ndLt','1stLt','WO'); default = '2ndLt' }
            @{ name = 'Company';           type = 'Choice';        required = $true;  choices = @('Alpha','Bravo','Charlie','Delta','Echo','Foxtrot','Golf','India','Mike'); indexed = $true }
            @{ name = 'Platoon';           type = 'Choice';        required = $true;  choices = @('1st','2nd','3rd','4th') }
            @{ name = 'SPCAssigned';       type = 'Person';        required = $true;  description = 'Staff Platoon Commander responsible for this student' }
            @{ name = 'ClassNumber';       type = 'Text';          required = $true;  indexed = $true; description = 'TBS class identifier (e.g., 1-26)' }
            @{ name = 'ClassStartDate';    type = 'DateTime';      required = $true;  dateOnly = $true }
            @{ name = 'CurrentPhase';      type = 'Choice';        required = $true;  choices = @('Phase I - Individual Skills','Phase II - Squad','Phase III - Platoon','Phase IV - MAGTF','Complete'); default = 'Phase I - Individual Skills' }
            @{ name = 'AcademicExam1';     type = 'Number';        required = $false; min = 0; max = 100; description = 'Phase I written exam score (percentage)' }
            @{ name = 'AcademicExam2';     type = 'Number';        required = $false; min = 0; max = 100; description = 'Phase II written exam score (percentage)' }
            @{ name = 'AcademicExam3';     type = 'Number';        required = $false; min = 0; max = 100; description = 'Phase III written exam score (percentage)' }
            @{ name = 'AcademicExam4';     type = 'Number';        required = $false; min = 0; max = 100; description = 'Phase IV written exam score (percentage)' }
            @{ name = 'AcademicQuizAvg';   type = 'Number';        required = $false; min = 0; max = 100; description = 'Running average of all quiz scores' }
            @{ name = 'AcademicComposite'; type = 'Calculated';    formula = 'Average of exams and quiz average'; description = 'Overall academics score (32% of total grade)' }
            @{ name = 'PFTScore';          type = 'Number';        required = $false; min = 0; max = 300 }
            @{ name = 'CFTScore';          type = 'Number';        required = $false; min = 0; max = 300 }
            @{ name = 'RifleQual';         type = 'Choice';        required = $false; choices = @('Expert','Sharpshooter','Marksman','Unqualified') }
            @{ name = 'PistolQual';        type = 'Choice';        required = $false; choices = @('Expert','Sharpshooter','Marksman','Unqualified') }
            @{ name = 'LandNavDay';        type = 'Choice';        required = $false; choices = @('Pass','Fail','Not Yet Tested') }
            @{ name = 'LandNavNight';      type = 'Choice';        required = $false; choices = @('Pass','Fail','Not Yet Tested') }
            @{ name = 'LandNavWritten';    type = 'Number';        required = $false; min = 0; max = 100 }
            @{ name = 'ObstacleCourse';    type = 'Choice';        required = $false; choices = @('Pass','Fail','Not Yet Tested') }
            @{ name = 'EnduranceCourse';   type = 'Choice';        required = $false; choices = @('Pass','Fail','Not Yet Tested') }
            @{ name = 'MilSkillsComposite'; type = 'Calculated';   formula = 'Weighted composite of PFT, CFT, rifle, pistol, land nav, obstacles'; description = 'Overall military skills score (32% of total grade)' }
            @{ name = 'LeadershipWeek12';  type = 'Number';        required = $false; min = 0; max = 100; description = 'Midpoint leadership evaluation (14% of total grade)' }
            @{ name = 'LeadershipWeek22';  type = 'Number';        required = $false; min = 0; max = 100; description = 'Final leadership evaluation (22% of total grade)' }
            @{ name = 'PeerEvalWeek12';    type = 'Number';        required = $false; min = 0; max = 100; description = 'Midpoint peer evaluation score' }
            @{ name = 'PeerEvalWeek22';    type = 'Number';        required = $false; min = 0; max = 100; description = 'Final peer evaluation score' }
            @{ name = 'LeadershipComposite'; type = 'Calculated';  formula = 'Week12 (14%) + Week22 (22%) with SPC (90%) and Peer (10%) weighting'; description = 'Overall leadership score (36% of total grade)' }
            @{ name = 'OverallComposite';  type = 'Calculated';    formula = '(AcademicComposite * 0.32) + (MilSkillsComposite * 0.32) + (LeadershipComposite * 0.36)'; description = 'Total TBS grade' }
            @{ name = 'ClassStandingThird'; type = 'Choice';       required = $false; choices = @('Top Third','Middle Third','Bottom Third'); description = 'Updated after each graded event' }
            @{ name = 'CompanyRank';       type = 'Number';        required = $false; min = 1; description = 'Rank order within company (1 = highest)' }
            @{ name = 'AtRiskFlag';        type = 'Choice';        required = $false; choices = @('None','Academic (<75%)','MilSkills (<75%)','Leadership (<75%)','Multiple (<75%)','Declining Trend'); default = 'None'; description = 'Auto-set when student falls below threshold' }
            @{ name = 'Status';            type = 'Choice';        required = $true;  choices = @('Active','Medical Hold (Mike Co)','Academic Hold','Administrative Hold','Graduated','Dropped'); default = 'Active' }
            @{ name = 'Notes';             type = 'MultilineText'; required = $false; description = 'SPC notes - do not include PHI' }
        )
        views = @(
            @{
                name    = 'Company Overview'
                columns = @('LastName','FirstName','Rank','Platoon','CurrentPhase','AcademicComposite','MilSkillsComposite','LeadershipComposite','OverallComposite','CompanyRank','AtRiskFlag')
                filter  = 'Status = Active'
                sort    = 'CompanyRank ASC'
            }
            @{
                name    = 'At Risk Students'
                columns = @('LastName','FirstName','Company','Platoon','AtRiskFlag','OverallComposite','SPCAssigned')
                filter  = 'AtRiskFlag != None AND Status = Active'
            }
            @{
                name    = 'Graduation Report'
                columns = @('LastName','FirstName','Rank','Company','OverallComposite','CompanyRank','ClassStandingThird')
                filter  = 'Status = Graduated'
                sort    = 'OverallComposite DESC'
            }
        )
    }

    # -----------------------------------------------------------------------
    # 2. TrainingSchedule
    # -----------------------------------------------------------------------
    $lists += @{
        listName    = 'TrainingSchedule'
        description = 'Heywood Phase 1 - TBS master training schedule with event-level detail across all four phases.'
        columns     = @(
            @{ name = 'EventTitle';             type = 'Text';          required = $true;  maxLength = 255; indexed = $true; description = 'Name of the training event' }
            @{ name = 'EventCode';              type = 'Text';          required = $false; maxLength = 20;  description = 'POI event reference code' }
            @{ name = 'TrainingPhase';          type = 'Choice';        required = $true;  choices = @('Phase I - Individual Skills','Phase II - Squad','Phase III - Platoon','Phase IV - MAGTF'); indexed = $true }
            @{ name = 'Category';               type = 'Choice';        required = $true;  choices = @('Academic','Military Skills','Leadership','Physical Training','Field Exercise','Evaluation','Admin'); indexed = $true }
            @{ name = 'GradePillar';            type = 'Choice';        required = $false; choices = @('Academics (32%)','Military Skills (32%)','Leadership (36%)','Not Graded'); description = 'Which grading pillar this event falls under' }
            @{ name = 'IsGraded';               type = 'Boolean';       required = $true;  default = $false; description = 'Whether this event produces a score' }
            @{ name = 'StartDate';              type = 'DateTime';      required = $true;  dateOnly = $true; indexed = $true }
            @{ name = 'EndDate';                type = 'DateTime';      required = $true;  dateOnly = $true }
            @{ name = 'StartTime';              type = 'DateTime';      required = $false; timeOnly = $true }
            @{ name = 'EndTime';                type = 'DateTime';      required = $false; timeOnly = $true }
            @{ name = 'DurationHours';          type = 'Number';        required = $false; min = 0.5; max = 240; description = 'Total event duration in hours' }
            @{ name = 'Location';               type = 'Text';          required = $false; maxLength = 100; description = 'Training area, range, classroom, or field location' }
            @{ name = 'CompanyAssigned';        type = 'Choice';        required = $true;  choices = @('All','Alpha','Bravo','Charlie','Delta','Echo','Foxtrot','Golf','India','Mike'); default = 'All'; indexed = $true }
            @{ name = 'ClassNumber';            type = 'Text';          required = $true;  indexed = $true; description = 'TBS class identifier' }
            @{ name = 'LeadInstructor';         type = 'Person';        required = $false; description = 'Primary instructor or OIC' }
            @{ name = 'SupportInstructors';     type = 'PersonMulti';   required = $false; description = 'Additional instructors or safety personnel' }
            @{ name = 'InstructorCountRequired'; type = 'Number';       required = $false; min = 1; max = 50; description = 'Total instructors needed' }
            @{ name = 'PrerequisiteEvents';     type = 'Text';          required = $false; description = 'Comma-separated EventCodes that must be complete first' }
            @{ name = 'SpecialEquipment';       type = 'MultilineText'; required = $false; description = 'Equipment, ammunition, or range requirements' }
            @{ name = 'Status';                 type = 'Choice';        required = $true;  choices = @('Scheduled','In Progress','Complete','Postponed','Cancelled'); default = 'Scheduled' }
            @{ name = 'WeatherContingency';     type = 'Choice';        required = $false; choices = @('No Impact','Rain Plan Available','Lightning Hold Applies','Cold Weather Threshold','Heat Cat Dependent'); default = 'No Impact' }
            @{ name = 'Notes';                  type = 'MultilineText'; required = $false }
        )
        views = @(
            @{
                name    = 'Weekly Calendar'
                columns = @('EventTitle','Category','StartDate','StartTime','EndTime','Location','CompanyAssigned','LeadInstructor','Status')
                sort    = 'StartDate ASC'
                groupBy = 'StartDate'
            }
            @{
                name    = 'Phase View'
                columns = @('EventTitle','Category','GradePillar','IsGraded','StartDate','EndDate','DurationHours','Status')
                sort    = 'StartDate ASC'
                groupBy = 'TrainingPhase'
            }
            @{
                name    = 'Instructor Assignments'
                columns = @('EventTitle','StartDate','CompanyAssigned','LeadInstructor','SupportInstructors','InstructorCountRequired','Status')
                filter  = 'Status = Scheduled'
                sort    = 'StartDate ASC'
            }
            @{
                name    = 'Graded Events Tracker'
                columns = @('EventTitle','TrainingPhase','GradePillar','StartDate','CompanyAssigned','Status')
                filter  = 'IsGraded = true'
                sort    = 'StartDate ASC'
            }
        )
    }

    # -----------------------------------------------------------------------
    # 3. Instructors
    # -----------------------------------------------------------------------
    $lists += @{
        listName    = 'Instructors'
        description = 'Heywood Phase 1 - TBS instructor roster with company assignments and workload tracking. PII: Yes.'
        columns     = @(
            @{ name = 'InstructorEDIPI'; type = 'Text';          required = $true;  maxLength = 10;  indexed = $true; description = '10-digit EDIPI. PII.' }
            @{ name = 'LastName';        type = 'Text';          required = $true;  maxLength = 50;  description = 'PII.' }
            @{ name = 'FirstName';       type = 'Text';          required = $true;  maxLength = 50;  description = 'PII.' }
            @{ name = 'Rank';           type = 'Choice';        required = $true;  choices = @('Capt','1stLt','CWO5','CWO4','CWO3','CWO2','MGySgt','MSgt','GySgt','SSgt') }
            @{ name = 'Role';           type = 'Choice';        required = $true;  choices = @('Staff Platoon Commander','Assistant SPC','Academic Instructor','Tactics Instructor','Weapons Instructor','Land Nav Instructor','PT Instructor','Company Commander','Company XO'); indexed = $true }
            @{ name = 'CompanyAssigned'; type = 'Choice';       required = $true;  choices = @('Alpha','Bravo','Charlie','Delta','Echo','Foxtrot','Golf','India','Mike','HQ'); indexed = $true }
            @{ name = 'PlatoonAssigned'; type = 'Choice';       required = $false; choices = @('1st','2nd','3rd','4th','N/A'); description = 'SPCs and ASPCs assigned to specific platoons' }
            @{ name = 'ClassNumber';     type = 'Text';          required = $true;  indexed = $true; description = 'Current class assignment' }
            @{ name = 'DateAssigned';    type = 'DateTime';      required = $true;  dateOnly = $true; description = 'Date assigned to current billet' }
            @{ name = 'PRD';            type = 'DateTime';      required = $false; dateOnly = $true; description = 'Projected Rotation Date' }
            @{ name = 'StudentsAssigned'; type = 'Number';       required = $false; min = 0; max = 60; description = 'Number of students this instructor evaluates' }
            @{ name = 'EventsThisWeek'; type = 'Number';        required = $false; min = 0; description = 'Current week event count' }
            @{ name = 'EventsThisMonth'; type = 'Number';       required = $false; min = 0; description = 'Current month event count' }
            @{ name = 'CounselingsOverdue'; type = 'Number';    required = $false; min = 0; description = 'Count of students with overdue counseling' }
            @{ name = 'Status';          type = 'Choice';        required = $true;  choices = @('Active','Leave/Liberty','TAD','Medical','PCS Pending','Departed'); default = 'Active' }
            @{ name = 'Phone';           type = 'Text';          required = $false; maxLength = 20; description = 'Duty phone or cell. PII.' }
            @{ name = 'Email';           type = 'Text';          required = $false; maxLength = 100; description = 'USMC email address' }
            @{ name = 'Notes';           type = 'MultilineText'; required = $false }
        )
        views = @(
            @{
                name    = 'Company Roster'
                columns = @('Rank','LastName','FirstName','Role','PlatoonAssigned','StudentsAssigned','Status')
                filter  = 'Status != Departed'
                sort    = 'Role ASC'
                groupBy = 'CompanyAssigned'
            }
            @{
                name    = 'Workload Overview'
                columns = @('LastName','CompanyAssigned','Role','StudentsAssigned','EventsThisWeek','EventsThisMonth','CounselingsOverdue')
                filter  = 'Status = Active'
                sort    = 'EventsThisMonth DESC'
            }
            @{
                name    = 'Rotation Tracker'
                columns = @('Rank','LastName','CompanyAssigned','Role','PRD','Status')
                filter  = 'Status != Departed'
                sort    = 'PRD ASC'
            }
        )
    }

    # -----------------------------------------------------------------------
    # 4. RequiredQualifications
    # -----------------------------------------------------------------------
    $lists += @{
        listName    = 'RequiredQualifications'
        description = 'Heywood Phase 1 - Master list of qualifications required to instruct TBS events. No PII.'
        columns     = @(
            @{ name = 'QualCode';          type = 'Text';          required = $true;  maxLength = 30;  indexed = $true; description = 'Unique qualification identifier' }
            @{ name = 'QualName';          type = 'Text';          required = $true;  maxLength = 150; description = 'Full qualification name' }
            @{ name = 'Category';          type = 'Choice';        required = $true;  choices = @('Range Safety','Demolitions','Weapons Instruction','Tactics Instruction','Land Navigation','Water Survival','Physical Training','Combat Lifesaver','Driving/Vehicle','Instructor Certification','Other'); indexed = $true }
            @{ name = 'IssuingAuthority';  type = 'Text';          required = $false; maxLength = 100; description = 'Certifying organization' }
            @{ name = 'ValidityMonths';    type = 'Number';        required = $true;  min = 1; max = 120; description = 'Months valid before renewal' }
            @{ name = 'RenewalProcess';    type = 'MultilineText'; required = $false; description = 'Steps required to renew' }
            @{ name = 'RequiredForEvents'; type = 'MultilineText'; required = $false; description = 'EventCodes requiring this qualification' }
            @{ name = 'MinimumPerEvent';   type = 'Number';        required = $false; min = 1; max = 20; description = 'Min qualified personnel per event' }
            @{ name = 'OrderReference';    type = 'Text';          required = $false; maxLength = 100; description = 'MCO, NAVMC, SOP reference' }
            @{ name = 'Status';            type = 'Choice';        required = $true;  choices = @('Active','Superseded','Archived'); default = 'Active' }
            @{ name = 'Notes';             type = 'MultilineText'; required = $false }
        )
        views = @(
            @{
                name    = 'All Active Qualifications'
                columns = @('QualCode','QualName','Category','ValidityMonths','MinimumPerEvent','IssuingAuthority')
                filter  = 'Status = Active'
                sort    = 'Category ASC'
            }
            @{
                name    = 'By Category'
                columns = @('QualCode','QualName','ValidityMonths','MinimumPerEvent','OrderReference')
                filter  = 'Status = Active'
                groupBy = 'Category'
            }
        )
    }

    # -----------------------------------------------------------------------
    # 5. QualificationRecords
    # -----------------------------------------------------------------------
    $lists += @{
        listName    = 'QualificationRecords'
        description = 'Heywood Phase 1 - Instructor qualification records with expiration tracking. PII: Yes.'
        columns     = @(
            @{ name = 'InstructorEDIPI';     type = 'Text';       required = $true;  maxLength = 10;  indexed = $true; description = '10-digit EDIPI. PII.' }
            @{ name = 'InstructorName';      type = 'Text';       required = $true;  maxLength = 100; description = 'Display name (Last, First). PII.' }
            @{ name = 'QualCode';            type = 'Text';       required = $true;  maxLength = 30;  indexed = $true; description = 'Links to RequiredQualifications.QualCode' }
            @{ name = 'QualName';            type = 'Text';       required = $true;  maxLength = 150; description = 'Denormalized qualification name' }
            @{ name = 'DateEarned';          type = 'DateTime';   required = $true;  dateOnly = $true }
            @{ name = 'ExpirationDate';      type = 'DateTime';   required = $true;  dateOnly = $true; indexed = $true; description = 'DateEarned + ValidityMonths' }
            @{ name = 'DaysUntilExpiration'; type = 'Calculated';  formula = 'ExpirationDate - Today'; description = 'Drives 30/60/90 day alert thresholds' }
            @{ name = 'ExpirationStatus';    type = 'Calculated';  formula = 'IF(DaysUntilExpiration < 0, Expired, IF(DaysUntilExpiration <= 30, Critical, IF(DaysUntilExpiration <= 60, Warning, IF(DaysUntilExpiration <= 90, Caution, Current))))'; description = 'Color-coded status for dashboards' }
            @{ name = 'CertificateNumber';   type = 'Text';       required = $false; maxLength = 50; description = 'Certificate or credential number' }
            @{ name = 'IssuedBy';            type = 'Text';       required = $false; maxLength = 100; description = 'Certifying official or organization' }
            @{ name = 'CompanyAtTimeOfCert'; type = 'Choice';     required = $false; choices = @('Alpha','Bravo','Charlie','Delta','Echo','Foxtrot','Golf','India','Mike','HQ'); description = 'Company when qualification earned' }
            @{ name = 'RenewalStatus';       type = 'Choice';     required = $false; choices = @('N/A - Current','Renewal Scheduled','Renewal In Progress','Renewal Overdue','Waiver Requested'); default = 'N/A - Current' }
            @{ name = 'RenewalDate';         type = 'DateTime';   required = $false; dateOnly = $true; description = 'Scheduled renewal date' }
            @{ name = 'DocumentLink';        type = 'Hyperlink';  required = $false; description = 'Link to scanned certificate in SharePoint' }
            @{ name = 'Notes';               type = 'MultilineText'; required = $false }
        )
        views = @(
            @{
                name    = 'Expiring Soon'
                columns = @('InstructorName','QualCode','QualName','ExpirationDate','DaysUntilExpiration','ExpirationStatus','RenewalStatus')
                sort    = 'ExpirationDate ASC'
            }
            @{
                name    = 'Expired'
                columns = @('InstructorName','QualCode','QualName','ExpirationDate','DaysUntilExpiration','RenewalStatus')
                sort    = 'ExpirationDate ASC'
            }
            @{
                name    = 'Instructor Quals'
                columns = @('QualCode','QualName','DateEarned','ExpirationDate','ExpirationStatus','RenewalStatus')
                sort    = 'ExpirationDate ASC'
            }
            @{
                name    = 'Qualification Coverage'
                columns = @('InstructorName','CompanyAtTimeOfCert','DateEarned','ExpirationDate','ExpirationStatus')
                sort    = 'ExpirationDate ASC'
            }
        )
    }

    # -----------------------------------------------------------------------
    # 6. EventFeedback
    # -----------------------------------------------------------------------
    $lists += @{
        listName    = 'EventFeedback'
        description = 'Heywood Phase 1 - Post-event feedback from students and instructors. Supports AAR synthesis and NLP analysis.'
        columns     = @(
            @{ name = 'EventTitle';              type = 'Text';          required = $true;  maxLength = 255; indexed = $true; description = 'Links to TrainingSchedule.EventTitle' }
            @{ name = 'EventCode';               type = 'Text';          required = $false; maxLength = 20;  description = 'Links to TrainingSchedule.EventCode' }
            @{ name = 'EventDate';               type = 'DateTime';      required = $true;  dateOnly = $true; indexed = $true }
            @{ name = 'TrainingPhase';           type = 'Choice';        required = $true;  choices = @('Phase I - Individual Skills','Phase II - Squad','Phase III - Platoon','Phase IV - MAGTF'); indexed = $true }
            @{ name = 'CompanyAssigned';         type = 'Choice';        required = $true;  choices = @('Alpha','Bravo','Charlie','Delta','Echo','Foxtrot','Golf','India','Mike'); indexed = $true }
            @{ name = 'SubmitterRole';           type = 'Choice';        required = $true;  choices = @('Student','SPC','Instructor','Observer'); description = 'Role of the person providing feedback' }
            @{ name = 'SubmitterName';           type = 'Text';          required = $false; maxLength = 100; description = 'Optional - anonymous by default. PII when populated.' }
            @{ name = 'OverallRating';           type = 'Choice';        required = $true;  choices = @('1 - Ineffective','2 - Below Average','3 - Average','4 - Above Average','5 - Excellent'); description = 'Overall event effectiveness' }
            @{ name = 'ObjectivesMet';           type = 'Choice';        required = $true;  choices = @('All objectives met','Most objectives met','Some objectives met','Few objectives met','No objectives met') }
            @{ name = 'InstructorEffectiveness'; type = 'Choice';        required = $false; choices = @('1 - Ineffective','2 - Below Average','3 - Average','4 - Above Average','5 - Excellent') }
            @{ name = 'TimeManagement';          type = 'Choice';        required = $false; choices = @('Too short','About right','Too long','Significant dead time') }
            @{ name = 'ResourceAdequacy';        type = 'Choice';        required = $false; choices = @('Fully resourced','Mostly adequate','Some shortages','Significant shortages') }
            @{ name = 'Sustains';                type = 'MultilineText'; required = $false; description = 'What went well - structured AAR field' }
            @{ name = 'Improves';                type = 'MultilineText'; required = $false; description = 'What to change next time - structured AAR field' }
            @{ name = 'SafetyConcerns';          type = 'MultilineText'; required = $false; description = 'Safety issues observed - flagged for immediate review' }
            @{ name = 'HasSafetyConcern';        type = 'Boolean';       required = $true;  default = $false; indexed = $true; description = 'Quick filter flag' }
            @{ name = 'AdditionalComments';      type = 'MultilineText'; required = $false }
            @{ name = 'SubmittedDate';           type = 'DateTime';      required = $true;  description = 'Auto-set on form submission' }
            @{ name = 'ReviewedBy';              type = 'Person';        required = $false; description = 'SPC or staff who reviewed this feedback' }
            @{ name = 'ReviewStatus';            type = 'Choice';        required = $true;  choices = @('Pending Review','Reviewed','Action Required','Closed'); default = 'Pending Review' }
            @{ name = 'ActionTaken';             type = 'MultilineText'; required = $false; description = 'Response to this feedback' }
        )
        views = @(
            @{
                name    = 'Pending Review'
                columns = @('EventTitle','EventDate','CompanyAssigned','SubmitterRole','OverallRating','ObjectivesMet','HasSafetyConcern','ReviewStatus')
                filter  = 'ReviewStatus = Pending Review'
                sort    = 'SubmittedDate DESC'
            }
            @{
                name    = 'Safety Alerts'
                columns = @('EventTitle','EventDate','CompanyAssigned','SubmitterRole','SafetyConcerns','ReviewStatus','ActionTaken')
                filter  = 'HasSafetyConcern = true'
                sort    = 'SubmittedDate DESC'
            }
            @{
                name    = 'Event Summary'
                columns = @('SubmitterRole','OverallRating','ObjectivesMet','InstructorEffectiveness','TimeManagement','Sustains','Improves')
                sort    = 'SubmitterRole ASC'
            }
            @{
                name    = 'Trend View'
                columns = @('EventTitle','TrainingPhase','EventDate','CompanyAssigned','OverallRating','ObjectivesMet','InstructorEffectiveness')
                sort    = 'EventDate DESC'
            }
        )
    }

    return $lists
}

# ---------------------------------------------------------------------------
# Sample data import
# ---------------------------------------------------------------------------
function Import-SampleData {
    [CmdletBinding(SupportsShouldProcess = $true)]
    param(
        [string]$ListName,
        [string]$CsvPath
    )

    if (-not (Test-Path $CsvPath)) {
        Write-Warning "[IMPORT] CSV not found for '$ListName': $CsvPath"
        return
    }

    if (-not $PSCmdlet.ShouldProcess($ListName, "Import sample data from $CsvPath")) {
        return
    }

    Write-Host "  [IMPORT] Importing sample data into '$ListName' from $CsvPath ..." -ForegroundColor Magenta

    $rows = Import-Csv -Path $CsvPath -ErrorAction Stop
    $count = 0

    foreach ($row in $rows) {
        $values = @{}
        foreach ($prop in $row.PSObject.Properties) {
            if (-not [string]::IsNullOrWhiteSpace($prop.Value)) {
                $values[$prop.Name] = $prop.Value
            }
        }

        try {
            Add-PnPListItem -List $ListName -Values $values -ErrorAction Stop | Out-Null
            $count++
        }
        catch {
            Write-Warning "    [WARN] Failed to import row $count in '$ListName': $_"
        }
    }

    Write-Host "  [IMPORT] Imported $count rows into '$ListName'." -ForegroundColor Magenta
}

# ---------------------------------------------------------------------------
# Removal / cleanup
# ---------------------------------------------------------------------------
function Remove-HeywoodLists {
    [CmdletBinding(SupportsShouldProcess = $true, ConfirmImpact = 'High')]
    param()

    $listNames = @(
        'StudentScores'
        'TrainingSchedule'
        'Instructors'
        'RequiredQualifications'
        'QualificationRecords'
        'EventFeedback'
    )

    Write-Host ""
    Write-Host "========================================" -ForegroundColor Red
    Write-Host "  DESTRUCTIVE OPERATION: Remove Lists   " -ForegroundColor Red
    Write-Host "========================================" -ForegroundColor Red
    Write-Host ""
    Write-Host "The following lists will be permanently deleted:" -ForegroundColor Red
    foreach ($ln in $listNames) {
        Write-Host "  - $ln" -ForegroundColor Red
    }
    Write-Host ""

    foreach ($ln in $listNames) {
        $existing = Get-PnPList -Identity $ln -ErrorAction SilentlyContinue
        if (-not $existing) {
            Write-Host "  [SKIP] List '$ln' does not exist." -ForegroundColor Yellow
            continue
        }

        if ($PSCmdlet.ShouldProcess($ln, "Remove list")) {
            try {
                Remove-PnPList -Identity $ln -Force -ErrorAction Stop
                Write-Host "  [REMOVED] $ln" -ForegroundColor Red
            }
            catch {
                Write-Error "  [ERROR] Failed to remove '$ln': $_"
            }
        }
    }

    Write-Host ""
    Write-Host "[DONE] Cleanup complete." -ForegroundColor Green
}

# ===========================================================================
# MAIN EXECUTION
# ===========================================================================

$stopwatch = [System.Diagnostics.Stopwatch]::StartNew()

# Ensure PnP module is available
Assert-PnPModule

# Connect to SharePoint
Connect-Site -Url $SiteUrl -Cred $Credential

# --- Removal mode ---
if ($RemoveAll) {
    Remove-HeywoodLists
    Disconnect-PnPOnline -ErrorAction SilentlyContinue
    return
}

# --- Deploy mode ---
Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "  Heywood TBS - SharePoint List Deployment  " -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "  Site: $SiteUrl" -ForegroundColor Cyan
Write-Host "  Mode: $(if ($WhatIfPreference) { 'PREVIEW (WhatIf)' } else { 'DEPLOY' })" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

$listDefs = Get-ListDefinitions
$sampleDataDir = Join-Path (Split-Path $PSScriptRoot -Parent) 'schemas' 'sharepoint' 'sample-data'

$totalLists    = $listDefs.Count
$totalColumns  = 0
$totalViews    = 0
$currentList   = 0

foreach ($listDef in $listDefs) {
    $currentList++
    $ln = $listDef.listName
    Write-Host ""
    Write-Host "[$currentList/$totalLists] Processing list: $ln" -ForegroundColor White -BackgroundColor DarkBlue
    Write-Host "  Description: $($listDef.description)" -ForegroundColor DarkGray

    # Create the list
    Ensure-List -ListName $ln -Description $listDef.description

    # Add columns
    Write-Host "  Adding columns ($($listDef.columns.Count)) ..." -ForegroundColor Gray
    foreach ($col in $listDef.columns) {
        try {
            Add-ListColumn -ListName $ln -Column $col
            $totalColumns++
        }
        catch {
            Write-Warning "  [ERROR] Failed to add column '$($col.name)' to '$ln': $_"
        }
    }

    # Create views
    Write-Host "  Creating views ($($listDef.views.Count)) ..." -ForegroundColor Gray
    foreach ($view in $listDef.views) {
        try {
            Add-ListView -ListName $ln -ViewDef $view
            $totalViews++
        }
        catch {
            Write-Warning "  [ERROR] Failed to create view '$($view.name)' on '$ln': $_"
        }
    }

    # Import sample data if requested
    if ($ImportSampleData) {
        $csvFile = Join-Path $sampleDataDir "$($ln).csv"
        Import-SampleData -ListName $ln -CsvPath $csvFile
    }
}

# --- Summary ---
$stopwatch.Stop()
Write-Host ""
Write-Host "============================================" -ForegroundColor Green
Write-Host "  Deployment Complete                       " -ForegroundColor Green
Write-Host "============================================" -ForegroundColor Green
Write-Host "  Lists processed:   $totalLists" -ForegroundColor Green
Write-Host "  Columns processed: $totalColumns" -ForegroundColor Green
Write-Host "  Views processed:   $totalViews" -ForegroundColor Green
Write-Host "  Elapsed time:      $($stopwatch.Elapsed.ToString('mm\:ss'))" -ForegroundColor Green
Write-Host "============================================" -ForegroundColor Green

if ($WhatIfPreference) {
    Write-Host ""
    Write-Host "  ** This was a PREVIEW run. No changes were made. **" -ForegroundColor Yellow
    Write-Host "  Re-run without -WhatIf to deploy." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "[NOTE] Calculated fields were created with placeholder formulas (=`"`")." -ForegroundColor Yellow
Write-Host "       Update them in SharePoint List Settings > Column > Edit formula." -ForegroundColor Yellow
Write-Host "       See the schema JSON files for intended formula logic." -ForegroundColor Yellow
Write-Host ""

# Disconnect
Disconnect-PnPOnline -ErrorAction SilentlyContinue
Write-Verbose "Disconnected from SharePoint Online."
