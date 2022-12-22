import json
from sql_storage import getPrograms, getVuzes


ProgramIdColumn = 0
ProgramVuzIdColumn = 2
ProgramHasProfessionsColumn = -2

VuzId = 0
VuzName = 1

def getProgramsIdWhoHasProfessions(vuzId:str) -> list[str]:
    programs = getPrograms()
    return [i[ProgramIdColumn] for i in programs if i[ProgramHasProfessionsColumn] == 1 and i[ProgramVuzIdColumn] == vuzId]

def saveToJson(data: list[str], filename:str):
    with open(filename, "a", encoding="utf-8") as file:
        json.dump(data, file, ensure_ascii=False, indent=2)
    print(filename, "Saved!")


def find_vuzes_without_professions(vuzes:list[tuple[str]]):
    vuzes_without_professions = []
    for vuz in vuzes:
        vuzHasProfessions = getProgramsIdWhoHasProfessions(vuz[VuzId])
        if not vuzHasProfessions:
            vuzes_without_professions.append(vuz[VuzName])
    return vuzes_without_professions

def find_programs_without_professions(programs:list[tuple[str]]):
    programs_without_professions = [i[ProgramIdColumn] for i in programs if i[ProgramHasProfessionsColumn] == 0]
    return programs_without_professions

if __name__ == "__main__":
    vuzes = getVuzes()
    vuzes_without_professions = find_vuzes_without_professions(vuzes)
    saveToJson(vuzes_without_professions, "vuzes_without_professions.json")
    programs = getPrograms()
    programs_without_professions = find_programs_without_professions(programs)
    saveToJson(programs_without_professions, "programs_without_professions.json")