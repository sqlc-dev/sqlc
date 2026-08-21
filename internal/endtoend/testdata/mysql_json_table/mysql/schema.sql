CREATE TABLE person (
  emails TEXT NOT NULL
);

CREATE TABLE person_email (
  email TEXT NOT NULL
);

INSERT INTO person_email (`email`)
  SELECT pe.email
  FROM person AS p
  JOIN JSON_TABLE(p.emails, '$[*]' COLUMNS(email TEXT PATH '$')) AS pe;
