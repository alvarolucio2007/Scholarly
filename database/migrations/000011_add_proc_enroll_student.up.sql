CREATE PROCEDURE sp_enroll_student(
      IN p_student_id BIGINT,
      IN p_course_id  BIGINT,
      OUT p_enrollment_id BIGINT
  )
  LANGUAGE plpgsql
  AS $$
  BEGIN
    IF NOT EXISTS (SELECT 1 FROM students where user_id=p_student_id)  THEN
      RAISE EXCEPTION 'student not found' USING ERRCODE = 'P0S01';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM courses WHERE id=p_course_id) THEN
      RAISE EXCEPTION 'course not found' USING ERRCODE = 'P0C01';
    END IF;
    IF EXISTS (SELECT 1 FROM enrollments WHERE student_id=p_student_id AND course_id=p_course_id)  THEN
      RAISE EXCEPTION 'student already enrolled in course' USING ERRCODE = 'P0E01';
    END IF;
    IF (SELECT COUNT(*) FROM enrollments WHERE course_id = p_course_id AND status='active' ) >= 
        (SELECT max_students FROM courses WHERE id = p_course_id ) THEN
      RAISE EXCEPTION 'course is full' USING ERRCODE='P0C02';
    END IF;
    INSERT INTO enrollments (student_id,course_id) VALUES (p_student_id,p_course_id) RETURNING id INTO p_enrollment_id;
  END;
  $$
