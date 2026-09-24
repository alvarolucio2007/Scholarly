CREATE TABLE courses(
  id BIGSERIAL PRIMARY KEY,
  teacher_id BIGINT NOT NULL REFERENCES teachers(user_id) ON DELETE RESTRICT,
  name VARCHAR(255) NOT NULL,
  code VARCHAR(50) NOT NULL,
  semester VARCHAR(50) NOT NULL,
  max_students SMALLINT NOT NULL DEFAULT 30,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ,

  CONSTRAINT courses_name_not_empty CHECK (LENGTH(TRIM(name))>0),
  CONSTRAINT courses_code_not_empty CHECK (LENGTH(TRIM(code))>0),
  CONSTRAINT courses_semester_not_empty CHECK (LENGTH(TRIM(semester))>0),
  CONSTRAINT courses_max_students_not_zero CHECK (max_students>0)
  CONSTRAINT courses_code_semester_unique UNIQUE (code, semester)
);
  CREATE INDEX idx_courses_teacher ON courses(teacher_id);
