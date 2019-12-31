CREATE TABLE teachers (
  id int NOT NULL,
  school_id VARCHAR(255) NOT NULL,
  school_lat FLOAT,
  school_lng FLOAT,
  department ENUM("English", "Math"),
  PRIMARY KEY (ID)
);

/* name: GetTeachersByID :one */
SELECT * FROM teachers WHERE id = ?

/* name: GetSomeTeachers :one */
SELECT school_id, id FROM teachers WHERE school_lng > ? AND school_lat < ?;

/* name: TeachersByID :one */
SELECT id, school_lat FROM teachers WHERE id = ? LIMIT 10 