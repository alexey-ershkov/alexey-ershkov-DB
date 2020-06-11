ALTER SYSTEM SET checkpoint_completion_target = '0.9';
ALTER SYSTEM SET wal_buffers = '6912kB';
ALTER SYSTEM SET default_statistics_target = '100';
ALTER SYSTEM SET random_page_cost = '1.1';
ALTER SYSTEM SET effective_io_concurrency = '200';
ALTER SYSTEM SET seq_page_cost = '0.1';
ALTER SYSTEM SET random_page_cost = '0.1';
ALTER SYSTEM SET max_worker_processes = '4';
ALTER SYSTEM SET max_parallel_workers_per_gather = '2';
ALTER SYSTEM SET max_parallel_workers = '4';
ALTER SYSTEM SET max_parallel_maintenance_workers = '2';

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

create unlogged table usr
(
    email    citext collate "C" not null,
    fullname text               not null,
    nickname citext collate "C" not null,
    about    text,
    constraint user_pk primary key (email)
);

create index index_usr_nickname on usr (nickname);
create index index_usr_all on usr (nickname, fullname, email, about);
cluster usr using index_usr_all;


create unique index usr_nickname_uindex
    on usr (nickname);

create unlogged table forum
(
    slug    citext collate "C" not null
        constraint forum_pk
            primary key,
    title   text               not null,
    usr     citext collate "C" not null
        constraint forum_user_email_fk
            references usr (nickname)
            on update cascade on delete cascade,
    threads bigint default 0,
    posts   bigint default 0
);

create index index_forum_slug on forum (slug);
create index index_forum_slug_hash on forum using hash (slug);
cluster forum using index_forum_slug_hash;
create index index_usr_fk on forum (usr);
create index index_forum_all on forum (slug, title, usr, posts, threads);

create unlogged table thread
(
    id      serial             not null
        constraint thread_pk
            primary key,
    title   text               not null,
    message text               not null,
    created timestamp with time zone,
    slug    citext collate "C",
    votes   int default 0,
    usr     citext collate "C" not null
        constraint thread_user_email_fk
            references usr (nickname)
            on update cascade on delete cascade,
    forum   citext collate "C" not null
        constraint thread_forum_slug_fk
            references forum
            on update cascade on delete cascade
);

create index index_thread_forum_created on thread (forum, created);
-- cluster thread using index_thread_forum_created;
create index index_thread_id_and_slug on thread (CITEXT(id), slug);
create index index_thread_id on thread (id);
create index index_thread_slug on thread (slug);
create index index_thread_slug_hash on thread using hash (slug);
create index index_thread_all on thread (usr, forum, message, title);
create index index_thread_usr_fk on thread (usr);
create index index_thread_forum_fk on thread (forum);


create unique index thread_slug_uindex
    on thread (slug);

create unlogged table post
(
    id       bigserial
        constraint post_pk
            primary key,
    message  text                  not null,
    isedited boolean default false not null,
    parent   integer default 0,
    created  timestamp,
    usr      citext collate "C"    not null
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

-- create index index_post_thread_path on post (thread, path);
create index index_post_thread_parent_path on post (thread, parent, path);
-- create index index_post_path1_path on post ((path[1]), path);
-- cluster post using index_post_path1_path;
create index index_post_thread_created_id on post (thread, created, id);

create index index_post_usr_fk on post (usr);
create index index_post_forum_fk on post (forum);



create unlogged table vote
(
    id     serial             not null
        constraint vote_pk
            primary key,
    vote   integer            not null,
    usr    citext collate "C" not null
        constraint vote_usr_nickname_fk
            references usr (nickname)
            on update cascade on delete cascade,
    thread integer            not null
        constraint vote_thread_id_fk
            references thread
            on update cascade on delete cascade,
    constraint vote_pk_2
        unique (usr, thread)
);


create index index_vote_thread on vote (thread);


create unlogged table forum_users
(
    forum    citext collate "C" not null
        constraint forum_users_forum_slug_fk
            references forum
            on update cascade on delete cascade,
    nickname citext collate "C" not null
        constraint forum_users_usr_nickname_fk
            references usr (nickname)
            on update cascade on delete cascade

);

create unique index index_forum_nickname on forum_users (forum, nickname);
-- cluster forum_users using index_forum_user_nickname;

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

create or replace function insert_into_forum_users()
    returns trigger as
$insert_into_forum_users$
begin
    insert into forum_users (nickname, forum)
    values (new.usr, new.forum)
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


create or replace function insert_thread_votes()
    returns trigger as
$update_thread_votes$
declare
    prev_vote int;
begin
    select vote
    from vote
    where thread = new.thread
      and usr = new.usr
    into prev_vote;
    if not FOUND then
        update thread set votes = (votes + new.vote) where id = new.thread;
    else
        if prev_vote != new.vote then
            update thread set votes = (votes + 2 * new.vote) where id = new.thread;
        end if;
    end if;
    return new;
end;
$update_thread_votes$ language plpgsql;


create trigger update_thread_votes
    before insert
    on vote
    for each row
execute procedure insert_thread_votes();

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