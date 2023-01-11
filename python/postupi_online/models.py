from typing import NamedTuple

class Scores(NamedTuple):
    budgetPoints:int
    paymentPoints:int
    budgetPlaces:str
    paymentPlaces:str

class Base(NamedTuple):
    name: str
    url: str
    description:str
    direction:str
    image:str
    logo:str
    cost:str
    scores:Scores

class Vuz(NamedTuple):
    vuzId:str
    description:str
    base:Base

class Contacts(NamedTuple):
    vuzId:str
    website:str
    email:str
    phone:str
    address:str

class Profession(NamedTuple):
    programId:str
    name:str
    image:str


class Specialization(NamedTuple):
    specId:str
    vuzId:str
    description:str
    base:Base

class Program(NamedTuple):
    programId:str
    specId:str
    vuzId:str
    hasProfession:bool
    description:str
    form:str
    exams:list[str]
    base:Base

class Education(NamedTuple):
    title:str
    direction:str
    year:int

class Training(NamedTuple):
    title:str
    direction:str
    year:int

