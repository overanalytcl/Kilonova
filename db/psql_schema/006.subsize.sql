ALTER TABLE submissions ADD COLUMN code_size INTEGER NOT NULL GENERATED ALWAYS AS (length(code)) STORED;
