CREATE TABLE account (
    id SERIAL UNIQUE NOT NULL,
    created_at timestamp without time zone NOT NULL,
    name character varying(255) NOT NULL,
    stripe_id character varying(64),
    plan character varying(64),
    slug character varying(255) NOT NULL,
    owner_id integer NOT NULL
);

CREATE TABLE billing_email (
    id SERIAL UNIQUE NOT NULL,
    created_at timestamp without time zone NOT NULL,
    address character varying(255) NOT NULL,
    account_id integer NOT NULL
);

CREATE TABLE identity (
    id SERIAL UNIQUE NOT NULL,
    created_at timestamp without time zone NOT NULL,
    provider_id character varying(64) NOT NULL,
    provider character varying(64) NOT NULL,
    provider_data text,
    user_id integer NOT NULL
);

ALTER TABLE ONLY identity
    ADD CONSTRAINT identity_provider_id_provider_key UNIQUE (provider_id, provider);


CREATE TABLE invitation (
    id SERIAL UNIQUE NOT NULL,
    created_at timestamp without time zone NOT NULL,
    token character varying(64) NOT NULL,
    email character varying(255) NOT NULL,
    expiration timestamp without time zone NOT NULL,
    account_id integer NOT NULL,
    inviter_id integer NOT NULL,
    user_id integer
);

ALTER TABLE ONLY invitation
    ADD CONSTRAINT invitation_token_key UNIQUE (token);

CREATE TABLE mailing_list (
    id SERIAL UNIQUE NOT NULL,
    created_at timestamp without time zone NOT NULL,
    email character varying(255) NOT NULL,
    ipaddr character varying(255) NOT NULL
);

CREATE TABLE membership (
    id SERIAL UNIQUE NOT NULL,
    created_at timestamp without time zone NOT NULL,
    account_id integer NOT NULL,
    user_id integer NOT NULL
);

ALTER TABLE ONLY membership
    ADD CONSTRAINT membership_user_id_account_id_key UNIQUE (user_id, account_id);

CREATE TABLE password_reset (
    id SERIAL UNIQUE NOT NULL,
    created_at timestamp without time zone NOT NULL,
    token character varying(64) NOT NULL,
    expiration timestamp without time zone NOT NULL,
    user_id integer NOT NULL,
    redeemed timestamp without time zone
);

ALTER TABLE ONLY password_reset
    ADD CONSTRAINT password_reset_token_key UNIQUE (token);

CREATE TABLE subscription (
    id SERIAL UNIQUE NOT NULL,
    created_at timestamp without time zone NOT NULL,
    product_id character varying(255) NOT NULL,
    product_group_id character varying(255) NOT NULL,
    stripe_id character varying(255),
    account_id integer NOT NULL
);

ALTER TABLE ONLY subscription
    ADD CONSTRAINT subscription_account_id_product_group_id_key UNIQUE (account_id, product_group_id);

CREATE TABLE "user" (
    id SERIAL UNIQUE NOT NULL,
    created_at timestamp without time zone NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    login character varying(255),
    password character varying(64),
    default_membership_id integer
);

ALTER TABLE ONLY "user"
    ADD CONSTRAINT user_login_key UNIQUE (login);

ALTER TABLE ONLY account
    ADD CONSTRAINT account_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES "user"(id);

ALTER TABLE ONLY billing_email
    ADD CONSTRAINT billing_email_account_id_fkey FOREIGN KEY (account_id) REFERENCES account(id);

ALTER TABLE ONLY "user"
    ADD CONSTRAINT fk_user_membership_id FOREIGN KEY (default_membership_id) REFERENCES membership(id);

ALTER TABLE ONLY identity
    ADD CONSTRAINT identity_user_id_fkey FOREIGN KEY (user_id) REFERENCES "user"(id);

ALTER TABLE ONLY invitation
    ADD CONSTRAINT invitation_account_id_fkey FOREIGN KEY (account_id) REFERENCES account(id);

ALTER TABLE ONLY invitation
    ADD CONSTRAINT invitation_inviter_id_fkey FOREIGN KEY (inviter_id) REFERENCES "user"(id);

ALTER TABLE ONLY invitation
    ADD CONSTRAINT invitation_user_id_fkey FOREIGN KEY (user_id) REFERENCES "user"(id);

ALTER TABLE ONLY membership
    ADD CONSTRAINT membership_account_id_fkey FOREIGN KEY (account_id) REFERENCES account(id);

ALTER TABLE ONLY membership
    ADD CONSTRAINT membership_user_id_fkey FOREIGN KEY (user_id) REFERENCES "user"(id);

ALTER TABLE ONLY password_reset
    ADD CONSTRAINT password_reset_user_id_fkey FOREIGN KEY (user_id) REFERENCES "user"(id);

ALTER TABLE ONLY subscription
    ADD CONSTRAINT subscription_account_id_fkey FOREIGN KEY (account_id) REFERENCES account(id);
