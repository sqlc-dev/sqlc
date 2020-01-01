CREATE TABLE teachers (
  id int NOT NULL,
  first_name varchar(255),
  last_name varchar(255),
  school_id VARCHAR(255) NOT NULL,
  school_lat FLOAT,
  school_lng FLOAT,
  department ENUM("English", "Math"),
  PRIMARY KEY (id)
);

CREATE TABLE students (
  id int NOT NULL,
  student_id varchar(10),
  first_name varchar(255),
  last_name varchar(255),
  PRIMARY KEY (id)
)

/* name: GetTeachersByID :one */
SELECT * FROM teachers WHERE id = ?

/* name: GetSomeTeachers :one */
SELECT school_id, id FROM teachers WHERE school_lng > ? AND school_lat < ?;

/* name: TeachersByID :one */
SELECT id, school_lat FROM teachers WHERE id = ? LIMIT 10 

/* name: GetStudents :many */
SELECT students.first_name, students.last_name, teachers.first_name teacherFirstName
  FROM students 
  LEFT JOIN teachers on teachers.school_id = students.school_id