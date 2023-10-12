CREATE TABLE public.articles (
    id integer NOT NULL,
    company integer NOT NULL,
    name text NOT NULL,
    sort_order integer DEFAULT 0 NOT NULL,
    object_type smallint NOT NULL,
    image text,
    "group" integer,
    notes text DEFAULT ''::text NOT NULL
);