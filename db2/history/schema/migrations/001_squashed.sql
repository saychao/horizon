-- noinspection SqlNoDataSourceInspectionForFile

-- +migrate Up


--
-- Name: history_accounts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE history_accounts (
    id bigserial NOT NULL,
    address character varying(64)
);


--
-- Name: history_balances; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE history_balances (
    id bigint NOT NULL,
    balance_id character varying(56) NOT NULL,
    asset character varying(4) NOT NULL,
    account_id character varying(56) NOT NULL,
    kyc jsonb NOT NULL
);


--
-- Name: history_balances_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE history_balances_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: history_balances_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE history_balances_id_seq OWNED BY history_balances.id;


--
-- Name: history_forfeit_requests; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE history_forfeit_requests (
    id bigint NOT NULL,
    target character varying(64) NOT NULL,
    amount character varying(64) NOT NULL,
    initiated_by_user boolean NOT NULL,
    accepted boolean,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    forfeit_type integer DEFAULT 0 NOT NULL
);


--
-- Name: history_forfeit_requests_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE history_forfeit_requests_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: history_forfeit_requests_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE history_forfeit_requests_id_seq OWNED BY history_forfeit_requests.id;


--
-- Name: history_ledgers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE history_ledgers (
    sequence integer NOT NULL,
    ledger_hash character varying(64) NOT NULL,
    previous_ledger_hash character varying(64),
    transaction_count integer DEFAULT 0 NOT NULL,
    operation_count integer DEFAULT 0 NOT NULL,
    closed_at timestamp without time zone NOT NULL,
    id bigint,
    importer_version integer DEFAULT 1 NOT NULL,
    total_coins bigint NOT NULL,
    fee_pool bigint NOT NULL,
    base_fee integer NOT NULL,
    base_reserve integer NOT NULL,
    max_tx_set_size integer NOT NULL
);


--
-- Name: history_operation_participants; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE history_operation_participants (
    id integer NOT NULL,
    history_operation_id bigint NOT NULL,
    history_account_id bigint NOT NULL,
    balance_id character varying(64) DEFAULT ''::character varying NOT NULL,
    effects jsonb
);


--
-- Name: history_operation_participants_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE history_operation_participants_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: history_operation_participants_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE history_operation_participants_id_seq OWNED BY history_operation_participants.id;


--
-- Name: history_operations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE history_operations (
    id bigint NOT NULL,
    transaction_id bigint NOT NULL,
    application_order integer NOT NULL,
    type integer NOT NULL,
    details jsonb,
    source_account character varying(64) DEFAULT ''::character varying NOT NULL,
    ledger_close_time timestamp without time zone DEFAULT now() NOT NULL,
    identifier bigint NOT NULL,
    state integer
);


--
-- Name: history_operations_identifier_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE history_operations_identifier_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: history_operations_identifier_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE history_operations_identifier_seq OWNED BY history_operations.identifier;


--
-- Name: history_payment_requests; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE history_payment_requests (
    id bigint NOT NULL,
    payment_id bigint NOT NULL,
    accepted boolean,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    details jsonb NOT NULL
);


--
-- Name: history_payment_requests_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE history_payment_requests_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: history_payment_requests_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE history_payment_requests_id_seq OWNED BY history_payment_requests.id;


--
-- Name: history_payment_requests_payment_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE history_payment_requests_payment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: history_payment_requests_payment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE history_payment_requests_payment_id_seq OWNED BY history_payment_requests.payment_id;


--
-- Name: history_transaction_participants; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE history_transaction_participants (
    id integer NOT NULL,
    history_transaction_id bigint NOT NULL,
    history_account_id bigint NOT NULL
);


--
-- Name: history_transaction_participants_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE history_transaction_participants_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: history_transaction_participants_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE history_transaction_participants_id_seq OWNED BY history_transaction_participants.id;


--
-- Name: history_transactions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE history_transactions (
    transaction_hash character varying(64) NOT NULL,
    ledger_sequence integer NOT NULL,
    application_order integer NOT NULL,
    account character varying(64) NOT NULL,
    account_sequence bigint NOT NULL,
    fee_paid integer NOT NULL,
    operation_count integer NOT NULL,
    ledger_close_time timestamp without time zone,
    id bigint,
    tx_envelope text NOT NULL,
    tx_result text NOT NULL,
    tx_meta text NOT NULL,
    tx_fee_meta text NOT NULL,
    signatures character varying(96)[] DEFAULT '{}'::character varying[] NOT NULL,
    memo_type character varying DEFAULT 'none'::character varying NOT NULL,
    memo character varying,
    time_bounds int8range
);


--
-- Name: history_balances id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY history_balances ALTER COLUMN id SET DEFAULT nextval('history_balances_id_seq'::regclass);

--
-- Name: history_forfeit_requests id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY history_forfeit_requests ALTER COLUMN id SET DEFAULT nextval('history_forfeit_requests_id_seq'::regclass);


--
-- Name: history_operation_participants id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY history_operation_participants ALTER COLUMN id SET DEFAULT nextval('history_operation_participants_id_seq'::regclass);


--
-- Name: history_operations identifier; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY history_operations ALTER COLUMN identifier SET DEFAULT nextval('history_operations_identifier_seq'::regclass);


--
-- Name: history_payment_requests id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY history_payment_requests ALTER COLUMN id SET DEFAULT nextval('history_payment_requests_id_seq'::regclass);


--
-- Name: history_payment_requests payment_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY history_payment_requests ALTER COLUMN payment_id SET DEFAULT nextval('history_payment_requests_payment_id_seq'::regclass);


--
-- Name: history_transaction_participants id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY history_transaction_participants ALTER COLUMN id SET DEFAULT nextval('history_transaction_participants_id_seq'::regclass);


--
-- Name: history_balances history_balances_balance_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY history_balances
    ADD CONSTRAINT history_balances_balance_id_key UNIQUE (balance_id);


--
-- Name: history_balances history_balances_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY history_balances
    ADD CONSTRAINT history_balances_pkey PRIMARY KEY (id);


--
-- Name: history_forfeit_requests history_forfeit_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY history_forfeit_requests
    ADD CONSTRAINT history_forfeit_requests_pkey PRIMARY KEY (id);


--
-- Name: history_operation_participants history_operation_participants_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY history_operation_participants
    ADD CONSTRAINT history_operation_participants_pkey PRIMARY KEY (id);


--
-- Name: history_payment_requests history_payment_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY history_payment_requests
    ADD CONSTRAINT history_payment_requests_pkey PRIMARY KEY (id);


--
-- Name: history_transaction_participants history_transaction_participants_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY history_transaction_participants
    ADD CONSTRAINT history_transaction_participants_pkey PRIMARY KEY (id);


--
-- Name: by_account; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX by_account ON history_transactions USING btree (account, account_sequence);


--
-- Name: by_hash; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX by_hash ON history_transactions USING btree (transaction_hash);


--
-- Name: by_ledger; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX by_ledger ON history_transactions USING btree (ledger_sequence, application_order);


--
-- Name: hist_op_p_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX hist_op_p_id ON history_operation_participants USING btree (history_account_id, history_operation_id, balance_id);


--
-- Name: hist_tx_p_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX hist_tx_p_id ON history_transaction_participants USING btree (history_account_id, history_transaction_id);


--
-- Name: hop_by_hoid; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hop_by_hoid ON history_operation_participants USING btree (history_operation_id);


--
-- Name: hs_ledger_by_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX hs_ledger_by_id ON history_ledgers USING btree (id);


--
-- Name: hs_transaction_by_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX hs_transaction_by_id ON history_transactions USING btree (id);


--
-- Name: htp_by_htid; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX htp_by_htid ON history_transaction_participants USING btree (history_transaction_id);


--
-- Name: index_history_accounts_on_address; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX index_history_accounts_on_address ON history_accounts USING btree (address);


--
-- Name: index_history_accounts_on_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX index_history_accounts_on_id ON history_accounts USING btree (id);


--
-- Name: index_history_ledgers_on_closed_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX index_history_ledgers_on_closed_at ON history_ledgers USING btree (closed_at);


--
-- Name: index_history_ledgers_on_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX index_history_ledgers_on_id ON history_ledgers USING btree (id);


--
-- Name: index_history_ledgers_on_importer_version; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX index_history_ledgers_on_importer_version ON history_ledgers USING btree (importer_version);


--
-- Name: index_history_ledgers_on_ledger_hash; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX index_history_ledgers_on_ledger_hash ON history_ledgers USING btree (ledger_hash);


--
-- Name: index_history_ledgers_on_previous_ledger_hash; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX index_history_ledgers_on_previous_ledger_hash ON history_ledgers USING btree (previous_ledger_hash);


--
-- Name: index_history_ledgers_on_sequence; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX index_history_ledgers_on_sequence ON history_ledgers USING btree (sequence);


--
-- Name: index_history_operations_on_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX index_history_operations_on_id ON history_operations USING btree (id);


--
-- Name: index_history_operations_on_transaction_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX index_history_operations_on_transaction_id ON history_operations USING btree (transaction_id);


--
-- Name: index_history_operations_on_type; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX index_history_operations_on_type ON history_operations USING btree (type);


--
-- Name: index_history_transactions_on_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX index_history_transactions_on_id ON history_transactions USING btree (id);


-- +migrate Down

drop table history_accounts cascade;
drop table history_balances cascade;
drop table history_forfeit_requests cascade;
drop table history_ledgers cascade;
drop table history_operation_participants cascade;
drop table history_operations cascade;
drop table history_payment_requests cascade;
drop table history_transaction_participants cascade;
drop table history_transactions cascade;