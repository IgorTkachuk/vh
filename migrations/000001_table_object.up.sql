CREATE TABLE object(
    id serial PRIMARY KEY,
    storage_name TEXT,
    orig_name TEXT,
    orig_date DATE,
    add_date DATE,
    billing_pn TEXT,
    user_name TEXT,
    notes TEXT
);