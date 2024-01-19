

-- CREATE TABLE IF NOT EXISTS first.entry_supplier (
--     id serial NOT NULL PRIMARY KEY
--     , supplier_name varchar(100) UNIQUE NOT NULL 
--     , contact_info text NOT NULL
--     , is_active boolean NOT NULL 
    
--     , created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now())
--     , updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now())
-- );

ALTER TABLE first.entry 
    ADD COLUMN supplier_name varchar(200) ,
    ADD COLUMN supplier_contact_info varchar(200) ,
    ADD COLUMN fix varchar(50) NOT NULL DEFAULT '',
    ADD COLUMN chemical_name varchar(50) NOT NULL DEFAULT '';