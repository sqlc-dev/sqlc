CREATE TABLE campus (id text not null);
CREATE TABLE students (id text not null);
CREATE TABLE product_meta (id text not null);
CREATE TABLE calories (id text not null);

-- name: ListCampuses :many
SELECT * FROM campus;

-- name: ListStudents :many
SELECT * FROM students;

-- name: ListMetadata :many
SELECT * FROM product_meta;

-- name: ListCalories :many
SELECT * FROM calories;
