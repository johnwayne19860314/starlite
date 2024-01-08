CREATE SCHEMA IF NOT EXISTS first;

CREATE TYPE  first.user_role  AS ENUM ('admin', 'power', 'internal');
CREATE TABLE IF NOT EXISTS first.user (
    id serial NOT NULL  PRIMARY KEY
   
    , user_name varchar(200) NOT NULL
    , user_email varchar(100) NOT NULL
    , user_password varchar(100) NOT NULL
    , user_role first.user_role NOT NULL
    , is_active boolean NOT NULL 
    
    , created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now())
    , updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL   DEFAULT (now())
);

CREATE UNIQUE INDEX on first.user(user_name);

CREATE TABLE IF NOT EXISTS first.entry_category (
    id serial NOT NULL  PRIMARY KEY
    , category varchar(10) UNIQUE NOT NULL 
    , note text NOT NULL
    , is_active boolean NOT NULL 
    
    , created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now())
    , updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS first.entry (
    id serial NOT NULL  PRIMARY KEY
   
    , entry_code varchar(50) UNIQUE NOT NULL 
    , entry_category varchar(10) REFERENCES first.entry_category(category) NOT NULL 
    , entry_name varchar(100) UNIQUE NOT NULL
    , entry_amount int NOT NULL
    , entry_weight float NOT NULL
    , entry_note text NOT NULL
    , is_active boolean NOT NULL 
    
    , created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now())
    , updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX on first.entry(entry_code);
CREATE UNIQUE INDEX on first.entry(entry_name);
CREATE UNIQUE INDEX on first.entry(entry_code,entry_name);

