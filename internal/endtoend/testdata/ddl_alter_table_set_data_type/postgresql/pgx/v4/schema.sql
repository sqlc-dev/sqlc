CREATE TABLE foo (bar text);
ALTER TABLE foo ALTER bar SET DATA TYPE bool USING bar::boolean;
