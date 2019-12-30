CREATE TABLE students (
  id int NOT NULL,
  school_id VARCHAR(255) NOT NULL,
  school_lat FLOAT,
  department ENUM("English", "Math"),
  PRIMARY KEY (ID)
);

/* name: GetAllStudents :many */
SELECT school_id, id FROM students WHERE id = ? AND school_id = ?

/* name: GetSomeStudents :one */
SELECT school_id, id FROM students WHERE school_id = ?

/* name: StudentByID :one */
SELECT id, school_lat FROM students WHERE id = ? LIMIT 10 