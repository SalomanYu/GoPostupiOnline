from rich.progress import track

from models import *
from mongo_storage import *
from sql_storage import *


def saveVuzesToSQL(vuzes:list[Vuz]):
    for item in track(range(len(vuzes)), description="[green]Adding vuzes..."):
        add_institution(vuzes[item], DATABASE_POSTUPI)

def saveSpecsToSQL(specs:list[Specialization]):
    for item in track(range(len(specs)), description="[green]Adding specs..."):
        add_spec(specs[item], DATABASE_POSTUPI)

def saveProgramToSQL(programs:list[Program]):
    for item in track(range(len(programs)), description="[green]Adding programs..."):
        add_program(programs[item], DATABASE_POSTUPI)

def saveProfessionsToSQL(professions:list[Profession]):
    for item in track(range(len(professions)), description="[green]Adding professions..."):
        add_profession(professions[item], DATABASE_POSTUPI)

def saveContactsToSQL(contacts:list[Contacts]):
    for item in track(range(len(contacts)), description="[green]Adding contacts..."):
        add_contact(contacts[item], DATABASE_POSTUPI)

if __name__ == "__main__":
    create_table_for_contacts(DATABASE_POSTUPI)
    contacts = getContacts()
    programs = getAllPrograms()
    specs = getAllSpecs()
    vuzes = getAllVuzes()
    professions = getProfessions()
    saveProgramToSQL(programs)
    saveContactsToSQL(contacts)
    saveVuzesToSQL(vuzes)
    saveSpecsToSQL(specs)
    saveProfessionsToSQL(professions)