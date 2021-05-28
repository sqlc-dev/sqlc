CREATE TABLE campus (id text not null);
CREATE TABLE students (id text not null);

-- name: ListCampuses :many
SELECT * FROM campus;

-- name: ListStudents :many
SELECT * FROM students;
