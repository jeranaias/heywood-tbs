package ai

import (
	"fmt"
	"strings"

	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/models"
)

// MockGreeting returns an appropriate greeting based on role.
func MockGreeting(role string, stats models.StudentStats) string {
	switch role {
	case auth.RoleStaff:
		return fmt.Sprintf("Good morning. I'm Heywood, your digital staff officer for TBS.\n\n"+
			"Current status: %d active students, average composite %.1f. "+
			"%d students (%.1f%%) are flagged at-risk.\n\n"+
			"How can I help you today?",
			stats.ActiveStudents, stats.AvgComposite,
			stats.AtRiskCount, stats.AtRiskPercent)
	case auth.RoleSPC:
		return fmt.Sprintf("Good morning, sir/ma'am. I'm Heywood, here to help you manage your students.\n\n"+
			"Your company currently has %d students with an average composite of %.1f. "+
			"%d (%.1f%%) are flagged at-risk.\n\n"+
			"I can help with counseling prep, performance tracking, AAR analysis, or scenario generation. What do you need?",
			stats.ActiveStudents, stats.AvgComposite,
			stats.AtRiskCount, stats.AtRiskPercent)
	case auth.RoleStudent:
		return "Hey there. I'm Heywood, your study assistant at TBS. " +
			"I can help you understand your scores, identify areas to focus on, and prepare for upcoming evaluations. " +
			"What would you like to know?"
	default:
		return "I'm Heywood, the TBS digital staff officer. How can I help?"
	}
}

// MockAtRiskResponse generates a mock response for at-risk student queries.
func MockAtRiskResponse(students []models.Student) string {
	if len(students) == 0 {
		return "No students are currently flagged at-risk. All students are performing within acceptable parameters."
	}

	var b strings.Builder
	fmt.Fprintf(&b, "**%d students are currently flagged at-risk.** Here's the breakdown:\n\n", len(students))

	for i, s := range students {
		if i >= 10 {
			fmt.Fprintf(&b, "\n...and %d more. Would you like the full list?\n", len(students)-10)
			break
		}
		flags := strings.Join(s.RiskFlags, ", ")
		if flags == "" {
			if s.AcademicComposite < 75 {
				flags = "academic below 75"
			} else if s.LeadershipComposite < 75 {
				flags = "leadership below 75"
			} else if s.MilSkillsComposite < 75 {
				flags = "mil skills below 75"
			} else {
				flags = "composite/trend concern"
			}
		}
		fmt.Fprintf(&b, "- **%s** (%s): Overall %.1f, Trend: %s — %s\n",
			s.ID, s.Rank, s.OverallComposite, s.Trend, flags)
	}

	b.WriteString("\n**Recommended actions:**\n")
	b.WriteString("1. Prioritize counseling sessions for students with declining trends\n")
	b.WriteString("2. Review academic support options for those below 75 in academics\n")
	b.WriteString("3. Coordinate with range/field instructors for mil skills deficiencies\n")
	b.WriteString("\n*This is AI-generated analysis. Verify all data before taking action.*")
	return b.String()
}

// MockCounselingResponse generates a mock counseling outline.
func MockCounselingResponse(s *models.Student) string {
	var b strings.Builder
	b.WriteString("# Counseling Outline — Draft\n\n")

	b.WriteString("## 1. Opening Statement\n")
	fmt.Fprintf(&b, "This counseling addresses %s's performance through the current training phase (%s). "+
		"Overall composite stands at %.1f.\n\n", s.ID, s.Phase, s.OverallComposite)

	b.WriteString("## 2. Strengths Observed\n")
	if s.AcademicComposite >= 85 {
		fmt.Fprintf(&b, "- Strong academic performance with a %.1f composite\n", s.AcademicComposite)
	}
	if s.MilSkillsComposite >= 85 {
		fmt.Fprintf(&b, "- Excellent military skills composite of %.1f\n", s.MilSkillsComposite)
	}
	if s.LeadershipComposite >= 85 {
		fmt.Fprintf(&b, "- Outstanding leadership composite of %.1f\n", s.LeadershipComposite)
	}
	if s.PFTScore >= 280 {
		fmt.Fprintf(&b, "- First-class PFT score of %d\n", s.PFTScore)
	}
	if s.Trend == "up" {
		b.WriteString("- Positive performance trend indicates strong trajectory\n")
	}

	b.WriteString("\n## 3. Areas for Improvement\n")
	if s.AcademicComposite < 80 {
		fmt.Fprintf(&b, "- Academic composite of %.1f needs attention (target: 80+)\n", s.AcademicComposite)
	}
	if s.MilSkillsComposite < 80 {
		fmt.Fprintf(&b, "- Military skills composite of %.1f below target\n", s.MilSkillsComposite)
	}
	if s.LeadershipComposite < 80 {
		fmt.Fprintf(&b, "- Leadership composite of %.1f requires focus\n", s.LeadershipComposite)
	}
	if s.Trend == "down" {
		b.WriteString("- Negative performance trend — needs course correction\n")
	}

	b.WriteString("\n## 4. Specific Actions\n")
	b.WriteString("1. Schedule weekly check-ins to track improvement\n")
	if s.AcademicComposite < 80 {
		b.WriteString("2. Attend academic study sessions and form a study group\n")
	}
	if s.LeadershipComposite < 80 {
		b.WriteString("3. Seek additional billet opportunities and peer mentoring\n")
	}
	b.WriteString("4. Maintain physical fitness standards as a foundation for all pillars\n")

	b.WriteString("\n## 5. Timeline\n")
	b.WriteString("Reassess in 2 weeks. Expect measurable improvement on next graded event.\n")

	b.WriteString("\n## 6. Closing Guidance\n")
	fmt.Fprintf(&b, "%s has the potential to succeed at TBS. Focused effort on the identified areas "+
		"will put them on track for a strong finish. The Marine Corps invested in them for a reason.\n", s.ID)

	b.WriteString("\n---\n*AI-generated draft. Review and modify before use. No PII was sent to the AI service.*")
	return b.String()
}

// MockScenarioResponse generates a mock tactical scenario.
func MockScenarioResponse(phase, objective, terrain string) string {
	return fmt.Sprintf(`# Tactical Scenario — %s

## 1. Situation

**Enemy:** A reinforced enemy squad (8-10 personnel) has established a patrol base in the %s terrain vicinity of Grid 123456. They have been observed conducting route reconnaissance along MSR Copper. Enemy is equipped with small arms and at least one crew-served weapon. Last observed at 0430 this morning.

**Friendly:** Your platoon (reinforced, 43 Marines) is the main effort for 2nd Company. 1st Platoon is the supporting effort to your north. Company weapons section has one section of machine guns attached to your platoon.

**Terrain & Weather:** %s terrain with limited visibility in draws and ravines. Temperature 45°F, wind 5-10 mph from the NW. Sunset at 1745. No moon tonight — limited natural illumination.

## 2. Mission
At 0600, 2nd Platoon attacks to destroy the enemy patrol base vicinity Grid 123456 in order to deny the enemy freedom of movement along MSR Copper.

## 3. Execution
- Main effort: 2nd Squad conducts the assault through the objective
- Supporting effort: 1st Squad establishes a support-by-fire position on the high ground to the south
- 3rd Squad establishes a blocking position to the east to prevent enemy withdrawal
- Attached machine gun section reinforces the support-by-fire position
- Reserve: Platoon Sergeant retains one fire team as reserve

## 4. Service Support
- Ammunition: Combat load + one additional magazine per Marine
- CASEVAC: CCP established at the assault position; ambulance on standby at the company CP (Grid 122458)
- Priority of resupply: ammunition, water, medical

## 5. Command and Signal
- Platoon Commander with main effort (2nd Squad)
- Succession: Platoon Sergeant, 1st Squad Leader, 2nd Squad Leader
- Primary: Squad Net (Freq 42.350) / Alternate: Platoon Net (Freq 42.425)
- Challenge/Password: THUNDER/BOLT

## 6. Injects
1. **T+15 min:** 1st Squad reports contact — enemy OP engaged their movement. How do you adjust?
2. **T+25 min:** CASEVAC — one Marine from 2nd Squad hit during the assault. Continue mission or pause?
3. **T+40 min:** Enemy reinforcements (one fire team) observed approaching from the east toward 3rd Squad's blocking position.

## 7. Assessment Criteria
- Appropriate use of supporting arms and crew-served weapons
- CASEVAC procedures executed correctly under fire
- Decision-making under pressure when injects disrupt the plan
- Communication between squads during the assault
- %s

---
*AI-generated scenario. Instructor must validate tactical accuracy before use.*`, phase, terrain, terrain, objective)
}

// MockGeneralResponse generates a general mock response.
func MockGeneralResponse(message string) string {
	msg := strings.ToLower(message)

	if strings.Contains(msg, "hello") || strings.Contains(msg, "hi") || strings.Contains(msg, "hey") {
		return "Hello. I'm Heywood, the TBS digital staff officer. I have data on students, instructors, training schedules, and qualifications. How can I help you today?"
	}
	if strings.Contains(msg, "help") || strings.Contains(msg, "what can you do") {
		return "I can help with:\n\n" +
			"1. **Student Performance** — \"Who are the at-risk students?\" or \"Tell me about STU-042\"\n" +
			"2. **Counseling Prep** — \"Prepare a counseling outline for STU-042\"\n" +
			"3. **AAR Analysis** — \"Analyze these AAR notes: [paste notes]\"\n" +
			"4. **Scenario Generation** — \"Generate a Phase 2 deliberate attack scenario\"\n" +
			"5. **Instructor Quals** — \"Any qualifications expiring soon?\"\n" +
			"6. **Training Schedule** — \"What's on the schedule this week?\"\n\n" +
			"Just ask in natural language. I'll pull the relevant data and give you a data-driven response."
	}
	if strings.Contains(msg, "thank") {
		return "Roger that. Let me know if you need anything else."
	}

	return "I can help with student performance, counseling prep, AAR analysis, scenario generation, instructor qualifications, and training schedules. Could you be more specific about what you need?"
}
