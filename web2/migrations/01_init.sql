create table post
(
    id      int auto_increment primary key,
    text    text,
    user_id int not null default 0
);

create table user
(
    id    int auto_increment primary key,
    name  varchar(200) not null default '',
    vk_id int,

    constraint user_vk_id_uindex
        unique (vk_id)
);
