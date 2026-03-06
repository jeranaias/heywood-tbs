#!/usr/bin/env node
/**
 * generate-exam-data.js
 *
 * Reads students.json and generates realistic TBS Phase I Exam 1
 * question-level results for all 200 students.
 *
 * Output: app/data/exam-results.json
 */

const fs = require("fs");
const path = require("path");

// ─── Exam Blueprint ───────────────────────────────────────────────────────────
const EXAM_QUESTIONS = [
  // Marine Corps History (3)
  { questionNum: 1, topic: "Marine Corps History" },
  { questionNum: 2, topic: "Marine Corps History" },
  { questionNum: 3, topic: "Marine Corps History" },

  // Leadership Principles (3)
  { questionNum: 4, topic: "Leadership Principles" },
  { questionNum: 5, topic: "Leadership Principles" },
  { questionNum: 6, topic: "Leadership Principles" },

  // Warfighting (MCDP 1) (3)
  { questionNum: 7, topic: "Warfighting (MCDP 1)" },
  { questionNum: 8, topic: "Warfighting (MCDP 1)" },
  { questionNum: 9, topic: "Warfighting (MCDP 1)" },

  // Tactics (MCWP 3-01) (3)
  { questionNum: 10, topic: "Tactics (MCWP 3-01)" },
  { questionNum: 11, topic: "Tactics (MCWP 3-01)" },
  { questionNum: 12, topic: "Tactics (MCWP 3-01)" },

  // Land Navigation (2)
  { questionNum: 13, topic: "Land Navigation" },
  { questionNum: 14, topic: "Land Navigation" },

  // Weapons Systems (2)
  { questionNum: 15, topic: "Weapons Systems" },
  { questionNum: 16, topic: "Weapons Systems" },

  // Fire Support (2)
  { questionNum: 17, topic: "Fire Support" },
  { questionNum: 18, topic: "Fire Support" },

  // Communications (2)
  { questionNum: 19, topic: "Communications" },
  { questionNum: 20, topic: "Communications" },
];

const TOTAL_QUESTIONS = EXAM_QUESTIONS.length; // 20

// ─── Seeded random for reproducibility ────────────────────────────────────────
// Simple mulberry32 PRNG seeded from student ID
function mulberry32(seed) {
  return function () {
    seed |= 0;
    seed = (seed + 0x6d2b79f5) | 0;
    let t = Math.imul(seed ^ (seed >>> 15), 1 | seed);
    t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
  };
}

function seedFromId(id) {
  // "STU-042" -> 42 * 7919 (a prime) for decent spread
  const num = parseInt(id.replace("STU-", ""), 10);
  return num * 7919 + 31337;
}

// ─── Topic-area weakness profiles for at-risk / low-scoring students ──────────
// These define the probability of getting a question WRONG in a given topic
// when a student is "weak" in that area.  Higher = more likely to miss.
const WEAKNESS_PROFILES = [
  // Profile A: weak on Tactics + Fire Support + Land Nav
  {
    "Tactics (MCWP 3-01)": 0.7,
    "Fire Support": 0.65,
    "Land Navigation": 0.6,
    "Communications": 0.4,
    "Warfighting (MCDP 1)": 0.35,
    "Weapons Systems": 0.3,
    "Leadership Principles": 0.2,
    "Marine Corps History": 0.15,
  },
  // Profile B: weak on Warfighting + Leadership + Communications
  {
    "Warfighting (MCDP 1)": 0.7,
    "Leadership Principles": 0.65,
    "Communications": 0.6,
    "Tactics (MCWP 3-01)": 0.4,
    "Fire Support": 0.35,
    "Land Navigation": 0.25,
    "Weapons Systems": 0.2,
    "Marine Corps History": 0.15,
  },
  // Profile C: weak on Land Nav + Weapons + History
  {
    "Land Navigation": 0.7,
    "Weapons Systems": 0.65,
    "Marine Corps History": 0.6,
    "Fire Support": 0.4,
    "Communications": 0.35,
    "Tactics (MCWP 3-01)": 0.3,
    "Warfighting (MCDP 1)": 0.2,
    "Leadership Principles": 0.15,
  },
  // Profile D: weak across technical topics
  {
    "Fire Support": 0.7,
    "Communications": 0.65,
    "Weapons Systems": 0.6,
    "Tactics (MCWP 3-01)": 0.5,
    "Land Navigation": 0.4,
    "Warfighting (MCDP 1)": 0.25,
    "Leadership Principles": 0.2,
    "Marine Corps History": 0.15,
  },
];

// ─── Main generation logic ────────────────────────────────────────────────────

function generateStudentExam(student) {
  const rng = mulberry32(seedFromId(student.id));

  // Target number of correct answers from the student's exam1 percentage
  const targetCorrect = Math.round((student.exam1 / 100) * TOTAL_QUESTIONS);

  // Clamp to [0, 20]
  const clampedTarget = Math.max(0, Math.min(TOTAL_QUESTIONS, targetCorrect));

  // Determine if this student should exhibit a weakness pattern
  const isWeak = student.atRisk || student.exam1 < 75;

  // Pick a weakness profile deterministically based on student ID
  const profileIndex =
    parseInt(student.id.replace("STU-", ""), 10) % WEAKNESS_PROFILES.length;
  const weaknessProfile = WEAKNESS_PROFILES[profileIndex];

  // Generate per-question miss probabilities
  // For strong students: roughly uniform low miss probability
  // For weak students: shaped by the weakness profile
  const baseMissRate = 1 - clampedTarget / TOTAL_QUESTIONS; // overall miss rate

  const questionMissProbs = EXAM_QUESTIONS.map((q) => {
    if (isWeak) {
      // Blend the weakness profile with the base miss rate
      const topicBias = weaknessProfile[q.topic] || 0.3;
      // Scale the profile so the average miss rate matches the target
      return topicBias * (baseMissRate / 0.4); // 0.4 is ~average profile value
    } else {
      // Strong students: slight random variation around base miss rate
      return baseMissRate * (0.7 + rng() * 0.6);
    }
  });

  // We need exactly (TOTAL_QUESTIONS - clampedTarget) wrong answers.
  // Strategy: assign a "miss score" to each question, pick the top N to be wrong.
  const numWrong = TOTAL_QUESTIONS - clampedTarget;

  // Generate a weighted random score for each question
  const missScores = questionMissProbs.map((prob, i) => ({
    index: i,
    score: prob + rng() * 0.3, // add noise so it's not perfectly deterministic
  }));

  // Sort descending by score — highest scores are most likely to be wrong
  missScores.sort((a, b) => b.score - a.score);

  // The top numWrong questions are wrong
  const wrongSet = new Set(missScores.slice(0, numWrong).map((m) => m.index));

  // Build the question results
  const questions = EXAM_QUESTIONS.map((q, i) => ({
    questionNum: q.questionNum,
    topic: q.topic,
    correct: !wrongSet.has(i),
  }));

  const correct = TOTAL_QUESTIONS - numWrong;
  const score = parseFloat(((correct / TOTAL_QUESTIONS) * 100).toFixed(1));

  return {
    studentId: student.id,
    examNum: 1,
    score,
    correct,
    total: TOTAL_QUESTIONS,
    questions,
  };
}

// ─── Run ──────────────────────────────────────────────────────────────────────

const studentsPath = path.join(__dirname, "app", "data", "students.json");
const outputPath = path.join(__dirname, "app", "data", "exam-results.json");

const students = JSON.parse(fs.readFileSync(studentsPath, "utf8"));

console.log(`Read ${students.length} students from ${studentsPath}`);
console.log(`Generating Phase I Exam 1 results (${TOTAL_QUESTIONS} questions)...`);

const results = students.map((s) => generateStudentExam(s));

// Verify
const scoreDiffs = results.map((r, i) => {
  const diff = Math.abs(r.score - students[i].exam1);
  return { id: r.studentId, generated: r.score, original: students[i].exam1, diff };
});

const maxDiff = Math.max(...scoreDiffs.map((d) => d.diff));
const avgDiff = scoreDiffs.reduce((sum, d) => sum + d.diff, 0) / scoreDiffs.length;
const within2 = scoreDiffs.filter((d) => d.diff <= 2.5).length;
const outliers = scoreDiffs.filter((d) => d.diff > 5);

console.log(`\nScore accuracy:`);
console.log(`  Max difference from original exam1: ${maxDiff.toFixed(1)} pts`);
console.log(`  Average difference: ${avgDiff.toFixed(2)} pts`);
console.log(`  Within ±2.5 pts: ${within2}/${students.length} (${((within2 / students.length) * 100).toFixed(1)}%)`);

if (outliers.length > 0) {
  console.log(`  Outliers (>5 pts off):`);
  outliers.forEach((o) =>
    console.log(`    ${o.id}: generated=${o.generated}, original=${o.original}, diff=${o.diff.toFixed(1)}`)
  );
}

// Show topic miss patterns for at-risk students
const atRiskResults = results.filter((r, i) => students[i].atRisk);
if (atRiskResults.length > 0) {
  const topicMissCounts = {};
  atRiskResults.forEach((r) => {
    r.questions.forEach((q) => {
      if (!q.correct) {
        topicMissCounts[q.topic] = (topicMissCounts[q.topic] || 0) + 1;
      }
    });
  });
  console.log(`\nAt-risk students (${atRiskResults.length}) — topic miss distribution:`);
  Object.entries(topicMissCounts)
    .sort((a, b) => b[1] - a[1])
    .forEach(([topic, count]) => {
      const pct = ((count / atRiskResults.length) * 100).toFixed(1);
      console.log(`  ${topic}: ${count} misses (avg ${pct}% of at-risk students missed a Q in this topic)`);
    });
}

// Write output
fs.writeFileSync(outputPath, JSON.stringify(results, null, 2), "utf8");
console.log(`\nWrote ${results.length} exam results to ${outputPath}`);
