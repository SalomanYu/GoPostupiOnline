import json
from sql_storage import getPrograms, getSpecs


ProgramIdColumn = 0
ProgramSpecIdColumn = 1
ProgramVuzIdColumn = 2
ProgramTitleColumn = 3
ProgramHasProfessionsColumn = -2

SpecIdColumn = 0
SpecVuzIdColumn = 1 
SpecNameColumn = 2
VuzId = 0
VuzName = 1


def saveToJson(data: list[str], filename:str):
    with open(filename, "w", encoding="utf-8") as file:
        json.dump(data, file, ensure_ascii=False, indent=2)
    print(filename, "Saved!")


def find_specs_without_professions(specs:list[tuple[str]]):
    specs_without_professions = []
    programs = getPrograms()
    programsWithProfessions = [i[ProgramSpecIdColumn] for i in programs if i[ProgramHasProfessionsColumn] == 1]
    for spec in specs:
        if spec[SpecVuzIdColumn] not in programsWithProfessions:
            specs_without_professions.append((spec[SpecNameColumn], spec[SpecIdColumn]))
    return specs_without_professions

def find_programs_without_professions(programs:list[tuple[str]]):
    programs_without_professions = [(i[ProgramIdColumn], i[ProgramTitleColumn], i[ProgramSpecIdColumn], i[ProgramVuzIdColumn]) for i in programs if i[ProgramHasProfessionsColumn] == 0]
    return programs_without_professions

if __name__ == "__main__":
    programs_without_professions = find_programs_without_professions(getPrograms())    
    specs = find_specs_without_professions(getSpecs())
    data = []
    for spec in specs:
        for prog in programs_without_professions:
            if prog[2] == spec[0]:
                data.append({
                    "spec_id": spec[1],
                    "program_id": prog[0],
                    "vuz_id": prog[3],
                    "program_title": prog[1],
                    "spec_title": spec[0],
                })
    saveToJson(data, "specs_without_professions.json")
    