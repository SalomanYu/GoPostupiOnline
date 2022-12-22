import pymongo

from models import *

def getAllVuzes() -> list[Vuz]:
    client = pymongo.MongoClient("mongodb://localhost:27017/")
    db = client["PostupiOnline"]
    collection = db["Vuz"]
    vuzes:list[Vuz] = []
    for item in collection.find():
        vuzes.append(Vuz(
            vuzId=item["vuz_id"],
            description=item["description"],
            base=getBase(item["base"])
        ))
    return vuzes

def getAllSpecs() -> list[Specialization]:
    client = pymongo.MongoClient("mongodb://localhost:27017/")
    db = client["PostupiOnline"]
    collection = db["Specialization"]
    specs:list[Specialization] = []
    for item in collection.find():
        specs.append(Specialization(
            specId=item["spec_id"],
            vuzId=item["vuz_id"],
            description=item["description"],
            base=getBase(item["base"])
        ))
    return specs

def getAllPrograms() -> list[Program]:
    client = pymongo.MongoClient("mongodb://localhost:27017/")
    db = client["PostupiOnline"]
    collection = db["Program"]
    programs:list[Program] = [] 
    for item in collection.find():
        try:exams = "|".join(i for i in item["exams"])
        except TypeError: exams = ""
        programs.append(Program(
            programId=item["program_id"],
            hasProfession=item["has_professions"],
            form=item["form"],
            exams=exams,
            specId=item["spec_id"],
            vuzId=item["vuz_id"],
            description=item["description"],
            base=getBase(item["base"])
        ))
    return programs

def getProfessions() -> list[Profession]:
    client = pymongo.MongoClient("mongodb://localhost:27017/")
    db = client["PostupiOnline"]
    collection = db["Profession"]
    professions:list[Profession] = []
    for item in collection.find():
        professions.append(Profession(
            programId=item["program_id"],
            name=item["name"],
            image=item["image"]
        ))
    return professions

def getContacts() -> list[Contacts]:
    client = pymongo.MongoClient("mongodb://localhost:27017/")
    db = client["PostupiOnline"]
    collection = db["Contacts"]
    contacts:list[Contacts] = []
    for item in collection.find():
        try:phone = "|".join(item["phone"].split(","))
        except:phone = ""
        contacts.append(Contacts(
            vuzId=item["vuz_id"],
            website=item["website"],
            email=item["email"],
            phone=phone,
            address=item["address"]
        ))
    return contacts

def getBase(base: dict) -> Base:
    return Base(
        name=base["name"],
        url=base["url"],
        description=base["description"],
        direction=base["direction"],
        image=base["image"],
        logo=base["logo"],
        cost=base["cost"],
        scores=getScores(base["scores"])

    )

def getScores(scores: dict) -> Scores:
    return Scores(
        budgetPlaces=scores["budget_places"],
        paymentPlaces=scores["payment_places"],
        budgetPoints=scores["budget_points"],
        paymentPoints=scores["payment_points"]
    )