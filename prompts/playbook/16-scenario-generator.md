# Prompt 16: Tactical Scenario Generator

**Use:** Generate tactical training scenarios calibrated to TBS phase and difficulty
**Platform:** GenAI.mil
**Complexity:** Intermediate
**Source:** Custom for TBS (no direct EDD analog)

---

## Prompt

```
Generate a tactical training scenario for The Basic School.

PARAMETERS:
- Phase: [Phase I Individual / Phase II Squad / Phase III Platoon / Phase IV MAGTF]
- Unit size: [Fire team / Squad / Platoon / Company]
- Training objective: [SPECIFIC_T&R_OBJECTIVE — e.g., "Conduct a squad attack," "Execute a platoon defense," "Plan and execute a company movement to contact"]
- Terrain: [DESCRIBE_OR_USE_TBS_DEFAULT — e.g., "Quantico training area, mixed hardwood forest, rolling terrain, limited road network"]
- Difficulty: [Introductory / Standard / Advanced / Stress test]
- Event type: [Tactical Decision Game (TDG) / Sand Table Exercise (STEX) / Field Exercise (FEX) / Classroom discussion]

CONSTRAINTS:
- Duration: [TIME_AVAILABLE — e.g., "30-minute TDG" or "3-day FEX"]
- Equipment available: [LIST_OR_DEFAULT — e.g., "Standard infantry weapons, no vehicles, no air support" or "Full platoon equipment plus SMAW and 60mm mortar"]
- Safety considerations: [ANY_RANGE_RESTRICTIONS_OR_TRAINING_AREA_LIMITS]

Generate the scenario with:

1. **Situation** (METT-TC format):
   - Mission: Assigned mission for the student unit
   - Enemy: OPFOR composition, disposition, and most likely/most dangerous course of action. Scale enemy capability to the difficulty level.
   - Terrain & Weather: Key terrain, observation/fields of fire, cover/concealment, obstacles, avenues of approach
   - Troops Available: Student unit composition and attachments
   - Time: Timeline and time constraints
   - Civil: Any civilian considerations

2. **Execution Requirements:**
   - What the student leader must produce (OPORD, FRAGO, verbal order, etc.)
   - Decision points the scenario forces
   - Dilemmas built into the scenario (there should be at least one hard tradeoff with no perfect answer)

3. **Evaluation Criteria:**
   - What "good" looks like at this difficulty level
   - Common mistakes to watch for
   - Key indicators the SPC should observe

4. **Instructor Notes:**
   - How to inject friction (change the situation mid-exercise)
   - Scaling: how to make it harder or easier on the fly
   - Post-exercise discussion questions for the AAR

Do NOT generate anything classified. Use generic OPFOR (e.g., "Centralian Revolutionary Force" — TBS's standard fictional adversary). All map references should be fictional grid coordinates.
```

---

## Notes
- Always have an SPC or instructor review generated scenarios before use
- Scenarios should reinforce the current POI block — don't introduce Phase III concepts in Phase I
- The "dilemma" requirement is key — TBS evaluates decision-making under ambiguity
- Difficulty calibration: Introductory = clear problem, one good answer. Stress test = ambiguous, time-pressured, multiple valid approaches.
