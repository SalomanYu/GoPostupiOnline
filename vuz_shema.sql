 create table contacts(
 id serial,
 institution_id text not null,
 address text,
 email text,
 phones text[],
 website text,
 primary key (id),
 foreign key (institution_id) references institution(id)
 );

  create table profession(
 id serial,
 program_id integer not null,
 name text not null,
 image text,
 primary key (id)
 );

 create table program(
 id integer not null,
 specialization_id varchar(8) not null,
 institution_id text not null,
 name text not null,
 description text,
 direction text not null,
 form text not null,
 subjects text[],
 cost integer,
 budget_places integer,
 payment_places integer,
 budget_points float,
 payment_points float,
 image text,
 url text not null,
 has_professions boolean,
 primary key (id, specialization_id, institution_id),
 foreign key (institution_id) references institution(id),
 check (direction in ('Бакалавриат', 'Специалитет', 'Магистратура'))
 );

 create table specialization(
 id varchar(8) not null,
 institution_id text not null,
 name text not null,
 description text,
 direction text not null,
 cost integer,
 budget_places integer,
 payment_places integer,
 budget_points float,
 payment_points float,
 image text,
 url text not null,
 primary key (id, institution_id),
 foreign key (institution_id) references institution(id)
 );

 create table institution(
 id text not null,
 name text not null,
 description text,
 cost integer,
 budget_places integer,
 payment_places integer,
 budget_points float,
 payment_points float,
 image text,
 logo text,
 url text not null,
 primary key (id));