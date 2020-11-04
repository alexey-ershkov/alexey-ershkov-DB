drop trigger IF EXISTS path_updater ON post;
drop trigger if exists forum_users_clear on forum_users;
drop trigger if exists forum_user_insert_after_thread on thread;
drop trigger if exists forum_user_insert_after_post on post;
drop trigger if exists insert_into_thread_votes on thread;
drop trigger if exists update_thread_votes on thread;
drop trigger if exists upd_forum_threads on thread;
drop function IF EXISTS updater;
drop function IF EXISTS insert_thread_votes;
drop function IF EXISTS insert_into_forum_users;
drop function if exists update_forum_threads;
drop table IF EXISTS usr CASCADE;
drop table IF EXISTS forum CASCADE;
drop table IF EXISTS thread CASCADE;
drop table IF EXISTS post CASCADE;
drop table IF EXISTS vote CASCADE;
drop table IF EXISTS forum_users CASCADE;

create EXTENSION IF NOT EXISTS CITEXT;

------------------------ USER ------------------------------------------
create unlogged table usr
(
    id       serial primary key,
    email    citext collate "C" not null unique,
    fullname text               not null,
    nickname citext collate "C" not null unique,
    about    text
);



------------------------- FORUM ------------------------------------
create unlogged table forum
(
    id      serial primary key,
    slug    citext collate "C" not null unique,
    title   text               not null,
    usr     citext collate "C" not null
        references usr (nickname)
            on delete cascade,
    threads bigint default 0,
    posts   bigint default 0
);



------------------------- THREAD ---------------------------------------
create unlogged table thread
(
    id      serial primary key,
    title   text               not null,
    message text               not null,
    created timestamp with time zone,
    slug    citext collate "C" unique,
    votes   int default 0,
    usr     citext collate "C" not null
        references usr (nickname)
            on delete cascade,
    forum   citext collate "C" not null
        references forum (slug)
            on delete cascade
);

------------------------- POST --------------------------------------------------------------
create unlogged table post
(
    id       bigserial          primary key,
    message  text                  not null,
    isedited boolean default false not null,
    parent   integer default 0,
    created  timestamp,
    usr      citext collate "C"    not null
             references usr (nickname)
             on delete cascade,
    thread   integer               not null
             references thread
             on delete cascade,
    forum    citext                not null
             references forum (slug)
             on delete cascade,
    path     bigint[]
);



--------------------------- VOTE ------------------------------------------
create unlogged table vote
(
    id     serial           primary key,
    vote   integer            not null,
    usr    citext collate "C" not null
            references usr (nickname)
            on delete cascade,
    thread integer            not null
            references thread
            on delete cascade
);


------------------------------ FORUM USERS -----------------------------------
create unlogged table forum_users
(
    forum    citext collate "C" not null
            references forum (slug)  on delete cascade,
    nickname citext collate "C" not null
            references usr (nickname) on delete cascade
);




---------------------- UPDATE PATH AND CHECK PARENT ---------------------------
create or replace function updater()
    RETURNS trigger AS
$BODY$
declare
    parent_path         bigint[];
    first_parent_thread int;
begin
    if (NEW.parent = 0) then
        NEW.path := array_append(NEW.path, NEW.id);
    else
        select thread, path
        from post
        where thread = NEW.thread
          and id = NEW.parent
        into first_parent_thread , parent_path;
        if not FOUND or first_parent_thread != NEW.thread then
            raise exception 'Parent post was created in another thread' using errcode = '00404';
        end if;

        NEW.path := parent_path || NEW.id;
    end if;
    return NEW;
end;
$BODY$ language plpgsql;

create trigger path_updater
    before insert
    on post
    for each row
EXECUTE procedure updater();


-------------------------------- INSERT THREAD VOTES -----------------------
create or replace function insert_thread_votes()
    returns trigger as
$insert_thread_votes$
declare
begin
    update thread set votes = (votes + new.vote) where id = new.thread;
    return new;
end;
$insert_thread_votes$ language plpgsql;

create trigger insert_thread_votes
    before insert
    on vote
    for each row
execute procedure insert_thread_votes();


------------------------------- UPDATE THREAD VOTES -------------------------
create or replace function update_thread_votes()
    returns trigger as
$update_thread_votes$
begin
    if new.vote > 0 then
        update thread set votes = (votes + 2) where id = new.thread;
    else
        update thread set votes = (votes - 2) where id = new.thread;
    end if;
    return new;
end;
$update_thread_votes$ language plpgsql;

create trigger update_thread_votes
    before update
    on vote
    for each row
execute procedure update_thread_votes();


------------------------------- UPDATE FORUM THREADS -------------------
create or replace function update_forum_threads()
    returns trigger as
$update_forum_threads$
begin
    update forum set threads = (threads + 1) where slug = new.forum;
    return new;
end;
$update_forum_threads$ language plpgsql;

create trigger upd_forum_threads
    after insert
    on thread
    for each row
execute procedure update_forum_threads();