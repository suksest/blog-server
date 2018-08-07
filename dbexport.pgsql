--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.13
-- Dumped by pg_dump version 9.5.13

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: books; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.books (
    isbn character(14) NOT NULL,
    title character varying(255) NOT NULL,
    author character varying(255) NOT NULL,
    price numeric(5,2) NOT NULL
);


ALTER TABLE public.books OWNER TO postgres;

--
-- Name: posts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.posts (
    id integer NOT NULL,
    title character varying(50) NOT NULL,
    body text NOT NULL,
    published_at timestamp without time zone NOT NULL,
    author_id integer NOT NULL
);


ALTER TABLE public.posts OWNER TO postgres;

--
-- Name: posts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.posts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.posts_id_seq OWNER TO postgres;

--
-- Name: posts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.posts_id_seq OWNED BY public.posts.id;


--
-- Name: posts_tags; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.posts_tags (
    tag_id integer,
    post_id integer
);


ALTER TABLE public.posts_tags OWNER TO postgres;

--
-- Name: tags; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tags (
    id integer NOT NULL,
    name character varying(50) NOT NULL
);


ALTER TABLE public.tags OWNER TO postgres;

--
-- Name: tags_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.tags_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.tags_id_seq OWNER TO postgres;

--
-- Name: tags_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.tags_id_seq OWNED BY public.tags.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(50) NOT NULL,
    email character varying(255) NOT NULL,
    password_hash text NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts ALTER COLUMN id SET DEFAULT nextval('public.posts_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tags ALTER COLUMN id SET DEFAULT nextval('public.tags_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: books; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.books (isbn, title, author, price) FROM stdin;
978-1505255607	The Time Machine	H. G. Wells	6.99
978-1505255601	Redion	Ridwan Fathin	699.00
978-1503379640	The Princess	Niccolò Machiavelli	6.99
\.


--
-- Data for Name: posts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.posts (id, title, body, published_at, author_id) FROM stdin;
1	My first gomidway post	Golang rocks!	2018-07-20 07:08:27.805582	1
4	My gomidway post	Go golang rocks! 	2018-07-27 08:12:29.42214	1
6	My gomidway post1	Go golang rocks! 	2018-07-28 09:38:42.305815	1
9	My gomidway post2	Go golang rocks! 	2018-07-28 10:22:45.757798	1
12	My gomidway post223	Go golang rocks! 	2018-07-29 10:46:03.175204	1
14	My gomidway post3	Go golang rocks! 	2018-07-29 11:40:47.275026	1
16	My gomidway post4	Go golang rocks! 	2018-07-29 11:49:04.393163	1
21	My gomidway post5	Go golang rocks! 	2018-07-30 03:53:29.909682	1
23	My gomidway post6	Go golang rocks! 	2018-07-30 04:16:16.914154	1
24	Give this awesome post a title	Go golang rocks! 	2018-07-31 06:50:05.51035	1
11	Updated Title	Go golang rocks! 	2018-07-29 10:39:28.077217	1
26	Give this awesome post a title2	Go golang rocks! 	2018-07-31 10:05:00.638538	1
27	Give this awesome post a title21	Go golang rocks! 	2018-08-02 08:05:42.387747	1
\.


--
-- Name: posts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.posts_id_seq', 27, true);


--
-- Data for Name: posts_tags; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.posts_tags (tag_id, post_id) FROM stdin;
1	1
2	1
1	4
2	4
1	6
2	6
1	9
2	9
1	26
2	26
1	27
2	27
\.


--
-- Data for Name: tags; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tags (id, name) FROM stdin;
1	intro
2	golang
\.


--
-- Name: tags_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.tags_id_seq', 2, true);


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, username, email, password_hash) FROM stdin;
1	foo	foo@bar.com	$2a$10$gCfbY8aXOxIqdh1JuPhr6.2V4ra6TkLJeo/t6RcPZ.81PPVTjjjM2
3	wan	wan@bar.com	$2a$10$MB8fVUcMehB/afkpdscTzeZ2p8AqZs77fsf/dOeh10vK6992dQUTm
14	Rid	rid@shade.com	$2a$10$Z8pw//OUP5af9UqGyxao.uepDgd493Z//sRAAqw8o2ImgVENFWlh2
10	airin	airin@shade.com	$2a$10$2fKoKCEPRd8mib0Jo8F5Fenr6OlW7BO7RzUi//PU96PNa6nwHiaXi
15	Rida	rida@shade.com	$2a$10$l0Ai9rqb.01PX4wRQFmHe.LM.MZf0LMHDhtmRjQ2zNAiuINu1UqIy
\.


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 15, true);


--
-- Name: books_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.books
    ADD CONSTRAINT books_pkey PRIMARY KEY (isbn);


--
-- Name: posts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);


--
-- Name: posts_title_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_title_key UNIQUE (title);


--
-- Name: tags_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);


--
-- Name: users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: posts_author_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_author_id_fkey FOREIGN KEY (author_id) REFERENCES public.users(id);


--
-- Name: posts_tags_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts_tags
    ADD CONSTRAINT posts_tags_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(id);


--
-- Name: posts_tags_tag_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts_tags
    ADD CONSTRAINT posts_tags_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES public.tags(id);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- Name: TABLE books; Type: ACL; Schema: public; Owner: postgres
--

REVOKE ALL ON TABLE public.books FROM PUBLIC;
REVOKE ALL ON TABLE public.books FROM postgres;
GRANT ALL ON TABLE public.books TO postgres;
GRANT ALL ON TABLE public.books TO airin;


--
-- Name: TABLE posts; Type: ACL; Schema: public; Owner: postgres
--

REVOKE ALL ON TABLE public.posts FROM PUBLIC;
REVOKE ALL ON TABLE public.posts FROM postgres;
GRANT ALL ON TABLE public.posts TO postgres;
GRANT ALL ON TABLE public.posts TO airin;


--
-- Name: SEQUENCE posts_id_seq; Type: ACL; Schema: public; Owner: postgres
--

REVOKE ALL ON SEQUENCE public.posts_id_seq FROM PUBLIC;
REVOKE ALL ON SEQUENCE public.posts_id_seq FROM postgres;
GRANT ALL ON SEQUENCE public.posts_id_seq TO postgres;
GRANT ALL ON SEQUENCE public.posts_id_seq TO airin;


--
-- Name: TABLE posts_tags; Type: ACL; Schema: public; Owner: postgres
--

REVOKE ALL ON TABLE public.posts_tags FROM PUBLIC;
REVOKE ALL ON TABLE public.posts_tags FROM postgres;
GRANT ALL ON TABLE public.posts_tags TO postgres;
GRANT ALL ON TABLE public.posts_tags TO airin;


--
-- Name: TABLE tags; Type: ACL; Schema: public; Owner: postgres
--

REVOKE ALL ON TABLE public.tags FROM PUBLIC;
REVOKE ALL ON TABLE public.tags FROM postgres;
GRANT ALL ON TABLE public.tags TO postgres;
GRANT ALL ON TABLE public.tags TO airin;


--
-- Name: SEQUENCE tags_id_seq; Type: ACL; Schema: public; Owner: postgres
--

REVOKE ALL ON SEQUENCE public.tags_id_seq FROM PUBLIC;
REVOKE ALL ON SEQUENCE public.tags_id_seq FROM postgres;
GRANT ALL ON SEQUENCE public.tags_id_seq TO postgres;
GRANT ALL ON SEQUENCE public.tags_id_seq TO airin;


--
-- Name: TABLE users; Type: ACL; Schema: public; Owner: postgres
--

REVOKE ALL ON TABLE public.users FROM PUBLIC;
REVOKE ALL ON TABLE public.users FROM postgres;
GRANT ALL ON TABLE public.users TO postgres;
GRANT ALL ON TABLE public.users TO airin;


--
-- Name: SEQUENCE users_id_seq; Type: ACL; Schema: public; Owner: postgres
--

REVOKE ALL ON SEQUENCE public.users_id_seq FROM PUBLIC;
REVOKE ALL ON SEQUENCE public.users_id_seq FROM postgres;
GRANT ALL ON SEQUENCE public.users_id_seq TO postgres;
GRANT ALL ON SEQUENCE public.users_id_seq TO airin;


--
-- PostgreSQL database dump complete
--

