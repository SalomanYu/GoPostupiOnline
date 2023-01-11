import sqlite3
from models import *

DATABASE_POSTUPI = "postupi_online_colleges.db"
TABLE_VUZES = "vuzes"
TABLE_SPECS = "specializations"
TABLE_PROGRAMS = "programs"
TABLE_PROFESSIONS = "professions"
TABLE_CONTACTS = "contacts"
COLUMNS_COUNT_FOR_VUZ = 11
COLUMNS_COUNT_FOR_SPEC = 12
COLUMNS_COUNT_FOR_PROGRAM = 16
COLUMNS_COUNT_FOR_PROFESSION = 3
COLUMNS_COUNT_FOR_CONTACTS = 5

def connect_to_sql(db_name:str):
    db = sqlite3.connect(db_name)
    cursor = db.cursor()
    return db, cursor


def create_table_for_institution(db_name: str):
    db, cursor = connect_to_sql(db_name)
    cursor.execute(f"""CREATE TABLE IF NOT EXISTS {TABLE_VUZES}(
        institutionID VARCHAR(50),
        name TEXT,
        description TEXT,
        img TEXT,
        logo TEXT,
        cost INTEGER,
        budget_places INTEGER,
        payment_places INTEGER,
        budget_points FLOAT,
        payment_points FLOAT,
        url TEXT
    )""")
    db.commit()
    db.close()

def create_table_for_specialization(db_name: str):
    db, cursor = connect_to_sql(db_name)
    cursor.execute(f"""CREATE TABLE IF NOT EXISTS {TABLE_SPECS}(
        specID VARCHAR(20),
        institutionID VARCHAR(30) REFERENCES institution(institutionID),
        name TEXT,
        description TEXT,
        direction TEXT,
        img TEXT,
        cost INTEGER,
        budget_places INTEGER,
        payment_places INTEGER,
        budget_points FLOAT,
        payment_points FLOAT,
        url TEXT
    )""")
    db.commit()
    db.close()

def create_table_for_programs(db_name: str):
    db, cursor = connect_to_sql(db_name)
    cursor.execute(f"""CREATE TABLE IF NOT EXISTS {TABLE_PROGRAMS}(
        programID INTEGER,
        specID INTEGER REFERENCES specialization(specID),
        institutionID VARCHAR(50) REFERENCES institution(institutionID),
        name TEXT,
        description TEXT,
        direction TEXT,
        form VARCHAR(50),
        img TEXT,
        cost INTEGER,
        budget_places INTEGER,
        payment_places INTEGER,
        budget_points FLOAT,
        payment_points FLOAT,
        subjects VARCHAR(255),
        has_professions BOOLEAN,
        url TEXT
    )""")
    db.commit()
    db.close()

def create_table_for_profession(db_name: str):
    db, cursor = connect_to_sql(db_name)
    cursor.execute(f"""CREATE TABLE IF NOT EXISTS {TABLE_PROFESSIONS}(
        professionID INTEGER PRIMARY KEY AUTOINCREMENT,
        programID INTEGER REFERENCES program(programID),
        name TEXT,
        img TEXT
    )""")
    db.commit()
    db.close()


def create_table_for_contacts(db_name: str):
    db, cursor = connect_to_sql(db_name)
    cursor.execute(f"""CREATE TABLE IF NOT EXISTS {TABLE_CONTACTS}(
        contactID INTEGER PRIMARY KEY AUTOINCREMENT,
        website VARCHAR(100),
        email VARCHAR(100),
        phones VARCHAR(100),
        address TEXT,
        institutionID VARCHAR(50) REFERENCES {TABLE_VUZES}(institutionID)
    )""")
    db.commit()
    db.close()

def add_institution(data:Vuz, db_name: str = ""):
    create_table_for_institution(db_name)
    db, cursor = connect_to_sql(db_name)
    cursor.execute(f"INSERT INTO {TABLE_VUZES}(institutionID, name , description, img, logo, cost, budget_places, payment_places, budget_points ,payment_points, url) VALUES({','.join(['?' for i in range(COLUMNS_COUNT_FOR_VUZ)])})",
    (data.vuzId, data.base.name,data.description, data.base.image, data.base.logo, data.base.cost, data.base.scores.budgetPlaces, data.base.scores.paymentPlaces, data.base.scores.budgetPoints, data.base.scores.paymentPoints, data.base.url))
    db.commit()
    db.close()

def add_spec(data:Specialization, db_name: str = ""):
    create_table_for_specialization(db_name)
    db, cursor = connect_to_sql(db_name)
    cursor.execute(f"INSERT INTO {TABLE_SPECS}(specID, institutionID, name , description, direction, img, cost, budget_places, payment_places, budget_points ,payment_points, url) VALUES({','.join(['?' for i in range(COLUMNS_COUNT_FOR_SPEC)])})", 
    (data.specId, data.vuzId, data.base.name, data.base.description, data.base.direction, data.base.image, data.base.cost, data.base.scores.budgetPlaces, data.base.scores.paymentPlaces, data.base.scores.budgetPoints, data.base.scores.paymentPoints, data.base.url))
    db.commit()
    db.close()

def add_contact(data:Contacts, db_name: str = ""):
    create_table_for_contacts(db_name)
    db, cursor = connect_to_sql(db_name)
    cursor.execute(f"INSERT INTO {TABLE_CONTACTS}(website, email, phones, address, institutionID) VALUES({','.join(['?' for i in range(COLUMNS_COUNT_FOR_CONTACTS)])})",
    (data.website, data.email, data.phone, data.address, data.vuzId))
    db.commit()
    db.close() 

def add_program(data:Program, db_name: str = ""):
    create_table_for_programs(db_name)
    db, cursor = connect_to_sql(db_name)
    cursor.execute(f"INSERT INTO {TABLE_PROGRAMS}(programID, specID, institutionID, name, description, direction, form, img, cost, budget_places, payment_places, budget_points ,payment_points, subjects, has_professions, url) VALUES({','.join(['?' for i in range(COLUMNS_COUNT_FOR_PROGRAM)])})",
    (data.programId, data.specId, data.vuzId, data.base.name, data.base.description, data.base.direction, data.form, data.base.image, data.base.cost, data.base.scores.budgetPlaces, data.base.scores.paymentPlaces, data.base.scores.budgetPoints, data.base.scores.paymentPoints, data.exams, data.hasProfession, data.base.url))
    db.commit()
    db.close()

def add_profession(data:Profession, db_name: str = ""):
    create_table_for_profession(db_name)
    db, cursor = connect_to_sql(db_name)
    cursor.execute(f"INSERT INTO {TABLE_PROFESSIONS}(programID, name, img) VALUES({','.join(['?' for i in range(COLUMNS_COUNT_FOR_PROFESSION)])})",
    (data.programId, data.name, data.image))
    db.commit()
    db.close()

def getPrograms() -> list[tuple[str]]:
    db, cursor = connect_to_sql(DATABASE_POSTUPI)
    cursor.execute(f"SELECT * FROM {TABLE_PROGRAMS}")
    programs = cursor.fetchall()
    db.close()
    return programs

def getVuzes() -> list[tuple[str]]:
    db, cursor = connect_to_sql(DATABASE_POSTUPI)
    cursor.execute(f"SELECT * FROM {TABLE_VUZES}")
    vuzes = cursor.fetchall()
    db.close()
    return vuzes

def getSpecs() -> list[tuple[str]]:
    db, cursor = connect_to_sql(DATABASE_POSTUPI)
    cursor.execute(f"SELECT * FROM {TABLE_SPECS}")
    specs = cursor.fetchall()
    db.close()
    return specs