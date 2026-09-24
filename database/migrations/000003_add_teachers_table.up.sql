CREATE TABLE teachers(
  user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  department VARCHAR(100) NOT NULL,

  CONSTRAINT teachers_department_not_empty CHECK (LENGTH(TRIM(department))>0)
)
