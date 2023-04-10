--
-- PostgreSQL database dump
--

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

--
-- TOC entry 211 (class 1259 OID 16404)
-- Name: products; Type: TABLE; Schema: public; Owner: postgres
--

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

--
-- TOC entry 3337 (class 0 OID 0)
-- Dependencies: 211
-- Name: COLUMN products.product_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.products.product_name IS 'Product name';

--
-- TOC entry 3338 (class 0 OID 0)
-- Dependencies: 211
-- Name: COLUMN products.product_amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.products.product_amount IS 'Product amount';

--
-- TOC entry 3339 (class 0 OID 0)
-- Dependencies: 211
-- Name: COLUMN products.product_user_created; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.products.product_user_created IS 'User create a product';

--
-- TOC entry 3340 (class 0 OID 0)
-- Dependencies: 211
-- Name: COLUMN products.product_date_created; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.products.product_date_created IS 'Date of product created';

--
-- TOC entry 3341 (class 0 OID 0)
-- Dependencies: 211
-- Name: COLUMN products.product_user_modify; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.products.product_user_modify IS 'Date of user modification';

--
-- TOC entry 3342 (class 0 OID 0)
-- Dependencies: 211
-- Name: COLUMN products.product_date_modify; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.products.product_date_modify IS 'Date of product modification';

--
-- TOC entry 213 (class 1259 OID 16429)
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.products_id_seq OWNER TO postgres;

--
-- TOC entry 3343 (class 0 OID 0)
-- Dependencies: 213
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.product_id;

--
-- TOC entry 209 (class 1259 OID 16385)
-- Name: types_identifiers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.types_identifiers (
                                          type_id integer NOT NULL,
                                          type_description character varying NOT NULL
);

ALTER TABLE public.types_identifiers OWNER TO postgres;

--
-- TOC entry 3344 (class 0 OID 0)
-- Dependencies: 209
-- Name: COLUMN types_identifiers.type_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.types_identifiers.type_id IS 'Identification type of character';

--
-- TOC entry 3345 (class 0 OID 0)
-- Dependencies: 209
-- Name: COLUMN types_identifiers.type_description; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.types_identifiers.type_description IS 'Description of the identification type';


--
-- TOC entry 210 (class 1259 OID 16392)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
                              user_name character varying(255) NOT NULL,
                              user_identifier integer NOT NULL,
                              user_email character varying(255) NOT NULL,
                              user_password character varying(255) NOT NULL,
                              user_type_identifier integer NOT NULL,
                              user_id integer NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 3346 (class 0 OID 0)
-- Dependencies: 210
-- Name: COLUMN users.user_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users.user_name IS 'Full name of the Character';

--
-- TOC entry 3347 (class 0 OID 0)
-- Dependencies: 210
-- Name: COLUMN users.user_identifier; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users.user_identifier IS 'User identification number';

--
-- TOC entry 3348 (class 0 OID 0)
-- Dependencies: 210
-- Name: COLUMN users.user_email; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users.user_email IS 'User email for login';

--
-- TOC entry 3349 (class 0 OID 0)
-- Dependencies: 210
-- Name: COLUMN users.user_password; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users.user_password IS 'User password for login';

--
-- TOC entry 3350 (class 0 OID 0)
-- Dependencies: 210
-- Name: COLUMN users.user_type_identifier; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users.user_type_identifier IS 'User identification document type';

--
-- TOC entry 212 (class 1259 OID 16419)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- TOC entry 3351 (class 0 OID 0)
-- Dependencies: 212
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.user_id;

--
-- TOC entry 3176 (class 2604 OID 16430)
-- Name: products product_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products ALTER COLUMN product_id SET DEFAULT nextval('public.products_id_seq'::regclass);

--
-- TOC entry 3175 (class 2604 OID 16420)
-- Name: users user_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN user_id SET DEFAULT nextval('public.users_id_seq'::regclass);

--
-- TOC entry 3329 (class 0 OID 16404)
-- Dependencies: 211
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.products (product_name, product_amount, product_user_created, product_date_created, product_user_modify, product_date_modify, product_id, product_price) FROM stdin;
BTC	10	1	2023-04-06 18:08:19	1	2023-04-06 18:08:19	5	28500

--
-- TOC entry 3327 (class 0 OID 16385)
-- Dependencies: 209
-- Data for Name: types_identifiers; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.types_identifiers (type_id, type_description) FROM stdin;
1	cc

--
-- TOC entry 3328 (class 0 OID 16392)
-- Dependencies: 210
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (user_name, user_identifier, user_email, user_password, user_type_identifier, user_id) FROM stdin;
aaa	1	aaa@aa.com	$2a$10$FZB/RkDK1rRjnbubgI58UeVGNE6ZuEsjRHwHNZMnLCa5fYd/1mzny	1	3
ivan	1	ivan-salazar@gmail.com	$2a$10$nPA4hNTB9jSx1Hr69eli5.Ln/LO1wB9hs7nlB8eAP/SmmedmI8bCi	1	4
ivanandres	1	ivanssalazar14@gmail.com	$2a$10$XENQXIRCQ.JYlqXOW8ZpaezVLj8R5Dp0xEPyZbb4Dv7e9MGhAc7im	1	5


--
-- TOC entry 3352 (class 0 OID 0)
-- Dependencies: 213
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.products_id_seq', 5, true);

--
-- TOC entry 3353 (class 0 OID 0)
-- Dependencies: 212
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 8, true);

--
-- TOC entry 3186 (class 2606 OID 16436)
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (product_id);

--
-- TOC entry 3178 (class 2606 OID 16391)
-- Name: types_identifiers types_identifiers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.types_identifiers
    ADD CONSTRAINT types_identifiers_pkey PRIMARY KEY (type_id);

--
-- TOC entry 3180 (class 2606 OID 16442)
-- Name: users unique_email; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT unique_email UNIQUE (user_email);


--
-- TOC entry 3182 (class 2606 OID 16438)
-- Name: users unique_username; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT unique_username UNIQUE (user_name);

--
-- TOC entry 3184 (class 2606 OID 16428)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);

--
-- TOC entry 3187 (class 2606 OID 16399)
-- Name: users fk_type_identifier; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_type_identifier FOREIGN KEY (user_type_identifier) REFERENCES public.types_identifiers(type_id);
-- Completed on 2023-04-06 13:16:46 -05
--
-- PostgreSQL database dump complete
--
