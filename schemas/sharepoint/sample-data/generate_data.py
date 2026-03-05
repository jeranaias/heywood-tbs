#!/usr/bin/env python3
"""
Generate realistic synthetic test data for the Heywood Initiative TBS project.
All names and EDIPIs are fictional. No real people.
"""
import csv
import random
import os
from datetime import date, timedelta, datetime

random.seed(42)  # Reproducible

OUTPUT_DIR = os.path.dirname(os.path.abspath(__file__))

# ─────────────────────────────────────────────
# Name pools (common American surnames/first names - all fictional combinations)
# ─────────────────────────────────────────────
LAST_NAMES = [
    "Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
    "Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson",
    "Thomas", "Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson",
    "White", "Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson",
    "Walker", "Young", "Allen", "King", "Wright", "Scott", "Torres", "Nguyen",
    "Hill", "Flores", "Green", "Adams", "Nelson", "Baker", "Hall", "Rivera",
    "Campbell", "Mitchell", "Carter", "Roberts", "Gomez", "Phillips", "Evans",
    "Turner", "Diaz", "Parker", "Cruz", "Edwards", "Collins", "Reyes", "Stewart",
    "Morris", "Morales", "Murphy", "Cook", "Rogers", "Gutierrez", "Ortiz",
    "Morgan", "Cooper", "Peterson", "Bailey", "Reed", "Kelly", "Howard", "Ramos",
    "Kim", "Cox", "Ward", "Richardson", "Watson", "Brooks", "Chavez", "Wood",
    "James", "Bennett", "Gray", "Mendoza", "Ruiz", "Hughes", "Price", "Alvarez",
    "Castillo", "Sanders", "Patel", "Myers", "Long", "Ross", "Foster", "Jimenez",
    "Powell", "Jenkins", "Perry", "Russell", "Sullivan", "Bell", "Coleman",
    "Butler", "Henderson", "Barnes", "Gonzales", "Fisher", "Vasquez", "Simmons",
    "Griffin", "Aguilar", "Stone", "Meyer", "Boyd", "Mills", "Warren", "Fox",
    "Rose", "Rice", "Moreno", "Schmidt", "Hicks", "Dunn", "Hunt", "Steele",
    "Burns", "Carr", "Olson", "Hale", "Marsh", "Burke", "Watkins", "West",
    "Barker", "Lynch", "Stephens", "Vargas", "Crawford", "Garza", "Wolfe",
    "Vega", "Hardy", "Keller", "Hoffman", "Medina", "Romero", "Delgado",
    "Tucker", "Brewer", "Bishop", "Sutton", "Singh", "Frank", "Dawson",
    "Bates", "Lamb", "Moran", "Holt", "Day", "Beard", "Owen", "Duffy",
    "Rojas", "Pace", "Stark", "Galvan", "Fields", "Hubbard", "Erickson",
    "Sharp", "Brock", "Conner", "Mcbride", "Coffey", "Quinn", "Pruitt",
    "Cabrera", "Acosta", "Ochoa", "Villarreal", "Strickland", "Bonilla",
    "Donovan", "Mcconnell", "Tran", "Pham", "Chung", "Orozco", "Ibarra",
    "Blackwell", "Whitehead", "Roth", "Benson", "Harmon", "Gillespie",
    "Mosley", "Chambers", "Mccoy", "Brennan", "Figueroa", "Mclaughlin",
]

FIRST_NAMES_MALE = [
    "James", "Robert", "John", "Michael", "David", "William", "Richard", "Joseph",
    "Thomas", "Christopher", "Charles", "Daniel", "Matthew", "Anthony", "Mark",
    "Donald", "Steven", "Paul", "Andrew", "Joshua", "Kenneth", "Kevin", "Brian",
    "George", "Timothy", "Ronald", "Edward", "Jason", "Jeffrey", "Ryan",
    "Jacob", "Gary", "Nicholas", "Eric", "Jonathan", "Stephen", "Larry",
    "Justin", "Scott", "Brandon", "Benjamin", "Samuel", "Raymond", "Gregory",
    "Frank", "Alexander", "Patrick", "Jack", "Dennis", "Jerry", "Tyler",
    "Aaron", "Jose", "Nathan", "Henry", "Peter", "Douglas", "Zachary",
    "Kyle", "Noah", "Ethan", "Logan", "Mason", "Lucas", "Dylan", "Caleb",
    "Ian", "Connor", "Elijah", "Blake", "Luke", "Gavin", "Jared", "Cole",
    "Marcus", "Trent", "Derek", "Max", "Chase", "Garrett", "Wesley", "Dalton",
    "Bryce", "Cody", "Grant", "Miguel", "Dillon", "Travis", "Clay", "Reid",
    "Carter", "Tucker", "Colton", "Wyatt", "Spencer", "Tristan", "Owen",
    "Parker", "Dominic", "Brayden", "Cameron", "Hunter", "Austin",
]

FIRST_NAMES_FEMALE = [
    "Mary", "Patricia", "Jennifer", "Linda", "Barbara", "Elizabeth", "Susan",
    "Jessica", "Sarah", "Karen", "Lisa", "Nancy", "Betty", "Margaret", "Sandra",
    "Ashley", "Emily", "Donna", "Michelle", "Kimberly", "Carol", "Amanda",
    "Melissa", "Deborah", "Stephanie", "Rebecca", "Sharon", "Laura", "Cynthia",
    "Kathleen", "Amy", "Angela", "Shirley", "Anna", "Brenda", "Pamela",
    "Nicole", "Samantha", "Katherine", "Christine", "Catherine", "Alicia",
    "Rachel", "Heather", "Diane", "Megan", "Brooke", "Andrea", "Hannah",
    "Olivia", "Emma", "Sophia", "Ava", "Lauren", "Grace", "Madison",
    "Chloe", "Victoria", "Natalie", "Lily", "Morgan", "Alexis", "Kayla",
    "Jasmine", "Taylor", "Abigail", "Hailey", "Julia", "Danielle", "Brittany",
    "Kelsey", "Paige", "Courtney", "Lindsey", "Molly", "Shelby", "Brianna",
    "Jordan", "Bailey", "Savannah", "Riley", "Allison", "Mackenzie", "Claire",
]

def random_edipi():
    """Generate a random 10-digit EDIPI (not starting with 0)."""
    return str(random.randint(1000000000, 9999999999))

def used_name_set():
    return set()

_used_names = set()

def unique_name(gender=None):
    """Generate a unique first/last name combo."""
    while True:
        if gender == "F":
            first = random.choice(FIRST_NAMES_FEMALE)
        elif gender == "M":
            first = random.choice(FIRST_NAMES_MALE)
        else:
            first = random.choice(FIRST_NAMES_MALE + FIRST_NAMES_FEMALE)
        last = random.choice(LAST_NAMES)
        key = (last, first)
        if key not in _used_names:
            _used_names.add(key)
            return last, first

# ─────────────────────────────────────────────
# 1. INSTRUCTORS (generate first — students reference SPCs)
# ─────────────────────────────────────────────
CLASS_START = date(2026, 1, 5)  # Monday, Jan 5, 2026
TODAY = date(2026, 3, 4)
CLASS_NUMBER = "1-26"

instructors = []

def make_instructor(role, rank, platoon, date_assigned, prd, students=0):
    last, first = unique_name("M" if random.random() < 0.80 else "F")
    edipi = random_edipi()
    # Events counts
    events_week = random.randint(2, 6)
    events_month = random.randint(8, 20)
    counselings_overdue = random.randint(0, 3) if students > 0 else 0
    phone = f"703-784-{random.randint(1000,9999)}"
    email = f"{first[0].lower()}.{last.lower()}@usmc.mil"
    status = "Active"
    notes = ""
    inst = {
        "InstructorEDIPI": edipi,
        "LastName": last,
        "FirstName": first,
        "Rank": rank,
        "Role": role,
        "CompanyAssigned": "Alpha",
        "PlatoonAssigned": platoon,
        "ClassNumber": CLASS_NUMBER,
        "DateAssigned": date_assigned.isoformat(),
        "PRD": prd.isoformat() if prd else "",
        "StudentsAssigned": students,
        "EventsThisWeek": events_week,
        "EventsThisMonth": events_month,
        "CounselingsOverdue": counselings_overdue,
        "Status": status,
        "Phone": phone,
        "Email": email,
        "Notes": notes,
    }
    instructors.append(inst)
    return inst

# Company Commander
co_cmd = make_instructor("Company Commander", "Capt", "N/A",
                         date(2025, 6, 15), date(2027, 6, 15), 0)
# Company XO
co_xo = make_instructor("Company XO", "1stLt", "N/A",
                         date(2025, 9, 1), date(2027, 3, 1), 0)
# 4 SPCs (Capt) — one per platoon
spc_names = []
for i, plat in enumerate(["1st", "2nd", "3rd", "4th"]):
    da = date(2025, 3 + i*2, 1)
    prd_date = date(2027, 3 + i*2, 1)
    spc = make_instructor("Staff Platoon Commander", "Capt", plat,
                           da, prd_date, 50)
    spc_names.append(f"{spc['Rank']} {spc['LastName']}")

# 2 Assistant SPCs (1stLt)
for plat in ["1st", "3rd"]:
    da = date(2025, 8, random.randint(1, 28))
    prd = date(2027, 8, 1)
    make_instructor("Assistant SPC", "1stLt", plat, da, prd, 0)

# 2 Tactics Instructors (GySgt, SSgt)
for rank in ["GySgt", "SSgt"]:
    da = date(2024, random.randint(6, 12), random.randint(1, 28))
    prd = date(2026, random.randint(4, 8), 1)  # within rotation window
    make_instructor("Tactics Instructor", rank, "N/A", da, prd, 0)

# 1 Weapons Instructor (SSgt)
make_instructor("Weapons Instructor", "SSgt", "N/A",
                date(2025, 2, 10), date(2027, 2, 10), 0)

# 1 PT Instructor (SSgt) — PRD within 90 days (rotation risk)
make_instructor("PT Instructor", "SSgt", "N/A",
                date(2024, 5, 15), date(2026, 5, 1), 0)

# Mark 2 instructors with PRD within 90 days of TODAY (2026-03-04)
# The tactics GySgt and PT instructor already have close PRDs.
# Let's force one to be very close:
instructors[-1]["PRD"] = date(2026, 4, 15).isoformat()  # PT instructor — 42 days
instructors[-3]["PRD"] = date(2026, 5, 20).isoformat()  # Tactics GySgt — 77 days

def write_instructors():
    cols = ["InstructorEDIPI", "LastName", "FirstName", "Rank", "Role",
            "CompanyAssigned", "PlatoonAssigned", "ClassNumber", "DateAssigned",
            "PRD", "StudentsAssigned", "EventsThisWeek", "EventsThisMonth",
            "CounselingsOverdue", "Status", "Phone", "Email", "Notes"]
    path = os.path.join(OUTPUT_DIR, "instructors.csv")
    with open(path, "w", newline="", encoding="utf-8") as f:
        w = csv.DictWriter(f, fieldnames=cols)
        w.writeheader()
        w.writerows(instructors)
    print(f"  instructors.csv: {len(instructors)} rows")

# ─────────────────────────────────────────────
# 2. STUDENT SCORES (~200 rows)
# ─────────────────────────────────────────────

def clamp(v, lo, hi):
    return max(lo, min(hi, v))

def rifle_pistol_qual():
    r = random.random()
    if r < 0.60:
        return "Expert"
    elif r < 0.85:
        return "Sharpshooter"
    elif r < 0.95:
        return "Marksman"
    else:
        return "Unqualified"

def land_nav_result():
    r = random.random()
    if r < 0.85:
        return "Pass"
    elif r < 0.95:
        return "Fail"
    else:
        return "Not Yet Tested"

students = []

# Generate 200 students
n_students = 200
platoon_assign = ["1st"] * 50 + ["2nd"] * 50 + ["3rd"] * 50 + ["4th"] * 50
random.shuffle(platoon_assign)

# Decide special statuses
status_pool = (["Active"] * 190 +
               ["Medical Hold (Mike Co)"] * 5 +
               ["Academic Hold"] * 3 +
               ["Dropped"] * 2)
random.shuffle(status_pool)

# Ranks
rank_pool = ["2ndLt"] * 193 + ["1stLt"] * 5 + ["WO"] * 2
random.shuffle(rank_pool)

for i in range(n_students):
    last, first = unique_name()
    edipi = random_edipi()
    rank = rank_pool[i]
    plat = platoon_assign[i]
    status = status_pool[i]

    # Map platoon to SPC
    plat_idx = ["1st", "2nd", "3rd", "4th"].index(plat)
    spc = spc_names[plat_idx]

    # --- Academic scores ---
    # Base academic ability for this student
    acad_ability = random.gauss(82, 8)

    # Phase I exam (the only one completed so far)
    exam1 = clamp(round(acad_ability + random.gauss(0, 5), 1), 45, 100)

    # Quiz average
    quiz_avg = clamp(round(acad_ability + random.gauss(0, 4), 1), 50, 100)

    # Academic composite (only exam1 and quiz avg available; simple average for Phase I)
    acad_composite = round((exam1 + quiz_avg) / 2, 1)

    # --- Military Skills ---
    pft = clamp(round(random.gauss(255, 25)), 200, 295)
    cft = clamp(round(random.gauss(265, 20)), 220, 300)
    rifle = rifle_pistol_qual()
    pistol = rifle_pistol_qual()
    land_nav_day = land_nav_result()
    land_nav_night = land_nav_result()
    land_nav_written = clamp(round(random.gauss(80, 10), 1), 40, 100)

    # Obstacle and endurance — most not yet tested in Phase I
    if random.random() < 0.3:
        obstacle = random.choice(["Pass", "Fail"])
    else:
        obstacle = "Not Yet Tested"
    if random.random() < 0.2:
        endurance = random.choice(["Pass", "Fail"])
    else:
        endurance = "Not Yet Tested"

    # MilSkills composite — weighted estimate
    # PFT/CFT normalized to 100-scale
    pft_norm = (pft / 300) * 100
    cft_norm = (cft / 300) * 100
    rifle_score = {"Expert": 95, "Sharpshooter": 85, "Marksman": 75, "Unqualified": 50}[rifle]
    pistol_score = {"Expert": 95, "Sharpshooter": 85, "Marksman": 75, "Unqualified": 50}[pistol]
    lnav_day_score = {"Pass": 90, "Fail": 55, "Not Yet Tested": 0}[land_nav_day]
    lnav_night_score = {"Pass": 90, "Fail": 55, "Not Yet Tested": 0}[land_nav_night]

    # Count tested events for averaging
    tested_scores = [pft_norm, cft_norm, rifle_score, pistol_score, land_nav_written]
    if land_nav_day != "Not Yet Tested":
        tested_scores.append(lnav_day_score)
    if land_nav_night != "Not Yet Tested":
        tested_scores.append(lnav_night_score)
    if obstacle != "Not Yet Tested":
        tested_scores.append(90 if obstacle == "Pass" else 55)
    if endurance != "Not Yet Tested":
        tested_scores.append(90 if endurance == "Pass" else 55)

    milskills_composite = round(sum(tested_scores) / len(tested_scores), 1) if tested_scores else 0

    # --- Leadership ---
    # Only Week 12 populated (midpoint eval); Week 22 blank
    lead_w12 = clamp(round(random.gauss(80, 9), 1), 55, 100)
    peer_w12 = clamp(round(random.gauss(78, 10), 1), 50, 100)

    # Leadership composite (only midpoint so far = 14% weight portion)
    # SPC 90%, Peer 10%
    lead_composite = round(lead_w12 * 0.90 + peer_w12 * 0.10, 1)

    # --- Overall Composite ---
    # Formula: Acad*0.32 + MilSkills*0.32 + Leadership*0.36
    overall = round(acad_composite * 0.32 + milskills_composite * 0.32 + lead_composite * 0.36, 1)

    # --- At Risk Flag ---
    at_risk_flags = []
    if acad_composite < 75:
        at_risk_flags.append("Academic (<75%)")
    if milskills_composite < 75:
        at_risk_flags.append("MilSkills (<75%)")
    if lead_composite < 75:
        at_risk_flags.append("Leadership (<75%)")

    if len(at_risk_flags) >= 2:
        at_risk = "Multiple (<75%)"
    elif len(at_risk_flags) == 1:
        at_risk = at_risk_flags[0]
    else:
        at_risk = "None"

    # If dropped or on hold, override phase
    if status == "Dropped":
        current_phase = "Phase I - Individual Skills"
        at_risk = "Multiple (<75%)" if at_risk == "None" else at_risk
    elif status in ("Medical Hold (Mike Co)", "Academic Hold"):
        current_phase = "Phase I - Individual Skills"
    else:
        current_phase = "Phase I - Individual Skills"

    # Notes for special cases
    notes = ""
    if status == "Medical Hold (Mike Co)":
        notes = random.choice([
            "Stress fracture right tibia - pending medical board",
            "ACL injury during obstacle course - surgery scheduled",
            "Heat injury recovery - cleared for limited duty",
            "Shoulder dislocation during grappling - PT in progress",
            "Plantar fasciitis - non-weight bearing profile",
        ])
    elif status == "Academic Hold":
        notes = random.choice([
            "Failed Phase I exam - remediation in progress",
            "Below 75% academic composite - counseled 15 Feb",
            "Missed 3 consecutive quizzes - admin review pending",
        ])
    elif status == "Dropped":
        notes = random.choice([
            "DOR - personal reasons",
            "Failed to meet minimum standards after remediation",
        ])

    student = {
        "StudentEDIPI": edipi,
        "LastName": last,
        "FirstName": first,
        "Rank": rank,
        "Company": "Alpha",
        "Platoon": plat,
        "SPCAssigned": spc,
        "ClassNumber": CLASS_NUMBER,
        "ClassStartDate": CLASS_START.isoformat(),
        "CurrentPhase": current_phase,
        "AcademicExam1": exam1,
        "AcademicExam2": "",
        "AcademicExam3": "",
        "AcademicExam4": "",
        "AcademicQuizAvg": quiz_avg,
        "AcademicComposite": acad_composite,
        "PFTScore": pft,
        "CFTScore": cft,
        "RifleQual": rifle,
        "PistolQual": pistol,
        "LandNavDay": land_nav_day,
        "LandNavNight": land_nav_night,
        "LandNavWritten": land_nav_written,
        "ObstacleCourse": obstacle,
        "EnduranceCourse": endurance,
        "MilSkillsComposite": milskills_composite,
        "LeadershipWeek12": lead_w12,
        "LeadershipWeek22": "",
        "PeerEvalWeek12": peer_w12,
        "PeerEvalWeek22": "",
        "LeadershipComposite": lead_composite,
        "OverallComposite": overall,
        "ClassStandingThird": "",  # filled after sorting
        "CompanyRank": 0,          # filled after sorting
        "AtRiskFlag": at_risk,
        "Status": status,
        "Notes": notes,
    }
    students.append(student)

# Compute class standing and company rank among Active students
active_students = [s for s in students if s["Status"] == "Active"]
active_students.sort(key=lambda s: s["OverallComposite"], reverse=True)
n_active = len(active_students)
for rank_pos, s in enumerate(active_students, 1):
    s["CompanyRank"] = rank_pos
    if rank_pos <= n_active // 3:
        s["ClassStandingThird"] = "Top Third"
    elif rank_pos <= 2 * n_active // 3:
        s["ClassStandingThird"] = "Middle Third"
    else:
        s["ClassStandingThird"] = "Bottom Third"

# Non-active students: no rank
for s in students:
    if s["Status"] != "Active":
        s["CompanyRank"] = ""
        s["ClassStandingThird"] = ""

def write_students():
    cols = [
        "StudentEDIPI", "LastName", "FirstName", "Rank", "Company", "Platoon",
        "SPCAssigned", "ClassNumber", "ClassStartDate", "CurrentPhase",
        "AcademicExam1", "AcademicExam2", "AcademicExam3", "AcademicExam4",
        "AcademicQuizAvg", "AcademicComposite",
        "PFTScore", "CFTScore", "RifleQual", "PistolQual",
        "LandNavDay", "LandNavNight", "LandNavWritten",
        "ObstacleCourse", "EnduranceCourse", "MilSkillsComposite",
        "LeadershipWeek12", "LeadershipWeek22", "PeerEvalWeek12", "PeerEvalWeek22",
        "LeadershipComposite", "OverallComposite",
        "ClassStandingThird", "CompanyRank", "AtRiskFlag", "Status", "Notes",
    ]
    path = os.path.join(OUTPUT_DIR, "student-scores.csv")
    with open(path, "w", newline="", encoding="utf-8") as f:
        w = csv.DictWriter(f, fieldnames=cols)
        w.writeheader()
        w.writerows(students)
    print(f"  student-scores.csv: {len(students)} rows")

# ─────────────────────────────────────────────
# 3. TRAINING SCHEDULE (~60 rows)
# ─────────────────────────────────────────────

schedule_events = []

# Lead instructors we can reference
lead_instructors_pool = [f"{i['Rank']} {i['LastName']}" for i in instructors]

def add_event(title, code, phase, category, pillar, graded, start, end,
              start_time, end_time, hours, location, status,
              weather="No Impact", lead=None, support_count=1, prereqs="",
              equipment="", notes=""):
    if lead is None:
        lead = random.choice(lead_instructors_pool)
    support = random.sample(lead_instructors_pool, min(support_count - 1, len(lead_instructors_pool)))
    support_str = "; ".join(support) if support else ""
    schedule_events.append({
        "EventTitle": title,
        "EventCode": code,
        "TrainingPhase": phase,
        "Category": category,
        "GradePillar": pillar,
        "IsGraded": graded,
        "StartDate": start.isoformat(),
        "EndDate": end.isoformat(),
        "StartTime": start_time,
        "EndTime": end_time,
        "DurationHours": hours,
        "Location": location,
        "CompanyAssigned": "Alpha",
        "ClassNumber": CLASS_NUMBER,
        "LeadInstructor": lead,
        "SupportInstructors": support_str,
        "InstructorCountRequired": support_count,
        "PrerequisiteEvents": prereqs,
        "SpecialEquipment": equipment,
        "Status": status,
        "WeatherContingency": weather,
        "Notes": notes,
    })

P1 = "Phase I - Individual Skills"
P2 = "Phase II - Squad"

# Week 1 (Jan 5-9)
add_event("TBS Check-In and Orientation", "PI-001", P1, "Admin", "Not Graded", False,
          date(2026,1,5), date(2026,1,5), "0700", "1700", 10,
          "TBS Main Building, Camp Barrett", "Complete")
add_event("History of TBS / Warrior Ethos", "PI-002", P1, "Academic", "Academics (32%)", False,
          date(2026,1,6), date(2026,1,6), "0800", "1100", 3,
          "Gruber Auditorium, Camp Barrett", "Complete")
add_event("Initial PFT", "PI-003", P1, "Physical Training", "Military Skills (32%)", True,
          date(2026,1,7), date(2026,1,7), "0600", "1000", 4,
          "Camp Barrett Track and Pull-up Bars", "Complete",
          weather="Cold Weather Threshold")
add_event("Initial CFT", "PI-004", P1, "Physical Training", "Military Skills (32%)", True,
          date(2026,1,8), date(2026,1,8), "0600", "1000", 4,
          "Camp Barrett MANUF Course", "Complete",
          weather="Heat Cat Dependent")
add_event("Small Unit Leadership Lecture", "PI-005", P1, "Academic", "Academics (32%)", False,
          date(2026,1,9), date(2026,1,9), "0800", "1200", 4,
          "Gruber Auditorium, Camp Barrett", "Complete")

# Week 2 (Jan 12-16)
add_event("Marine Corps Planning Process Overview", "PI-006", P1, "Academic", "Academics (32%)", False,
          date(2026,1,12), date(2026,1,12), "0800", "1200", 4,
          "TBS Classroom B-204", "Complete")
add_event("Land Navigation Theory", "PI-007", P1, "Academic", "Academics (32%)", False,
          date(2026,1,13), date(2026,1,13), "0800", "1600", 8,
          "TBS Classroom B-204", "Complete")
add_event("Map Reading Practical Application", "PI-008", P1, "Military Skills", "Military Skills (32%)", False,
          date(2026,1,14), date(2026,1,14), "0800", "1600", 8,
          "OCS/TBS Training Area", "Complete",
          weather="Rain Plan Available")
add_event("Compass and Protractor Lab", "PI-009", P1, "Military Skills", "Military Skills (32%)", False,
          date(2026,1,15), date(2026,1,15), "0800", "1200", 4,
          "TBS Classroom B-204", "Complete")
add_event("Academic Quiz 1 - Map Reading", "PI-010", P1, "Evaluation", "Academics (32%)", True,
          date(2026,1,16), date(2026,1,16), "0800", "0930", 1.5,
          "TBS Classroom B-204", "Complete")

# Week 3 (Jan 20-23) — MLK Day Jan 19
add_event("Land Navigation Day Practical", "PI-011", P1, "Military Skills", "Military Skills (32%)", True,
          date(2026,1,20), date(2026,1,20), "0600", "1800", 12,
          "Chopawamsic Training Area", "Complete",
          weather="Lightning Hold Applies", support_count=6,
          equipment="Lensatic compasses (200), protractors, map sheets 1:50000")
add_event("Land Navigation Night Practical", "PI-012", P1, "Military Skills", "Military Skills (32%)", True,
          date(2026,1,21), date(2026,1,21), "1800", "0200", 8,
          "Chopawamsic Training Area", "Complete",
          weather="Lightning Hold Applies", support_count=6, prereqs="PI-011")
add_event("Land Navigation Written Exam", "PI-013", P1, "Evaluation", "Military Skills (32%)", True,
          date(2026,1,22), date(2026,1,22), "0800", "1000", 2,
          "TBS Classroom B-204", "Complete", prereqs="PI-011,PI-012")
add_event("Rifle Marksmanship Fundamentals", "PI-014", P1, "Military Skills", "Military Skills (32%)", False,
          date(2026,1,23), date(2026,1,23), "0800", "1600", 8,
          "Indoor Simulated Marksmanship Trainer (ISMT), Camp Barrett", "Complete")

# Week 4 (Jan 26-30)
add_event("Known Distance Rifle Qualification - Day 1", "PI-015", P1, "Military Skills", "Military Skills (32%)", False,
          date(2026,1,26), date(2026,1,26), "0600", "1700", 11,
          "Range 8, MCB Quantico", "Complete",
          weather="Lightning Hold Applies", support_count=8,
          equipment="M27 IAR/M4 (200), 5.56mm ball (40,000 rds), targets, pasters")
add_event("Known Distance Rifle Qualification - Day 2 (Qual Day)", "PI-016", P1, "Military Skills", "Military Skills (32%)", True,
          date(2026,1,27), date(2026,1,27), "0600", "1700", 11,
          "Range 8, MCB Quantico", "Complete",
          weather="Lightning Hold Applies", support_count=8, prereqs="PI-015",
          equipment="M27 IAR/M4 (200), 5.56mm ball (10,000 rds), scorecards")
add_event("Combat Marksmanship Program Day 1", "PI-017", P1, "Military Skills", "Military Skills (32%)", False,
          date(2026,1,28), date(2026,1,28), "0700", "1700", 10,
          "Range 4, MCB Quantico", "Complete",
          support_count=6, equipment="M27/M4, 5.56mm ball (20,000 rds)")
add_event("Combat Marksmanship Program Day 2", "PI-018", P1, "Military Skills", "Military Skills (32%)", False,
          date(2026,1,29), date(2026,1,29), "0700", "1700", 10,
          "Range 4, MCB Quantico", "Complete",
          support_count=6, equipment="M27/M4, 5.56mm ball (20,000 rds)")
add_event("Warfighting Fundamentals Lecture", "PI-019", P1, "Academic", "Academics (32%)", False,
          date(2026,1,30), date(2026,1,30), "0800", "1200", 4,
          "Gruber Auditorium, Camp Barrett", "Complete")

# Week 5 (Feb 2-6)
add_event("Pistol Marksmanship Fundamentals", "PI-020", P1, "Military Skills", "Military Skills (32%)", False,
          date(2026,2,2), date(2026,2,2), "0800", "1600", 8,
          "ISMT, Camp Barrett", "Complete")
add_event("Pistol Qualification", "PI-021", P1, "Military Skills", "Military Skills (32%)", True,
          date(2026,2,3), date(2026,2,3), "0700", "1600", 9,
          "Range 2, MCB Quantico", "Complete",
          support_count=6, prereqs="PI-020",
          equipment="M18 pistol (200), 9mm ball (8,000 rds), targets")
add_event("Patrol Base Operations", "PI-022", P1, "Military Skills", "Military Skills (32%)", False,
          date(2026,2,4), date(2026,2,4), "0800", "1700", 9,
          "TBS Training Area", "Complete",
          weather="Rain Plan Available")
add_event("Combat Lifesaver Course", "PI-023", P1, "Military Skills", "Military Skills (32%)", False,
          date(2026,2,5), date(2026,2,6), "0800", "1700", 18,
          "TBS Medical Training Facility", "Complete",
          support_count=4, equipment="CLS bags (50), tourniquets, IV trainers")

# Week 6 (Feb 9-13)
add_event("Operations Order Format Lecture", "PI-024", P1, "Academic", "Academics (32%)", False,
          date(2026,2,9), date(2026,2,9), "0800", "1200", 4,
          "TBS Classroom B-204", "Complete")
add_event("Obstacle Course Assessment", "PI-025", P1, "Physical Training", "Military Skills (32%)", True,
          date(2026,2,10), date(2026,2,10), "0700", "1100", 4,
          "TBS Obstacle Course, Camp Barrett", "Complete",
          weather="Cold Weather Threshold")
add_event("Fire Team Tactics Lecture", "PI-026", P1, "Academic", "Academics (32%)", False,
          date(2026,2,11), date(2026,2,11), "0800", "1200", 4,
          "TBS Classroom B-204", "Complete")
add_event("Academic Quiz 2 - Tactics Fundamentals", "PI-027", P1, "Evaluation", "Academics (32%)", True,
          date(2026,2,12), date(2026,2,12), "0800", "0930", 1.5,
          "TBS Classroom B-204", "Complete")
add_event("PT Session - Endurance Course", "PI-028", P1, "Physical Training", "Military Skills (32%)", True,
          date(2026,2,13), date(2026,2,13), "0600", "0900", 3,
          "TBS Endurance Course, Camp Barrett", "Complete",
          weather="Heat Cat Dependent")

# Week 7 (Feb 16-20)
add_event("Offensive Operations Lecture", "PI-029", P1, "Academic", "Academics (32%)", False,
          date(2026,2,16), date(2026,2,16), "0800", "1200", 4,
          "Gruber Auditorium, Camp Barrett", "Complete")
add_event("Defensive Operations Lecture", "PI-030", P1, "Academic", "Academics (32%)", False,
          date(2026,2,17), date(2026,2,17), "0800", "1200", 4,
          "Gruber Auditorium, Camp Barrett", "Complete")
add_event("Call for Fire Lecture", "PI-031", P1, "Academic", "Academics (32%)", False,
          date(2026,2,18), date(2026,2,18), "0800", "1200", 4,
          "TBS Classroom B-204", "Complete")
add_event("Call for Fire Practical Application", "PI-032", P1, "Military Skills", "Academics (32%)", False,
          date(2026,2,19), date(2026,2,19), "0800", "1600", 8,
          "OP4 Observation Post, MCB Quantico", "Complete",
          weather="Rain Plan Available")
add_event("Phase I Written Exam", "PI-033", P1, "Evaluation", "Academics (32%)", True,
          date(2026,2,20), date(2026,2,20), "0800", "1100", 3,
          "Gruber Auditorium, Camp Barrett", "Complete")

# Week 8 (Feb 23-27)
add_event("Midpoint Leadership Evaluation (Week 12 Equivalent)", "PI-034", P1, "Leadership", "Leadership (36%)", True,
          date(2026,2,23), date(2026,2,24), "0800", "1700", 16,
          "Camp Barrett", "Complete",
          notes="SPC evaluations and peer evals due NLT 1700 24 Feb")
add_event("Midpoint Peer Evaluations", "PI-035", P1, "Leadership", "Leadership (36%)", True,
          date(2026,2,25), date(2026,2,25), "0800", "1200", 4,
          "TBS Classroom B-204", "Complete",
          notes="Anonymous peer ranking forms distributed")
add_event("IED Awareness and Counter-IED Ops", "PI-036", P1, "Academic", "Academics (32%)", False,
          date(2026,2,26), date(2026,2,26), "0800", "1200", 4,
          "TBS Classroom B-204", "Complete")
add_event("Radio Communications Practical", "PI-037", P1, "Military Skills", "Military Skills (32%)", False,
          date(2026,2,27), date(2026,2,27), "0800", "1600", 8,
          "TBS Training Area", "Complete",
          equipment="AN/PRC-152 radios (50), SINCGARS, batteries")

# Week 9 (Mar 2-6) — current week, some in progress / scheduled
add_event("Squad Movement Techniques", "PII-001", P2, "Military Skills", "Military Skills (32%)", False,
          date(2026,3,2), date(2026,3,2), "0800", "1700", 9,
          "TBS Training Area", "Complete",
          weather="Rain Plan Available")
add_event("Squad Attack - Dry Run", "PII-002", P2, "Military Skills", "Military Skills (32%)", False,
          date(2026,3,3), date(2026,3,3), "0800", "1700", 9,
          "TBS Training Area", "Complete")
add_event("Squad Attack - Live Fire", "PII-003", P2, "Military Skills", "Military Skills (32%)", True,
          date(2026,3,4), date(2026,3,4), "0600", "1800", 12,
          "Range 10, MCB Quantico", "In Progress",
          weather="Lightning Hold Applies", support_count=10,
          equipment="M27/M4 (200), 5.56 ball/tracer (30,000 rds), smoke grenades (100), sim grenades (200)")
add_event("Academic Quiz 3 - Offensive/Defensive Ops", "PII-004", P2, "Evaluation", "Academics (32%)", True,
          date(2026,3,5), date(2026,3,5), "0800", "0930", 1.5,
          "TBS Classroom B-204", "Scheduled")
add_event("Squad Patrol - Day", "PII-005", P2, "Military Skills", "Military Skills (32%)", True,
          date(2026,3,6), date(2026,3,6), "0600", "1800", 12,
          "Chopawamsic Training Area", "Scheduled",
          weather="Lightning Hold Applies", support_count=6)

# Week 10 (Mar 9-13)
add_event("Squad Defense", "PII-006", P2, "Military Skills", "Military Skills (32%)", True,
          date(2026,3,9), date(2026,3,9), "0800", "1700", 9,
          "TBS Training Area", "Scheduled",
          weather="Rain Plan Available", support_count=6)
add_event("Military Law and Rules of Engagement", "PII-007", P2, "Academic", "Academics (32%)", False,
          date(2026,3,10), date(2026,3,10), "0800", "1200", 4,
          "TBS Classroom B-204", "Scheduled")
add_event("Close Air Support Lecture", "PII-008", P2, "Academic", "Academics (32%)", False,
          date(2026,3,11), date(2026,3,11), "0800", "1200", 4,
          "Gruber Auditorium, Camp Barrett", "Scheduled")
add_event("Combined Arms Planning Exercise", "PII-009", P2, "Academic", "Academics (32%)", False,
          date(2026,3,12), date(2026,3,12), "0800", "1700", 9,
          "TBS Classroom B-204", "Scheduled",
          notes="Sand table exercise — platoon billets assigned")
add_event("PT Assessment - 5-Mile Conditioning Hike", "PII-010", P2, "Physical Training", "Military Skills (32%)", False,
          date(2026,3,13), date(2026,3,13), "0600", "1000", 4,
          "Camp Barrett to Range Complex Road", "Scheduled",
          weather="Heat Cat Dependent")

# Week 11-12 (Mar 16-27)
add_event("Squad OPORD Practical - Graded", "PII-011", P2, "Evaluation", "Academics (32%)", True,
          date(2026,3,16), date(2026,3,17), "0800", "1700", 18,
          "TBS Training Area", "Scheduled",
          support_count=4, notes="Each student presents 5-paragraph order")
add_event("Mechanized Operations Lecture", "PII-012", P2, "Academic", "Academics (32%)", False,
          date(2026,3,18), date(2026,3,18), "0800", "1200", 4,
          "Gruber Auditorium, Camp Barrett", "Scheduled")
add_event("Convoy Operations Practical", "PII-013", P2, "Military Skills", "Military Skills (32%)", False,
          date(2026,3,19), date(2026,3,19), "0700", "1700", 10,
          "Main Service Road, MCB Quantico", "Scheduled",
          weather="Rain Plan Available", support_count=6,
          equipment="HMMWV (12), radio sets (12), BFT displays")
add_event("Night Squad Patrol", "PII-014", P2, "Military Skills", "Military Skills (32%)", True,
          date(2026,3,20), date(2026,3,20), "1800", "0600", 12,
          "Chopawamsic Training Area", "Scheduled",
          weather="Lightning Hold Applies", support_count=6)
add_event("Field Exercise 1 - FINEX Alpha", "PII-015", P2, "Field Exercise", "Military Skills (32%)", True,
          date(2026,3,23), date(2026,3,27), "0001", "2359", 120,
          "Quantico Training Areas (Multiple)", "Scheduled",
          weather="Lightning Hold Applies", support_count=12,
          equipment="Full combat load, C-rats (5 days), blank ammo, pyrotechnics",
          notes="5-day field exercise — final Phase II graded event. Liberty upon return.")

# Additional Phase I events to reach ~60 total
add_event("CBRN Defense Lecture", "PI-038", P1, "Academic", "Academics (32%)", False,
          date(2026,1,28), date(2026,1,28), "1300", "1600", 3,
          "TBS Classroom B-206", "Complete")
add_event("CBRN Gas Chamber Practical", "PI-039", P1, "Military Skills", "Military Skills (32%)", False,
          date(2026,1,29), date(2026,1,29), "0700", "1200", 5,
          "CBRN Training Facility, MCB Quantico", "Complete",
          support_count=4, equipment="M50 JSGPM masks (200), CS tablets, decon kits")
add_event("Water Survival Basic", "PI-040", P1, "Physical Training", "Military Skills (32%)", True,
          date(2026,2,3), date(2026,2,3), "1300", "1700", 4,
          "TBS Pool Facility, Camp Barrett", "Complete",
          support_count=4)
add_event("Military Briefing Techniques", "PI-041", P1, "Academic", "Academics (32%)", False,
          date(2026,2,10), date(2026,2,10), "1300", "1600", 3,
          "TBS Classroom B-204", "Complete")
add_event("Crew-Served Weapons Familiarization", "PI-042", P1, "Military Skills", "Military Skills (32%)", False,
          date(2026,2,16), date(2026,2,16), "1300", "1700", 4,
          "Weapons Demonstration Facility, Camp Barrett", "Complete",
          support_count=3, equipment="M240B (4), M2 .50 cal (2), Mk19 (2) - display only")
add_event("Grappling and MCMAP", "PI-043", P1, "Physical Training", "Military Skills (32%)", False,
          date(2026,2,18), date(2026,2,18), "0600", "0900", 3,
          "TBS Gym, Camp Barrett", "Complete")
add_event("Ethics and Core Values Seminar", "PI-044", P1, "Academic", "Academics (32%)", False,
          date(2026,2,24), date(2026,2,24), "0800", "1200", 4,
          "Gruber Auditorium, Camp Barrett", "Complete")
add_event("Night Compass Course Remediation", "PI-045", P1, "Military Skills", "Military Skills (32%)", False,
          date(2026,2,25), date(2026,2,25), "1800", "2200", 4,
          "Chopawamsic Training Area", "Complete",
          support_count=3, notes="Remediation for students who failed PI-012")

def write_schedule():
    cols = [
        "EventTitle", "EventCode", "TrainingPhase", "Category", "GradePillar",
        "IsGraded", "StartDate", "EndDate", "StartTime", "EndTime", "DurationHours",
        "Location", "CompanyAssigned", "ClassNumber", "LeadInstructor",
        "SupportInstructors", "InstructorCountRequired", "PrerequisiteEvents",
        "SpecialEquipment", "Status", "WeatherContingency", "Notes",
    ]
    path = os.path.join(OUTPUT_DIR, "training-schedule.csv")
    with open(path, "w", newline="", encoding="utf-8") as f:
        w = csv.DictWriter(f, fieldnames=cols)
        w.writeheader()
        w.writerows(schedule_events)
    print(f"  training-schedule.csv: {len(schedule_events)} rows")

# ─────────────────────────────────────────────
# 4. REQUIRED QUALIFICATIONS (~25 rows)
# ─────────────────────────────────────────────

required_quals = [
    {"QualCode": "RSO-RIFLE", "QualName": "Range Safety Officer - Rifle", "Category": "Range Safety",
     "IssuingAuthority": "WTBN Quantico", "ValidityMonths": 12,
     "RenewalProcess": "Complete RSO recertification course (8 hrs) and pass written exam",
     "RequiredForEvents": "PI-015,PI-016,PI-017,PI-018,PII-003",
     "MinimumPerEvent": 2, "OrderReference": "MCO 3570.1C", "Status": "Active", "Notes": ""},
    {"QualCode": "RSO-PISTOL", "QualName": "Range Safety Officer - Pistol", "Category": "Range Safety",
     "IssuingAuthority": "WTBN Quantico", "ValidityMonths": 12,
     "RenewalProcess": "Complete RSO recertification course (8 hrs) and pass written exam",
     "RequiredForEvents": "PI-021",
     "MinimumPerEvent": 2, "OrderReference": "MCO 3570.1C", "Status": "Active", "Notes": ""},
    {"QualCode": "RSO-DEMO", "QualName": "Range Safety Officer - Demolitions", "Category": "Demolitions",
     "IssuingAuthority": "WTBN Quantico", "ValidityMonths": 12,
     "RenewalProcess": "Complete demolitions RSO course and live practical",
     "RequiredForEvents": "",
     "MinimumPerEvent": 1, "OrderReference": "MCO 3570.1C", "Status": "Active", "Notes": ""},
    {"QualCode": "OIC-DEMO", "QualName": "Officer in Charge - Demolitions Range", "Category": "Demolitions",
     "IssuingAuthority": "TECOM", "ValidityMonths": 24,
     "RenewalProcess": "TECOM-approved refresher course with live demo exercise",
     "RequiredForEvents": "",
     "MinimumPerEvent": 1, "OrderReference": "MCO 3570.1C Ch.4", "Status": "Active", "Notes": ""},
    {"QualCode": "OIC-RANGE", "QualName": "Officer in Charge - Live Fire Range", "Category": "Range Safety",
     "IssuingAuthority": "WTBN Quantico", "ValidityMonths": 12,
     "RenewalProcess": "Complete OIC recertification and supervise minimum 2 live-fire events",
     "RequiredForEvents": "PI-015,PI-016,PI-017,PI-018,PI-021,PII-003",
     "MinimumPerEvent": 1, "OrderReference": "MCO 3574.2L", "Status": "Active", "Notes": ""},
    {"QualCode": "INST-LANDNAV", "QualName": "Land Navigation Instructor", "Category": "Land Navigation",
     "IssuingAuthority": "TBS", "ValidityMonths": 24,
     "RenewalProcess": "Requalify land nav course and demonstrate teaching proficiency",
     "RequiredForEvents": "PI-007,PI-008,PI-009,PI-011,PI-012,PI-013",
     "MinimumPerEvent": 4, "OrderReference": "TBS SOP 3-2", "Status": "Active", "Notes": ""},
    {"QualCode": "INST-WATSURVIVAL", "QualName": "Water Survival Instructor", "Category": "Water Survival",
     "IssuingAuthority": "MCIWS", "ValidityMonths": 24,
     "RenewalProcess": "MCIWS recertification including rescue swimmer demonstration",
     "RequiredForEvents": "",
     "MinimumPerEvent": 3, "OrderReference": "MCO 1500.52D", "Status": "Active", "Notes": ""},
    {"QualCode": "CLS", "QualName": "Combat Lifesaver", "Category": "Combat Lifesaver",
     "IssuingAuthority": "FMTB East", "ValidityMonths": 12,
     "RenewalProcess": "Complete CLS refresher (4 hrs) with practical evaluation",
     "RequiredForEvents": "PI-023",
     "MinimumPerEvent": 2, "OrderReference": "MCO 1500.56A", "Status": "Active", "Notes": ""},
    {"QualCode": "PMI-RIFLE", "QualName": "Primary Marksmanship Instructor - Rifle", "Category": "Weapons Instruction",
     "IssuingAuthority": "WTBN Quantico", "ValidityMonths": 24,
     "RenewalProcess": "PMI refresher course and qualification with assigned weapon",
     "RequiredForEvents": "PI-014,PI-015,PI-016,PI-017,PI-018",
     "MinimumPerEvent": 4, "OrderReference": "MCO 3574.2L", "Status": "Active", "Notes": ""},
    {"QualCode": "PMI-PISTOL", "QualName": "Primary Marksmanship Instructor - Pistol", "Category": "Weapons Instruction",
     "IssuingAuthority": "WTBN Quantico", "ValidityMonths": 24,
     "RenewalProcess": "PMI refresher course and pistol qualification",
     "RequiredForEvents": "PI-020,PI-021",
     "MinimumPerEvent": 4, "OrderReference": "MCO 3574.2L", "Status": "Active", "Notes": ""},
    {"QualCode": "INST-CQB", "QualName": "Close Quarters Battle Instructor", "Category": "Tactics Instruction",
     "IssuingAuthority": "SOTG / MCTOG", "ValidityMonths": 36,
     "RenewalProcess": "Attend CQB instructor update (40 hrs) with live-fire qual",
     "RequiredForEvents": "",
     "MinimumPerEvent": 2, "OrderReference": "MCRP 3-01A", "Status": "Active", "Notes": ""},
    {"QualCode": "INST-MOUT", "QualName": "Military Operations in Urban Terrain Instructor", "Category": "Tactics Instruction",
     "IssuingAuthority": "TBS", "ValidityMonths": 24,
     "RenewalProcess": "MOUT instructor recertification and practical evaluation",
     "RequiredForEvents": "",
     "MinimumPerEvent": 4, "OrderReference": "MCWP 3-35.3", "Status": "Active", "Notes": ""},
    {"QualCode": "INST-DEMO", "QualName": "Demolitions Instructor", "Category": "Demolitions",
     "IssuingAuthority": "EOTG", "ValidityMonths": 24,
     "RenewalProcess": "EOTG demolitions instructor course with live practical",
     "RequiredForEvents": "",
     "MinimumPerEvent": 2, "OrderReference": "MCO 8023.3B", "Status": "Active", "Notes": ""},
    {"QualCode": "INST-PATROLLING", "QualName": "Patrolling Instructor", "Category": "Tactics Instruction",
     "IssuingAuthority": "TBS", "ValidityMonths": 24,
     "RenewalProcess": "Conduct 2 patrol evals under senior instructor supervision",
     "RequiredForEvents": "PII-005,PII-014",
     "MinimumPerEvent": 4, "OrderReference": "TBS SOP 4-1", "Status": "Active", "Notes": ""},
    {"QualCode": "CPR-AED", "QualName": "CPR/AED Certified", "Category": "Combat Lifesaver",
     "IssuingAuthority": "Red Cross / NAEMT", "ValidityMonths": 24,
     "RenewalProcess": "Complete CPR/AED recertification course",
     "RequiredForEvents": "PI-003,PI-004,PI-025,PI-028",
     "MinimumPerEvent": 2, "OrderReference": "MCO 6200.1F", "Status": "Active", "Notes": ""},
    {"QualCode": "INST-PT", "QualName": "Physical Training Instructor", "Category": "Physical Training",
     "IssuingAuthority": "HQMC Fitness", "ValidityMonths": 36,
     "RenewalProcess": "Complete HQMC PT instructor recertification and fitness assessment",
     "RequiredForEvents": "PI-003,PI-004,PI-025,PI-028,PII-010",
     "MinimumPerEvent": 1, "OrderReference": "MCO 6100.13A", "Status": "Active", "Notes": ""},
    {"QualCode": "INST-SWIM", "QualName": "Swim Qualification Instructor", "Category": "Water Survival",
     "IssuingAuthority": "MCIWS", "ValidityMonths": 12,
     "RenewalProcess": "Annual recertification with MCIWS including rescue swimmer drill",
     "RequiredForEvents": "",
     "MinimumPerEvent": 2, "OrderReference": "MCO 1500.52D", "Status": "Active", "Notes": ""},
    {"QualCode": "DRV-HMMWV", "QualName": "HMMWV Operator License", "Category": "Driving/Vehicle",
     "IssuingAuthority": "Motor Transport", "ValidityMonths": 24,
     "RenewalProcess": "Annual license renewal with road test and PMCS evaluation",
     "RequiredForEvents": "PII-013",
     "MinimumPerEvent": 12, "OrderReference": "MCO 11240.66D", "Status": "Active", "Notes": ""},
    {"QualCode": "DRV-7TON", "QualName": "7-Ton (MTVR) Operator License", "Category": "Driving/Vehicle",
     "IssuingAuthority": "Motor Transport", "ValidityMonths": 24,
     "RenewalProcess": "Annual license renewal with road test and PMCS evaluation",
     "RequiredForEvents": "",
     "MinimumPerEvent": 2, "OrderReference": "MCO 11240.66D", "Status": "Active", "Notes": ""},
    {"QualCode": "INST-COMMS", "QualName": "Communications Instructor", "Category": "Instructor Certification",
     "IssuingAuthority": "TBS / MCTOG", "ValidityMonths": 24,
     "RenewalProcess": "Demonstrate proficiency on current radio systems and complete refresher",
     "RequiredForEvents": "PI-037",
     "MinimumPerEvent": 2, "OrderReference": "TBS SOP 5-3", "Status": "Active", "Notes": ""},
    {"QualCode": "INST-IED", "QualName": "Counter-IED Instructor", "Category": "Tactics Instruction",
     "IssuingAuthority": "JIEDDO / TECOM", "ValidityMonths": 24,
     "RenewalProcess": "Complete annual C-IED update brief and lane certification",
     "RequiredForEvents": "PI-036",
     "MinimumPerEvent": 1, "OrderReference": "MCO 3502.7A", "Status": "Active", "Notes": ""},
    {"QualCode": "INST-CFF", "QualName": "Call for Fire Instructor", "Category": "Tactics Instruction",
     "IssuingAuthority": "TBS / AITB", "ValidityMonths": 24,
     "RenewalProcess": "Complete CFF instructor refresher with practical at OP4",
     "RequiredForEvents": "PI-031,PI-032",
     "MinimumPerEvent": 2, "OrderReference": "MCWP 3-16.6", "Status": "Active", "Notes": ""},
    {"QualCode": "INST-CBRN", "QualName": "CBRN Defense Instructor", "Category": "Instructor Certification",
     "IssuingAuthority": "CBRN School", "ValidityMonths": 24,
     "RenewalProcess": "Annual recertification with live agent training",
     "RequiredForEvents": "",
     "MinimumPerEvent": 2, "OrderReference": "MCO 3400.3J", "Status": "Active", "Notes": ""},
    {"QualCode": "SAFE-HEAT", "QualName": "Heat Injury Prevention Monitor", "Category": "Other",
     "IssuingAuthority": "BAS / Medical", "ValidityMonths": 12,
     "RenewalProcess": "Complete annual HIPM course (4 hrs) before 1 May",
     "RequiredForEvents": "PI-003,PI-004,PI-025,PI-028,PII-010",
     "MinimumPerEvent": 1, "OrderReference": "MCO 6200.1F", "Status": "Active", "Notes": ""},
    {"QualCode": "SAFE-COLD", "QualName": "Cold Weather Injury Prevention Monitor", "Category": "Other",
     "IssuingAuthority": "BAS / Medical", "ValidityMonths": 12,
     "RenewalProcess": "Complete annual cold weather injury prevention course before 1 Oct",
     "RequiredForEvents": "PI-003,PI-025",
     "MinimumPerEvent": 1, "OrderReference": "MCO 6200.1F", "Status": "Active", "Notes": ""},
]

def write_required_quals():
    cols = ["QualCode", "QualName", "Category", "IssuingAuthority", "ValidityMonths",
            "RenewalProcess", "RequiredForEvents", "MinimumPerEvent", "OrderReference",
            "Status", "Notes"]
    path = os.path.join(OUTPUT_DIR, "required-qualifications.csv")
    with open(path, "w", newline="", encoding="utf-8") as f:
        w = csv.DictWriter(f, fieldnames=cols)
        w.writeheader()
        w.writerows(required_quals)
    print(f"  required-qualifications.csv: {len(required_quals)} rows")

# ─────────────────────────────────────────────
# 5. QUALIFICATION RECORDS (~80 rows)
# ─────────────────────────────────────────────

qual_records = []

# Build a lookup for validity months
validity_lookup = {q["QualCode"]: q["ValidityMonths"] for q in required_quals}
qual_name_lookup = {q["QualCode"]: q["QualName"] for q in required_quals}

# Assign qualifications to instructors based on their roles
def instructor_qual_map():
    """Map each instructor to a set of relevant qualifications."""
    mapping = {}
    for inst in instructors:
        role = inst["Role"]
        edipi = inst["InstructorEDIPI"]
        quals = set()

        # Everyone gets CPR-AED and at least one safety cert
        quals.add("CPR-AED")
        quals.add("SAFE-HEAT")

        if role == "Company Commander":
            quals.update(["OIC-RANGE", "OIC-DEMO", "RSO-RIFLE", "INST-PATROLLING"])
        elif role == "Company XO":
            quals.update(["OIC-RANGE", "RSO-RIFLE", "RSO-PISTOL"])
        elif role == "Staff Platoon Commander":
            quals.update(["RSO-RIFLE", "RSO-PISTOL", "INST-LANDNAV", "INST-PATROLLING", "CLS"])
        elif role == "Assistant SPC":
            quals.update(["RSO-RIFLE", "INST-LANDNAV", "CLS"])
        elif role == "Tactics Instructor":
            quals.update(["RSO-RIFLE", "RSO-PISTOL", "INST-PATROLLING", "INST-MOUT", "INST-CQB", "CLS", "SAFE-COLD"])
        elif role == "Weapons Instructor":
            quals.update(["RSO-RIFLE", "RSO-PISTOL", "RSO-DEMO", "PMI-RIFLE", "PMI-PISTOL", "OIC-RANGE"])
        elif role == "PT Instructor":
            quals.update(["INST-PT", "INST-SWIM", "SAFE-COLD"])

        # Some random extras
        if random.random() < 0.3:
            quals.add("DRV-HMMWV")
        if random.random() < 0.15:
            quals.add("INST-COMMS")

        mapping[edipi] = quals
    return mapping

inst_quals = instructor_qual_map()

# Now generate records with expiration status distribution:
# ~60% Current (>90 days), ~15% Caution (60-90), ~10% Warning (30-60), ~10% Critical (<30), ~5% Expired
STATUS_WEIGHTS = [
    ("current", 0.60),
    ("caution", 0.15),
    ("warning", 0.10),
    ("critical", 0.10),
    ("expired", 0.05),
]

def pick_expiration_status():
    r = random.random()
    cum = 0
    for status, weight in STATUS_WEIGHTS:
        cum += weight
        if r < cum:
            return status
    return "current"

cert_counter = 1000

for inst in instructors:
    edipi = inst["InstructorEDIPI"]
    name_display = f"{inst['LastName']}, {inst['FirstName']}"
    quals_for_this = inst_quals.get(edipi, set())

    for qcode in quals_for_this:
        validity_months = validity_lookup.get(qcode, 12)
        exp_status = pick_expiration_status()

        # Work backwards from desired expiration status to set DateEarned
        if exp_status == "current":
            days_until = random.randint(91, validity_months * 30)
            expiration = TODAY + timedelta(days=days_until)
        elif exp_status == "caution":
            days_until = random.randint(61, 90)
            expiration = TODAY + timedelta(days=days_until)
        elif exp_status == "warning":
            days_until = random.randint(31, 60)
            expiration = TODAY + timedelta(days=days_until)
        elif exp_status == "critical":
            days_until = random.randint(1, 30)
            expiration = TODAY + timedelta(days=days_until)
        else:  # expired
            days_until = random.randint(-60, -1)
            expiration = TODAY + timedelta(days=days_until)

        date_earned = expiration - timedelta(days=validity_months * 30)

        # Renewal status
        if days_until < 0:
            renewal = random.choice(["Renewal Overdue", "Renewal In Progress", "Waiver Requested"])
            renewal_date_val = ""
        elif days_until <= 60:
            renewal = random.choice(["Renewal Scheduled", "Renewal In Progress"])
            renewal_date_val = (expiration - timedelta(days=random.randint(1, 14))).isoformat()
        elif days_until <= 90:
            renewal = random.choice(["Renewal Scheduled", "N/A - Current"])
            renewal_date_val = (expiration - timedelta(days=random.randint(5, 30))).isoformat() if renewal == "Renewal Scheduled" else ""
        else:
            renewal = "N/A - Current"
            renewal_date_val = ""

        cert_counter += 1
        cert_num = f"TBS-{cert_counter:05d}"

        issuer = ""
        for rq in required_quals:
            if rq["QualCode"] == qcode:
                issuer = rq["IssuingAuthority"]
                break

        qual_records.append({
            "InstructorEDIPI": edipi,
            "InstructorName": name_display,
            "QualCode": qcode,
            "QualName": qual_name_lookup.get(qcode, qcode),
            "DateEarned": date_earned.isoformat(),
            "ExpirationDate": expiration.isoformat(),
            "DaysUntilExpiration": days_until,
            "ExpirationStatus": (
                "Expired" if days_until < 0 else
                "Critical (30 days)" if days_until <= 30 else
                "Warning (60 days)" if days_until <= 60 else
                "Caution (90 days)" if days_until <= 90 else
                "Current"
            ),
            "CertificateNumber": cert_num,
            "IssuedBy": issuer,
            "CompanyAtTimeOfCert": "Alpha",
            "RenewalStatus": renewal,
            "RenewalDate": renewal_date_val,
            "DocumentLink": "",
            "Notes": "",
        })

def write_qual_records():
    cols = ["InstructorEDIPI", "InstructorName", "QualCode", "QualName",
            "DateEarned", "ExpirationDate", "DaysUntilExpiration", "ExpirationStatus",
            "CertificateNumber", "IssuedBy", "CompanyAtTimeOfCert",
            "RenewalStatus", "RenewalDate", "DocumentLink", "Notes"]
    path = os.path.join(OUTPUT_DIR, "qualification-records.csv")
    with open(path, "w", newline="", encoding="utf-8") as f:
        w = csv.DictWriter(f, fieldnames=cols)
        w.writeheader()
        w.writerows(qual_records)
    print(f"  qualification-records.csv: {len(qual_records)} rows")

# ─────────────────────────────────────────────
# 6. EVENT FEEDBACK (~40 rows)
# ─────────────────────────────────────────────

feedback = []

# Only completed events get feedback
completed_events = [e for e in schedule_events if e["Status"] == "Complete"]

# Pick ~15 events to have feedback, with 2-4 feedback entries each
feedback_events = random.sample(completed_events, min(15, len(completed_events)))

rating_choices = ["1 - Ineffective", "2 - Below Average", "3 - Average", "4 - Above Average", "5 - Excellent"]
objectives_choices = ["All objectives met", "Most objectives met", "Some objectives met", "Few objectives met", "No objectives met"]
time_choices = ["Too short", "About right", "Too long", "Significant dead time"]
resource_choices = ["Fully resourced", "Mostly adequate", "Some shortages", "Significant shortages"]

SUSTAIN_COMMENTS = [
    "Instructor clearly explained the fundamentals before hands-on practice.",
    "Good ratio of instructors to students on the lanes.",
    "Well-organized flow from classroom to field application.",
    "Excellent use of realistic scenarios that reinforced earlier lectures.",
    "Safety brief was thorough without being excessive.",
    "Good integration of lessons learned from recent operational deployments.",
    "Effective use of time - minimal dead time between iterations.",
    "AAR conducted at the end was very valuable for consolidating learning.",
    "Instructors were patient and provided individual feedback.",
    "Range setup was efficient - minimal time waiting to fire.",
    "Night iteration was challenging but realistic.",
    "Good progressive build from individual to team tasks.",
    "Supporting materials (handouts, map sheets) were well prepared.",
    "Equipment was pre-staged and ready when students arrived.",
    "Peer leadership opportunities were well distributed.",
]

IMPROVE_COMMENTS = [
    "Could use more repetitions before the graded event.",
    "Need better lighting for the classroom portion - hard to see slides.",
    "Would benefit from more time on the practical application phase.",
    "Some students didn't get a chance to lead during iterations.",
    "Should provide the reading material further in advance.",
    "Range had too much downtime between relays - suggest adding a concurrent training station.",
    "Need more map sheets - several students had to share.",
    "Consider breaking into smaller groups for the hands-on portions.",
    "The timeline was too compressed - rushed through the last two objectives.",
    "Would help to have a demonstration before students attempt the task.",
    "Feedback forms should be distributed immediately after the event, not the next day.",
    "Water and shade at the range were insufficient for the heat conditions.",
    "Radio equipment had multiple deadlined items - need PMCS before the event.",
    "Recommend adding a sand table exercise before going to the field.",
    "Some instructors gave contradictory guidance on the patrol order format.",
]

SAFETY_COMMENTS = [
    "One student flagged another Marine during magazine change drill. Corrected immediately by RSO.",
    "Water bull ran dry at 1400 - resupply took 45 minutes in high heat conditions.",
    "Loose concertina wire found along the boundary of the obstacle course lane.",
    "Two students showed signs of heat exhaustion - corpsman treated on site.",
    "Night land nav course had an unmarked ravine not shown on the map overlay.",
]

ADDITIONAL_COMMENTS = [
    "Best training event so far in the cycle.",
    "Recommend scheduling this earlier in the phase.",
    "Would benefit from a follow-up practical next week.",
    "Good event but needs more ammo allocation next time.",
    "",
    "",
    "",
    "",
]

feedback_id = 0
for event in feedback_events:
    n_feedback = random.randint(2, 4)
    for _ in range(n_feedback):
        feedback_id += 1
        submitter_role = random.choice(["Student", "Student", "Student", "SPC", "Instructor"])

        if submitter_role == "Student":
            submitter_name = ""  # anonymous
        else:
            submitter_name = f"{random.choice(lead_instructors_pool)}"

        # Rating distribution: mostly 3-4, some 5, few 1-2
        rating_weights = [0.03, 0.07, 0.30, 0.40, 0.20]
        rating = random.choices(rating_choices, weights=rating_weights, k=1)[0]

        # Objectives — correlated with rating
        rating_num = int(rating[0])
        if rating_num >= 4:
            obj = random.choice(objectives_choices[:2])
        elif rating_num == 3:
            obj = random.choice(objectives_choices[:3])
        else:
            obj = random.choice(objectives_choices[2:])

        # Instructor effectiveness (students only)
        if submitter_role == "Student":
            ie_weights = [0.02, 0.08, 0.25, 0.40, 0.25]
            inst_eff = random.choices(rating_choices, weights=ie_weights, k=1)[0]
        else:
            inst_eff = ""

        time_mgmt = random.choices(time_choices, weights=[0.10, 0.60, 0.15, 0.15], k=1)[0]
        resources = random.choices(resource_choices, weights=[0.40, 0.35, 0.20, 0.05], k=1)[0]

        sustain = random.choice(SUSTAIN_COMMENTS)
        improve = random.choice(IMPROVE_COMMENTS)

        # Safety concerns: ~10% of feedback
        has_safety = random.random() < 0.10
        safety_text = random.choice(SAFETY_COMMENTS) if has_safety else ""

        additional = random.choice(ADDITIONAL_COMMENTS)

        # Submitted within 1-3 days of event
        event_end = date.fromisoformat(event["EndDate"])
        submitted = event_end + timedelta(days=random.randint(0, 2))

        # Review status
        if has_safety:
            review_status = random.choice(["Action Required", "Reviewed"])
            reviewed_by = random.choice(lead_instructors_pool[:3])
            action_taken = "Safety concern routed to S-3 and reviewed at weekly training meeting." if review_status == "Reviewed" else ""
        elif submitted < TODAY - timedelta(days=7):
            review_status = random.choice(["Reviewed", "Reviewed", "Closed"])
            reviewed_by = random.choice(lead_instructors_pool[:3])
            action_taken = ""
        else:
            review_status = "Pending Review"
            reviewed_by = ""
            action_taken = ""

        feedback.append({
            "EventTitle": event["EventTitle"],
            "EventCode": event["EventCode"],
            "EventDate": event["EndDate"],
            "TrainingPhase": event["TrainingPhase"],
            "CompanyAssigned": "Alpha",
            "SubmitterRole": submitter_role,
            "SubmitterName": submitter_name,
            "OverallRating": rating,
            "ObjectivesMet": obj,
            "InstructorEffectiveness": inst_eff,
            "TimeManagement": time_mgmt,
            "ResourceAdequacy": resources,
            "Sustains": sustain,
            "Improves": improve,
            "SafetyConcerns": safety_text,
            "HasSafetyConcern": has_safety,
            "AdditionalComments": additional,
            "SubmittedDate": datetime.combine(submitted, datetime.min.time()).isoformat(),
            "ReviewedBy": reviewed_by,
            "ReviewStatus": review_status,
            "ActionTaken": action_taken,
        })

def write_feedback():
    cols = [
        "EventTitle", "EventCode", "EventDate", "TrainingPhase", "CompanyAssigned",
        "SubmitterRole", "SubmitterName", "OverallRating", "ObjectivesMet",
        "InstructorEffectiveness", "TimeManagement", "ResourceAdequacy",
        "Sustains", "Improves", "SafetyConcerns", "HasSafetyConcern",
        "AdditionalComments", "SubmittedDate", "ReviewedBy", "ReviewStatus", "ActionTaken",
    ]
    path = os.path.join(OUTPUT_DIR, "event-feedback.csv")
    with open(path, "w", newline="", encoding="utf-8") as f:
        w = csv.DictWriter(f, fieldnames=cols)
        w.writeheader()
        w.writerows(feedback)
    print(f"  event-feedback.csv: {len(feedback)} rows")

# ─────────────────────────────────────────────
# GENERATE ALL FILES
# ─────────────────────────────────────────────
if __name__ == "__main__":
    print("Generating Heywood TBS sample data...")
    write_instructors()
    write_students()
    write_schedule()
    write_required_quals()
    write_qual_records()
    write_feedback()
    print("Done. All files in:", OUTPUT_DIR)
