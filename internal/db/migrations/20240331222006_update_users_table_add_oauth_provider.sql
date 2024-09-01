-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS thunderdome.auth_nonce (
    nonce_id character varying(64) NOT NULL PRIMARY KEY,
    created_date timestamp with time zone DEFAULT now(),
    expire_date timestamp with time zone DEFAULT (now() + '10 minutes'::interval)
);
CREATE TABLE thunderdome.auth_credential (
    user_id uuid UNIQUE NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    email character varying(320),
    password text,
    verified boolean DEFAULT false,
    mfa_enabled boolean DEFAULT false NOT NULL,
    created_date timestamp with time zone NOT NULL DEFAULT now(),
    updated_date timestamp with time zone NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX cred_email_unique_idx ON thunderdome.auth_credential USING btree (lower((email)::text));
CREATE TABLE thunderdome.auth_identity (
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    provider character varying(64) NOT NULL,
    sub TEXT NOT NULL,
    email character varying(320) NOT NULL,
    picture TEXT,
    verified boolean NOT NULL DEFAULT false,
    created_date timestamp with time zone NOT NULL DEFAULT now(),
    updated_date timestamp with time zone NOT NULL DEFAULT now(),
    UNIQUE(provider, sub)
);
DO $$
BEGIN
    INSERT INTO thunderdome.auth_credential (
        user_id, email, password, verified, mfa_enabled, created_date, updated_date
    )
    SELECT id, LOWER(email), password, verified, mfa_enabled, created_date, updated_date
        FROM thunderdome.users
        WHERE type <> 'GUEST';
END $$;
ALTER TABLE thunderdome.users ADD COLUMN picture TEXT;
DROP INDEX thunderdome.email_unique_idx;
CREATE OR REPLACE FUNCTION thunderdome.prune_auth_nonces() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
  row_count int;
BEGIN
  DELETE FROM thunderdome.auth_nonce WHERE expire_date < NOW();
  IF found THEN
    GET DIAGNOSTICS row_count = ROW_COUNT;
    RAISE NOTICE 'DELETED % row(s) FROM auth_nonce', row_count;
  END IF;
  RETURN NULL;
END;
$$;
CREATE TRIGGER prune_auth_nonces AFTER INSERT ON thunderdome.auth_nonce FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.prune_auth_nonces();
TRUNCATE TABLE thunderdome.user_verify;
CREATE OR REPLACE FUNCTION thunderdome.user_register(username character varying, useremail character varying, hashedpassword text, usertype character varying, OUT userid uuid, OUT verifyid uuid)
 RETURNS record
 LANGUAGE plpgsql
AS $function$
BEGIN
    INSERT INTO thunderdome.users (name, email, type)
    VALUES (username, useremail, usertype)
    RETURNING id INTO userid;

    INSERT INTO thunderdome.auth_credential (user_id, email, password)
    VALUES (userid, useremail, hashedpassword);

    INSERT INTO thunderdome.user_verify (user_id) VALUES (userid) RETURNING verify_id INTO verifyid;
END;
$function$;
CREATE OR REPLACE FUNCTION thunderdome.user_register_existing(activeuserid uuid, username character varying, useremail character varying, hashedpassword text, usertype character varying, OUT userid uuid, OUT verifyid uuid)
 RETURNS record
 LANGUAGE plpgsql
AS $function$
BEGIN
    UPDATE thunderdome.users
    SET
        name = userName,
        email = userEmail,
        type = userType,
        last_active = NOW(),
        updated_date = NOW()
    WHERE id = activeUserId
    RETURNING id INTO userId;

    INSERT INTO thunderdome.auth_credential (user_id, email, password) VALUES (userid, useremail, hashedpassword);

    INSERT INTO thunderdome.user_verify (user_id) VALUES (userId) RETURNING verify_id INTO verifyId;
END;
$function$;
CREATE OR REPLACE FUNCTION thunderdome.user_reset_create(useremail character varying, OUT resetid uuid, OUT userid uuid, OUT username character varying)
 RETURNS record
 LANGUAGE plpgsql
AS $function$
BEGIN
    SELECT u.id, u.name INTO userId, userName
    FROM thunderdome.auth_credential ac
    JOIN thunderdome.users u ON u.id = ac.user_id
    WHERE ac.email = userEmail;

    IF FOUND THEN
        INSERT INTO thunderdome.user_reset (user_id) VALUES (userId) RETURNING reset_id INTO resetId;
    ELSE
        RAISE EXCEPTION 'Nonexistent User --> %', userEmail USING HINT = 'Please check your Email';
    END IF;
END;
$function$;
CREATE OR REPLACE PROCEDURE thunderdome.user_password_reset(IN resetid uuid, IN userpassword text)
 LANGUAGE plpgsql
AS $procedure$
DECLARE matchedUserId UUID;
BEGIN
	matchedUserId := (
        SELECT w.id
        FROM thunderdome.user_reset wr
        LEFT JOIN thunderdome.users w ON w.id = wr.user_id
        WHERE wr.reset_id = resetId AND NOW() < wr.expire_date
    );

    IF matchedUserId IS NULL THEN
        -- attempt delete in case reset record expired
        DELETE FROM thunderdome.user_reset WHERE reset_id = resetId;
        RAISE 'Valid Reset ID not found';
    END IF;

    UPDATE thunderdome.auth_credential SET password = userPassword, updated_date = NOW() WHERE user_id = matchedUserId;
    UPDATE thunderdome.users SET last_active = NOW(), updated_date = NOW() WHERE id = matchedUserId;
    DELETE FROM thunderdome.user_reset WHERE reset_id = resetId;

    COMMIT;
END;
$procedure$;
CREATE OR REPLACE PROCEDURE thunderdome.user_mfa_remove(IN userid uuid)
 LANGUAGE plpgsql
AS $procedure$
BEGIN
    DELETE FROM thunderdome.user_mfa WHERE user_id = userId;
    UPDATE thunderdome.auth_credential SET mfa_enabled = false, updated_date = NOW() WHERE user_id = userId;

    COMMIT;
END;
$procedure$;
CREATE OR REPLACE PROCEDURE thunderdome.user_mfa_enable(IN userid uuid, IN mfasecret text)
 LANGUAGE plpgsql
AS $procedure$
BEGIN
    INSERT INTO thunderdome.user_mfa (user_id, secret) VALUES (userId, mfaSecret);
    UPDATE thunderdome.auth_credential SET mfa_enabled = true, updated_date = NOW() WHERE user_id = userId;

    COMMIT;
END;
$procedure$;
CREATE OR REPLACE PROCEDURE thunderdome.user_account_verify(IN verifyid uuid)
 LANGUAGE plpgsql
AS $procedure$
DECLARE matchedUserId UUID;
BEGIN
	matchedUserId := (
        SELECT w.id
        FROM thunderdome.user_verify wv
        JOIN thunderdome.users w ON w.id = wv.user_id
        WHERE wv.verify_id = verifyId AND NOW() < wv.expire_date
    );

    IF matchedUserId IS NULL THEN
        -- attempt delete in case verify record expired
        DELETE FROM thunderdome.user_verify WHERE verify_id = verifyId;
        RAISE 'Valid Verify ID not found';
    END IF;

    UPDATE thunderdome.users SET verified = 'TRUE', last_active = NOW(), updated_date = NOW() WHERE id = matchedUserId;
    UPDATE thunderdome.auth_credential SET verified = 'TRUE', updated_date = NOW() WHERE user_id = matchedUserId;
    DELETE FROM thunderdome.user_verify WHERE verify_id = verifyId;

    COMMIT;
END;
$procedure$;
ALTER TABLE thunderdome.users DROP COLUMN password;
ALTER TABLE thunderdome.users DROP COLUMN mfa_enabled;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION thunderdome.user_register(username character varying, useremail character varying, hashedpassword text, usertype character varying, OUT userid uuid, OUT verifyid uuid)
 RETURNS record
 LANGUAGE plpgsql
AS $function$
BEGIN
    INSERT INTO thunderdome.users (name, email, type, password)
    VALUES (username, useremail, usertype, hashedpassword)
    RETURNING id INTO userid;

    INSERT INTO thunderdome.user_verify (user_id) VALUES (userid) RETURNING verify_id INTO verifyid;
END;
$function$;
CREATE OR REPLACE FUNCTION thunderdome.user_register_existing(activeuserid uuid, username character varying, useremail character varying, hashedpassword text, usertype character varying, OUT userid uuid, OUT verifyid uuid)
 RETURNS record
 LANGUAGE plpgsql
AS $function$
BEGIN
    UPDATE thunderdome.users
    SET
        name = userName,
        email = userEmail,
        password = hashedPassword,
        type = userType,
        last_active = NOW(),
        updated_date = NOW()
    WHERE id = activeUserId
    RETURNING id INTO userId;

    INSERT INTO thunderdome.user_verify (user_id) VALUES (userId) RETURNING verify_id INTO verifyId;
END;
$function$;
CREATE OR REPLACE FUNCTION thunderdome.user_reset_create(useremail character varying, OUT resetid uuid, OUT userid uuid, OUT username character varying)
 RETURNS record
 LANGUAGE plpgsql
AS $function$
BEGIN
    SELECT id, name INTO userId, userName FROM thunderdome.users WHERE email = userEmail;
    IF FOUND THEN
        INSERT INTO thunderdome.user_reset (user_id) VALUES (userId) RETURNING reset_id INTO resetId;
    ELSE
        RAISE EXCEPTION 'Nonexistent User --> %', userEmail USING HINT = 'Please check your Email';
    END IF;
END;
$function$;
CREATE OR REPLACE PROCEDURE thunderdome.user_password_reset(IN resetid uuid, IN userpassword text)
 LANGUAGE plpgsql
AS $procedure$
DECLARE matchedUserId UUID;
BEGIN
	matchedUserId := (
        SELECT w.id
        FROM thunderdome.user_reset wr
        LEFT JOIN thunderdome.users w ON w.id = wr.user_id
        WHERE wr.reset_id = resetId AND NOW() < wr.expire_date
    );

    IF matchedUserId IS NULL THEN
        -- attempt delete in case reset record expired
        DELETE FROM thunderdome.user_reset WHERE reset_id = resetId;
        RAISE 'Valid Reset ID not found';
    END IF;

    UPDATE thunderdome.users SET password = userPassword, last_active = NOW(), updated_date = NOW()
        WHERE id = matchedUserId;
    DELETE FROM thunderdome.user_reset WHERE reset_id = resetId;

    COMMIT;
END;
$procedure$;
CREATE OR REPLACE PROCEDURE thunderdome.user_mfa_remove(IN userid uuid)
 LANGUAGE plpgsql
AS $procedure$
BEGIN
    DELETE FROM thunderdome.user_mfa WHERE user_id = userId;
    UPDATE thunderdome.users SET mfa_enabled = false, updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$procedure$;
CREATE OR REPLACE PROCEDURE thunderdome.user_mfa_enable(IN userid uuid, IN mfasecret text)
 LANGUAGE plpgsql
AS $procedure$
BEGIN
    INSERT INTO thunderdome.user_mfa (user_id, secret) VALUES (userId, mfaSecret);
    UPDATE thunderdome.users SET mfa_enabled = true, updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$procedure$;
CREATE OR REPLACE PROCEDURE thunderdome.user_account_verify(IN verifyid uuid)
 LANGUAGE plpgsql
AS $procedure$
DECLARE matchedUserId UUID;
BEGIN
	matchedUserId := (
        SELECT w.id
        FROM thunderdome.user_verify wv
        LEFT JOIN thunderdome.users w ON w.id = wv.user_id
        WHERE wv.verify_id = verifyId AND NOW() < wv.expire_date
    );

    IF matchedUserId IS NULL THEN
        -- attempt delete in case verify record expired
        DELETE FROM thunderdome.user_verify WHERE verify_id = verifyId;
        RAISE 'Valid Verify ID not found';
    END IF;

    UPDATE thunderdome.users SET verified = 'TRUE', last_active = NOW(), updated_date = NOW() WHERE id = matchedUserId;
    DELETE FROM thunderdome.user_verify WHERE verify_id = verifyId;

    COMMIT;
END;
$procedure$;
ALTER TABLE thunderdome.users ADD COLUMN password TEXT;
ALTER TABLE thunderdome.users ADD COLUMN mfa_enabled boolean NOT NULL DEFAULT false;
CREATE UNIQUE INDEX IF NOT EXISTS email_unique_idx ON thunderdome.users USING btree (lower((email)::text));
DROP TRIGGER prune_auth_nonces ON thunderdome.auth_nonce;
DROP FUNCTION thunderdome.prune_auth_nonces();
ALTER TABLE thunderdome.users DROP COLUMN picture;
DROP TABLE thunderdome.auth_nonce;
DROP TABLE thunderdome.auth_credential;
DROP TABLE thunderdome.auth_identity;
-- +goose StatementEnd
