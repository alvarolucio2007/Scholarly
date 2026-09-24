CREATE TABLE enrollments(
  id BIGSERIAL PRIMARY KEY,
  student_id BIGINT NOT NULL REFERENCES students(user_id) ON DELETE CASCADE,
  course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  enrolled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  status VARCHAR(20) NOT NULL DEFAULT 'active',

  CONSTRAINT enrollments_status_valid CHECK (status IN ('active','cancelled','completed')),
  CONSTRAINT enrollments_student_course_unique UNIQUE (student_id, course_id)
);
CREATE INDEX idx_enrollments_course  ON enrollments(course_id);
