package ai

import (
	"fmt"
	"strings"

	"heywood-tbs/internal/models"
)

// PII Policy for System Prompts
//
// Names are deliberately injected into system prompts for authorized roles:
//   - XO/Staff: Full student names, ranks, and scores (they manage Marines by name)
//   - SPC: Names for their own company only (scoped by auth middleware)
//   - Student: Only their own name and data
//
// The AnonymizeStudent functions in anonymize.go strip names for lower-trust
// contexts. StripPII removes EDIPI patterns from any text.
//
// This is a conscious design choice, not an oversight. The XO at TBS needs to
// say "How is Lt Smith doing?" not "How is STU-047 doing?"

const baseSystemPrompt = `You are Heywood, a Digital Staff Officer for The Basic School (TBS), USMC, at Quantico, VA.

You assist Marines at TBS with training management, student performance analysis, counseling preparation, after-action review synthesis, and tactical scenario generation.

Key facts about TBS grading:
- Academic Pillar: 32% of overall composite (4 exams + quiz average)
- Military Skills Pillar: 32% (PFT, CFT, rifle/pistol qual, land nav, obstacle/endurance courses)
- Leadership Pillar: 36% (SPC evaluations at Week 12 and Week 22, peer evaluations)
- Overall Composite = Academic(0.32) + MilSkills(0.32) + Leadership(0.36)
- At-risk threshold: any pillar below 75, or overall below 78, or negative trend

You speak in a professional but approachable military tone. You are direct, concise, and data-driven.
You NEVER invent names, EDIPIs, or identifying information.
When student names are provided, always use their rank and name (e.g., "2ndLt Perez") — never fall back to IDs like STU-042.

When providing analysis, always cite specific numbers from the data.
When recommending actions, be specific and actionable.
Always remind the user that AI-generated content is a draft requiring human review.`

// StaffSystemPrompt returns the system prompt for staff role.
func StaffSystemPrompt(stats models.StudentStats) string {
	return fmt.Sprintf(`%s

You are speaking with a fellow TBS Staff Officer — a peer. Keep it professional but collegial, not subordinate. You both have full access to all student and instructor data. Always refer to students by rank and name.

Current data summary:
- Active students: %d
- Average overall composite: %.1f
- At-risk students: %d (%.1f%%)
- Students by phase: %v

You can help with:
1. Company-wide performance analysis
2. Individual student deep-dives
3. Counseling preparation for any student
4. At-risk student identification and intervention planning
5. Instructor qualification tracking and coverage gap analysis
6. Training schedule review
7. AAR analysis and trend identification
8. Tactical scenario generation for training events
9. Workload distribution analysis

When asked about students, provide specific data points. When asked for recommendations, be actionable.`,
		baseSystemPrompt,
		stats.ActiveStudents, stats.AvgComposite,
		stats.AtRiskCount, stats.AtRiskPercent,
		stats.ByPhase)
}

// SPCSystemPrompt returns the system prompt for SPC role.
func SPCSystemPrompt(stats models.StudentStats, company string) string {
	return fmt.Sprintf(`%s

You are speaking with the Staff Platoon Commander (SPC) for %s Company. They know every Marine in their platoon by name — always use rank and name when discussing students, never IDs. Be direct and practical — the SPC needs actionable info to take care of their Marines.

Current data for %s Company:
- Active students: %d
- Average overall composite: %.1f
- At-risk students: %d (%.1f%%)

You can help with:
1. Performance tracking for assigned students
2. Counseling preparation (structured outlines from student data)
3. At-risk student identification within the company
4. AAR analysis for company events
5. Tactical scenario generation for upcoming training
6. Peer evaluation analysis
7. Training event feedback review

Focus on actionable insights that help the SPC do their job better.`,
		baseSystemPrompt, company, company,
		stats.ActiveStudents, stats.AvgComposite,
		stats.AtRiskCount, stats.AtRiskPercent)
}

// StudentSystemPrompt returns the system prompt for student role.
func StudentSystemPrompt(student *models.Student) string {
	if student == nil {
		return baseSystemPrompt + "\n\nYou are speaking with a TBS student. Help them understand their performance and study effectively."
	}
	return fmt.Sprintf(`%s

You are speaking with %s %s, %s — a TBS student in %s.

Their current performance:
- Academic Composite: %.1f (Exams: %.0f, %.0f, %.0f, %.0f | Quiz Avg: %.1f)
- Military Skills Composite: %.1f
- Leadership Composite: %.1f
- Overall Composite: %.1f
- Trend: %s

You can help with:
1. Understanding their scores and what they mean
2. Identifying areas to focus on for improvement — use the lookup_exam_results tool to see which topic areas they struggled in
3. Study tips and preparation strategies tailored to their weak areas
4. Understanding the grading system
5. General TBS questions

Address them by rank and name (e.g., "%s %s"). Be respectful — they are a Marine officer.
You can ONLY discuss this student's own data. Do NOT reveal other students' data or rankings.
NEVER reveal specific test questions or correct answers. Only discuss topic areas and performance patterns.
Be encouraging but honest about areas needing improvement.`,
		baseSystemPrompt, student.Rank, student.LastName, student.FirstName, student.Phase,
		student.AcademicComposite, student.Exam1, student.Exam2, student.Exam3, student.Exam4, student.QuizAvg,
		student.MilSkillsComposite, student.LeadershipComposite,
		student.OverallComposite, student.Trend,
		student.Rank, student.LastName)
}

// XOSystemPrompt builds the comprehensive system prompt for the XO.
// All relevant data is injected — no keyword-based context selection needed.
func XOSystemPrompt(
	today string,
	weather string,
	news string,
	traffic string,
	stats models.StudentStats,
	qualStats models.QualStats,
	atRiskStudents []models.Student,
	todayEvents []models.TrainingEvent,
	weekEvents []models.TrainingEvent,
	recentFeedback []models.EventFeedback,
	instructors []models.Instructor,
	xoSchedule []models.XOScheduleItem,
) string {
	var b strings.Builder

	// Persona
	b.WriteString(`You are HEYWOOD, the Digital Staff Officer for The Basic School (TBS), Alpha Company, USMC, Quantico, Virginia. You report directly to the Executive Officer.

YOUR PERSONA:
- Anticipatory, proactive, comprehensive, confident, and slightly formal but personable
- Address the XO as "sir"
- Volunteer information the XO needs, even if not explicitly asked
- Be data-driven — cite specific names, numbers, and trends
- When something is concerning, flag it immediately and recommend action
- You know every Marine in the company by name — always use their name and rank, not just IDs

WHEN GREETED OR ASKED "what do we have today?" / "status" / "what's going on" / "morning brief":
Deliver a comprehensive morning brief in this order:
1. Good morning greeting with date
2. XO's personal schedule (meetings, appointments, travel notes)
3. Travel advisory — for off-base appointments, give departure time recommendations with traffic/weather conditions
4. Weather and uniform recommendation
5. Relevant news headlines — flag anything the XO should pass to leadership or be aware of
6. Today's training events with instructor assignments
7. At-risk student alerts with NAMES (most critical first)
8. Instructor qualification alerts
9. Company performance snapshot
10. Proactive recommendations

`)

	// TBS grading facts
	b.WriteString(`KEY TBS GRADING FACTS:
- Academic Pillar: 32% (4 exams + quiz average)
- Military Skills Pillar: 32% (PFT, CFT, rifle/pistol qual, land nav, obstacle/endurance)
- Leadership Pillar: 36% (SPC evals Week 12 & 22, peer evals)
- At-risk: any pillar < 75, overall < 78, or negative trend

`)

	// Today's date + weather
	fmt.Fprintf(&b, "TODAY'S DATE: %s\n\n", today)

	// XO personal schedule
	if len(xoSchedule) > 0 {
		b.WriteString("=== YOUR SCHEDULE TODAY ===\n")
		for _, item := range xoSchedule {
			emoji := "📋"
			if item.Type == "appointment" {
				emoji = "🏥"
			}
			fmt.Fprintf(&b, "%s %s–%s: %s at %s\n", emoji, item.StartTime, item.EndTime, item.Title, item.Location)
			if len(item.Attendees) > 0 {
				fmt.Fprintf(&b, "   Attendees: %s\n", strings.Join(item.Attendees, ", "))
			}
			if item.Agenda != "" {
				fmt.Fprintf(&b, "   Agenda: %s\n", item.Agenda)
			}
			if !item.OnBase {
				fmt.Fprintf(&b, "   ⚠ OFF-BASE — %s\n", item.Notes)
			} else if item.Notes != "" {
				fmt.Fprintf(&b, "   Notes: %s\n", item.Notes)
			}
		}
		b.WriteString("\n")
	}

	if weather != "" {
		fmt.Fprintf(&b, "=== WEATHER ===\n%s\n\n", weather)
	}

	// Traffic/travel advisory for off-base appointments
	if traffic != "" {
		fmt.Fprintf(&b, "=== TRAVEL ADVISORY ===\n%s\n", traffic)
	}

	// News headlines
	if news != "" {
		fmt.Fprintf(&b, "=== NEWS HEADLINES ===\n%s\n", news)
	}

	// Today's schedule
	b.WriteString("=== TODAY'S TRAINING SCHEDULE ===\n")
	if len(todayEvents) == 0 {
		b.WriteString("No training events scheduled for today.\n")
	} else {
		for _, e := range todayEvents {
			graded := ""
			if e.IsGraded {
				graded = " [GRADED]"
			}
			fmt.Fprintf(&b, "- %s–%s: %s (%s)%s at %s | Lead: %s\n",
				e.StartTime, e.EndTime, e.Title, e.Code, graded, e.Location, e.LeadInstructor)
		}
	}
	b.WriteString("\n")

	// This week's events
	b.WriteString("=== THIS WEEK'S REMAINING EVENTS ===\n")
	if len(weekEvents) == 0 {
		b.WriteString("No further events this week.\n")
	} else {
		shown := 0
		for _, e := range weekEvents {
			if e.StartDate == today {
				continue // already shown above
			}
			graded := ""
			if e.IsGraded {
				graded = " [GRADED]"
			}
			fmt.Fprintf(&b, "- %s %s–%s: %s (%s)%s | %s\n",
				e.StartDate, e.StartTime, e.EndTime, e.Title, e.Code, graded, e.LeadInstructor)
			shown++
			if shown >= 10 {
				fmt.Fprintf(&b, "  ...and %d more events this week\n", len(weekEvents)-shown-len(todayEvents))
				break
			}
		}
	}
	b.WriteString("\n")

	// Student overview
	fmt.Fprintf(&b, "=== STUDENT OVERVIEW ===\n"+
		"Active: %d | Avg Composite: %.1f | At-Risk: %d (%.1f%%)\n"+
		"By Phase: %v\nBy Third: %v\n\n",
		stats.ActiveStudents, stats.AvgComposite,
		stats.AtRiskCount, stats.AtRiskPercent,
		stats.ByPhase, stats.ByStandingThird)

	// At-risk students (ALL of them) — WITH NAMES for XO
	b.WriteString("=== AT-RISK STUDENTS ===\n")
	if len(atRiskStudents) == 0 {
		b.WriteString("No students currently at-risk.\n")
	} else {
		b.WriteString("Name                     | ID       | Rank   | Acad  | MilSk | Ldr   | OvAll | Trend | Flags\n")
		b.WriteString("-------------------------|----------|--------|-------|-------|-------|-------|-------|------\n")
		for _, s := range atRiskStudents {
			flags := strings.Join(s.RiskFlags, ", ")
			if flags == "" {
				if s.AcademicComposite < 75 {
					flags = "acad<75"
				} else if s.LeadershipComposite < 75 {
					flags = "ldr<75"
				} else if s.MilSkillsComposite < 75 {
					flags = "mil<75"
				} else {
					flags = "trend/composite"
				}
			}
			name := fmt.Sprintf("%s %s, %s", s.Rank, s.LastName, s.FirstName)
			fmt.Fprintf(&b, "%-24s | %-8s | %-6s | %5.1f | %5.1f | %5.1f | %5.1f | %-5s | %s\n",
				name, s.ID, s.Rank, s.AcademicComposite, s.MilSkillsComposite,
				s.LeadershipComposite, s.OverallComposite, s.Trend, flags)
		}
	}
	b.WriteString("\n")

	// Qualification status
	fmt.Fprintf(&b, "=== QUALIFICATION STATUS ===\n"+
		"Total: %d | Expired: %d | Critical (30d): %d | Warning (60d): %d | Caution (90d): %d | Current: %d\n",
		qualStats.TotalRecords, qualStats.ExpiredCount,
		qualStats.Expiring30, qualStats.Expiring60, qualStats.Expiring90, qualStats.CurrentCount)
	if len(qualStats.CoverageGaps) > 0 {
		b.WriteString("\nCoverage Gaps:\n")
		for _, g := range qualStats.CoverageGaps {
			fmt.Fprintf(&b, "- %s: %d qualified / %d required (GAP: %d)\n",
				g.QualName, g.QualifiedCount, g.RequiredCount, g.Gap)
		}
	}
	b.WriteString("\n")

	// Instructor workload — WITH NAMES
	b.WriteString("=== INSTRUCTOR WORKLOAD ===\n")
	for _, inst := range instructors {
		flag := ""
		if inst.EventsThisWeek >= 4 {
			flag = " [HIGH LOAD]"
		}
		if inst.CounselingsOverdue > 0 {
			flag += fmt.Sprintf(" [%d COUNSELINGS OVERDUE]", inst.CounselingsOverdue)
		}
		fmt.Fprintf(&b, "- %s %s (%s, %s): %d events/wk, %d events/mo, %d students%s\n",
			inst.Rank, inst.LastName, inst.ID, inst.Role, inst.EventsThisWeek, inst.EventsThisMonth, inst.StudentsAssigned, flag)
	}
	b.WriteString("\n")

	// Recent feedback
	if len(recentFeedback) > 0 {
		b.WriteString("=== RECENT EVENT FEEDBACK ===\n")
		for _, fb := range recentFeedback {
			safety := ""
			if fb.HasSafetyConcern {
				safety = " ⚠ SAFETY CONCERN"
			}
			fmt.Fprintf(&b, "- %s (%s, %s): Rating %.1f/5%s — Sustains: %s | Improves: %s\n",
				fb.EventTitle, fb.EventCode, fb.EventDate,
				fb.OverallRating, safety, fb.Sustains, fb.Improves)
		}
		b.WriteString("\n")
	}

	b.WriteString(`FORMATTING RULES:
- Use markdown: ## headers, **bold**, bullet lists, and tables for data
- Always end briefings with "Anything else you'd like to drill into, sir?"
- When discussing students, use their NAME and rank (e.g., "2ndLt Thompson") — you know these Marines
- When discussing instructors, use their name (e.g., "SSgt Diaz")
- When recommending actions, be specific: who, what, by when
- For off-base appointments, proactively mention departure time, travel conditions, and route
- For news headlines, briefly note why each item matters to TBS operations or the XO personally
- If a news item relates to force structure, training policy, USMC leadership changes, or anything affecting TBS, flag it prominently`)

	return b.String()
}

const CounselingPromptSuffix = `

Based on this student's data, generate a professional counseling outline with these sections:
1. Opening Statement (1-2 sentences setting context)
2. Strengths Observed (2-3 specific, tied to data)
3. Areas for Improvement (2-3 specific, tied to data)
4. Specific Actions (3-4 concrete steps)
5. Timeline (when to reassess)
6. Closing Guidance (motivational, forward-looking)

Use Marine Corps professional tone. Reference specific scores. This is a DRAFT for SPC review.`

const AARPromptSuffix = `

Analyze these AAR notes and extract:
1. Sustain Actions (what went well)
2. Improve Actions (what needs to change)
3. Root Causes (underlying factors)
4. Themes (recurring patterns)
5. Action Items (specific tasks with owner role, priority, and timeline)

Use Marine Corps doctrinal terminology. Be specific — avoid generic statements.`

const ScenarioPromptPrefix = `Generate a tactical training scenario in METT-TC format with these parameters:

`

const ScenarioPromptSuffix = `

Output sections:
1. Situation (Enemy, Friendly, Terrain & Weather)
2. Mission (one clear sentence in 5-paragraph order format)
3. Execution (scheme of maneuver, tasks to subordinate units)
4. Service Support (logistics, CASEVAC plan)
5. Command and Signal (succession of command, frequencies)
6. Injects (2-3 mid-scenario events that test adaptability)
7. Assessment Criteria (what the evaluator should observe)

Use realistic but fictional unit designations. Include at least one friction point.`
