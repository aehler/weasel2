-- run as srm_monitor
create schema service_logs;

create table service_logs.logs (
    log_id          bigserial not null primary key,
    service         varchar(100) not null default '',
    severity_level  integer not null default 0,
    info            varchar not null default '',
    details         jsonb not null default '{}',
    occured         timestamp not null default current_timestamp,
    created         timestamp not null default current_timestamp
);

create index idx_service_logs_service on service_logs.logs (service);