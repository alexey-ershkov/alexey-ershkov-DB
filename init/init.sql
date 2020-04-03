-- Unknown how to generate base type type

CREATE EXTENSION IF NOT EXISTS CITEXT;

create table usr
(
    email    citext not null
        constraint user_pk
            primary key,
    fullname text   not null,
    nickname citext not null,
    about    text
);


create unique index usr_nickname_uindex
    on usr (nickname);

create table forum
(
    slug  citext not null
        constraint forum_pk
            primary key,
    title text   not null,
    usr   citext not null
        constraint forum_user_email_fk
            references usr (nickname)
            on update cascade on delete cascade
);


create table thread
(
    id      serial not null
        constraint thread_pk
            primary key,
    title   text   not null,
    message text   not null,
    created timestamp with time zone,
    slug    citext,
    usr     citext not null
        constraint thread_user_email_fk
            references usr (nickname)
            on update cascade on delete cascade,
    forum   citext not null
        constraint thread_forum_slug_fk
            references forum
            on update cascade on delete cascade
);


create unique index thread_slug_uindex
    on thread (slug);

create table post
(
    id       serial                not null
        constraint post_pk
            primary key,
    message  text                  not null,
    isedited boolean default false not null,
    parent   integer default 0,
    created  timestamp,
    usr      citext                not null
        constraint post_usr_nickname_fk
            references usr (nickname)
            on update cascade on delete cascade,
    thread   integer               not null
        constraint post_thread_id_fk
            references thread
            on update cascade on delete cascade
);


create table vote
(
    id     serial  not null
        constraint vote_pk
            primary key,
    vote   integer not null,
    usr    citext  not null
        constraint vote_usr_nickname_fk
            references usr (nickname)
            on update cascade on delete cascade,
    thread integer not null
        constraint vote_thread_id_fk
            references thread
            on update cascade on delete cascade
);



