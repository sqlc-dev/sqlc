CREATE TABLE teachers (
  id int NOT NULL,
  first_name varchar(255),
  last_name varchar(255),
  school_id int NOT NULL,
  class_id int NOT NULL, 
  school_lat FLOAT,
  school_lng FLOAT,
  department ENUM("English", "Math"),
  PRIMARY KEY (id)
);

CREATE TABLE students (
  id int NOT NULL,
  class_id int NOT NULL,
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

/* name: GetStudentsTeacher :one */
SELECT students.first_name, students.last_name, teachers.first_name teacherFirstName,
teachers.id teacher_id
  FROM students 
  Left JOIN teachers on teachers.class_id = students.class_id
  WHERE students.id = :studentID
