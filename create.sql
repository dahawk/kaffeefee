CREATE TABLE users
(
    id bigserial NOT NULL,
    name text,
    CONSTRAINT id_pkey PRIMARY KEY (id)
);

CREATE TABLE log
(
    id bigserial NOT NULL,
    userid bigint,
    ts bigint,
    cnt integer,
    CONSTRAINT log_pkey PRIMARY KEY (id),
    CONSTRAINT "log_userId_fkey" FOREIGN KEY (userid)
        REFERENCES users (id) MATCH SIMPLE
        ON UPDATE NO ACTION ON DELETE NO ACTION
);