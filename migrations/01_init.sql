create database meme9;
use meme9;

create table edges
(
    id        int auto_increment
        primary key,
    from_id   int not null,
    to_id     int not null,
    edge_type int not null,
    date      int not null,
    constraint edges_from_id_to_id_edge_type_uindex
        unique (from_id, to_id, edge_type)
);

create table objects
(
    id       bigint auto_increment
        primary key,
    obj_type smallint not null,
    data     json     not null
);

create table meme9.uniques
(
    id        int auto_increment
        primary key,
    type      int          not null,
    `key`     varchar(200) not null,
    object_id int          not null,
    constraint table_name_pk
        unique (type, `key`)
);



insert into objects(id, obj_type, data)
values (-5, 2, '{
  "SaveSecret": "test"
}');
