CREATE TABLE IF NOT EXIST public.currency
(
    id character varying COLLATE pg_catalog."default" NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    rate numeric NOT NULL,
    insert_dt timestamp with time zone NOT NULL DEFAULT timezone('utc'::text, now()),
    CONSTRAINT currency_pkey PRIMARY KEY (id)
)
    TABLESPACE pg_default;

ALTER TABLE public.currency
    OWNER to igor;