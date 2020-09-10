-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS public.currency
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE public.currency;
-- +goose StatementEnd
