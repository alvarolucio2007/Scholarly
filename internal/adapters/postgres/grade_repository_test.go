package postgres_test

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/postgres"
	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/go-openapi/testify/require"
)

type gradeScenario struct {
	User       *domain.User
	User2      *domain.User
	Teacher    *domain.Teacher
	Student    *domain.Student
	Course     *domain.Course
	Enrollment *domain.Enrollment
	Test       *domain.Test
}

func createGradeTest(t *testing.T, grade *domain.Grade) {
	r := postgres.NewGradeRepo(testDB, testLogger)
	err := r.Create(t.Context(), grade)
	require.NoError(t, err)
}

func cleanupGrade(t *testing.T, gradeID, testID, courseID, enrollmentID, teacherID, studentID, user1ID, user2ID int64) {
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM grades WHERE id = $1", gradeID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM tests WHERE id = $1", testID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM courses WHERE id = $1", courseID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM enrollments WHERE id = $1", enrollmentID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM teachers WHERE user_id=$1", teacherID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM students WHERE user_id=$1", studentID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user1ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user2ID)
		require.NoError(t, err)
	})
}

func newGradeScenario(t *testing.T) gradeScenario {
	user := &domain.User{Name: "Create test", Email: "createtest@email.com", CPF: "00000000000", PasswordHash: []byte("123456789")}
	createUserTest(t, user)
	user2 := &domain.User{Name: "Create student test", Email: "createstudent@email.com", CPF: "00000000002", PasswordHash: []byte("123456789")}
	createUserTest(t, user2)
	teacher := &domain.Teacher{UserID: user.ID, Department: "CourseCreation"}
	createTeacherTest(t, teacher)
	student := &domain.Student{UserID: user2.ID, EnrollmentNumber: "Test"}
	createStudentTest(t, student)
	course := &domain.Course{TeacherID: teacher.UserID, Name: "GetByID", Code: "GETBYID", Semester: "GETBYIDSEMESTER", MaxStudents: 100}
	createCourseTest(t, course)
	enrollment := &domain.Enrollment{StudentID: student.UserID, CourseID: course.ID, Status: "active"}
	createEnrollmentTest(t, enrollment)
	test := &domain.Test{CourseID: course.ID, Name: "Successful test creation", Weight: 6}
	createTestTest(t, test)
	return gradeScenario{user, user2, teacher, student, course, enrollment, test}
}

func TestGradeRepo_Create(t *testing.T) {
	s := newGradeScenario(t)
	grade := &domain.Grade{TestID: s.Test.ID, EnrollmentID: s.Enrollment.ID, Value: 10}
	r := postgres.NewGradeRepo(testDB, testLogger)
	tests := []struct {
		name    string // description of this test case
		grade   *domain.Grade
		wantErr bool
	}{
		{
			name:    "should create grade with valid credentials",
			grade:   grade,
			wantErr: false,
		}, {
			name:    "shouldn't create grade with invalid testID",
			grade:   &domain.Grade{TestID: 0, EnrollmentID: s.Enrollment.ID, Value: 10},
			wantErr: true,
		}, {
			name:    "shouldn't create grade with invalid enrollmentID",
			grade:   &domain.Grade{TestID: s.Test.ID, EnrollmentID: 0, Value: 10},
			wantErr: true,
		}, {
			name:    "shouldn't create grade with invalid value",
			grade:   &domain.Grade{TestID: s.Test.ID, EnrollmentID: 0, Value: 1000},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := r.Create(t.Context(), tt.grade)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Create() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Create() succeeded unexpectedly")
			}
		})
	}
	cleanupGrade(t, grade.ID, s.Test.ID, s.Course.ID, s.Enrollment.ID, s.Teacher.UserID, s.Student.UserID, s.User.ID, s.User2.ID)
}

func TestGradeRepo_GetByID(t *testing.T) {
	s := newGradeScenario(t)
	grade := &domain.Grade{TestID: s.Test.ID, EnrollmentID: s.Enrollment.ID, Value: 10}
	createGradeTest(t, grade)
	r := postgres.NewGradeRepo(testDB, testLogger)
	tests := []struct {
		name    string // description of this test case
		gradeID int64
		want    *domain.Grade
		wantErr bool
	}{
		{
			name:    "should read valid grade IDs",
			gradeID: grade.ID,
			want:    grade,
			wantErr: false,
		}, {
			name:    "should not return anything if grade ID is invalid",
			gradeID: -1,
			want:    nil,
			wantErr: true,
		}, {
			name:    "should not return anythign if grade ID doesn't exist",
			gradeID: grade.ID + 1,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := r.GetByID(t.Context(), tt.gradeID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetByID() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetByID() succeeded unexpectedly")
			}
			require.EqualValues(t, got, tt.want)
		})
	}
	cleanupGrade(t, grade.ID, s.Test.ID, s.Course.ID, s.Enrollment.ID, s.Teacher.UserID, s.Student.UserID, s.User.ID, s.User2.ID)
}

func TestGradeRepo_List(t *testing.T) {
	s := newGradeScenario(t)
	grade := &domain.Grade{TestID: s.Test.ID, EnrollmentID: s.Enrollment.ID, Value: 0}
	createGradeTest(t, grade)
	r := postgres.NewGradeRepo(testDB, testLogger)
	tests := []struct {
		name    string // description of this test case
		filter  ports.GradeFilter
		want    []*domain.Grade
		wantErr bool
	}{
		{
			name:    "should list all grades if no parameters are given",
			filter:  ports.GradeFilter{},
			want:    []*domain.Grade{grade},
			wantErr: false,
		}, {
			name:    "should'nt list any grades if no enrollment IDs match",
			filter:  ports.GradeFilter{EnrollmentID: &s.Test.ID},
			want:    []*domain.Grade{},
			wantErr: false,
		}, {
			name:    "shouldn't list any grades if no test IDs match",
			filter:  ports.GradeFilter{EnrollmentID: &s.Teacher.UserID},
			want:    []*domain.Grade{},
			wantErr: false,
		}, {
			name:    "should list only select grades if right parameters are given",
			filter:  ports.GradeFilter{TestID: &s.Test.ID, EnrollmentID: &s.Enrollment.ID},
			want:    []*domain.Grade{grade},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := r.List(t.Context(), tt.filter)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("List() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("List() succeeded unexpectedly")
			}
			require.Equal(t, len(got), len(tt.want))
			for i := range got {
				require.EqualValues(t, tt.want[i], got[i])
			}
		})
	}
	cleanupGrade(t, grade.ID, s.Test.ID, s.Course.ID, s.Enrollment.ID, s.Teacher.UserID, s.Student.UserID, s.User.ID, s.User2.ID)
}

func TestGradeRepo_Update(t *testing.T) {
	s := newGradeScenario(t)
	grade := &domain.Grade{TestID: s.Test.ID, EnrollmentID: s.Enrollment.ID, Value: 0}
	createGradeTest(t, grade)
	r := postgres.NewGradeRepo(testDB, testLogger)
	tests := []struct {
		name    string // description of this test case
		grade   *domain.Grade
		wantErr bool
	}{
		{
			name:    "shouldn't edit anything if all params are empty",
			grade:   &domain.Grade{ID: grade.ID},
			wantErr: false,
		}, {
			name:    "should try to edit only testID (and throw error because FK constraint)",
			grade:   &domain.Grade{ID: grade.ID, TestID: 999},
			wantErr: true,
		}, {
			name:    "should try to edit only enrollmentID (and throw error because FK constraint)",
			grade:   &domain.Grade{ID: grade.ID, EnrollmentID: 999},
			wantErr: true,
		}, {
			name:    "should edit only value",
			grade:   &domain.Grade{ID: grade.ID, Value: 5},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := r.Update(t.Context(), tt.grade)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Update() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Update() succeeded unexpectedly")
			}
			updatedGrade, err := r.GetByID(t.Context(), grade.ID)
			require.NoError(t, err)
			if tt.grade.Value != 0 {
				require.NotEqual(t, grade.Value, updatedGrade.Value)
			}
			if tt.grade.TestID != 0 {
				require.NotEqual(t, grade.TestID, updatedGrade.TestID)
			}
			if tt.grade.EnrollmentID != 0 {
				require.NotEqual(t, grade.EnrollmentID, updatedGrade.EnrollmentID)
			}
		})
	}
	cleanupGrade(t, grade.ID, s.Test.ID, s.Course.ID, s.Enrollment.ID, s.Teacher.UserID, s.Student.UserID, s.User.ID, s.User2.ID)
}

func TestGradeRepo_Delete(t *testing.T) {
	s := newGradeScenario(t)
	grade := &domain.Grade{TestID: s.Test.ID, EnrollmentID: s.Enrollment.ID, Value: 0}
	createGradeTest(t, grade)
	r := postgres.NewGradeRepo(testDB, testLogger)
	tests := []struct {
		name    string // description of this test case
		gradeID int64
		wantErr bool
	}{
		{
			name:    "should delete valid grade",
			gradeID: grade.ID,
			wantErr: false,
		}, {
			name:    "shouldn't delete invalid gradeID",
			gradeID: -1,
			wantErr: true,
		}, {
			name:    "shouldn't delete gradeID that doesn't exist",
			gradeID: 999,
			wantErr: true,
		}, {
			name:    "should return error if grade was already deleted",
			gradeID: grade.ID,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := r.Delete(t.Context(), tt.gradeID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Delete() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Delete() succeeded unexpectedly")
			}
		})
	}
	cleanupGrade(t, grade.ID, s.Test.ID, s.Course.ID, s.Enrollment.ID, s.Teacher.UserID, s.Student.UserID, s.User.ID, s.User2.ID)
}

func TestGradeRepo_ListAverages(t *testing.T) {
	s := newGradeScenario(t)
	grade := &domain.Grade{TestID: s.Test.ID, EnrollmentID: s.Enrollment.ID, Value: 0}
	createGradeTest(t, grade)
	r := postgres.NewGradeRepo(testDB, testLogger)
	tests := []struct {
		name    string // description of this test case
		filter  ports.ReportCardFilter
		want    []*domain.ReportCard
		wantErr bool
	}{
		{
			name:    "should list all averages",
			filter:  ports.ReportCardFilter{StudentID: &s.Student.UserID},
			want:    []*domain.ReportCard{{StudentID: s.Student.UserID, StudentName: s.User2.Name, CourseID: s.Course.ID, CourseName: s.Course.Name, CourseCode: s.Course.Code, TeacherName: s.User.Name, Average: 0, Situation: domain.SituationStatus("failed")}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := r.ListAverages(t.Context(), tt.filter)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ListAverages() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ListAverages() succeeded unexpectedly")
			}
			require.EqualValues(t, tt.want, got)
		})
	}
	cleanupGrade(t, grade.ID, s.Test.ID, s.Course.ID, s.Enrollment.ID, s.Teacher.UserID, s.Student.UserID, s.User.ID, s.User2.ID)
}

func TestGradeRepo_GetAverage(t *testing.T) {
	s := newGradeScenario(t)
	grade := &domain.Grade{TestID: s.Test.ID, EnrollmentID: s.Enrollment.ID, Value: 10}
	createGradeTest(t, grade)
	r := postgres.NewGradeRepo(testDB, testLogger)
	tests := []struct {
		name      string // description of this test case
		studentID int64
		courseID  int64
		want      float64
		wantErr   bool
	}{
		{
			name:      "should return average",
			studentID: s.Student.UserID,
			courseID:  s.Course.ID,
			want:      10,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := r.GetAverage(t.Context(), tt.studentID, tt.courseID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetAverage() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetAverage() succeeded unexpectedly")
			}
			if tt.want != got {
				t.Errorf("GetAverage() = %v, want %v", got, tt.want)
			}
		})
	}
	cleanupGrade(t, grade.ID, s.Test.ID, s.Course.ID, s.Enrollment.ID, s.Teacher.UserID, s.Student.UserID, s.User.ID, s.User2.ID)
}
