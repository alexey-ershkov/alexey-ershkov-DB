drop trigger IF EXISTS path_updater ON post;
drop trigger if exists forum_users_clear on forum_users;
drop trigger if exists forum_user_insert_after_thread on thread;
drop trigger if exists forum_user_insert_after_post on post;
drop function IF EXISTS updater;
drop function IF EXISTS clear_forum_users;
drop function IF EXISTS insert_into_forum_users;
drop table IF EXISTS usr CASCADE;
drop table IF EXISTS forum CASCADE;
drop table IF EXISTS thread CASCADE;
drop table IF EXISTS post CASCADE;
drop table IF EXISTS vote CASCADE;
drop table IF EXISTS forum_users CASCADE;

create EXTENSION IF NOT EXISTS CITEXT;

create unlogged table usr
(
    email    citext collate "C" not null,
    fullname text   not null,
    nickname citext collate "C" not null,
    about    text,
    constraint user_pk primary key (email)
);

create index index_usr_nickname on usr (nickname);


create unique index usr_nickname_uindex
    on usr (nickname);

create unlogged table forum
(
    slug  citext collate "C" not null
        constraint forum_pk
            primary key,
    title text   not null,
    usr   citext collate "C" not null
        constraint forum_user_email_fk
            references usr (nickname)
            on update cascade on delete cascade
);


create index index_usr_fk on forum (usr);


create unlogged table thread
(
    id      serial not null
        constraint thread_pk
            primary key,
    title   text   not null,
    message text   not null,
    created timestamp with time zone,
    slug    citext collate "C",
    usr     citext collate "C" not null
        constraint thread_user_email_fk
            references usr (nickname)
            on update cascade on delete cascade,
    forum   citext collate "C" not null
        constraint thread_forum_slug_fk
            references forum
            on update cascade on delete cascade
);


create index index_thread_usr_fk on thread (usr);
create index index_thread_forum_fk on thread (forum);


create unique index thread_slug_uindex
    on thread (slug);

create unlogged table post
(
    id       bigserial             not null
        constraint post_pk
            primary key,
    message  text                  not null,
    isedited boolean default false not null,
    parent   integer default 0,
    created  timestamp,
    usr      citext  collate "C" not null
        constraint post_usr_nickname_fk
            references usr (nickname)
            on update cascade on delete cascade,
    thread   integer               not null
        constraint post_thread_id_fk
            references thread
            on update cascade on delete cascade,
    forum    citext                not null
        constraint post_forum_slug_fk
            references forum
            on update cascade on delete cascade,
    path     bigint[]
);


create index index_post_usr_fk on post (usr);
create index index_post_forum_fk on post(forum);
create index index_post_thread_fk on post(thread);


create unlogged table vote
(
    id     serial  not null
        constraint vote_pk
            primary key,
    vote   integer not null,
    usr    citext collate "C" not null
        constraint vote_usr_nickname_fk
            references usr (nickname)
            on update cascade on delete cascade,
    thread integer not null
        constraint vote_thread_id_fk
            references thread
            on update cascade on delete cascade,
    constraint vote_pk_2
        unique (usr, thread)
);




create unlogged table forum_users
(
    id bigserial primary key ,
    forum    citext collate "C" not null
        constraint forum_users_forum_slug_fk
            references forum
            on update cascade on delete cascade,
    nickname citext collate "C" not null
        constraint forum_users_usr_nickname_fk
            references usr (nickname)
            on update cascade on delete cascade

);

create or replace function updater()
    RETURNS trigger AS
$BODY$
begin
    update post set path = path || NEW.id WHERE thread = NEW.thread AND id = NEW.id;
    RETURN NEW;
END;
$BODY$ LANGUAGE plpgsql;

create trigger path_updater
    after insert
    on post
    for each row
EXECUTE procedure updater();

create or replace function insert_into_forum_users()
    returns trigger as
$insert_into_forum_users$
    begin
        insert into forum_users (nickname, forum) values (new.usr, new.forum)
        on conflict do nothing;
        return new;
    exception
        when SQLSTATE '40P01' then
            return new;
    end;
$insert_into_forum_users$ LANGUAGE plpgsql;

create trigger forum_user_insert_after_post
    after insert
    on post
    for each row
EXECUTE procedure insert_into_forum_users();

create trigger forum_user_insert_after_thread
    after insert
    on thread
    for each row
EXECUTE procedure insert_into_forum_users();

