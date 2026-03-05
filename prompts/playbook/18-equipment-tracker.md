# Prompt 18: Equipment/Supply Tracker

**Use:** Design a Power App for tracking company-level equipment and supply readiness
**Platform:** GenAI.mil
**Complexity:** Intermediate
**Source:** Adapted from EDD Equipment Tracker (Admin Tools #1)

---

## Prompt

```
I need a Power Apps canvas app for tracking TBS company equipment and supply readiness. The app should have:

1. A SharePoint list as the data source with columns:
   - ItemName (text)
   - NSN (text — National Stock Number, if applicable)
   - SerialNumber (text)
   - Category (choice: Communications, Optics, Weapons, IT/Electronics, Field Equipment, Vehicles, Training Aids, Other)
   - AssignedTo (choice: Company HQ, 1st Platoon, 2nd Platoon, 3rd Platoon, 4th Platoon, Armory, Supply)
   - Condition (choice: Serviceable, Unserviceable, In Repair, Deadline, Missing)
   - LastInventoryDate (date)
   - NextInventoryDue (calculated: LastInventoryDate + 30 days)
   - Location (text)
   - CustodianRank (text)
   - Notes (multiline text)

2. Three screens:
   - **Inventory List:** Gallery view with search and filter by Category and Condition. Show color-coded condition badges (green=serviceable, red=deadline/missing, yellow=in repair).
   - **Item Detail:** View/edit form for a selected item with photo capability (take picture of equipment condition)
   - **Add New Item:** Form with validation (SerialNumber required, no duplicates within Category)

3. A banner at the top displaying "[COMPANY] Company Equipment Tracker" and the current user.

4. Dashboard summary at the top of Inventory List:
   - Total items: [N]
   - Serviceable: [N] (green)
   - Unserviceable/Deadline: [N] (red)
   - In Repair: [N] (yellow)
   - Missing: [N] (red, bold)

5. Export capability: Generate a formatted equipment status report for the company commander.

Start with the SharePoint list schema, then build each screen.
```
