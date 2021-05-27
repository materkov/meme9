create table likes
(
    id int auto_increment
        primary key,
    post_id int not null,
    user_id int not null,
    time int not null,
    constraint likes_post_id_user_id_uindex
        unique (post_id, user_id)
);
create table followers
(
    id int auto_increment
        primary key,
    user1_id int not null,
    user2_id int not null,
    follow_date int not null,
    constraint followers_user1_id_user2_id_uindex
        unique (user1_id, user2_id)
);
create table comment
(
    id int auto_increment
        primary key,
    post_id int not null,
    text varchar(800) null,
    date int not null,
    user_id int not null
);



create table friend
(
    id int auto_increment
        primary key,
    user1 int not null,
    user2 int null
);



create table photo
(
    id int auto_increment
        primary key,
    url varchar(500) null
);

create table post
(
    id int auto_increment
        primary key,
    user_id int not null,
    date int null,
    text text null,
    photo_id int null
);

create table token
(
    id int auto_increment
        primary key,
    token varchar(50) null,
    user_id int not null,
    constraint token_token_uindex
        unique (token)
);

create table user
(
    id int auto_increment
        primary key,
    name varchar(200) null,
    avatar_id int null,
    vk_id int null,
    vk_avatar varchar(500) null,
    constraint user_vk_id_uindex
        unique (vk_id)
);

create table user_mentions
(
    id int auto_increment
        primary key,
    user_id int null,
    mention varchar(200) null,
    mention_time int null
);

create table user_profile
(
    ud int not null
        primary key,
    full_name varchar(500) null
);

create table objects
(
    id int auto_increment
        primary key,
    object_type smallint not null
);

