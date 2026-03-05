// Client-side mock AI response engine
// Replicates Go backend's mock.go responses for static/offline mode

import type { Student, ChatMessage } from './types'
import * as data from './staticData'

export function generateMockResponse(message: string, _history: ChatMessage[], role: string): string {
  const msg = message.toLowerCase()

  // Greeting / short messages
  if (msg.length < 20 && !msg.includes('risk') && !msg.includes('stu-')) {
    return mockGreeting(role)
  }

  // At-risk queries
  if (msg.includes('at risk') || msg.includes('at-risk') || msg.includes('struggling') || msg.includes('failing')) {
    const { students } = data.getAtRiskStudents()
    return mockAtRiskResponse(students)
  }

  // Counseling
  if (msg.includes('counsel')) {
    const id = extractStudentId(msg)
    if (id) {
      const student = data.getStudent(id)
      if (student) return mockCounselingResponse(student)
    }
    return 'I can prepare a counseling outline. Please specify a student ID (e.g., "counseling for STU-042").'
  }

  // Specific student lookup
  const studentId = extractStudentId(msg)
  if (studentId) {
    const student = data.getStudent(studentId)
    if (student) return mockStudentResponse(student)
    return `I couldn't find a student with ID ${studentId}. Student IDs range from STU-001 to STU-200.`
  }

  // Scenario generation
  if (msg.includes('scenario') || msg.includes('mett-tc') || msg.includes('attack') || msg.includes('defense') || msg.includes('patrol')) {
    return mockScenarioResponse()
  }

  // Schedule
  if (msg.includes('schedule') || msg.includes('this week') || msg.includes('training')) {
    return mockScheduleResponse()
  }

  // Qualifications
  if (msg.includes('qual') || msg.includes('expir') || msg.includes('instructor')) {
    return mockQualResponse()
  }

  // Help
  if (msg.includes('help') || msg.includes('what can you do') || msg.includes('capabilities')) {
    return mockHelpResponse()
  }

  // Thanks
  if (msg.includes('thank')) {
    return 'Roger that. Let me know if you need anything else.'
  }

  // General fallback
  return mockGeneralResponse()
}

function mockGreeting(role: string): string {
  const stats = data.getStudentStats()
  switch (role) {
    case 'staff':
      return `Good morning. I'm Heywood, your digital staff officer for TBS.\n\n` +
        `Current status: ${stats.activeStudents} active students, average composite ${stats.avgComposite.toFixed(1)}. ` +
        `${stats.atRiskCount} students (${stats.atRiskPercent.toFixed(1)}%) are flagged at-risk.\n\n` +
        `How can I help you today?`
    case 'spc':
      return `Good morning, sir/ma'am. I'm Heywood, here to help you manage your students.\n\n` +
        `Your company currently has ${stats.activeStudents} students with an average composite of ${stats.avgComposite.toFixed(1)}. ` +
        `${stats.atRiskCount} are flagged at-risk.\n\n` +
        `I can help with counseling prep, performance tracking, AAR analysis, or scenario generation. What do you need?`
    case 'student':
      return `Hey there. I'm Heywood, your study assistant at TBS. ` +
        `I can help you understand your scores, identify areas to focus on, and prepare for upcoming evaluations. ` +
        `What would you like to know?`
    default:
      return "I'm Heywood, the TBS digital staff officer. How can I help?"
  }
}

function mockAtRiskResponse(students: Student[]): string {
  if (students.length === 0) {
    return 'No students are currently flagged at-risk. All students are performing within acceptable parameters.'
  }

  let text = `**${students.length} students are currently flagged at-risk.** Here's the breakdown:\n\n`

  const shown = students.slice(0, 10)
  for (const s of shown) {
    let flags = s.riskFlags.join(', ')
    if (!flags) {
      if (s.academicComposite < 75) flags = 'academic below 75'
      else if (s.leadershipComposite < 75) flags = 'leadership below 75'
      else if (s.milSkillsComposite < 75) flags = 'mil skills below 75'
      else flags = 'composite/trend concern'
    }
    text += `- **${s.id}** (${s.rank}): Overall ${s.overallComposite.toFixed(1)}, Trend: ${s.trend} — ${flags}\n`
  }

  if (students.length > 10) {
    text += `\n...and ${students.length - 10} more. Would you like the full list?\n`
  }

  text += '\n**Recommended actions:**\n'
  text += '1. Prioritize counseling sessions for students with declining trends\n'
  text += '2. Review academic support options for those below 75 in academics\n'
  text += '3. Coordinate with range/field instructors for mil skills deficiencies\n'
  text += '\n*This is AI-generated analysis. Verify all data before taking action.*'
  return text
}

function mockCounselingResponse(s: Student): string {
  let text = '# Counseling Outline — Draft\n\n'

  text += '## 1. Opening Statement\n'
  text += `This counseling addresses ${s.id}'s performance through the current training phase (${s.phase}). `
  text += `Overall composite stands at ${s.overallComposite.toFixed(1)}.\n\n`

  text += '## 2. Strengths Observed\n'
  if (s.academicComposite >= 85) text += `- Strong academic performance with a ${s.academicComposite.toFixed(1)} composite\n`
  if (s.milSkillsComposite >= 85) text += `- Excellent military skills composite of ${s.milSkillsComposite.toFixed(1)}\n`
  if (s.leadershipComposite >= 85) text += `- Outstanding leadership composite of ${s.leadershipComposite.toFixed(1)}\n`
  if (s.pftScore >= 280) text += `- First-class PFT score of ${s.pftScore}\n`
  if (s.trend === 'up') text += '- Positive performance trend indicates strong trajectory\n'

  text += '\n## 3. Areas for Improvement\n'
  if (s.academicComposite < 80) text += `- Academic composite of ${s.academicComposite.toFixed(1)} needs attention (target: 80+)\n`
  if (s.milSkillsComposite < 80) text += `- Military skills composite of ${s.milSkillsComposite.toFixed(1)} below target\n`
  if (s.leadershipComposite < 80) text += `- Leadership composite of ${s.leadershipComposite.toFixed(1)} requires focus\n`
  if (s.trend === 'down') text += '- Negative performance trend — needs course correction\n'

  text += '\n## 4. Specific Actions\n'
  text += '1. Schedule weekly check-ins to track improvement\n'
  if (s.academicComposite < 80) text += '2. Attend academic study sessions and form a study group\n'
  if (s.leadershipComposite < 80) text += '3. Seek additional billet opportunities and peer mentoring\n'
  text += '4. Maintain physical fitness standards as a foundation for all pillars\n'

  text += '\n## 5. Timeline\n'
  text += 'Reassess in 2 weeks. Expect measurable improvement on next graded event.\n'

  text += '\n## 6. Closing Guidance\n'
  text += `${s.id} has the potential to succeed at TBS. Focused effort on the identified areas `
  text += 'will put them on track for a strong finish. The Marine Corps invested in them for a reason.\n'

  text += '\n---\n*AI-generated draft. Review and modify before use. No PII was sent to the AI service.*'
  return text
}

function mockStudentResponse(s: Student): string {
  return `## ${s.id} — ${s.rank} ${s.lastName}, ${s.firstName}\n\n` +
    `**Phase:** ${s.phase} | **Company:** ${s.company} | **SPC:** ${s.spc}\n\n` +
    `| Pillar | Score | Weight |\n|--------|-------|--------|\n` +
    `| Academic | ${s.academicComposite.toFixed(1)} | 32% |\n` +
    `| Mil Skills | ${s.milSkillsComposite.toFixed(1)} | 32% |\n` +
    `| Leadership | ${s.leadershipComposite.toFixed(1)} | 36% |\n` +
    `| **Overall** | **${s.overallComposite.toFixed(1)}** | |\n\n` +
    `**Trend:** ${s.trend} | **Standing:** ${s.classStandingThird} | ` +
    `**At Risk:** ${s.atRisk ? 'Yes — ' + (s.riskFlags.join(', ') || 'flagged') : 'No'}\n\n` +
    `**Exams:** E1: ${s.exam1.toFixed(1)}, E2: ${s.exam2.toFixed(1)}, E3: ${s.exam3.toFixed(1)}, E4: ${s.exam4.toFixed(1)} | Quiz Avg: ${s.quizAvg.toFixed(1)}\n` +
    `**PFT:** ${s.pftScore} | **CFT:** ${s.cftScore} | **Rifle:** ${s.rifleQual} | **Pistol:** ${s.pistolQual}\n\n` +
    `Would you like me to prepare a counseling outline for this student?`
}

function mockScenarioResponse(): string {
  return `# Tactical Scenario — Phase I Deliberate Attack

## 1. Situation

**Enemy:** A reinforced enemy squad (8-10 personnel) has established a patrol base in the wooded terrain vicinity of Grid 123456. They have been observed conducting route reconnaissance along MSR Copper. Enemy is equipped with small arms and at least one crew-served weapon. Last observed at 0430 this morning.

**Friendly:** Your platoon (reinforced, 43 Marines) is the main effort for 2nd Company. 1st Platoon is the supporting effort to your north. Company weapons section has one section of machine guns attached to your platoon.

**Terrain & Weather:** Wooded terrain with limited visibility in draws and ravines. Temperature 45°F, wind 5-10 mph from the NW. Sunset at 1745. No moon tonight — limited natural illumination.

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

---
*AI-generated scenario. Instructor must validate tactical accuracy before use.*`
}

function mockScheduleResponse(): string {
  const { events } = data.listSchedule()
  const upcoming = events.filter(e => e.status !== 'Complete').slice(0, 8)

  if (upcoming.length === 0) {
    return 'No upcoming events on the schedule. All scheduled events have been completed.'
  }

  let text = '## Upcoming Training Events\n\n'
  text += '| Date | Event | Category | Location |\n|------|-------|----------|----------|\n'
  for (const e of upcoming) {
    text += `| ${e.startDate} | ${e.title} | ${e.category} | ${e.location} |\n`
  }
  text += `\n${events.length} total events on the schedule. ${events.filter(e => e.status === 'Complete').length} completed.`
  return text
}

function mockQualResponse(): string {
  const stats = data.getQualStats()
  let text = '## Instructor Qualification Status\n\n'
  text += `- **Expired:** ${stats.expiredCount}\n`
  text += `- **Expiring within 30 days:** ${stats.expiring30}\n`
  text += `- **Expiring within 60 days:** ${stats.expiring60}\n`
  text += `- **Expiring within 90 days:** ${stats.expiring90}\n`
  text += `- **Current:** ${stats.currentCount}\n\n`

  if (stats.coverageGaps.length > 0) {
    text += '### Coverage Gaps\n\n'
    text += '| Qualification | Qualified | Required | Gap |\n|---------------|-----------|----------|-----|\n'
    for (const g of stats.coverageGaps) {
      text += `| ${g.qualName} | ${g.qualifiedCount} | ${g.requiredCount} | **-${g.gap}** |\n`
    }
    text += '\n**Action required:** Coordinate with WTBN for accelerated certifications on critical gaps.'
  } else {
    text += 'No critical coverage gaps identified.'
  }

  return text
}

function mockHelpResponse(): string {
  return `I can help with:\n\n` +
    `1. **Student Performance** — "Who are the at-risk students?" or "Tell me about STU-042"\n` +
    `2. **Counseling Prep** — "Prepare a counseling outline for STU-042"\n` +
    `3. **AAR Analysis** — "Analyze these AAR notes: [paste notes]"\n` +
    `4. **Scenario Generation** — "Generate a Phase 2 deliberate attack scenario"\n` +
    `5. **Instructor Quals** — "Any qualifications expiring soon?"\n` +
    `6. **Training Schedule** — "What's on the schedule this week?"\n\n` +
    `Just ask in natural language. I'll pull the relevant data and give you a data-driven response.`
}

function mockGeneralResponse(): string {
  return 'I can help with student performance, counseling prep, AAR analysis, scenario generation, ' +
    'instructor qualifications, and training schedules. Could you be more specific about what you need?'
}

function extractStudentId(msg: string): string | null {
  // Match "STU-042" or "student #42" or "student 42"
  const stuMatch = msg.match(/stu-(\d+)/i)
  if (stuMatch) {
    const num = parseInt(stuMatch[1], 10)
    return `STU-${String(num).padStart(3, '0')}`
  }
  const numMatch = msg.match(/student\s*#?\s*(\d+)/i)
  if (numMatch) {
    const num = parseInt(numMatch[1], 10)
    return `STU-${String(num).padStart(3, '0')}`
  }
  return null
}
