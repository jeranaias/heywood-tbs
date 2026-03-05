# Prompt 20: Post-Event Feedback Collection

**Use:** Design a Power App for collecting structured student feedback after training events
**Platform:** GenAI.mil
**Complexity:** Beginner
**Source:** Adapted from EDD After-Action Feedback Form (Training Tools #3)

---

## Prompt

```
Create a Power Apps canvas app for collecting structured feedback from TBS students after training events and exercises. Requirements:

1. SharePoint list "EventFeedback" columns:
   - EventName (text)
   - EventDate (date)
   - Phase (choice: Phase I, Phase II, Phase III, Phase IV)
   - EventType (choice: Classroom, Field Exercise, Physical Training, Evaluation, Sand Table, TDG)
   - Company (choice: Alpha through Golf, India)
   - SubmittedBy (person)
   - OverallRating (number 1-5)
   - ContentRelevance (number 1-5: "How relevant was this training to your development as an officer?")
   - InstructorEffectiveness (number 1-5: "How effective was the instructor/SPC?")
   - PaceRating (choice: Too Slow, Just Right, Too Fast)
   - Difficulty (choice: Too Easy, Appropriate, Too Hard)
   - MostValuable (multiline text: "What was the most valuable part of this event?")
   - LeastValuable (multiline text: "What was the least valuable or could be improved?")
   - Suggestions (multiline text: "What specific changes would you recommend?")
   - WouldRecommendToNextClass (yes/no)

2. Form layout:
   - Star-rating visual for the three numeric ratings (1-5 stars using icons)
   - Anonymous option: Toggle that replaces SubmittedBy with "Anonymous"
   - Required field validation: OverallRating and at least one text field must be filled
   - Confirmation screen: "Thank you — your feedback helps improve training for future classes."

3. Admin summary view (Staff/SPC access only):
   - Average ratings per event as bar charts
   - Response count and completion rate per event
   - Trend over time: Are ratings improving or declining through the phases?
   - Word frequency analysis from text responses (most common themes)
   - Comparison: Same event across different companies (are some SPCs getting different results?)

4. Data handling:
   - When "Anonymous" is selected, the record truly does not store the submitter identity
   - Individual responses are visible only to Staff role
   - SPCs see aggregate results for their company only

Start with the star-rating component, then the full form layout.
```

---

## Notes
- Anonymous feedback is critical for honest input — students must trust the system
- This data feeds into Prompt 15 (AAR Synthesis) for trend analysis
- Phase 3+: This data flows into Dataverse for cross-cycle analysis
