create table public.posts
(
    id       serial primary key,
    user_id  integer not null,
    title    varchar not null,
    body     varchar not null
);