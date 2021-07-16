package app

const migrations = `
create table assoc
(
	id int auto_increment
		primary key,
	id1 int not null,
	id2 int not null,
	type varchar(50) not null,
	data json not null,
	constraint assoc_id1_type_id2_uindex
		unique (id1, type, id2)
);

create table ids
(
	id int auto_increment
		primary key
);

create table object
(
	id int auto_increment
		primary key,
	type smallint not null,
	data json not null
);




`