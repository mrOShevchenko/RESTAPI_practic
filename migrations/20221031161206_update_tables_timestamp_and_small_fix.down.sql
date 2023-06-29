alter table if exists public.users
drop column created_date;

alter table if exists public.users
drop column deleted_date;

alter table if exists public.users
drop column updated_date;

alter table if exists public.posts
drop column created_date;

alter table if exists public.posts
drop column deleted_date;

alter table if exists public.posts
drop column updated_date;

alter table if exists public.commentses
drop column created_date;

alter table if exists public.commentses
drop column deleted_date;

alter table if exists public.commentses
drop column updated_date;
