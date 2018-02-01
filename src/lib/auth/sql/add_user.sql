create or replace function weasel_auth.add_user(
    _organization_id bigint,
    _user_firstname  varchar,
    _user_lastname   varchar,
    _user_middlename varchar,
    _user_job_title  varchar,
    _user_phone  varchar,
    _user_email varchar,
    _user_password varchar,
    _is_admin   boolean,
    _timezone_id int,
    out _user_id bigint
) returns bigint as $$

    declare
        _token_id varchar;
    begin

        PERFORM 1 FROM weasel_auth.users WHERE user_login = lower(_user_email) and user_password = _user_password;

        IF FOUND THEN
            RAISE EXCEPTION 'USER_EXISTS';
        END IF;

        PERFORM 1 FROM weasel_auth.users WHERE lower(user_email) = lower(_user_email);

        IF FOUND THEN
            RAISE EXCEPTION 'EMAIL_EXISTS';
        END IF;

        INSERT INTO weasel_auth.users (
            organization_id,
            user_firstname,
            user_lastname,
            user_middlename,
            user_job_title,
            user_phone,
            user_login,
            user_email,
            user_password,
            is_admin,
            is_active,
            is_deleted,
            timezone_id
        ) VALUES (
            _organization_id,
            _user_firstname,
            _user_lastname,
            _user_middlename,
            _user_job_title,
            _user_phone,
            _user_email,
            _user_email,
            _user_password,
            _is_admin,
            true,
            false,
            _timezone_id
        ) RETURNING weasel_auth.users.user_id INTO _user_id;

    end;
$$ language plpgsql;