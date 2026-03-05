# Prompt 03: Student Performance User Stories

**Use:** Define features for Heywood tools by role — what each user type needs
**Platform:** GenAI.mil
**Complexity:** Beginner
**Source:** Adapted from EDD User Story Generator (#3)

---

## Prompt

```
I'm building a training support tool for The Basic School called "Heywood." It will serve as a role-based AI agent for TBS.

Tool description: [DESCRIBE_THE_SPECIFIC_CAPABILITY — e.g., "a student performance dashboard," "an AAR processing system," "a counseling preparation tool"]

The users are:
- TBS Staff (S-3, Executive Officer, Commanding Officer): Oversee all 8 training companies, need aggregate visibility and command-level reporting
- Staff Platoon Commanders (SPCs): Responsible for evaluating and developing ~25 students in their platoon, need individual student data and counseling tools
- Students (Lieutenants/WOs): Need visibility into their own performance and study resources, should NOT see other students' data or instructor tools

For each user type, generate 3-5 user stories in this format:
"As a [user type], I want to [action] so that [benefit]."

Focus on the most common daily/weekly tasks, not edge cases. Prioritize each story as:
- **Must Have** — Tool is useless without this
- **Should Have** — Significantly improves value
- **Nice to Have** — Can wait for a later version

Also flag any stories that involve PII and note the data handling requirement.
```

---

## Notes
- Use this prompt when defining requirements for any new Heywood component
- The three-role structure (Staff / SPC / Student) is consistent across all Heywood features
