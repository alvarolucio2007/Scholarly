CREATE OR REPLACE FUNCTION fn_weighted_average(
    p_student_id BIGINT,
    p_course_id  BIGINT
) RETURNS NUMERIC
LANGUAGE plpgsql
AS $$
DECLARE
    v_result NUMERIC;
BEGIN
    IF NOT EXISTS (SELECT 1 FROM students WHERE user_id = p_student_id) THEN
        RAISE EXCEPTION 'student not found: %', p_student_id
            USING ERRCODE = 'P0S01';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM courses WHERE id = p_course_id) THEN
        RAISE EXCEPTION 'course not found: %', p_course_id
            USING ERRCODE = 'P0C01';
    END IF;

    SELECT COALESCE(
        SUM(t.weight * g.value) / SUM(t.weight),
        0
    )
    INTO v_result
    FROM grades g
    JOIN tests t       ON t.id = g.test_id
    JOIN enrollments e ON e.id = g.enrollment_id
    WHERE e.student_id = p_student_id
      AND e.course_id  = p_course_id;

    RETURN v_result;
END;
$$;
