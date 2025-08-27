ALTER TABLE results
ADD COLUMN status TEXT DEFAULT 'ok',
ADD COLUMN error_message TEXT;
