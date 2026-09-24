CREATE TABLE students(
  user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  enrollment_number VARCHAR(20) UNIQUE,

  CONSTRAINT students_enrollment_not_empty CHECK (LENGTH(TRIM(enrollment_number))>0)
)
