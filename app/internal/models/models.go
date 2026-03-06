package models

// Student represents a TBS student officer with full grading data.
type Student struct {
	ID                 string   `json:"id"`
	EDIPI              string   `json:"edipi"`
	Rank               string   `json:"rank"`
	LastName           string   `json:"lastName"`
	FirstName          string   `json:"firstName"`
	Company            string   `json:"company"`
	Platoon            string   `json:"platoon"`
	SPC                string   `json:"spc"`
	ClassNumber        string   `json:"classNumber"`
	ClassStartDate     string   `json:"classStartDate"`
	Phase              string   `json:"phase"`
	Exam1              float64  `json:"exam1"`
	Exam2              float64  `json:"exam2"`
	Exam3              float64  `json:"exam3"`
	Exam4              float64  `json:"exam4"`
	QuizAvg            float64  `json:"quizAvg"`
	AcademicComposite  float64  `json:"academicComposite"`
	PFTScore           int      `json:"pftScore"`
	CFTScore           int      `json:"cftScore"`
	RifleQual          string   `json:"rifleQual"`
	PistolQual         string   `json:"pistolQual"`
	LandNavDay         string   `json:"landNavDay"`
	LandNavNight       string   `json:"landNavNight"`
	LandNavWritten     float64  `json:"landNavWritten"`
	ObstacleCourse     string   `json:"obstacleCourse"`
	EnduranceCourse    string   `json:"enduranceCourse"`
	MilSkillsComposite float64  `json:"milSkillsComposite"`
	LeadershipWeek12   float64  `json:"leadershipWeek12"`
	LeadershipWeek22   float64  `json:"leadershipWeek22"`
	PeerEvalWeek12     float64  `json:"peerEvalWeek12"`
	PeerEvalWeek22     float64  `json:"peerEvalWeek22"`
	LeadershipComposite float64 `json:"leadershipComposite"`
	OverallComposite   float64  `json:"overallComposite"`
	ClassStandingThird string   `json:"classStandingThird"`
	CompanyRank        int      `json:"companyRank"`
	Trend              string   `json:"trend"`
	AtRisk             bool     `json:"atRisk"`
	RiskFlags          []string `json:"riskFlags"`
	Status             string   `json:"status"`
	Notes              string   `json:"notes"`
}

// Instructor represents a TBS instructor or staff member.
type Instructor struct {
	ID                 string `json:"id"`
	EDIPI              string `json:"edipi"`
	LastName           string `json:"lastName"`
	FirstName          string `json:"firstName"`
	Rank               string `json:"rank"`
	Role               string `json:"role"`
	Company            string `json:"company"`
	Platoon            string `json:"platoon"`
	ClassNumber        string `json:"classNumber"`
	DateAssigned       string `json:"dateAssigned"`
	PRD                string `json:"prd"`
	StudentsAssigned   int    `json:"studentsAssigned"`
	EventsThisWeek     int    `json:"eventsThisWeek"`
	EventsThisMonth    int    `json:"eventsThisMonth"`
	CounselingsOverdue int    `json:"counselingsOverdue"`
	Status             string `json:"status"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
	Notes              string `json:"notes"`
}

// TrainingEvent represents a scheduled training event.
type TrainingEvent struct {
	ID                      string  `json:"id"`
	Title                   string  `json:"title"`
	Code                    string  `json:"code"`
	Phase                   string  `json:"phase"`
	Category                string  `json:"category"`
	GradePillar             string  `json:"gradePillar"`
	IsGraded                bool    `json:"isGraded"`
	StartDate               string  `json:"startDate"`
	EndDate                 string  `json:"endDate"`
	StartTime               string  `json:"startTime"`
	EndTime                 string  `json:"endTime"`
	DurationHours           float64 `json:"durationHours"`
	Location                string  `json:"location"`
	Company                 string  `json:"company"`
	ClassNumber             string  `json:"classNumber"`
	LeadInstructor          string  `json:"leadInstructor"`
	SupportInstructors      string  `json:"supportInstructors"`
	InstructorCountRequired int     `json:"instructorCountRequired"`
	PrerequisiteEvents      string  `json:"prerequisiteEvents"`
	SpecialEquipment        string  `json:"specialEquipment"`
	Status                  string  `json:"status"`
	WeatherContingency      string  `json:"weatherContingency"`
	Notes                   string  `json:"notes"`
}

// Qualification is a reference entry for a qualification type.
type Qualification struct {
	ID                string `json:"id"`
	Code              string `json:"code"`
	Name              string `json:"name"`
	Category          string `json:"category"`
	IssuingAuthority  string `json:"issuingAuthority"`
	ValidityMonths    int    `json:"validityMonths"`
	RenewalProcess    string `json:"renewalProcess"`
	RequiredForEvents string `json:"requiredForEvents"`
	MinimumPerEvent   int    `json:"minimumPerEvent"`
	OrderReference    string `json:"orderReference"`
	Status            string `json:"status"`
	Notes             string `json:"notes"`
}

// QualRecord is an individual instructor's qualification record.
type QualRecord struct {
	ID                  string `json:"id"`
	InstructorEDIPI     string `json:"instructorEdipi"`
	InstructorName      string `json:"instructorName"`
	QualCode            string `json:"qualCode"`
	QualName            string `json:"qualName"`
	DateEarned          string `json:"dateEarned"`
	ExpirationDate      string `json:"expirationDate"`
	DaysUntilExpiration int    `json:"daysUntilExpiration"`
	ExpirationStatus    string `json:"expirationStatus"`
	CertificateNumber   string `json:"certificateNumber"`
	IssuedBy            string `json:"issuedBy"`
	RenewalStatus       string `json:"renewalStatus"`
	RenewalDate         string `json:"renewalDate"`
	Notes               string `json:"notes"`
}

// EventFeedback is feedback submitted for a training event.
type EventFeedback struct {
	ID                      string  `json:"id"`
	EventTitle              string  `json:"eventTitle"`
	EventCode               string  `json:"eventCode"`
	EventDate               string  `json:"eventDate"`
	Phase                   string  `json:"phase"`
	Company                 string  `json:"company"`
	SubmitterRole           string  `json:"submitterRole"`
	SubmitterName           string  `json:"submitterName"`
	OverallRating           float64 `json:"overallRating"`
	ObjectivesMet           string  `json:"objectivesMet"`
	InstructorEffectiveness string  `json:"instructorEffectiveness"`
	TimeManagement          string  `json:"timeManagement"`
	ResourceAdequacy        string  `json:"resourceAdequacy"`
	Sustains                string  `json:"sustains"`
	Improves                string  `json:"improves"`
	SafetyConcerns          string  `json:"safetyConcerns"`
	HasSafetyConcern        bool    `json:"hasSafetyConcern"`
	AdditionalComments      string  `json:"additionalComments"`
	SubmittedDate           string  `json:"submittedDate"`
	ReviewedBy              string  `json:"reviewedBy"`
	ReviewStatus            string  `json:"reviewStatus"`
	ActionTaken             string  `json:"actionTaken"`
}

// StudentStats holds aggregated KPI data for students.
type StudentStats struct {
	ActiveStudents    int                `json:"activeStudents"`
	AvgComposite      float64            `json:"avgComposite"`
	AtRiskCount       int                `json:"atRiskCount"`
	AtRiskPercent     float64            `json:"atRiskPercent"`
	ByPhase           map[string]int     `json:"byPhase"`
	ByStandingThird   map[string]int     `json:"byStandingThird"`
}

// QualStats holds aggregated qualification KPI data.
type QualStats struct {
	TotalRecords  int            `json:"totalRecords"`
	ExpiredCount  int            `json:"expiredCount"`
	Expiring30    int            `json:"expiring30"`
	Expiring60    int            `json:"expiring60"`
	Expiring90    int            `json:"expiring90"`
	CurrentCount  int            `json:"currentCount"`
	CoverageGaps  []CoverageGap `json:"coverageGaps"`
}

// CoverageGap identifies a qualification with insufficient coverage.
type CoverageGap struct {
	QualCode       string `json:"qualCode"`
	QualName       string `json:"qualName"`
	QualifiedCount int    `json:"qualifiedCount"`
	RequiredCount  int    `json:"requiredCount"`
	Gap            int    `json:"gap"`
}

// XOScheduleItem is a personal schedule entry for the XO (meetings, appointments).
type XOScheduleItem struct {
	ID        string   `json:"id"`
	Type      string   `json:"type"`      // "meeting" or "appointment"
	Title     string   `json:"title"`
	Date      string   `json:"date"`
	StartTime string   `json:"startTime"`
	EndTime   string   `json:"endTime"`
	Location  string   `json:"location"`
	OnBase    bool     `json:"onBase"`
	Latitude  float64  `json:"latitude,omitempty"`  // for off-base routing
	Longitude float64  `json:"longitude,omitempty"` // for off-base routing
	Attendees []string `json:"attendees"`
	Agenda    string   `json:"agenda"`
	Notes     string   `json:"notes"`
}

// ChatRequest is the request body for the chat endpoint.
type ChatRequest struct {
	Message string        `json:"message"`
	History []ChatMessage `json:"history"`
	Stream  bool          `json:"stream"`
}

// ChatMessage represents a single message in chat history.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse is the response from the chat endpoint.
type ChatResponse struct {
	Response string `json:"response"`
}

// AuthInfo holds the current user's role information.
type AuthInfo struct {
	Role      string `json:"role"`
	Company   string `json:"company"`
	StudentID string `json:"studentId,omitempty"`
	Name      string `json:"name"`
}
