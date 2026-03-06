package data

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"heywood-tbs/internal/models"

	"github.com/xuri/excelize/v2"
)

// ExcelReader reads structured data from Excel files using column mappings.
type ExcelReader struct {
	path string
}

// NewExcelReader creates a reader for the given .xlsx file.
func NewExcelReader(path string) *ExcelReader {
	return &ExcelReader{path: path}
}

// ReadHeaders returns the header row from the given sheet.
func (r *ExcelReader) ReadHeaders(sheet string) ([]string, error) {
	f, err := excelize.OpenFile(r.path)
	if err != nil {
		return nil, fmt.Errorf("open excel: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("read sheet %q: %w", sheet, err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("sheet %q is empty", sheet)
	}

	return rows[0], nil
}

// ListSheets returns all sheet names in the workbook.
func (r *ExcelReader) ListSheets() ([]string, error) {
	f, err := excelize.OpenFile(r.path)
	if err != nil {
		return nil, fmt.Errorf("open excel: %w", err)
	}
	defer f.Close()

	return f.GetSheetList(), nil
}

// PreviewSheet returns the first N data rows (after header) from a sheet.
func (r *ExcelReader) PreviewSheet(sheet string, n int) ([][]string, error) {
	f, err := excelize.OpenFile(r.path)
	if err != nil {
		return nil, fmt.Errorf("open excel: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("read sheet %q: %w", sheet, err)
	}

	// Skip header row
	if len(rows) <= 1 {
		return nil, nil
	}
	data := rows[1:]
	if len(data) > n {
		data = data[:n]
	}
	return data, nil
}

// ReadStudents reads student records from the given sheet using column mappings.
func (r *ExcelReader) ReadStudents(sheet string, mappings []ColumnMapping) ([]models.Student, error) {
	f, err := excelize.OpenFile(r.path)
	if err != nil {
		return nil, fmt.Errorf("open excel: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("read sheet %q: %w", sheet, err)
	}

	if len(rows) <= 1 {
		return nil, nil
	}

	// Build field->column index map
	fieldIdx := make(map[string]int)
	for _, m := range mappings {
		if m.FieldName != "" {
			fieldIdx[m.FieldName] = m.Column
		}
	}

	var students []models.Student
	for i, row := range rows[1:] { // skip header
		s := models.Student{}
		s.ID = cellStr(row, fieldIdx, "id")
		if s.ID == "" {
			s.ID = fmt.Sprintf("excel-%d", i+1)
		}
		s.EDIPI = cellStr(row, fieldIdx, "edipi")
		s.LastName = cellStr(row, fieldIdx, "lastName")
		s.FirstName = cellStr(row, fieldIdx, "firstName")
		s.Rank = cellStr(row, fieldIdx, "rank")
		s.Company = cellStr(row, fieldIdx, "company")
		s.Platoon = cellStr(row, fieldIdx, "platoon")
		s.SPC = cellStr(row, fieldIdx, "spc")
		s.ClassNumber = cellStr(row, fieldIdx, "classNumber")
		s.ClassStartDate = cellStr(row, fieldIdx, "classStartDate")
		s.Phase = cellStr(row, fieldIdx, "phase")
		s.Exam1 = cellFloat(row, fieldIdx, "exam1")
		s.Exam2 = cellFloat(row, fieldIdx, "exam2")
		s.Exam3 = cellFloat(row, fieldIdx, "exam3")
		s.Exam4 = cellFloat(row, fieldIdx, "exam4")
		s.QuizAvg = cellFloat(row, fieldIdx, "quizAvg")
		s.AcademicComposite = cellFloat(row, fieldIdx, "academicComposite")
		s.PFTScore = int(cellFloat(row, fieldIdx, "pftScore"))
		s.CFTScore = int(cellFloat(row, fieldIdx, "cftScore"))
		s.RifleQual = cellStr(row, fieldIdx, "rifleQual")
		s.PistolQual = cellStr(row, fieldIdx, "pistolQual")
		s.LandNavDay = cellStr(row, fieldIdx, "landNavDay")
		s.LandNavNight = cellStr(row, fieldIdx, "landNavNight")
		s.LandNavWritten = cellFloat(row, fieldIdx, "landNavWritten")
		s.ObstacleCourse = cellStr(row, fieldIdx, "obstacleCourse")
		s.EnduranceCourse = cellStr(row, fieldIdx, "enduranceCourse")
		s.MilSkillsComposite = cellFloat(row, fieldIdx, "milSkillsComposite")
		s.LeadershipWeek12 = cellFloat(row, fieldIdx, "leadershipWeek12")
		s.LeadershipWeek22 = cellFloat(row, fieldIdx, "leadershipWeek22")
		s.PeerEvalWeek12 = cellFloat(row, fieldIdx, "peerEvalWeek12")
		s.PeerEvalWeek22 = cellFloat(row, fieldIdx, "peerEvalWeek22")
		s.LeadershipComposite = cellFloat(row, fieldIdx, "leadershipComposite")
		s.OverallComposite = cellFloat(row, fieldIdx, "overallComposite")
		s.Trend = cellStr(row, fieldIdx, "trend")
		s.Status = cellStr(row, fieldIdx, "status")
		s.Notes = cellStr(row, fieldIdx, "notes")

		// Parse at-risk boolean
		riskStr := strings.ToLower(cellStr(row, fieldIdx, "atRisk"))
		s.AtRisk = riskStr == "true" || riskStr == "yes" || riskStr == "1" || riskStr == "y"

		// Parse risk flags (comma-separated)
		flagStr := cellStr(row, fieldIdx, "riskFlags")
		if flagStr != "" {
			for _, f := range strings.Split(flagStr, ",") {
				f = strings.TrimSpace(f)
				if f != "" {
					s.RiskFlags = append(s.RiskFlags, f)
				}
			}
		}

		if s.LastName != "" || s.FirstName != "" {
			students = append(students, s)
		}
	}

	slog.Info("Excel import: students parsed", "count", len(students), "sheet", sheet)
	return students, nil
}

// ReadInstructors reads instructor records from the given sheet using column mappings.
func (r *ExcelReader) ReadInstructors(sheet string, mappings []ColumnMapping) ([]models.Instructor, error) {
	f, err := excelize.OpenFile(r.path)
	if err != nil {
		return nil, fmt.Errorf("open excel: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("read sheet %q: %w", sheet, err)
	}

	if len(rows) <= 1 {
		return nil, nil
	}

	fieldIdx := make(map[string]int)
	for _, m := range mappings {
		if m.FieldName != "" {
			fieldIdx[m.FieldName] = m.Column
		}
	}

	var instructors []models.Instructor
	for i, row := range rows[1:] {
		inst := models.Instructor{}
		inst.ID = cellStr(row, fieldIdx, "id")
		if inst.ID == "" {
			inst.ID = fmt.Sprintf("excel-inst-%d", i+1)
		}
		inst.EDIPI = cellStr(row, fieldIdx, "edipi")
		inst.LastName = cellStr(row, fieldIdx, "lastName")
		inst.FirstName = cellStr(row, fieldIdx, "firstName")
		inst.Rank = cellStr(row, fieldIdx, "rank")
		inst.Role = cellStr(row, fieldIdx, "role")
		inst.Company = cellStr(row, fieldIdx, "company")
		inst.Platoon = cellStr(row, fieldIdx, "platoon")
		inst.ClassNumber = cellStr(row, fieldIdx, "classNumber")
		inst.DateAssigned = cellStr(row, fieldIdx, "dateAssigned")
		inst.PRD = cellStr(row, fieldIdx, "prd")
		inst.StudentsAssigned = int(cellFloat(row, fieldIdx, "studentsAssigned"))
		inst.Status = cellStr(row, fieldIdx, "status")
		inst.Phone = cellStr(row, fieldIdx, "phone")
		inst.Email = cellStr(row, fieldIdx, "email")
		inst.Notes = cellStr(row, fieldIdx, "notes")

		if inst.LastName != "" || inst.FirstName != "" {
			instructors = append(instructors, inst)
		}
	}

	slog.Info("Excel import: instructors parsed", "count", len(instructors), "sheet", sheet)
	return instructors, nil
}

// ---- helpers ----

func cellStr(row []string, fieldIdx map[string]int, field string) string {
	idx, ok := fieldIdx[field]
	if !ok || idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}

func cellFloat(row []string, fieldIdx map[string]int, field string) float64 {
	s := cellStr(row, fieldIdx, field)
	if s == "" {
		return 0
	}
	// Remove percent signs, commas
	s = strings.TrimSuffix(s, "%")
	s = strings.ReplaceAll(s, ",", "")
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
