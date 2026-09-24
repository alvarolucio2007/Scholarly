CREATE TABLE tests(
  id BIGSERIAL PRIMARY KEY,
  course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  weight NUMERIC(4,2) NOT NULL DEFAULT 1.0,
  test_date DATE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ,

  CONSTRAINT tests_weight_valid CHECK (weight>0),
  CONSTRAINT tests_name_not_empty CHECK (LENGTH(TRIM(name))>0)
);
CREATE INDEX idx_tests_course ON tests(course_id);
