# Prompt 06: PIA Screening for Student Data

**Use:** Determine if a Privacy Impact Assessment is required for a new TBS data collection
**Platform:** GenAI.mil
**Complexity:** Intermediate
**Source:** Adapted from EDD PIA Threshold Analysis (#2)

---

## Prompt

```
Answer these questions about a planned tool/data collection at The Basic School to determine PIA requirements:

Tool description: [TOOL_DESCRIPTION — e.g., "Power BI dashboard tracking student performance across Leadership, Military Skills, and Academics pillars"]

1. Will it collect PII from individuals? [YES_OR_NO]
   (PII at TBS includes: student names, EDIPIs, ranks, class standing, grades, peer evaluation data, fitness scores)
2. If yes, from how many individuals? [APPROXIMATE_NUMBER — e.g., "200 per company, up to 1,600 across TBS"]
3. Will it create a new database/repository containing PII? [YES_OR_NO]
4. Will PII be retrievable by individual identifier (name, EDIPI)? [YES_OR_NO]
5. Will it share PII with external parties outside TBS? [YES_OR_NO]
6. Will it use PII for a new purpose not covered by existing TBS training records? [YES_OR_NO]
7. Will it collect PII about minors (under 18)? [NO — TBS students are commissioned officers]
8. Will it collect sensitive PII (SSN, financial, medical/PHI, biometric)? [YES_OR_NO]
   (If medical/PHI — this requires MHS GENESIS coordination and DHA data sharing agreement)

Based on the answers, tell me:
- Whether a PIA is likely required (reference DoDI 5400.16 and DD Form 2930)
- What level of privacy review I need (threshold analysis only, full PIA, or SORN review)
- What I should discuss with the TBS Privacy Officer
- Recommended data minimization steps (can I use aggregate data instead of individual? Can I anonymize?)

Note: If the tool only handles aggregate, non-identifiable data (e.g., "Company A average academic score: 87%"), a PIA is likely NOT required. Individual-level data with identifiers triggers PIA requirements.
```

---

## Notes
- Complete this for EVERY new data source integration (Phase 1 through Phase 4)
- TBS student grades and class standing ARE PII — they are retrievable by individual
- PT scores linked to individuals ARE PII
- The goal is always to minimize PII: use aggregate data when possible
