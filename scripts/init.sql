BEGIN;

CREATE TABLE IF NOT EXISTS public.banner_feature_tag
(
    banner_id integer,
    feature_id integer NOT NULL,
    tag_id integer NOT NULL,
    CONSTRAINT feature_tag UNIQUE (feature_id, tag_id)
);

CREATE TABLE IF NOT EXISTS public.banners
(
    id serial NOT NULL,
    title text COLLATE pg_catalog."default" NOT NULL,
    text text COLLATE pg_catalog."default" NOT NULL,
    url text COLLATE pg_catalog."default" NOT NULL,
    visible boolean DEFAULT true,
    feature_id integer,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT banners_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.features
(
    id serial NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT features_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.tags
(
    id serial NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT tags_pkey PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.banner_feature_tag
    ADD CONSTRAINT banner_id_fkey FOREIGN KEY (banner_id)
        REFERENCES public.banners (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID;


ALTER TABLE IF EXISTS public.banner_feature_tag
    ADD CONSTRAINT banner_tag_tag_id_fkey FOREIGN KEY (tag_id)
        REFERENCES public.tags (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID;


ALTER TABLE IF EXISTS public.banner_feature_tag
    ADD CONSTRAINT feature_tag_feature_id_fkey FOREIGN KEY (feature_id)
        REFERENCES public.features (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID;


ALTER TABLE IF EXISTS public.banners
    ADD CONSTRAINT f_ref FOREIGN KEY (feature_id)
        REFERENCES public.features (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        DEFERRABLE
        NOT VALID;

create table if not exists public.users
(
    name     text not null,
    password text not null,
    role     text not null,
    id       serial
        constraint users_pk
            primary key
);

create or replace function update_timestamp_column()
    returns trigger as $$
begin
    new.updated_at = current_timestamp;
    return new;
end;
$$ language 'plpgsql';


create or replace trigger update_updated_at_column before update
    on banners for each row execute function update_timestamp_column();

END;