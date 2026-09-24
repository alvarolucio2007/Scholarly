
CREATE TABLE grades(
  id BIGSERIAL PRIMARY KEY,
  test_id BIGINT NOT NULL REFERENCES tests(id) ON DELETE CASCADE,
  enrollment_id BIGINT NOT NULL REFERENCES enrollments(id) ON DELETE CASCADE,
  value NUMERIC(4,2) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ,

  CONSTRAINT grades_value_valid CHECK (value>=0 AND value<=10),
 CONSTRAINT grades_unique UNIQUE (enrollment_id, test_id)
);
CREATE INDEX idx_grades_test ON grades(test_id);
