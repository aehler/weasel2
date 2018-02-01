create database srm_monitor;

create user srm_monitor with password 'IamAdiscoDancer' createdb;
grant all privileges on database srm_monitor to srm_monitor;

\c srm_monitor;

