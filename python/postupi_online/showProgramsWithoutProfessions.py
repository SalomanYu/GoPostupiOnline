from sql_storage import getPrograms, getVuzes


ProgramIdColumn = 0
ProgramVuzIdColumn = 2
ProgramHasProfessionsColumn = -2

VuzId = 0
VuzName = 1

def getProgramsIdWhoHasProfessions(vuzId:str) -> list[str]:
    programs = getPrograms()
    return [i[ProgramIdColumn] for i in programs if i[ProgramHasProfessionsColumn] == 1 and i[ProgramVuzIdColumn] == vuzId]

if __name__ == "__main__":
    vuzes = getVuzes()
    for vuz in vuzes:
        vuzHasProfessions = getProgramsIdWhoHasProfessions(vuz[VuzId])
        if not vuzHasProfessions:
            print(vuz[VuzName])