package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

var _ ports.ReportCardRepository = (*ReportCardRepo)(nil)

type ReportCardRepo struct {
	db *sql.DB
}

func NewReportCardRepo(db *sql.DB) *ReportCardRepo {
	return &ReportCardRepo{db: db}
}

func (r *ReportCardRepo) List(ctx context.Context, filter ports.ReportCardFilter) ([]*domain.ReportCard, error) {
	var query string
	var args []any
	if filter.StudentID != nil {
		query = `SELECT student_id, student_name, course_id, course_name, course_code, teacher_name, average, situation FROM vw_student_report_card WHERE student_id = $1 ORDER BY student_id,course_name`
		args = append(args, *filter.StudentID)
	} else {
		query = `SELECT student_id, student_name, course_id, course_name, course_code, teacher_name, average, situation FROM vw_student_report_card ORDER BY student_id,course_name`
	}
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("postgres: view report card: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("postgres: rows.close report card: %v", err)
		}
	}()
	var reportCards []*domain.ReportCard
	for rows.Next() {
		var rep domain.ReportCard
		if err := rows.Scan(&rep.StudentID, &rep.StudentName, &rep.CourseID, &rep.CourseName, &rep.CourseCode, &rep.TeacherName, &rep.Average, &rep.Situation); err != nil {
			return nil, fmt.Errorf("postgres: rows.scan report card: %w", err)
		}
		reportCards = append(reportCards, &rep)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: iterate report card: %w", err)
	}
	return reportCards, nil
}
