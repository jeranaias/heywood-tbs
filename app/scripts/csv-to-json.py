#!/usr/bin/env python3
"""Convert TBS sample CSV data to JSON for the Heywood Go API."""

import csv
import json
import os
import sys

SAMPLE_DIR = os.path.join(os.path.dirname(__file__), "..", "..", "schemas", "sharepoint", "sample-data")
OUT_DIR = os.path.join(os.path.dirname(__file__), "..", "data")


def safe_float(val, default=0.0):
    if val is None or val.strip() == "":
        return default
    try:
        return float(val)
    except ValueError:
        return default


def safe_int(val, default=0):
    if val is None or val.strip() == "":
        return default
    try:
        return int(float(val))
    except ValueError:
        return default


def convert_students(reader):
    students = []
    for i, row in enumerate(reader):
        at_risk = row.get("AtRiskFlag", "None")
        risk_flags = []
        if at_risk and at_risk != "None":
            risk_flags = [f.strip() for f in at_risk.split(",") if f.strip()]

        academic = safe_float(row.get("AcademicComposite"))
        milskills = safe_float(row.get("MilSkillsComposite"))
        leadership = safe_float(row.get("LeadershipComposite"))
        overall = safe_float(row.get("OverallComposite"))

        # Determine trend from exam progression
        exams = [safe_float(row.get(f"AcademicExam{n}")) for n in range(1, 5)]
        valid_exams = [e for e in exams if e > 0]
        trend = "flat"
        if len(valid_exams) >= 2:
            if valid_exams[-1] > valid_exams[-2] + 2:
                trend = "up"
            elif valid_exams[-1] < valid_exams[-2] - 2:
                trend = "down"

        is_at_risk = len(risk_flags) > 0 or academic < 75 or milskills < 75 or leadership < 75

        students.append({
            "id": f"STU-{i+1:03d}",
            "edipi": row.get("StudentEDIPI", ""),
            "rank": row.get("Rank", ""),
            "lastName": row.get("LastName", ""),
            "firstName": row.get("FirstName", ""),
            "company": row.get("Company", ""),
            "platoon": row.get("Platoon", ""),
            "spc": row.get("SPCAssigned", ""),
            "classNumber": row.get("ClassNumber", ""),
            "classStartDate": row.get("ClassStartDate", ""),
            "phase": row.get("CurrentPhase", ""),
            "exam1": safe_float(row.get("AcademicExam1")),
            "exam2": safe_float(row.get("AcademicExam2")),
            "exam3": safe_float(row.get("AcademicExam3")),
            "exam4": safe_float(row.get("AcademicExam4")),
            "quizAvg": safe_float(row.get("AcademicQuizAvg")),
            "academicComposite": academic,
            "pftScore": safe_int(row.get("PFTScore")),
            "cftScore": safe_int(row.get("CFTScore")),
            "rifleQual": row.get("RifleQual", ""),
            "pistolQual": row.get("PistolQual", ""),
            "landNavDay": row.get("LandNavDay", ""),
            "landNavNight": row.get("LandNavNight", ""),
            "landNavWritten": safe_float(row.get("LandNavWritten")),
            "obstacleCourse": row.get("ObstacleCourse", ""),
            "enduranceCourse": row.get("EnduranceCourse", ""),
            "milSkillsComposite": milskills,
            "leadershipWeek12": safe_float(row.get("LeadershipWeek12")),
            "leadershipWeek22": safe_float(row.get("LeadershipWeek22")),
            "peerEvalWeek12": safe_float(row.get("PeerEvalWeek12")),
            "peerEvalWeek22": safe_float(row.get("PeerEvalWeek22")),
            "leadershipComposite": leadership,
            "overallComposite": overall,
            "classStandingThird": row.get("ClassStandingThird", ""),
            "companyRank": safe_int(row.get("CompanyRank")),
            "trend": trend,
            "atRisk": is_at_risk,
            "riskFlags": risk_flags,
            "status": row.get("Status", ""),
            "notes": row.get("Notes", ""),
        })
    return students


def convert_instructors(reader):
    instructors = []
    for i, row in enumerate(reader):
        instructors.append({
            "id": f"INST-{i+1:03d}",
            "edipi": row.get("InstructorEDIPI", ""),
            "lastName": row.get("LastName", ""),
            "firstName": row.get("FirstName", ""),
            "rank": row.get("Rank", ""),
            "role": row.get("Role", ""),
            "company": row.get("CompanyAssigned", ""),
            "platoon": row.get("PlatoonAssigned", ""),
            "classNumber": row.get("ClassNumber", ""),
            "dateAssigned": row.get("DateAssigned", ""),
            "prd": row.get("PRD", ""),
            "studentsAssigned": safe_int(row.get("StudentsAssigned")),
            "eventsThisWeek": safe_int(row.get("EventsThisWeek")),
            "eventsThisMonth": safe_int(row.get("EventsThisMonth")),
            "counselingsOverdue": safe_int(row.get("CounselingsOverdue")),
            "status": row.get("Status", ""),
            "phone": row.get("Phone", ""),
            "email": row.get("Email", ""),
            "notes": row.get("Notes", ""),
        })
    return instructors


def convert_schedule(reader):
    events = []
    for i, row in enumerate(reader):
        events.append({
            "id": f"EVT-{i+1:03d}",
            "title": row.get("EventTitle", ""),
            "code": row.get("EventCode", ""),
            "phase": row.get("TrainingPhase", ""),
            "category": row.get("Category", ""),
            "gradePillar": row.get("GradePillar", ""),
            "isGraded": row.get("IsGraded", "False") == "True",
            "startDate": row.get("StartDate", ""),
            "endDate": row.get("EndDate", ""),
            "startTime": row.get("StartTime", ""),
            "endTime": row.get("EndTime", ""),
            "durationHours": safe_float(row.get("DurationHours")),
            "location": row.get("Location", ""),
            "company": row.get("CompanyAssigned", ""),
            "classNumber": row.get("ClassNumber", ""),
            "leadInstructor": row.get("LeadInstructor", ""),
            "supportInstructors": row.get("SupportInstructors", ""),
            "instructorCountRequired": safe_int(row.get("InstructorCountRequired")),
            "prerequisiteEvents": row.get("PrerequisiteEvents", ""),
            "specialEquipment": row.get("SpecialEquipment", ""),
            "status": row.get("Status", ""),
            "weatherContingency": row.get("WeatherContingency", ""),
            "notes": row.get("Notes", ""),
        })
    return events


def convert_qualifications(reader):
    quals = []
    for i, row in enumerate(reader):
        quals.append({
            "id": f"QUAL-{i+1:03d}",
            "code": row.get("QualCode", ""),
            "name": row.get("QualName", ""),
            "category": row.get("Category", ""),
            "issuingAuthority": row.get("IssuingAuthority", ""),
            "validityMonths": safe_int(row.get("ValidityMonths")),
            "renewalProcess": row.get("RenewalProcess", ""),
            "requiredForEvents": row.get("RequiredForEvents", ""),
            "minimumPerEvent": safe_int(row.get("MinimumPerEvent")),
            "orderReference": row.get("OrderReference", ""),
            "status": row.get("Status", ""),
            "notes": row.get("Notes", ""),
        })
    return quals


def convert_qual_records(reader):
    records = []
    for i, row in enumerate(reader):
        records.append({
            "id": f"QR-{i+1:03d}",
            "instructorEdipi": row.get("InstructorEDIPI", ""),
            "instructorName": row.get("InstructorName", ""),
            "qualCode": row.get("QualCode", ""),
            "qualName": row.get("QualName", ""),
            "dateEarned": row.get("DateEarned", ""),
            "expirationDate": row.get("ExpirationDate", ""),
            "daysUntilExpiration": safe_int(row.get("DaysUntilExpiration")),
            "expirationStatus": row.get("ExpirationStatus", ""),
            "certificateNumber": row.get("CertificateNumber", ""),
            "issuedBy": row.get("IssuedBy", ""),
            "renewalStatus": row.get("RenewalStatus", ""),
            "renewalDate": row.get("RenewalDate", ""),
            "notes": row.get("Notes", ""),
        })
    return records


def convert_feedback(reader):
    feedback = []
    for i, row in enumerate(reader):
        # Parse rating to numeric
        rating_str = row.get("OverallRating", "")
        rating = 3.0
        if rating_str:
            try:
                rating = float(rating_str[0])
            except (ValueError, IndexError):
                pass

        feedback.append({
            "id": f"FB-{i+1:03d}",
            "eventTitle": row.get("EventTitle", ""),
            "eventCode": row.get("EventCode", ""),
            "eventDate": row.get("EventDate", ""),
            "phase": row.get("TrainingPhase", ""),
            "company": row.get("CompanyAssigned", ""),
            "submitterRole": row.get("SubmitterRole", ""),
            "submitterName": row.get("SubmitterName", ""),
            "overallRating": rating,
            "objectivesMet": row.get("ObjectivesMet", ""),
            "instructorEffectiveness": row.get("InstructorEffectiveness", ""),
            "timeManagement": row.get("TimeManagement", ""),
            "resourceAdequacy": row.get("ResourceAdequacy", ""),
            "sustains": row.get("Sustains", ""),
            "improves": row.get("Improves", ""),
            "safetyConcerns": row.get("SafetyConcerns", ""),
            "hasSafetyConcern": row.get("HasSafetyConcern", "False") == "True",
            "additionalComments": row.get("AdditionalComments", ""),
            "submittedDate": row.get("SubmittedDate", ""),
            "reviewedBy": row.get("ReviewedBy", ""),
            "reviewStatus": row.get("ReviewStatus", ""),
            "actionTaken": row.get("ActionTaken", ""),
        })
    return feedback


CONVERSIONS = {
    "student-scores.csv": ("students.json", convert_students),
    "instructors.csv": ("instructors.json", convert_instructors),
    "training-schedule.csv": ("schedule.json", convert_schedule),
    "required-qualifications.csv": ("qualifications.json", convert_qualifications),
    "qualification-records.csv": ("qual-records.json", convert_qual_records),
    "event-feedback.csv": ("feedback.json", convert_feedback),
}


def main():
    os.makedirs(OUT_DIR, exist_ok=True)

    for csv_file, (json_file, converter) in CONVERSIONS.items():
        csv_path = os.path.join(SAMPLE_DIR, csv_file)
        json_path = os.path.join(OUT_DIR, json_file)

        if not os.path.exists(csv_path):
            print(f"SKIP: {csv_file} not found")
            continue

        with open(csv_path, "r", encoding="utf-8-sig") as f:
            reader = csv.DictReader(f)
            data = converter(reader)

        with open(json_path, "w", encoding="utf-8") as f:
            json.dump(data, f, indent=2)

        print(f"OK: {csv_file} -> {json_file} ({len(data)} records)")


if __name__ == "__main__":
    main()
