-- Creates a table for storing identification types
-- Table: public.type_identifiers

-- DROP TABLE public.type_identifiers;


CREATE TABLE public.types_identifiers
(
    type_id integer NOT NULL,
    type_description character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT types_identifiers_pkey PRIMARY KEY (type_id)
)

TABLESPACE pg_default;

-- Sets the owner of the table to 'postgres'
ALTER TABLE public.types_identifiers
    OWNER to postgres;

-- Adds a comment for the 'type_id' column
COMMENT ON COLUMN public.types_identifiers.type_id
    IS 'Identification type of character';

-- Adds a comment for the 'type_description' column
COMMENT ON COLUMN public.types_identifiers.type_description
    IS 'Description of the identification type';

-- Creates a table for storing user information
-- Table: public.users

-- DROP TABLE public.users;


CREATE TABLE public.users
(
    user_id integer NOT NULL,
    user_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    user_identifier integer NOT NULL,
    user_email character varying(255) COLLATE pg_catalog."default" NOT NULL,
    user_password character varying(255) COLLATE pg_catalog."default" NOT NULL,
    user_type_identifier integer NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (user_id),
    CONSTRAINT fk_type_identifier FOREIGN KEY (user_type_identifier)
        REFERENCES public.types_identifiers (type_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

-- Sets the owner of the table to 'postgres'
ALTER TABLE public.users
    OWNER to postgres;

-- Adds a comment for the 'user_id' column
COMMENT ON COLUMN public.users.user_id
    IS 'Primary key of Users';

-- Adds a comment for the 'user_name' column
COMMENT ON COLUMN public.users.user_name
    IS 'Full name of the Character';

-- Adds a comment for the 'user_identifier' column
COMMENT ON COLUMN public.users.user_identifier
    IS 'User identification number';

-- Adds a comment for the 'user_email' column
COMMENT ON COLUMN public.users.user_email
    IS 'User email for login';

-- Adds a comment for the 'user_password' column
COMMENT ON COLUMN public.users.user_password
    IS 'User password for login';

-- Adds a comment for the 'user_type_identifier' column
COMMENT ON COLUMN public.users.user_type_identifier
    IS 'User identification document type';

-- Creates a table for storing product information
-- -- Table: public.products

-- DROP TABLE public.products;


CREATE TABLE public.products
(
    product_id integer NOT NULL,
    product_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    product_amount integer,
    product_user_created integer NOT NULL,
    product_date_created timestamp(0) without time zone NOT NULL,
    product_user_modify integer NOT NULL,
    "product_date_modify" timestamp(0) without time zone NOT NULL,
    CONSTRAINT product_pkey PRIMARY KEY (product_id),
    CONSTRAINT fk_user_created FOREIGN KEY (product_user_created)
        REFERENCES public.users (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT fk_user_modify FOREIGN KEY (product_user_modify)
        REFERENCES public.users (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)

TABLESPACE pg_default;

-- Sets the owner of the table to 'postgres'
ALTER TABLE public.products
    OWNER to postgres;

-- Comments for the columns in the 'products' table
ALTER TABLE public.products
    OWNER to postgres;

COMMENT ON COLUMN public.products.product_id
    IS 'Product ID';

COMMENT ON COLUMN public.products.product_name
    IS 'Product name';

COMMENT ON COLUMN public.products.product_amount
    IS 'Product amount';

COMMENT ON COLUMN public.products.product_user_created
    IS 'User create a product';

COMMENT ON COLUMN public.products.product_date_created
    IS 'Date of product created';

COMMENT ON COLUMN public.products.product_date_modify
    IS 'Date of product modification';

COMMENT ON COLUMN public.products."product_user_modify"
    IS 'Date of user modification';
    