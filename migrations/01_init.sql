create database meme9;
use meme9;

create table edges
(
    id         int auto_increment
        primary key,
    from_id    int          not null,
    to_id      int          not null,
    edge_type  int          not null,
    unique_key varchar(400) null,
    constraint unique_edges_key
        unique (from_id, edge_type, to_id, unique_key)
);

create table objects
(
    id       bigint auto_increment
        primary key,
    obj_type smallint not null,
    data     json     not null
);


insert into objects(id, obj_type, data)
values (5, 2, '{
  "SaveSecret": "test"
}');
