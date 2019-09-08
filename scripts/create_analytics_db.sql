--
-- PostgreSQL database dump
--

-- Dumped from database version 10.6
-- Dumped by pg_dump version 10.9 (Ubuntu 10.9-0ubuntu0.18.04.1)

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

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

/* COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';*/


SET default_tablespace = '';

SET default_with_oids = false;

CREATE TABLE IF NOT EXISTS public.test_meta_data (
    "timestamp" timestamp without time zone,
    test_name TEXT    NOT NULL,
    values TEXT    NOT NULL

);
--
-- Name: ledger_db_end; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ledger_db_end (
    node character varying(5),
    "timestamp" timestamp without time zone,
    sequence integer
);


--
-- Name: ledger_db_start; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ledger_db_start (
    node character varying(5),
    "timestamp" timestamp without time zone,
    sequence integer
);


--
-- Name: ledger_tx_count; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ledger_tx_count (
    node character varying(5),
    count integer,
    sequence integer
);


--
-- Name: ledger_vote_end; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ledger_vote_end (
    node character varying(5),
    "timestamp" timestamp without time zone,
    sequence integer
);


--
-- Name: ledger_vote_start; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ledger_vote_start (
    node character varying(5),
    "timestamp" timestamp without time zone,
    sequence integer
);


--
-- Name: core_timings; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.core_timings AS
 SELECT db_s.node,
    db_s.sequence,
    tx_c.count,
    db_s."timestamp" AS db_start,
    db_e."timestamp" AS db_end,
    date_part('epoch'::text, (db_e."timestamp" - db_s."timestamp")) AS db_duration,
    v_s."timestamp" AS v_start,
    v_e."timestamp" AS v_end,
    date_part('epoch'::text, (v_e."timestamp" - v_s."timestamp")) AS v_duration,
    ts_md."test_name" AS test_name
   FROM public.ledger_db_start db_s,
    public.ledger_db_end db_e,
    public.ledger_vote_start v_s,
    public.ledger_vote_end v_e,
    public.ledger_tx_count tx_c,
    public.test_meta_data ts_md
  WHERE ((db_s.sequence = db_e.sequence) AND (db_s.sequence = v_s.sequence) AND (db_s.sequence = v_e.sequence) AND (db_s.sequence = tx_c.sequence) AND ((db_s.node)::text = (db_e.node)::text) AND ((db_s.node)::text = (v_s.node)::text) AND ((db_s.node)::text = (v_e.node)::text) AND ((db_s.node)::text = (tx_c.node)::text));


--
-- Name: ingestion; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ingestion (
    ledger integer NOT NULL,
    "timestamp" timestamp(6) without time zone NOT NULL
);


--
-- Name: submission; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submission (
    hash character varying(64) NOT NULL,
    "timestamp" timestamp(6) without time zone NOT NULL
);


--
-- Name: last_submission; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.last_submission AS
 SELECT submission.hash,
    max(submission."timestamp") AS ts
   FROM public.submission
  GROUP BY submission.hash;


--
-- Name: tx_ledger; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tx_ledger (
    hash character varying(64) NOT NULL,
    sequence integer NOT NULL
);


--
-- Name: tx_detail; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.tx_detail AS
 SELECT t.hash,
    t.sequence,
    s.ts AS submission_time,
    i."timestamp" AS ingestion_time,
    date_part('epoch'::text, (i."timestamp" - s.ts)) AS duration
   FROM public.tx_ledger t,
    public.last_submission s,
    public.ingestion i
  WHERE ((t.sequence = i.ledger) AND ((t.hash)::text = (s.hash)::text));


--
-- Name: ingestion ingestion_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ingestion
    ADD CONSTRAINT ingestion_pkey PRIMARY KEY (ledger);


--
-- Name: ledger_db_end.sequence; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "ledger_db_end.sequence" ON public.ledger_db_end USING btree (sequence);


--
-- Name: ledger_db_start.sequence; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "ledger_db_start.sequence" ON public.ledger_db_start USING btree (sequence);


--
-- Name: ledger_tx_count.sequence; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "ledger_tx_count.sequence" ON public.ledger_tx_count USING btree (sequence);


--
-- Name: ledger_vote_end.sequence; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "ledger_vote_end.sequence" ON public.ledger_vote_end USING btree (sequence);


--
-- Name: ledger_vote_start.sequence; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "ledger_vote_start.sequence" ON public.ledger_vote_start USING btree (sequence);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: -
--

REVOKE ALL ON SCHEMA public FROM rdsadmin;
REVOKE ALL ON SCHEMA public FROM PUBLIC;
/* GRANT ALL ON SCHEMA public TO stellar;*/
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

