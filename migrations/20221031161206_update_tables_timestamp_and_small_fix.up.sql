alter table if exists public.users
add created_date timestamp;

alter table if exists public.users
add deleted_date timestamp;

alter table if exists public.users
add updated_date timestamp;

alter table if exists public.posts
add created_date timestamp;

alter table if exists public.posts
add deleted_date timestamp;

alter table if exists public.posts
add updated_date timestamp;

alter table if exists public.commentses
add created_date timestamp;

alter table if exists public.commentses
add deleted_date timestamp;

alter table if exists public.commentses
add updated_date timestamp;