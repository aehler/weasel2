create schema weasel_auth;

create sequence weasel_auth.user_id start 1024;
create sequence weasel_auth.organization_id start 1024;

create table weasel_auth.organizations(
    organization_id bigint not null default nextval('weasel_auth.organization_id'),
    created_at timestamp not null default current_timestamp,
    organization_name varchar not null default '',
    primary key(organization_id)
);

create table weasel_auth.users(
    user_id bigint NOT NULL default nextval('weasel_auth.user_id'),
    user_password varchar(64) not null,
    user_firstname varchar(64) not null,
    user_lastname varchar(64) not null,
    user_middlename varchar(64) not null,
    user_groups     bigint[]  not null default '{}',
    user_roles      bigint[]  not null default '{}',
    user_login varchar(128) not null,
    user_job_title varchar(255) not null default '',
    user_phone varchar not null default '',
    user_email varchar(128) not null,
    email_is_confirmed boolean not null default false,
    password_expiration_date date not null default current_date,
    organization_id bigint  not null references weasel_auth.organizations,
    timezone_id  int not null default 0,
    is_admin     boolean not null default false,
    is_active boolean not null default false,
    is_deleted boolean not null default false,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    primary key(user_id)
) with(fillfactor = 90);

create unique index uidx_users_email on weasel_auth.users (lower(user_email));
create index idx_user_member_in on weasel_auth.users (organization_id);
create unique index uidx_users_password on weasel_auth.users (user_password, user_login);