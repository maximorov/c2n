create function set_updated_column() returns trigger
    language plpgsql
as
$$
BEGIN
    NEW.updated = now();
RETURN NEW;
END;
$$;
create trigger updated
    before update
    on tasks
    for each row
    execute procedure set_updated_column();
create type enum_task_activity_status as enum ('taken', 'completed', 'expired', 'refused');
create type enum_task_status as enum ('raw', 'new', 'in_progress', 'done', 'expired', 'cancelled', 'refused');

create table if not exists tasks
(
    id       bigserial
    constraint tasks_pk
    primary key,
    user_id  bigint                                                               not null,
    position point                                                                not null,
    status   enum_task_status default 'raw'::enum_task_status                     not null,
    text     varchar(255)                                                         not null,
    deadline timestamp        default (CURRENT_TIMESTAMP + '1 day'::interval day) not null,
    created  timestamp        default CURRENT_TIMESTAMP                           not null,
    updated  timestamp
    );

create unique index if not exists tasks_id_uindex
    on tasks (id);

create trigger updated
    before update
    on tasks
    for each row
    execute procedure set_updated_column();

create table if not exists users
(
    id           bigserial
    constraint users_pk
    primary key,
    phone_number bigint                              not null,
    created      timestamp default CURRENT_TIMESTAMP not null,
    updated      timestamp,
    deleted      timestamp
);

create trigger updated
    before update
    on users
    for each row
    execute procedure set_updated_column();

create table if not exists users_soc_nets
(
    id      bigserial
    constraint users_soc_nets_pk
    primary key,
    user_id bigint                              not null,
    created timestamp default CURRENT_TIMESTAMP not null,
    updated timestamp,
    deleted timestamp
);

create table if not exists tasks_activity
(
    task_id     bigint                                                        not null,
    executor_id bigint                                                        not null,
    status      timestamp default (CURRENT_TIMESTAMP + '1 day'::interval day) not null,
    created     timestamp default CURRENT_TIMESTAMP                           not null,
    updated     timestamp,
    deadline    timestamp default (CURRENT_TIMESTAMP + '1 day'::interval day)
    );

create trigger updated
    before update
    on tasks_activity
    for each row
    execute procedure set_updated_column();

create table if not exists users_executors
(
    user_id  bigint             not null,
    position point              not null,
    area     smallint default 1 not null,
    city     varchar(31)        not null
    );

create table if not exists tasks_appeales
(
    id      bigserial
    constraint tasks_appeales_pk
    primary key,
    user_id bigint       not null,
    task_id bigint       not null,
    text    varchar(255) not null
    );
