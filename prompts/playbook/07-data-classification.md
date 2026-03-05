# Prompt 07: TBS Data Classification

**Use:** Identify PII and data sensitivity in TBS data fields before building tools
**Platform:** GenAI.mil
**Complexity:** Beginner
**Source:** Adapted from EDD Data Classification (#1)

---

## Prompt

```
I'm building a tool for The Basic School that will collect/store/process the following data:

[LIST_ALL_DATA_FIELDS_YOUR_TOOL_WILL_HANDLE]

(List field names and data types only — do not paste actual data values or student records into this prompt.)

Example fields for TBS tools:
- Student class standing (rank order within company)
- Leadership evaluation score (Week 12, Week 22)
- PFT/CFT composite scores
- Land navigation pass/fail
- Peer evaluation rankings
- SPC narrative comments
- Company assignment
- Graduation date

For each field I listed, tell me:
1. Is this PII (Personally Identifiable Information)?
2. If yes, what category? (Direct identifier, quasi-identifier, sensitive PII)
3. Could this field be combined with others to identify an individual?
4. What's the minimum classification level required? (Unclassified, CUI, or higher)
5. Can this field be anonymized or aggregated without losing its usefulness?

Then summarize:
- Total PII fields identified
- Whether a Privacy Impact Assessment is likely required
- Recommended data minimization (fields I could remove, anonymize, or aggregate)
- Which platform this data can live on (SharePoint on MCEN for unclassified, GenAI.mil for CUI, neither for PHI/classified)

Reference DoDI 5400.16 and NIST 800-53 privacy controls for PII definitions.
```

---

## Notes
- Run this BEFORE designing any SharePoint list or Dataverse table
- Student names + grades = PII. Aggregated company averages without names = not PII.
- Consult the TBS Privacy Officer for edge cases
