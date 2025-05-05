CREATE SEQUENCE IF NOT EXISTS public.account_id_seq
   INCREMENT 1
   START 1
   MINVALUE 1
   MAXVALUE 9223372036854775807
   CACHE 1;

CREATE TABLE IF NOT EXISTS public.account
(
   id bigint NOT NULL DEFAULT nextval('account_id_seq'::regclass),
   username character varying(50) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
   password character varying(50) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
   enabled boolean NOT NULL DEFAULT false,
   CONSTRAINT account_pkey PRIMARY KEY (id),
   CONSTRAINT account_username_key UNIQUE (username)
);

ALTER SEQUENCE public.account_id_seq
   OWNED BY public.account.id;