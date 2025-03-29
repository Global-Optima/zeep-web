SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: admin_role; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.admin_role AS ENUM (
    'ADMIN',
    'OWNER'
);


--
-- Name: franchisee_employee_role; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.franchisee_employee_role AS ENUM (
    'FRANCHISEE_MANAGER',
    'FRANCHISEE_OWNER'
);


--
-- Name: http_method; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.http_method AS ENUM (
    'GET',
    'POST',
    'PUT',
    'PATCH',
    'DELETE'
);


--
-- Name: operation_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.operation_type AS ENUM (
    'GET',
    'CREATE',
    'UPDATE',
    'DELETE'
);


--
-- Name: region_employee_role; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.region_employee_role AS ENUM (
    'REGION_WAREHOUSE_MANAGER'
);


--
-- Name: store_employee_role; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.store_employee_role AS ENUM (
    'STORE_MANAGER',
    'BARISTA'
);


--
-- Name: valid_phone; Type: DOMAIN; Schema: public; Owner: -
--

CREATE DOMAIN public.valid_phone AS character varying(16)
	CONSTRAINT valid_phone_check CHECK (((VALUE)::text ~ '^\+[1-9]\d{1,14}$'::text));


--
-- Name: warehouse_employee_role; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.warehouse_employee_role AS ENUM (
    'WAREHOUSE_MANAGER',
    'WAREHOUSE_EMPLOYEE'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: additive_categories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.additive_categories (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    is_multiple_select boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: additive_categories_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.additive_categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: additive_categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.additive_categories_id_seq OWNED BY public.additive_categories.id;


--
-- Name: additive_ingredients; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.additive_ingredients (
    id integer NOT NULL,
    ingredient_id integer NOT NULL,
    additive_id integer NOT NULL,
    quantity numeric(10,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT additive_ingredients_quantity_check CHECK ((quantity > (0)::numeric))
);


--
-- Name: additive_ingredients_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.additive_ingredients_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: additive_ingredients_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.additive_ingredients_id_seq OWNED BY public.additive_ingredients.id;


--
-- Name: additives; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.additives (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    base_price numeric(10,2) NOT NULL,
    size numeric(10,2) NOT NULL,
    unit_id integer NOT NULL,
    additive_category_id integer NOT NULL,
    image_key character varying(2048),
    machine_id character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT additives_base_price_check CHECK ((base_price > (0)::numeric)),
    CONSTRAINT additives_size_check CHECK ((size > (0)::numeric))
);


--
-- Name: additives_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.additives_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: additives_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.additives_id_seq OWNED BY public.additives.id;


--
-- Name: admin_employees; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.admin_employees (
    id integer NOT NULL,
    employee_id integer NOT NULL,
    role public.admin_role NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: admin_employees_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.admin_employees_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: admin_employees_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.admin_employees_id_seq OWNED BY public.admin_employees.id;


--
-- Name: bonuses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bonuses (
    id integer NOT NULL,
    bonuses numeric(10,2),
    customer_id integer NOT NULL,
    expires_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT bonuses_bonuses_check CHECK ((bonuses >= (0)::numeric))
);


--
-- Name: bonuses_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bonuses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bonuses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bonuses_id_seq OWNED BY public.bonuses.id;


--
-- Name: customer_addresses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.customer_addresses (
    id integer NOT NULL,
    customer_id integer NOT NULL,
    address character varying(255) NOT NULL,
    longitude character varying(20),
    latitude character varying(20),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: customer_addresses_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.customer_addresses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: customer_addresses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.customer_addresses_id_seq OWNED BY public.customer_addresses.id;


--
-- Name: customers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.customers (
    id integer NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    phone public.valid_phone,
    is_verified boolean DEFAULT false,
    is_banned boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: customers_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.customers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: customers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.customers_id_seq OWNED BY public.customers.id;


--
-- Name: employee_audits; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.employee_audits (
    id integer NOT NULL,
    employee_id integer NOT NULL,
    operation_type public.operation_type NOT NULL,
    component_name character varying(255) NOT NULL,
    details jsonb,
    ip_address character varying(45) NOT NULL,
    resource_url text NOT NULL,
    method public.http_method NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: employee_audits_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.employee_audits_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: employee_audits_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.employee_audits_id_seq OWNED BY public.employee_audits.id;


--
-- Name: employee_notification_recipients; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.employee_notification_recipients (
    id integer NOT NULL,
    notification_id integer NOT NULL,
    employee_id integer NOT NULL,
    is_read boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone
);


--
-- Name: employee_notification_recipients_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.employee_notification_recipients_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: employee_notification_recipients_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.employee_notification_recipients_id_seq OWNED BY public.employee_notification_recipients.id;


--
-- Name: employee_notifications; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.employee_notifications (
    id integer NOT NULL,
    event_type character varying(255) NOT NULL,
    priority character varying(50) NOT NULL,
    details jsonb DEFAULT '{}'::jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone
);


--
-- Name: employee_notifications_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.employee_notifications_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: employee_notifications_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.employee_notifications_id_seq OWNED BY public.employee_notifications.id;


--
-- Name: employee_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.employee_tokens (
    id integer NOT NULL,
    token character varying(255) NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    employee_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: employee_tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.employee_tokens_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: employee_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.employee_tokens_id_seq OWNED BY public.employee_tokens.id;


--
-- Name: employee_work_tracks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.employee_work_tracks (
    id integer NOT NULL,
    start_work_at timestamp with time zone,
    end_work_at timestamp with time zone,
    employee_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: employee_work_tracks_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.employee_work_tracks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: employee_work_tracks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.employee_work_tracks_id_seq OWNED BY public.employee_work_tracks.id;


--
-- Name: employee_workdays; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.employee_workdays (
    id integer NOT NULL,
    day character varying(15) NOT NULL,
    start_at time without time zone NOT NULL,
    end_at time without time zone NOT NULL,
    employee_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: employee_workdays_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.employee_workdays_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: employee_workdays_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.employee_workdays_id_seq OWNED BY public.employee_workdays.id;


--
-- Name: employees; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.employees (
    id integer NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    phone public.valid_phone,
    email character varying(255),
    hashed_password character varying(255) NOT NULL,
    is_active boolean NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: employees_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.employees_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: employees_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.employees_id_seq OWNED BY public.employees.id;


--
-- Name: facility_addresses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.facility_addresses (
    id integer NOT NULL,
    address character varying(255) NOT NULL,
    longitude numeric(9,6),
    latitude numeric(9,6),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: facility_addresses_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.facility_addresses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: facility_addresses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.facility_addresses_id_seq OWNED BY public.facility_addresses.id;


--
-- Name: franchisee_employees; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.franchisee_employees (
    id integer NOT NULL,
    franchisee_id integer NOT NULL,
    employee_id integer NOT NULL,
    role public.franchisee_employee_role NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: franchisee_employees_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.franchisee_employees_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: franchisee_employees_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.franchisee_employees_id_seq OWNED BY public.franchisee_employees.id;


--
-- Name: franchisees; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.franchisees (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    description character varying(1024),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: franchisees_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.franchisees_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: franchisees_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.franchisees_id_seq OWNED BY public.franchisees.id;


--
-- Name: ingredient_categories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ingredient_categories (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: ingredient_categories_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ingredient_categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ingredient_categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ingredient_categories_id_seq OWNED BY public.ingredient_categories.id;


--
-- Name: ingredients; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ingredients (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    calories numeric(5,2),
    fat numeric(5,2),
    carbs numeric(5,2),
    proteins numeric(5,2),
    expiration_in_days integer,
    unit_id integer NOT NULL,
    category_id integer NOT NULL,
    is_allergen boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT ingredients_calories_check CHECK ((calories >= (0)::numeric)),
    CONSTRAINT ingredients_carbs_check CHECK ((carbs >= (0)::numeric)),
    CONSTRAINT ingredients_expiration_in_days_check CHECK ((expiration_in_days >= 0)),
    CONSTRAINT ingredients_fat_check CHECK ((fat >= (0)::numeric)),
    CONSTRAINT ingredients_proteins_check CHECK ((proteins >= (0)::numeric))
);


--
-- Name: ingredients_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ingredients_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ingredients_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ingredients_id_seq OWNED BY public.ingredients.id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.orders (
    id integer NOT NULL,
    customer_id integer,
    display_number integer NOT NULL,
    customer_name character varying(255) NOT NULL,
    store_employee_id integer,
    store_id integer NOT NULL,
    delivery_address_id integer,
    status character varying(50) NOT NULL,
    total numeric(10,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    completed_at timestamp with time zone,
    deleted_at timestamp with time zone,
    CONSTRAINT orders_total_check CHECK ((total >= (0)::numeric))
);


--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- Name: product_categories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_categories (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: product_categories_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.product_categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: product_categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.product_categories_id_seq OWNED BY public.product_categories.id;


--
-- Name: product_size_additives; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_size_additives (
    id integer NOT NULL,
    product_size_id integer NOT NULL,
    additive_id integer NOT NULL,
    is_default boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: product_size_additives_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.product_size_additives_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: product_size_additives_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.product_size_additives_id_seq OWNED BY public.product_size_additives.id;


--
-- Name: product_size_ingredients; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_size_ingredients (
    id integer NOT NULL,
    ingredient_id integer NOT NULL,
    product_size_id integer NOT NULL,
    quantity numeric(10,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT product_size_ingredients_quantity_check CHECK ((quantity > (0)::numeric))
);


--
-- Name: product_size_ingredients_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.product_size_ingredients_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: product_size_ingredients_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.product_size_ingredients_id_seq OWNED BY public.product_size_ingredients.id;


--
-- Name: product_sizes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_sizes (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    unit_id integer NOT NULL,
    base_price numeric(10,2) NOT NULL,
    size numeric(10,2) NOT NULL,
    product_id integer NOT NULL,
    discount_id integer,
    machine_id character varying(40) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT product_sizes_base_price_check CHECK ((base_price > (0)::numeric)),
    CONSTRAINT product_sizes_size_check CHECK ((size > (0)::numeric))
);


--
-- Name: product_sizes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.product_sizes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: product_sizes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.product_sizes_id_seq OWNED BY public.product_sizes.id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.products (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    image_key character varying(2048),
    video_key character varying(2048),
    category_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- Name: recipe_steps; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.recipe_steps (
    id integer NOT NULL,
    product_id integer NOT NULL,
    step integer NOT NULL,
    name character varying(100),
    description text,
    image_key character varying(2048),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: recipe_steps_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.recipe_steps_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: recipe_steps_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.recipe_steps_id_seq OWNED BY public.recipe_steps.id;


--
-- Name: referrals; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.referrals (
    id integer NOT NULL,
    customer_id integer NOT NULL,
    referee_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: referrals_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.referrals_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: referrals_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.referrals_id_seq OWNED BY public.referrals.id;


--
-- Name: region_employees; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.region_employees (
    id integer NOT NULL,
    employee_id integer NOT NULL,
    region_id integer NOT NULL,
    role public.region_employee_role NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: region_employees_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.region_employees_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: region_employees_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.region_employees_id_seq OWNED BY public.region_employees.id;


--
-- Name: regions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.regions (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: regions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.regions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: regions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.regions_id_seq OWNED BY public.regions.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(128) NOT NULL
);


--
-- Name: stock_material_categories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.stock_material_categories (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: stock_material_categories_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.stock_material_categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: stock_material_categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.stock_material_categories_id_seq OWNED BY public.stock_material_categories.id;


--
-- Name: stock_materials; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.stock_materials (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    ingredient_id integer NOT NULL,
    safety_stock numeric(10,2) NOT NULL,
    unit_id integer NOT NULL,
    size numeric(10,2) NOT NULL,
    category_id integer NOT NULL,
    barcode character varying(255),
    expiration_period_in_days integer DEFAULT 1095 NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT stock_materials_safety_stock_check CHECK ((safety_stock >= (0)::numeric))
);


--
-- Name: stock_materials_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.stock_materials_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: stock_materials_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.stock_materials_id_seq OWNED BY public.stock_materials.id;


--
-- Name: stock_request_ingredients; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.stock_request_ingredients (
    id integer NOT NULL,
    stock_request_id integer NOT NULL,
    stock_material_id integer NOT NULL,
    quantity numeric(10,2) NOT NULL,
    delivered_date timestamp with time zone,
    expiration_date timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT stock_request_ingredients_quantity_check CHECK ((quantity > (0)::numeric))
);


--
-- Name: stock_request_ingredients_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.stock_request_ingredients_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: stock_request_ingredients_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.stock_request_ingredients_id_seq OWNED BY public.stock_request_ingredients.id;


--
-- Name: stock_requests; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.stock_requests (
    id integer NOT NULL,
    store_id integer NOT NULL,
    warehouse_id integer NOT NULL,
    status character varying(50) NOT NULL,
    details jsonb,
    store_comment text,
    warehouse_comment text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: stock_requests_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.stock_requests_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: stock_requests_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.stock_requests_id_seq OWNED BY public.stock_requests.id;


--
-- Name: store_additives; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.store_additives (
    id integer NOT NULL,
    additive_id integer NOT NULL,
    store_id integer NOT NULL,
    store_price numeric(10,2),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT store_additives_store_price_check CHECK ((store_price > (0)::numeric))
);


--
-- Name: store_additives_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.store_additives_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: store_additives_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.store_additives_id_seq OWNED BY public.store_additives.id;


--
-- Name: store_employees; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.store_employees (
    id integer NOT NULL,
    employee_id integer NOT NULL,
    store_id integer NOT NULL,
    role public.store_employee_role NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: store_employees_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.store_employees_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: store_employees_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.store_employees_id_seq OWNED BY public.store_employees.id;


--
-- Name: store_product_sizes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.store_product_sizes (
    id integer NOT NULL,
    product_size_id integer NOT NULL,
    store_product_id integer NOT NULL,
    store_price numeric(10,2),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT store_product_sizes_store_price_check CHECK ((store_price > (0)::numeric))
);


--
-- Name: store_product_sizes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.store_product_sizes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: store_product_sizes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.store_product_sizes_id_seq OWNED BY public.store_product_sizes.id;


--
-- Name: store_products; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.store_products (
    id integer NOT NULL,
    product_id integer NOT NULL,
    store_id integer NOT NULL,
    is_available boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: store_products_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.store_products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: store_products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.store_products_id_seq OWNED BY public.store_products.id;


--
-- Name: store_stocks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.store_stocks (
    id integer NOT NULL,
    store_id integer NOT NULL,
    ingredient_id integer NOT NULL,
    low_stock_threshold numeric(10,2) NOT NULL,
    quantity numeric(10,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT store_stocks_low_stock_threshold_check CHECK ((low_stock_threshold > (0)::numeric)),
    CONSTRAINT store_stocks_quantity_check CHECK ((quantity >= (0)::numeric))
);


--
-- Name: store_stocks_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.store_stocks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: store_stocks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.store_stocks_id_seq OWNED BY public.store_stocks.id;


--
-- Name: stores; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.stores (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    facility_address_id integer,
    franchisee_id integer,
    warehouse_id integer NOT NULL,
    is_active boolean DEFAULT true,
    contact_phone public.valid_phone,
    contact_email character varying(255),
    store_hours character varying(255),
    last_inventory_sync_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: stores_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.stores_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: stores_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.stores_id_seq OWNED BY public.stores.id;


--
-- Name: suborder_additives; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.suborder_additives (
    id integer NOT NULL,
    suborder_id integer NOT NULL,
    store_additive_id integer NOT NULL,
    price numeric(10,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT suborder_additives_price_check CHECK ((price > (0)::numeric))
);


--
-- Name: suborder_additives_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.suborder_additives_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: suborder_additives_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.suborder_additives_id_seq OWNED BY public.suborder_additives.id;


--
-- Name: suborders; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.suborders (
    id integer NOT NULL,
    order_id integer NOT NULL,
    store_product_size_id integer NOT NULL,
    price numeric(10,2) NOT NULL,
    status character varying(50) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    completed_at timestamp with time zone,
    deleted_at timestamp with time zone,
    CONSTRAINT suborders_price_check CHECK ((price > (0)::numeric))
);


--
-- Name: suborders_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.suborders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: suborders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.suborders_id_seq OWNED BY public.suborders.id;


--
-- Name: supplier_materials; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.supplier_materials (
    id integer NOT NULL,
    stock_material_id integer NOT NULL,
    supplier_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: supplier_materials_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.supplier_materials_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: supplier_materials_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.supplier_materials_id_seq OWNED BY public.supplier_materials.id;


--
-- Name: supplier_prices; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.supplier_prices (
    id integer NOT NULL,
    supplier_material_id integer NOT NULL,
    base_price numeric(10,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT supplier_prices_base_price_check CHECK ((base_price > (0)::numeric))
);


--
-- Name: supplier_prices_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.supplier_prices_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: supplier_prices_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.supplier_prices_id_seq OWNED BY public.supplier_prices.id;


--
-- Name: supplier_warehouse_deliveries; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.supplier_warehouse_deliveries (
    id integer NOT NULL,
    supplier_id integer NOT NULL,
    warehouse_id integer NOT NULL,
    delivery_date timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: supplier_warehouse_deliveries_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.supplier_warehouse_deliveries_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: supplier_warehouse_deliveries_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.supplier_warehouse_deliveries_id_seq OWNED BY public.supplier_warehouse_deliveries.id;


--
-- Name: supplier_warehouse_delivery_materials; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.supplier_warehouse_delivery_materials (
    id integer NOT NULL,
    delivery_id integer NOT NULL,
    stock_material_id integer NOT NULL,
    barcode character varying(255) NOT NULL,
    quantity numeric(10,2) NOT NULL,
    expiration_date timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT supplier_warehouse_delivery_materials_quantity_check CHECK ((quantity > (0)::numeric))
);


--
-- Name: supplier_warehouse_delivery_materials_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.supplier_warehouse_delivery_materials_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: supplier_warehouse_delivery_materials_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.supplier_warehouse_delivery_materials_id_seq OWNED BY public.supplier_warehouse_delivery_materials.id;


--
-- Name: suppliers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.suppliers (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    contact_email character varying(255),
    contact_phone public.valid_phone,
    city character varying(100) NOT NULL,
    address character varying(255),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: suppliers_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.suppliers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: suppliers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.suppliers_id_seq OWNED BY public.suppliers.id;


--
-- Name: transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.transactions (
    id integer NOT NULL,
    type character varying(50) NOT NULL,
    order_id integer NOT NULL,
    bin character varying(20) NOT NULL,
    transaction_id character varying(50) NOT NULL,
    process_id character varying(50),
    payment_method character varying(50) NOT NULL,
    amount numeric(10,2) NOT NULL,
    currency character(3) NOT NULL,
    qr_number character varying(50),
    card_mask character varying(16),
    icc character varying(255),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.transactions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.transactions_id_seq OWNED BY public.transactions.id;


--
-- Name: units; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.units (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    conversion_factor numeric(10,4) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: units_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.units_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: units_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.units_id_seq OWNED BY public.units.id;


--
-- Name: verification_codes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.verification_codes (
    id integer NOT NULL,
    customer_id integer NOT NULL,
    code character varying(6) NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: verification_codes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.verification_codes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: verification_codes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.verification_codes_id_seq OWNED BY public.verification_codes.id;


--
-- Name: warehouse_employees; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.warehouse_employees (
    id integer NOT NULL,
    employee_id integer NOT NULL,
    warehouse_id integer NOT NULL,
    role public.warehouse_employee_role NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: warehouse_employees_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.warehouse_employees_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: warehouse_employees_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.warehouse_employees_id_seq OWNED BY public.warehouse_employees.id;


--
-- Name: warehouse_stocks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.warehouse_stocks (
    id integer NOT NULL,
    warehouse_id integer NOT NULL,
    stock_material_id integer NOT NULL,
    quantity numeric(10,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT warehouse_stocks_quantity_check CHECK ((quantity >= (0)::numeric))
);


--
-- Name: warehouse_stocks_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.warehouse_stocks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: warehouse_stocks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.warehouse_stocks_id_seq OWNED BY public.warehouse_stocks.id;


--
-- Name: warehouses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.warehouses (
    id integer NOT NULL,
    facility_address_id integer NOT NULL,
    region_id integer NOT NULL,
    name character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Name: warehouses_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.warehouses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: warehouses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.warehouses_id_seq OWNED BY public.warehouses.id;


--
-- Name: additive_categories id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.additive_categories ALTER COLUMN id SET DEFAULT nextval('public.additive_categories_id_seq'::regclass);


--
-- Name: additive_ingredients id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.additive_ingredients ALTER COLUMN id SET DEFAULT nextval('public.additive_ingredients_id_seq'::regclass);


--
-- Name: additives id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.additives ALTER COLUMN id SET DEFAULT nextval('public.additives_id_seq'::regclass);


--
-- Name: admin_employees id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.admin_employees ALTER COLUMN id SET DEFAULT nextval('public.admin_employees_id_seq'::regclass);


--
-- Name: bonuses id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bonuses ALTER COLUMN id SET DEFAULT nextval('public.bonuses_id_seq'::regclass);


--
-- Name: customer_addresses id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer_addresses ALTER COLUMN id SET DEFAULT nextval('public.customer_addresses_id_seq'::regclass);


--
-- Name: customers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customers ALTER COLUMN id SET DEFAULT nextval('public.customers_id_seq'::regclass);


--
-- Name: employee_audits id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_audits ALTER COLUMN id SET DEFAULT nextval('public.employee_audits_id_seq'::regclass);


--
-- Name: employee_notification_recipients id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_notification_recipients ALTER COLUMN id SET DEFAULT nextval('public.employee_notification_recipients_id_seq'::regclass);


--
-- Name: employee_notifications id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_notifications ALTER COLUMN id SET DEFAULT nextval('public.employee_notifications_id_seq'::regclass);


--
-- Name: employee_tokens id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_tokens ALTER COLUMN id SET DEFAULT nextval('public.employee_tokens_id_seq'::regclass);


--
-- Name: employee_work_tracks id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_work_tracks ALTER COLUMN id SET DEFAULT nextval('public.employee_work_tracks_id_seq'::regclass);


--
-- Name: employee_workdays id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_workdays ALTER COLUMN id SET DEFAULT nextval('public.employee_workdays_id_seq'::regclass);


--
-- Name: employees id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employees ALTER COLUMN id SET DEFAULT nextval('public.employees_id_seq'::regclass);


--
-- Name: facility_addresses id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.facility_addresses ALTER COLUMN id SET DEFAULT nextval('public.facility_addresses_id_seq'::regclass);


--
-- Name: franchisee_employees id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.franchisee_employees ALTER COLUMN id SET DEFAULT nextval('public.franchisee_employees_id_seq'::regclass);


--
-- Name: franchisees id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.franchisees ALTER COLUMN id SET DEFAULT nextval('public.franchisees_id_seq'::regclass);


--
-- Name: ingredient_categories id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ingredient_categories ALTER COLUMN id SET DEFAULT nextval('public.ingredient_categories_id_seq'::regclass);


--
-- Name: ingredients id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ingredients ALTER COLUMN id SET DEFAULT nextval('public.ingredients_id_seq'::regclass);


--
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- Name: product_categories id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_categories ALTER COLUMN id SET DEFAULT nextval('public.product_categories_id_seq'::regclass);


--
-- Name: product_size_additives id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_size_additives ALTER COLUMN id SET DEFAULT nextval('public.product_size_additives_id_seq'::regclass);


--
-- Name: product_size_ingredients id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_size_ingredients ALTER COLUMN id SET DEFAULT nextval('public.product_size_ingredients_id_seq'::regclass);


--
-- Name: product_sizes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_sizes ALTER COLUMN id SET DEFAULT nextval('public.product_sizes_id_seq'::regclass);


--
-- Name: products id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- Name: recipe_steps id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.recipe_steps ALTER COLUMN id SET DEFAULT nextval('public.recipe_steps_id_seq'::regclass);


--
-- Name: referrals id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.referrals ALTER COLUMN id SET DEFAULT nextval('public.referrals_id_seq'::regclass);


--
-- Name: region_employees id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.region_employees ALTER COLUMN id SET DEFAULT nextval('public.region_employees_id_seq'::regclass);


--
-- Name: regions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.regions ALTER COLUMN id SET DEFAULT nextval('public.regions_id_seq'::regclass);


--
-- Name: stock_material_categories id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock_material_categories ALTER COLUMN id SET DEFAULT nextval('public.stock_material_categories_id_seq'::regclass);


--
-- Name: stock_materials id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock_materials ALTER COLUMN id SET DEFAULT nextval('public.stock_materials_id_seq'::regclass);


--
-- Name: stock_request_ingredients id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock_request_ingredients ALTER COLUMN id SET DEFAULT nextval('public.stock_request_ingredients_id_seq'::regclass);


--
-- Name: stock_requests id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock_requests ALTER COLUMN id SET DEFAULT nextval('public.stock_requests_id_seq'::regclass);


--
-- Name: store_additives id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_additives ALTER COLUMN id SET DEFAULT nextval('public.store_additives_id_seq'::regclass);


--
-- Name: store_employees id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_employees ALTER COLUMN id SET DEFAULT nextval('public.store_employees_id_seq'::regclass);


--
-- Name: store_product_sizes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_product_sizes ALTER COLUMN id SET DEFAULT nextval('public.store_product_sizes_id_seq'::regclass);


--
-- Name: store_products id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_products ALTER COLUMN id SET DEFAULT nextval('public.store_products_id_seq'::regclass);


--
-- Name: store_stocks id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_stocks ALTER COLUMN id SET DEFAULT nextval('public.store_stocks_id_seq'::regclass);


--
-- Name: stores id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stores ALTER COLUMN id SET DEFAULT nextval('public.stores_id_seq'::regclass);


--
-- Name: suborder_additives id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suborder_additives ALTER COLUMN id SET DEFAULT nextval('public.suborder_additives_id_seq'::regclass);


--
-- Name: suborders id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suborders ALTER COLUMN id SET DEFAULT nextval('public.suborders_id_seq'::regclass);


--
-- Name: supplier_materials id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_materials ALTER COLUMN id SET DEFAULT nextval('public.supplier_materials_id_seq'::regclass);


--
-- Name: supplier_prices id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_prices ALTER COLUMN id SET DEFAULT nextval('public.supplier_prices_id_seq'::regclass);


--
-- Name: supplier_warehouse_deliveries id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_warehouse_deliveries ALTER COLUMN id SET DEFAULT nextval('public.supplier_warehouse_deliveries_id_seq'::regclass);


--
-- Name: supplier_warehouse_delivery_materials id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_warehouse_delivery_materials ALTER COLUMN id SET DEFAULT nextval('public.supplier_warehouse_delivery_materials_id_seq'::regclass);


--
-- Name: suppliers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers ALTER COLUMN id SET DEFAULT nextval('public.suppliers_id_seq'::regclass);


--
-- Name: transactions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transactions ALTER COLUMN id SET DEFAULT nextval('public.transactions_id_seq'::regclass);


--
-- Name: units id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.units ALTER COLUMN id SET DEFAULT nextval('public.units_id_seq'::regclass);


--
-- Name: verification_codes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.verification_codes ALTER COLUMN id SET DEFAULT nextval('public.verification_codes_id_seq'::regclass);


--
-- Name: warehouse_employees id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_employees ALTER COLUMN id SET DEFAULT nextval('public.warehouse_employees_id_seq'::regclass);


--
-- Name: warehouse_stocks id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_stocks ALTER COLUMN id SET DEFAULT nextval('public.warehouse_stocks_id_seq'::regclass);


--
-- Name: warehouses id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouses ALTER COLUMN id SET DEFAULT nextval('public.warehouses_id_seq'::regclass);


--
-- Name: additive_categories additive_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.additive_categories
    ADD CONSTRAINT additive_categories_pkey PRIMARY KEY (id);


--
-- Name: additive_ingredients additive_ingredients_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.additive_ingredients
    ADD CONSTRAINT additive_ingredients_pkey PRIMARY KEY (id);


--
-- Name: additives additives_machine_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.additives
    ADD CONSTRAINT additives_machine_id_key UNIQUE (machine_id);


--
-- Name: additives additives_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.additives
    ADD CONSTRAINT additives_pkey PRIMARY KEY (id);


--
-- Name: admin_employees admin_employees_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.admin_employees
    ADD CONSTRAINT admin_employees_pkey PRIMARY KEY (id);


--
-- Name: bonuses bonuses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bonuses
    ADD CONSTRAINT bonuses_pkey PRIMARY KEY (id);


--
-- Name: customer_addresses customer_addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer_addresses
    ADD CONSTRAINT customer_addresses_pkey PRIMARY KEY (id);


--
-- Name: customers customers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_pkey PRIMARY KEY (id);


--
-- Name: employee_audits employee_audits_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_audits
    ADD CONSTRAINT employee_audits_pkey PRIMARY KEY (id);


--
-- Name: employee_notification_recipients employee_notification_recipients_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_notification_recipients
    ADD CONSTRAINT employee_notification_recipients_pkey PRIMARY KEY (id);


--
-- Name: employee_notifications employee_notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_notifications
    ADD CONSTRAINT employee_notifications_pkey PRIMARY KEY (id);


--
-- Name: employee_tokens employee_tokens_employee_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_tokens
    ADD CONSTRAINT employee_tokens_employee_id_key UNIQUE (employee_id);


--
-- Name: employee_tokens employee_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_tokens
    ADD CONSTRAINT employee_tokens_pkey PRIMARY KEY (id);


--
-- Name: employee_tokens employee_tokens_token_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_tokens
    ADD CONSTRAINT employee_tokens_token_key UNIQUE (token);


--
-- Name: employee_work_tracks employee_work_tracks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_work_tracks
    ADD CONSTRAINT employee_work_tracks_pkey PRIMARY KEY (id);


--
-- Name: employee_workdays employee_workdays_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_workdays
    ADD CONSTRAINT employee_workdays_pkey PRIMARY KEY (id);


--
-- Name: employees employees_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT employees_pkey PRIMARY KEY (id);


--
-- Name: facility_addresses facility_addresses_address_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.facility_addresses
    ADD CONSTRAINT facility_addresses_address_key UNIQUE (address);


--
-- Name: facility_addresses facility_addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.facility_addresses
    ADD CONSTRAINT facility_addresses_pkey PRIMARY KEY (id);


--
-- Name: franchisee_employees franchisee_employees_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.franchisee_employees
    ADD CONSTRAINT franchisee_employees_pkey PRIMARY KEY (id);


--
-- Name: franchisees franchisees_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.franchisees
    ADD CONSTRAINT franchisees_pkey PRIMARY KEY (id);


--
-- Name: ingredient_categories ingredient_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ingredient_categories
    ADD CONSTRAINT ingredient_categories_pkey PRIMARY KEY (id);


--
-- Name: ingredients ingredients_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ingredients
    ADD CONSTRAINT ingredients_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: product_categories product_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_categories
    ADD CONSTRAINT product_categories_pkey PRIMARY KEY (id);


--
-- Name: product_size_additives product_size_additives_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_size_additives
    ADD CONSTRAINT product_size_additives_pkey PRIMARY KEY (id);


--
-- Name: product_size_ingredients product_size_ingredients_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_size_ingredients
    ADD CONSTRAINT product_size_ingredients_pkey PRIMARY KEY (id);


--
-- Name: product_sizes product_sizes_machine_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_sizes
    ADD CONSTRAINT product_sizes_machine_id_key UNIQUE (machine_id);


--
-- Name: product_sizes product_sizes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_sizes
    ADD CONSTRAINT product_sizes_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: recipe_steps recipe_steps_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.recipe_steps
    ADD CONSTRAINT recipe_steps_pkey PRIMARY KEY (id);


--
-- Name: referrals referrals_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.referrals
    ADD CONSTRAINT referrals_pkey PRIMARY KEY (id);


--
-- Name: region_employees region_employees_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.region_employees
    ADD CONSTRAINT region_employees_pkey PRIMARY KEY (id);


--
-- Name: regions regions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.regions
    ADD CONSTRAINT regions_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: stock_material_categories stock_material_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock_material_categories
    ADD CONSTRAINT stock_material_categories_pkey PRIMARY KEY (id);


--
-- Name: stock_materials stock_materials_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock_materials
    ADD CONSTRAINT stock_materials_pkey PRIMARY KEY (id);


--
-- Name: stock_request_ingredients stock_request_ingredients_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock_request_ingredients
    ADD CONSTRAINT stock_request_ingredients_pkey PRIMARY KEY (id);


--
-- Name: stock_requests stock_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock_requests
    ADD CONSTRAINT stock_requests_pkey PRIMARY KEY (id);


--
-- Name: store_additives store_additives_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_additives
    ADD CONSTRAINT store_additives_pkey PRIMARY KEY (id);


--
-- Name: store_employees store_employees_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_employees
    ADD CONSTRAINT store_employees_pkey PRIMARY KEY (id);


--
-- Name: store_product_sizes store_product_sizes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_product_sizes
    ADD CONSTRAINT store_product_sizes_pkey PRIMARY KEY (id);


--
-- Name: store_products store_products_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_products
    ADD CONSTRAINT store_products_pkey PRIMARY KEY (id);


--
-- Name: store_stocks store_stocks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_stocks
    ADD CONSTRAINT store_stocks_pkey PRIMARY KEY (id);


--
-- Name: stores stores_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stores
    ADD CONSTRAINT stores_pkey PRIMARY KEY (id);


--
-- Name: suborder_additives suborder_additives_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suborder_additives
    ADD CONSTRAINT suborder_additives_pkey PRIMARY KEY (id);


--
-- Name: suborders suborders_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suborders
    ADD CONSTRAINT suborders_pkey PRIMARY KEY (id);


--
-- Name: supplier_materials supplier_materials_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_materials
    ADD CONSTRAINT supplier_materials_pkey PRIMARY KEY (id);


--
-- Name: supplier_prices supplier_prices_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_prices
    ADD CONSTRAINT supplier_prices_pkey PRIMARY KEY (id);


--
-- Name: supplier_warehouse_deliveries supplier_warehouse_deliveries_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_warehouse_deliveries
    ADD CONSTRAINT supplier_warehouse_deliveries_pkey PRIMARY KEY (id);


--
-- Name: supplier_warehouse_delivery_materials supplier_warehouse_delivery_materials_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_warehouse_delivery_materials
    ADD CONSTRAINT supplier_warehouse_delivery_materials_pkey PRIMARY KEY (id);


--
-- Name: suppliers suppliers_contact_phone_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_contact_phone_key UNIQUE (contact_phone);


--
-- Name: suppliers suppliers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_pkey PRIMARY KEY (id);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: transactions transactions_process_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_process_id_key UNIQUE (process_id);


--
-- Name: transactions transactions_transaction_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_transaction_id_key UNIQUE (transaction_id);


--
-- Name: units units_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.units
    ADD CONSTRAINT units_pkey PRIMARY KEY (id);


--
-- Name: verification_codes verification_codes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.verification_codes
    ADD CONSTRAINT verification_codes_pkey PRIMARY KEY (id);


--
-- Name: warehouse_employees warehouse_employees_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_employees
    ADD CONSTRAINT warehouse_employees_pkey PRIMARY KEY (id);


--
-- Name: warehouse_stocks warehouse_stocks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_stocks
    ADD CONSTRAINT warehouse_stocks_pkey PRIMARY KEY (id);


--
-- Name: warehouses warehouses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouses_pkey PRIMARY KEY (id);


--
-- Name: idx_employee_audits_employee_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_employee_audits_employee_id ON public.employee_audits USING btree (employee_id);


--
-- Name: idx_employee_audits_timestamp; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_employee_audits_timestamp ON public.employee_audits USING btree (created_at);


--
-- Name: idx_employee_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_employee_id ON public.employee_notification_recipients USING btree (employee_id);


--
-- Name: idx_notification_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_notification_id ON public.employee_notification_recipients USING btree (notification_id);


--
-- Name: idx_orders_store_display; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_orders_store_display ON public.orders USING btree (store_id, display_number);


--
-- Name: unique_additive_category_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_additive_category_name ON public.additive_categories USING btree (name) WHERE (deleted_at IS NULL);


--
-- Name: unique_additive_ingredient; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_additive_ingredient ON public.additive_ingredients USING btree (ingredient_id, additive_id) WHERE (deleted_at IS NULL);


--
-- Name: unique_additive_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_additive_name ON public.additives USING btree (name) WHERE (deleted_at IS NULL);


--
-- Name: unique_admin_employee; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_admin_employee ON public.admin_employees USING btree (employee_id, deleted_at) WHERE (deleted_at IS NULL);


--
-- Name: unique_customer_address; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_customer_address ON public.customer_addresses USING btree (customer_id, lower((address)::text)) WHERE (deleted_at IS NULL);


--
-- Name: unique_customer_phone; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_customer_phone ON public.customers USING btree (phone) WHERE (deleted_at IS NULL);


--
-- Name: unique_employee_email; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_employee_email ON public.employees USING btree (email) WHERE (deleted_at IS NULL);


--
-- Name: unique_employee_phone; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_employee_phone ON public.employees USING btree (phone) WHERE (deleted_at IS NULL);


--
-- Name: unique_employee_workday; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_employee_workday ON public.employee_workdays USING btree (employee_id, day) WHERE (deleted_at IS NULL);


--
-- Name: unique_facility_coordinates; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_facility_coordinates ON public.facility_addresses USING btree (longitude, latitude) WHERE ((deleted_at IS NULL) AND (longitude IS NOT NULL) AND (latitude IS NOT NULL));


--
-- Name: unique_franchisee_employee; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_franchisee_employee ON public.franchisee_employees USING btree (employee_id, deleted_at) WHERE (deleted_at IS NULL);


--
-- Name: unique_franchisee_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_franchisee_name ON public.franchisees USING btree (name) WHERE (deleted_at IS NULL);


--
-- Name: unique_ingredient_category_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_ingredient_category_name ON public.ingredient_categories USING btree (name) WHERE (deleted_at IS NULL);


--
-- Name: unique_ingredient_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_ingredient_name ON public.ingredients USING btree (name) WHERE (deleted_at IS NULL);


--
-- Name: unique_product_category_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_product_category_name ON public.product_categories USING btree (name) WHERE (deleted_at IS NULL);


--
-- Name: unique_product_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_product_name ON public.products USING btree (name) WHERE (deleted_at IS NULL);


--
-- Name: unique_product_size_additive; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_product_size_additive ON public.product_size_additives USING btree (product_size_id, additive_id) WHERE (deleted_at IS NULL);


--
-- Name: unique_product_size_ingredient; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_product_size_ingredient ON public.product_size_ingredients USING btree (product_size_id, ingredient_id) WHERE (deleted_at IS NULL);


--
-- Name: unique_product_size_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_product_size_name ON public.product_sizes USING btree (product_id, name) WHERE (deleted_at IS NULL);


--
-- Name: unique_recipe_step_number; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_recipe_step_number ON public.recipe_steps USING btree (product_id, step) WHERE (deleted_at IS NULL);


--
-- Name: unique_referrals; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_referrals ON public.referrals USING btree (customer_id, referee_id) WHERE (deleted_at IS NULL);


--
-- Name: unique_region_employee; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_region_employee ON public.region_employees USING btree (employee_id, deleted_at) WHERE (deleted_at IS NULL);


--
-- Name: unique_stock_material_barcode; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_stock_material_barcode ON public.stock_materials USING btree (barcode) WHERE (deleted_at IS NULL);


--
-- Name: unique_stock_material_category_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_stock_material_category_name ON public.stock_material_categories USING btree (name) WHERE (deleted_at IS NULL);


--
-- Name: unique_stock_request_ingredient; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_stock_request_ingredient ON public.stock_request_ingredients USING btree (stock_request_id, stock_material_id) WHERE (deleted_at IS NULL);


--
-- Name: unique_store_additive; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_store_additive ON public.store_additives USING btree (store_id, additive_id) WHERE (deleted_at IS NULL);


--
-- Name: unique_store_contact_email; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_store_contact_email ON public.stores USING btree (contact_email) WHERE (deleted_at IS NULL);


--
-- Name: unique_store_contact_phone; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_store_contact_phone ON public.stores USING btree (contact_phone) WHERE (deleted_at IS NULL);


--
-- Name: unique_store_employee; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_store_employee ON public.store_employees USING btree (employee_id, deleted_at) WHERE (deleted_at IS NULL);


--
-- Name: unique_store_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_store_name ON public.stores USING btree (name) WHERE (deleted_at IS NULL);


--
-- Name: unique_store_product; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_store_product ON public.store_products USING btree (store_id, product_id) WHERE (deleted_at IS NULL);


--
-- Name: unique_store_product_size; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_store_product_size ON public.store_product_sizes USING btree (store_product_id, product_size_id) WHERE (deleted_at IS NULL);


--
-- Name: unique_store_stock; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_store_stock ON public.store_stocks USING btree (store_id, ingredient_id) WHERE (deleted_at IS NULL);


--
-- Name: unique_supplier_material; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_supplier_material ON public.supplier_materials USING btree (supplier_id, stock_material_id) WHERE (deleted_at IS NULL);


--
-- Name: unique_unit_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_unit_name ON public.units USING btree (name) WHERE (deleted_at IS NULL);


--
-- Name: unique_warehouse_employee; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_warehouse_employee ON public.warehouse_employees USING btree (employee_id, deleted_at) WHERE (deleted_at IS NULL);


--
-- Name: unique_warehouse_stock; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_warehouse_stock ON public.warehouse_stocks USING btree (warehouse_id, stock_material_id) WHERE (deleted_at IS NULL);


--
-- Name: additive_ingredients additive_ingredients_additive_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.additive_ingredients
    ADD CONSTRAINT additive_ingredients_additive_id_fkey FOREIGN KEY (additive_id) REFERENCES public.additives(id) ON DELETE CASCADE;


--
-- Name: additive_ingredients additive_ingredients_ingredient_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.additive_ingredients
    ADD CONSTRAINT additive_ingredients_ingredient_id_fkey FOREIGN KEY (ingredient_id) REFERENCES public.ingredients(id) ON DELETE CASCADE;


--
-- Name: additives additives_additive_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.additives
    ADD CONSTRAINT additives_additive_category_id_fkey FOREIGN KEY (additive_category_id) REFERENCES public.additive_categories(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: additives additives_unit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.additives
    ADD CONSTRAINT additives_unit_id_fkey FOREIGN KEY (unit_id) REFERENCES public.units(id) ON DELETE RESTRICT;


--
-- Name: admin_employees admin_employees_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.admin_employees
    ADD CONSTRAINT admin_employees_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.employees(id) ON DELETE CASCADE;


--
-- Name: bonuses bonuses_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bonuses
    ADD CONSTRAINT bonuses_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON DELETE CASCADE;


--
-- Name: customer_addresses customer_addresses_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer_addresses
    ADD CONSTRAINT customer_addresses_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON DELETE CASCADE;


--
-- Name: employee_audits employee_audits_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_audits
    ADD CONSTRAINT employee_audits_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.employees(id) ON DELETE CASCADE;


--
-- Name: employee_notification_recipients employee_notification_recipients_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_notification_recipients
    ADD CONSTRAINT employee_notification_recipients_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.employees(id) ON DELETE CASCADE;


--
-- Name: employee_notification_recipients employee_notification_recipients_notification_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_notification_recipients
    ADD CONSTRAINT employee_notification_recipients_notification_id_fkey FOREIGN KEY (notification_id) REFERENCES public.employee_notifications(id) ON DELETE CASCADE;


--
-- Name: employee_tokens employee_tokens_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_tokens
    ADD CONSTRAINT employee_tokens_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.employees(id) ON DELETE CASCADE;


--
-- Name: employee_work_tracks employee_work_tracks_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_work_tracks
    ADD CONSTRAINT employee_work_tracks_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.employees(id) ON DELETE CASCADE;


--
-- Name: employee_workdays employee_workdays_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_workdays
    ADD CONSTRAINT employee_workdays_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.employees(id) ON DELETE CASCADE;


--
-- Name: franchisee_employees franchisee_employees_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.franchisee_employees
    ADD CONSTRAINT franchisee_employees_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.employees(id) ON DELETE CASCADE;


--
-- Name: franchisee_employees franchisee_employees_franchisee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.franchisee_employees
    ADD CONSTRAINT franchisee_employees_franchisee_id_fkey FOREIGN KEY (franchisee_id) REFERENCES public.franchisees(id) ON DELETE CASCADE;


--
-- Name: ingredients ingredients_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ingredients
    ADD CONSTRAINT ingredients_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.ingredient_categories(id) ON DELETE RESTRICT;


--
-- Name: ingredients ingredients_unit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ingredients
    ADD CONSTRAINT ingredients_unit_id_fkey FOREIGN KEY (unit_id) REFERENCES public.units(id) ON DELETE RESTRICT;


--
-- Name: orders orders_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON DELETE SET NULL;


--
-- Name: orders orders_delivery_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_delivery_address_id_fkey FOREIGN KEY (delivery_address_id) REFERENCES public.customer_addresses(id) ON DELETE SET NULL;


--
-- Name: orders orders_store_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_store_employee_id_fkey FOREIGN KEY (store_employee_id) REFERENCES public.store_employees(id) ON DELETE SET NULL;


--
-- Name: orders orders_store_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_store_id_fkey FOREIGN KEY (store_id) REFERENCES public.stores(id);


--
-- Name: product_size_additives product_size_additives_additive_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_size_additives
    ADD CONSTRAINT product_size_additives_additive_id_fkey FOREIGN KEY (additive_id) REFERENCES public.additives(id) ON DELETE CASCADE;


--
-- Name: product_size_additives product_size_additives_product_size_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_size_additives
    ADD CONSTRAINT product_size_additives_product_size_id_fkey FOREIGN KEY (product_size_id) REFERENCES public.product_sizes(id) ON DELETE CASCADE;


--
-- Name: product_size_ingredients product_size_ingredients_ingredient_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_size_ingredients
    ADD CONSTRAINT product_size_ingredients_ingredient_id_fkey FOREIGN KEY (ingredient_id) REFERENCES public.ingredients(id) ON DELETE CASCADE;


--
-- Name: product_size_ingredients product_size_ingredients_product_size_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_size_ingredients
    ADD CONSTRAINT product_size_ingredients_product_size_id_fkey FOREIGN KEY (product_size_id) REFERENCES public.product_sizes(id) ON DELETE CASCADE;


--
-- Name: product_sizes product_sizes_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_sizes
    ADD CONSTRAINT product_sizes_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: product_sizes product_sizes_unit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_sizes
    ADD CONSTRAINT product_sizes_unit_id_fkey FOREIGN KEY (unit_id) REFERENCES public.units(id) ON DELETE RESTRICT;


--
-- Name: products products_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.product_categories(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: recipe_steps recipe_steps_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.recipe_steps
    ADD CONSTRAINT recipe_steps_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: referrals referrals_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.referrals
    ADD CONSTRAINT referrals_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON DELETE CASCADE;


--
-- Name: referrals referrals_referee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.referrals
    ADD CONSTRAINT referrals_referee_id_fkey FOREIGN KEY (referee_id) REFERENCES public.customers(id) ON DELETE CASCADE;


--
-- Name: region_employees region_employees_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.region_employees
    ADD CONSTRAINT region_employees_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.employees(id) ON DELETE CASCADE;


--
-- Name: region_employees region_employees_region_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.region_employees
    ADD CONSTRAINT region_employees_region_id_fkey FOREIGN KEY (region_id) REFERENCES public.regions(id) ON DELETE CASCADE;


--
-- Name: stock_materials stock_materials_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock_materials
    ADD CONSTRAINT stock_materials_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.stock_material_categories(id) ON DELETE RESTRICT;


--
-- Name: stock_materials stock_materials_ingredient_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock_materials
    ADD CONSTRAINT stock_materials_ingredient_id_fkey FOREIGN KEY (ingredient_id) REFERENCES public.ingredients(id) ON DELETE CASCADE;


--
-- Name: stock_materials stock_materials_unit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock_materials
    ADD CONSTRAINT stock_materials_unit_id_fkey FOREIGN KEY (unit_id) REFERENCES public.units(id) ON DELETE RESTRICT;


--
-- Name: stock_request_ingredients stock_request_ingredients_stock_material_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock_request_ingredients
    ADD CONSTRAINT stock_request_ingredients_stock_material_id_fkey FOREIGN KEY (stock_material_id) REFERENCES public.stock_materials(id) ON DELETE CASCADE;


--
-- Name: stock_request_ingredients stock_request_ingredients_stock_request_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock_request_ingredients
    ADD CONSTRAINT stock_request_ingredients_stock_request_id_fkey FOREIGN KEY (stock_request_id) REFERENCES public.stock_requests(id) ON DELETE CASCADE;


--
-- Name: stock_requests stock_requests_store_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock_requests
    ADD CONSTRAINT stock_requests_store_id_fkey FOREIGN KEY (store_id) REFERENCES public.stores(id) ON DELETE CASCADE;


--
-- Name: stock_requests stock_requests_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock_requests
    ADD CONSTRAINT stock_requests_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id) ON DELETE CASCADE;


--
-- Name: store_additives store_additives_additive_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_additives
    ADD CONSTRAINT store_additives_additive_id_fkey FOREIGN KEY (additive_id) REFERENCES public.additives(id) ON DELETE CASCADE;


--
-- Name: store_additives store_additives_store_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_additives
    ADD CONSTRAINT store_additives_store_id_fkey FOREIGN KEY (store_id) REFERENCES public.stores(id) ON DELETE CASCADE;


--
-- Name: store_employees store_employees_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_employees
    ADD CONSTRAINT store_employees_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.employees(id) ON DELETE CASCADE;


--
-- Name: store_employees store_employees_store_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_employees
    ADD CONSTRAINT store_employees_store_id_fkey FOREIGN KEY (store_id) REFERENCES public.stores(id) ON DELETE CASCADE;


--
-- Name: store_product_sizes store_product_sizes_product_size_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_product_sizes
    ADD CONSTRAINT store_product_sizes_product_size_id_fkey FOREIGN KEY (product_size_id) REFERENCES public.product_sizes(id) ON DELETE CASCADE;


--
-- Name: store_product_sizes store_product_sizes_store_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_product_sizes
    ADD CONSTRAINT store_product_sizes_store_product_id_fkey FOREIGN KEY (store_product_id) REFERENCES public.store_products(id);


--
-- Name: store_products store_products_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_products
    ADD CONSTRAINT store_products_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: store_products store_products_store_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_products
    ADD CONSTRAINT store_products_store_id_fkey FOREIGN KEY (store_id) REFERENCES public.stores(id) ON DELETE CASCADE;


--
-- Name: store_stocks store_stocks_ingredient_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_stocks
    ADD CONSTRAINT store_stocks_ingredient_id_fkey FOREIGN KEY (ingredient_id) REFERENCES public.ingredients(id) ON DELETE CASCADE;


--
-- Name: store_stocks store_stocks_store_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.store_stocks
    ADD CONSTRAINT store_stocks_store_id_fkey FOREIGN KEY (store_id) REFERENCES public.stores(id) ON DELETE CASCADE;


--
-- Name: stores stores_facility_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stores
    ADD CONSTRAINT stores_facility_address_id_fkey FOREIGN KEY (facility_address_id) REFERENCES public.facility_addresses(id);


--
-- Name: stores stores_franchisee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stores
    ADD CONSTRAINT stores_franchisee_id_fkey FOREIGN KEY (franchisee_id) REFERENCES public.franchisees(id) ON DELETE CASCADE;


--
-- Name: stores stores_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stores
    ADD CONSTRAINT stores_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id) ON DELETE CASCADE;


--
-- Name: suborder_additives suborder_additives_store_additive_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suborder_additives
    ADD CONSTRAINT suborder_additives_store_additive_id_fkey FOREIGN KEY (store_additive_id) REFERENCES public.store_additives(id) ON DELETE CASCADE;


--
-- Name: suborder_additives suborder_additives_suborder_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suborder_additives
    ADD CONSTRAINT suborder_additives_suborder_id_fkey FOREIGN KEY (suborder_id) REFERENCES public.suborders(id) ON DELETE CASCADE;


--
-- Name: suborders suborders_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suborders
    ADD CONSTRAINT suborders_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- Name: suborders suborders_store_product_size_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suborders
    ADD CONSTRAINT suborders_store_product_size_id_fkey FOREIGN KEY (store_product_size_id) REFERENCES public.store_product_sizes(id) ON DELETE RESTRICT;


--
-- Name: supplier_materials supplier_materials_stock_material_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_materials
    ADD CONSTRAINT supplier_materials_stock_material_id_fkey FOREIGN KEY (stock_material_id) REFERENCES public.stock_materials(id) ON DELETE CASCADE;


--
-- Name: supplier_materials supplier_materials_supplier_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_materials
    ADD CONSTRAINT supplier_materials_supplier_id_fkey FOREIGN KEY (supplier_id) REFERENCES public.suppliers(id) ON DELETE CASCADE;


--
-- Name: supplier_prices supplier_prices_supplier_material_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_prices
    ADD CONSTRAINT supplier_prices_supplier_material_id_fkey FOREIGN KEY (supplier_material_id) REFERENCES public.supplier_materials(id) ON DELETE CASCADE;


--
-- Name: supplier_warehouse_deliveries supplier_warehouse_deliveries_supplier_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_warehouse_deliveries
    ADD CONSTRAINT supplier_warehouse_deliveries_supplier_id_fkey FOREIGN KEY (supplier_id) REFERENCES public.suppliers(id) ON DELETE CASCADE;


--
-- Name: supplier_warehouse_deliveries supplier_warehouse_deliveries_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_warehouse_deliveries
    ADD CONSTRAINT supplier_warehouse_deliveries_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id) ON DELETE CASCADE;


--
-- Name: supplier_warehouse_delivery_materials supplier_warehouse_delivery_materials_delivery_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_warehouse_delivery_materials
    ADD CONSTRAINT supplier_warehouse_delivery_materials_delivery_id_fkey FOREIGN KEY (delivery_id) REFERENCES public.supplier_warehouse_deliveries(id) ON DELETE CASCADE;


--
-- Name: supplier_warehouse_delivery_materials supplier_warehouse_delivery_materials_stock_material_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_warehouse_delivery_materials
    ADD CONSTRAINT supplier_warehouse_delivery_materials_stock_material_id_fkey FOREIGN KEY (stock_material_id) REFERENCES public.stock_materials(id) ON DELETE CASCADE;


--
-- Name: transactions transactions_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- Name: verification_codes verification_codes_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.verification_codes
    ADD CONSTRAINT verification_codes_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON DELETE CASCADE;


--
-- Name: warehouse_employees warehouse_employees_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_employees
    ADD CONSTRAINT warehouse_employees_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.employees(id) ON DELETE CASCADE;


--
-- Name: warehouse_employees warehouse_employees_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_employees
    ADD CONSTRAINT warehouse_employees_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id) ON DELETE CASCADE;


--
-- Name: warehouse_stocks warehouse_stocks_stock_material_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_stocks
    ADD CONSTRAINT warehouse_stocks_stock_material_id_fkey FOREIGN KEY (stock_material_id) REFERENCES public.stock_materials(id) ON DELETE CASCADE;


--
-- Name: warehouse_stocks warehouse_stocks_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_stocks
    ADD CONSTRAINT warehouse_stocks_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id) ON DELETE CASCADE;


--
-- Name: warehouses warehouses_facility_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouses_facility_address_id_fkey FOREIGN KEY (facility_address_id) REFERENCES public.facility_addresses(id) ON DELETE RESTRICT;


--
-- Name: warehouses warehouses_region_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouses_region_id_fkey FOREIGN KEY (region_id) REFERENCES public.regions(id) ON DELETE RESTRICT;


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20250329033159');
