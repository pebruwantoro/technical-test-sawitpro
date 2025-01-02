-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.

-- THIS IS SCRIPT FOR CREATING ESTATES TABLE
CREATE TABLE estates (
	id UUID PRIMARY KEY,
	width INT NOT NULL CHECK ( width > 0 AND width <= 50000 ),
	length INT NOT NULL CHECK ( length > 0 AND length <= 50000 )
);

-- THIS IS SCRIPT FOR CREATING TREES TABLE
CREATE TABLE trees (
    id UUID PRIMARY KEY,
    estate_id UUID REFERENCES estates(id) ON DELETE CASCADE,
	x INT NOT NULL CHECK ( x > 0 ),
	y INT NOT NULL CHECK ( y > 0 ),
	height INT NOT NULL CHECK ( height >= 1 AND height <= 30 ),
	UNIQUE (estate_id, x, y)
);