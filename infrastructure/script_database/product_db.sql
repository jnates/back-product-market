-- PostgreSQL database dump

-- Dumped from database version 14.1
-- Dumped by pg_dump version 14.7 (Ubuntu 14.7-0ubuntu0.22.04.1)

-- Started on 2023-04-06 13:16:46 -05

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;
SET default_tablespace = '';
SET default_table_access_method = heap;

-- Table Definitions

CREATE TABLE public.products (
                                 product_name character varying(255) NOT NULL,
                                 product_amount integer,
                                 product_user_created integer NOT NULL,
                                 product_date_created timestamp(0) without time zone NOT NULL,
                                 product_user_modify integer NOT NULL,
                                 product_date_modify timestamp(0) without time zone NOT NULL,
                                 product_id integer NOT NULL,
                                 product_price double precision
);

ALTER TABLE public.products OWNER TO postgres;

COMMENT ON COLUMN public.products.product_name IS 'Product name';
COMMENT ON COLUMN public.products.product_amount IS 'Product amount';
COMMENT ON COLUMN public.products.product_user_created IS 'User create a product';
COMMENT ON COLUMN public.products.product_date_created IS 'Date of product created';
COMMENT ON COLUMN public.products.product_user_modify IS 'Date of user modification';
COMMENT ON COLUMN public.products.product_date_modify IS 'Date of product modification';

CREATE SEQUENCE public.products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.product_id;

CREATE TABLE public.types_identifiers (
                                          type_id integer NOT NULL,
                                          type_description character varying NOT NULL
);

ALTER TABLE public.types_identifiers OWNER TO postgres;

COMMENT ON COLUMN public.types_identifiers.type_id IS 'Identification type of character';
COMMENT ON COLUMN public.types_identifiers.type_description IS 'Description of the identification type';

CREATE TABLE public.users (
                              user_name character varying(255) NOT NULL,
                              user_identifier integer NOT NULL,
                              user_email character varying(255) NOT NULL,
                              user_password character varying(255) NOT NULL,
                              user_type_identifier integer NOT NULL,
                              user_id integer NOT NULL
);

ALTER TABLE public.users OWNER TO postgres;

COMMENT ON COLUMN public.users.user_name IS 'Full name of the Character';
COMMENT ON COLUMN public.users.user_identifier IS 'User identification number';
COMMENT ON COLUMN public.users.user_email IS 'User email for login';
COMMENT ON COLUMN public.users.user_password IS 'User password for login';
COMMENT ON COLUMN public.users.user_type_identifier IS 'User identification document type';

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.user_id;

-- Set Defaults

ALTER TABLE ONLY public.products ALTER COLUMN product_id SET DEFAULT nextval('public.products_id_seq'::regclass);
ALTER TABLE ONLY public.users ALTER COLUMN user_id SET DEFAULT nextval('public.users_id_seq'::regclass);

-- Data Inserts

INSERT INTO public.products (product_name, product_amount, product_user_created, product_date_created, product_user_modify, product_date_modify, product_id, product_price)
VALUES
    ('BTC', 10, 1, '2023-04-06 18:08:19', 1, '2023-04-06 18:08:19', 5, 28500);

INSERT INTO public.types_identifiers (type_id, type_description)
VALUES
    (1, 'cc');

INSERT INTO public.users (user_name, user_identifier, user_email, user_password, user_type_identifier, user_id)
VALUES
    ('aaa', 1, 'aaa@aa.com', '$2a$10$FZB/RkDK1rRjnbubgI58UeVGNE6ZuEsjRHwHNZMnLCa5fYd/1mzny', 1, 3),
    ('juan', 1, 'nates1999123@gmail.com', '$2a$10$nPA4hNTB9jSx1Hr69eli5.Ln/LO1wB9hs7nlB8eAP/SmmedmI8bCi', 1, 4),
    ('jnates', 1, 'natesjd@gmail.com', '$2a$10$XENQXIRCQ.JYlqXOW8ZpaezVLj8R5Dp0xEPyZbb4Dv7e9MGhAc7im', 1, 5);

-- Set Sequences

SELECT pg_catalog.setval('public.products_id_seq', 5, true);
SELECT pg_catalog.setval('public.users_id_seq', 8, true);

-- Constraints

ALTER TABLE ONLY public.products ADD CONSTRAINT products_pkey PRIMARY KEY (product_id);
ALTER TABLE ONLY public.types_identifiers ADD CONSTRAINT types_identifiers_pkey PRIMARY KEY (type_id);
ALTER TABLE ONLY public.users ADD CONSTRAINT unique_email UNIQUE (user_email);
ALTER TABLE ONLY public.users ADD CONSTRAINT unique_username UNIQUE (user_name);
ALTER TABLE ONLY public.users ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);
ALTER TABLE ONLY public.users ADD CONSTRAINT fk_type_identifier FOREIGN KEY (user_type_identifier) REFERENCES public.types_identifiers(type_id);

-- Completed on 2023-04-06 13:16:46 -05
-- PostgreSQL database dump complete


-- Data Inserts

INSERT INTO public.products (product_name, product_amount, product_user_created, product_date_created, product_user_modify, product_date_modify, product_id, product_price)
VALUES
    ('BTC', 10, 1, '2023-04-06 18:08:19', 1, '2023-04-06 18:08:19', 5, 28500);

INSERT INTO public.types_identifiers (type_id, type_description)
VALUES
    (1, 'cc');

INSERT INTO public.users (user_name, user_identifier, user_email, user_password, user_type_identifier, user_id)
VALUES
    ('aaa', 1, 'aaa@aa.com', '$2a$10$FZB/RkDK1rRjnbubgI58UeVGNE6ZuEsjRHwHNZMnLCa5fYd/1mzny', 1, 3),
    ('juan', 1, 'nates1999123@gmail.com', '$2a$10$nPA4hNTB9jSx1Hr69eli5.Ln/LO1wB9hs7nlB8eAP/SmmedmI8bCi', 1, 4),
    ('jnates', 1, 'natesjd@gmail.com', '$2a$10$XENQXIRCQ.JYlqXOW8ZpaezVLj8R5Dp0xEPyZbb4Dv7e9MGhAc7im', 1, 5);