CREATE VIEW vw_student_report_card AS
SELECT
    s.user_id         AS student_id,
    u.name            AS student_name,
    c.id              AS course_id,
    c.name            AS course_name,
    c.code            AS course_code,
    tu.name               AS teacher_name,
    fn_weighted_average(s.user_id,c.id)               AS average,
  CASE  
        WHEN fn_weighted_average(s.user_id, c.id) >= 7 THEN 'Approved'
        WHEN fn_weighted_average(s.user_id, c.id) >= 5 THEN 'Recovery'
        ELSE 'Failed'
  END AS situation
FROM students s
JOIN users u         ON s.user_id=u.id
JOIN enrollments e   ON e.student_id=s.user_id
JOIN courses c       ON c.id=e.course_id
JOIN teachers t      ON c.teacher_id=t.user_id
JOIN users tu        ON t.user_id=tu.id;
