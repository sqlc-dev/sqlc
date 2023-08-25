CREATE TABLE bar (
    id INT NOT NULL,
    "!!!nobody,_,-would-believe---this-...?!" INT,
    "parent        id" INT);

-- name: test :one
SELECT * from bar limit 1;
