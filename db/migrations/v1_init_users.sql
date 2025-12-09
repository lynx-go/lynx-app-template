create table if not exists public.users
(
    id                   bigint                                 not null
    constraint users_pk
    primary key,
    username             varchar(255)                           not null,
    display_name         varchar(255)                           not null,
    password_hash        varchar(255),
    avatar_url           varchar(512),
    phone                varchar(32),
    phone_confirmed_at   timestamp with time zone,
                                       email                varchar(256),
    email_confirmed_at   timestamp with time zone,
    status               integer                  default 0     not null,
    gender               integer                  default 0     not null,
    app_metadata         jsonb,
    user_metadata        jsonb,
    last_sign_in_at      timestamp with time zone default now() not null,
    banned_until         timestamp with time zone,
    confirmed_at         timestamp with time zone,
                                       confirmation_token   varchar(255),
    confirmation_sent_at timestamp with time zone,
                                       recovery_token       varchar(255),
    recovery_sent_at     timestamp with time zone,
                                       role                 varchar(255),
    is_super_admin       boolean                  default false,
    created_at           timestamp with time zone default now() not null,
    updated_at           timestamp with time zone default now() not null,
    created_by           bigint                                 not null,
    updated_by           bigint                                 not null
    );

comment on table public.users is '用户表';

comment on column public.users.id is '用户主键';

comment on column public.users.username is '用户名';

comment on column public.users.display_name is '显示用户名';

comment on column public.users.password_hash is '密码';

comment on column public.users.avatar_url is '头像';

comment on column public.users.phone_confirmed_at is '手机号验证时间戳';

comment on column public.users.email is '邮箱';

comment on column public.users.email_confirmed_at is '邮箱验证时间戳';

comment on column public.users.status is '状态，0-初始（未验证），1-正常（已验证），-1-被锁定，-2-已注销';

comment on column public.users.gender is '性别，0-未设置，1-男，2-女';

comment on column public.users.app_metadata is '应用元数据';

comment on column public.users.user_metadata is '用户元数据';

comment on column public.users.last_sign_in_at is '最近登录时间戳';

comment on column public.users.banned_until is '封禁到期时间戳';

comment on column public.users.confirmed_at is '验证时间戳';

comment on column public.users.role is '角色';

comment on column public.users.created_at is '创建时间戳';

comment on column public.users.updated_at is '更新时间戳';

alter table public.users
    owner to appsets;

create table if not exists public.identities
(
    id              bigint                                 not null
    constraint credentials_pk
    primary key,
    user_id         bigint                                 not null
    constraint identities_users_id_fk
    references public.users,
    provider        varchar(256)                           not null,
    provider_id     varchar(256)                           not null,
    identity_data   jsonb,
    last_sign_in_at timestamp with time zone default now() not null,
    created_at      timestamp with time zone default now() not null,
    updated_at      timestamp with time zone default now() not null,
    created_by      bigint                                 not null,
    updated_by      bigint                                 not null
    );

comment on table public.identities is '用户第三方认证';

comment on column public.identities.id is '主键';

comment on column public.identities.user_id is '用户 UID';

comment on column public.identities.provider is '认证提供方';

comment on column public.identities.provider_id is '第三方认证平台的用户唯一ID';

comment on column public.identities.identity_data is '认证的元数据';

comment on column public.identities.last_sign_in_at is '最近登录时间戳';

alter table public.identities
    owner to appsets;

create unique index if not exists credentials_user_id_provider_uindex
    on public.identities (user_id, provider);

create unique index if not exists users_username_uindex
    on public.users (username);

create index if not exists users_email_index
    on public.users (email);

create table if not exists public.refresh_tokens
(
    id         varchar(255)                           not null
    constraint refresh_tokens_pk
    primary key,
    user_id    bigint                                 not null,
    token      varchar(255)                           not null,
    revoked    boolean                  default false not null,
    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone default now() not null
    );

comment on table public.refresh_tokens is 'Refresh Token';

alter table public.refresh_tokens
    owner to appsets;

create index if not exists refresh_tokens_token_index
    on public.refresh_tokens (token);

create index if not exists refresh_tokens_user_id_index
    on public.refresh_tokens (user_id);
