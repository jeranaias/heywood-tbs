package data

import "strings"

// ColumnMapping maps a spreadsheet column header to a model field name.
type ColumnMapping struct {
	Column    int    `json:"column"`    // 0-based column index
	Header    string `json:"header"`    // original header text from file
	FieldName string `json:"fieldName"` // target model field name
	AutoMatch bool   `json:"autoMatch"` // true if auto-detected, false if manually mapped
}

// ColumnMapTemplate is a saved mapping template for reuse.
type ColumnMapTemplate struct {
	Name       string          `json:"name"`
	DataType   string          `json:"dataType"` // "students", "instructors", "schedule", etc.
	Mappings   []ColumnMapping `json:"mappings"`
	CreatedAt  string          `json:"createdAt"`
}

// fieldAliases maps model field names to known column header variations.
// All aliases are lowercase for case-insensitive matching.
var fieldAliases = map[string][]string{
	// Student fields
	"id":                  {"id", "student id", "studentid"},
	"edipi":               {"edipi", "dod id", "dodid", "dod_id"},
	"lastName":            {"last name", "lastname", "surname", "lname", "last"},
	"firstName":           {"first name", "firstname", "fname", "first"},
	"rank":                {"rank", "grade", "pay grade", "paygrade"},
	"company":             {"company", "co", "coy", "unit"},
	"platoon":             {"platoon", "plt", "pltn"},
	"spc":                 {"spc", "staff platoon commander"},
	"classNumber":         {"class", "class number", "classnumber", "class #", "class no"},
	"classStartDate":      {"class start", "start date", "classstart"},
	"phase":               {"phase", "training phase"},
	"exam1":               {"exam 1", "exam1", "test 1", "test1"},
	"exam2":               {"exam 2", "exam2", "test 2", "test2"},
	"exam3":               {"exam 3", "exam3", "test 3", "test3"},
	"exam4":               {"exam 4", "exam4", "test 4", "test4"},
	"quizAvg":             {"quiz avg", "quizavg", "quiz average", "quizzes"},
	"academicComposite":   {"academic", "academic composite", "academic score"},
	"pftScore":            {"pft", "pft score"},
	"cftScore":            {"cft", "cft score"},
	"rifleQual":           {"rifle", "rifle qual", "rifle qualification"},
	"pistolQual":          {"pistol", "pistol qual", "pistol qualification"},
	"landNavDay":          {"land nav day", "day land nav"},
	"landNavNight":        {"land nav night", "night land nav"},
	"landNavWritten":      {"land nav written", "written land nav"},
	"obstacleCourse":      {"obstacle course", "o-course", "ocourse"},
	"enduranceCourse":     {"endurance course", "e-course", "ecourse"},
	"milSkillsComposite":  {"mil skills", "military skills", "milskills composite"},
	"leadershipWeek12":    {"leadership wk12", "leadership week 12", "ldr wk12"},
	"leadershipWeek22":    {"leadership wk22", "leadership week 22", "ldr wk22"},
	"peerEvalWeek12":      {"peer eval wk12", "peer week 12"},
	"peerEvalWeek22":      {"peer eval wk22", "peer week 22"},
	"leadershipComposite": {"leadership", "leadership composite", "leadership score"},
	"overallComposite":    {"overall", "overall composite", "composite", "overall score"},
	"trend":               {"trend", "performance trend"},
	"atRisk":              {"at risk", "at-risk", "atrisk", "risk"},
	"riskFlags":           {"risk flags", "flags", "risk factors"},
	"status":              {"status"},
	"notes":               {"notes", "remarks", "comments"},

	// Instructor-specific fields
	"role":                {"role", "billet", "position"},
	"dateAssigned":        {"date assigned", "assigned", "start"},
	"prd":                 {"prd", "rotation date", "projected rotation date"},
	"studentsAssigned":    {"students assigned", "students"},
	"phone":               {"phone", "phone number", "telephone"},
	"email":               {"email", "email address", "mail"},
}

// AutoMapColumns takes a list of header strings and returns detected mappings.
// Headers that don't match any known alias are returned with FieldName = "".
func AutoMapColumns(headers []string) []ColumnMapping {
	mappings := make([]ColumnMapping, len(headers))
	usedFields := make(map[string]bool)

	for i, h := range headers {
		mappings[i] = ColumnMapping{
			Column: i,
			Header: h,
		}

		normalized := strings.TrimSpace(strings.ToLower(h))
		if normalized == "" {
			continue
		}

		// Try to find a matching field
		for field, aliases := range fieldAliases {
			if usedFields[field] {
				continue // don't double-map
			}
			for _, alias := range aliases {
				if normalized == alias {
					mappings[i].FieldName = field
					mappings[i].AutoMatch = true
					usedFields[field] = true
					break
				}
			}
			if mappings[i].FieldName != "" {
				break
			}
		}

		// Fuzzy fallback: check if header contains a field alias as substring
		if mappings[i].FieldName == "" {
			for field, aliases := range fieldAliases {
				if usedFields[field] {
					continue
				}
				for _, alias := range aliases {
					if len(alias) >= 4 && strings.Contains(normalized, alias) {
						mappings[i].FieldName = field
						mappings[i].AutoMatch = true
						usedFields[field] = true
						break
					}
				}
				if mappings[i].FieldName != "" {
					break
				}
			}
		}
	}

	return mappings
}

// AvailableFields returns the list of all known model field names
// that can be mapped to, grouped by data type.
func AvailableFields() map[string][]string {
	return map[string][]string{
		"students": {
			"id", "edipi", "lastName", "firstName", "rank", "company", "platoon",
			"spc", "classNumber", "classStartDate", "phase",
			"exam1", "exam2", "exam3", "exam4", "quizAvg", "academicComposite",
			"pftScore", "cftScore", "rifleQual", "pistolQual",
			"landNavDay", "landNavNight", "landNavWritten",
			"obstacleCourse", "enduranceCourse", "milSkillsComposite",
			"leadershipWeek12", "leadershipWeek22", "peerEvalWeek12", "peerEvalWeek22",
			"leadershipComposite", "overallComposite",
			"trend", "atRisk", "riskFlags", "status", "notes",
		},
		"instructors": {
			"id", "edipi", "lastName", "firstName", "rank", "role",
			"company", "platoon", "classNumber", "dateAssigned", "prd",
			"studentsAssigned", "status", "phone", "email", "notes",
		},
	}
}
