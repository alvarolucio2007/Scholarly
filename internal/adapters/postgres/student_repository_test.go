package postgres_test

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/postgres"
	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
)

type studentScenario struct {
	User *domain.User
}

func newStudentScenario(t *testing.T) studentScenario {
	user := &domain.User{Name: "Create Student", Email: "createstudent@email.com", CPF: "00000000000", PasswordHash: []byte("123456789")}
	createUserTest(t, user)
	return studentScenario{User: user}
}
func cleanupStudent(t *testing.T, studentID, userID int64) {
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM students WHERE user_id=$1", studentID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", userID)
		require.NoError(t, err)
	})
}
func createStudentTest(t *testing.T, student *domain.Student) {
	r := postgres.NewStudentRepo(testDB)
	err := r.Create(t.Context(), student)
	require.NoError(t, err)
}

func TestStudentRepo_Create(t *testing.T) {
	s := newStudentScenario(t)
	student := &domain.Student{UserID: s.User.ID, EnrollmentNumber: "TestingCreation"}
	r := postgres.NewStudentRepo(testDB)
	tests := []struct {
		name    string // description of this test case
		student *domain.Student
		wantErr bool
	}{
		{
			name:    "should accept valid userIDs",
			student: student,
			wantErr: false,
		}, {
			name:    "should not accept invalid userIDs",
			student: &domain.Student{UserID: 0, EnrollmentNumber: "Invalid"},
			wantErr: true,
		}, {
			name:    "should not accept duplicated userIDs",
			student: &domain.Student{UserID: s.User.ID, EnrollmentNumber: "DuplicatedCreation"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := r.Create(t.Context(), tt.student)
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
	cleanupStudent(t, student.UserID, s.User.ID)
}

func TestStudentRepo_GetByID(t *testing.T) {
	r := postgres.NewStudentRepo(testDB)
	s := newStudentScenario(t)
	student := domain.Student{UserID: s.User.ID, EnrollmentNumber: "Testing"}
	createStudentTest(t, &student)
	tests := []struct {
		name      string // description of this test case
		studentID int64
		want      domain.Student
		wantErr   bool
	}{
		{
			name:      "should read valid userIDs",
			studentID: s.User.ID,
			want:      student,
			wantErr:   false,
		},
		{
			name:      "should not return anything if userID is invalid",
			studentID: -1,
			want:      domain.Student{},
			wantErr:   true,
		},
		{
			name:      "should not return anything if userID doesn't exist",
			studentID: 100,
			want:      domain.Student{},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := r.GetByID(t.Context(), tt.studentID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetByID() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetByID() succeeded unexpectedly")
			}
			require.Equal(t, tt.want, got)
		})
	}
	cleanupStudent(t, student.UserID, s.User.ID)
}

func TestStudentRepo_GetByEnrollment(t *testing.T) {
	s := newStudentScenario(t)
	student := domain.Student{UserID: s.User.ID, EnrollmentNumber: "Testing"}
	createStudentTest(t, &student)
	tests := []struct {
		name       string // description of this test case
		enrollment string
		want       domain.Student
		wantErr    bool
	}{
		{
			name:       "should return valid user with valid enrollment number",
			enrollment: "Testing",
			want:       student,
			wantErr:    false,
		}, {
			name:       "should not return anything if userID doesn't exist",
			enrollment: "DOESN'T EXIST",
			want:       domain.Student{},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := postgres.NewStudentRepo(testDB)
			got, gotErr := r.GetByEnrollment(t.Context(), tt.enrollment)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetByEnrollment() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetByEnrollment() succeeded unexpectedly")
			}
			require.EqualValues(t, tt.want, got)
		})
	}
	cleanupStudent(t, student.UserID, s.User.ID)
}

func TestStudentRepo_Update(t *testing.T) {
	s := newStudentScenario(t)
	student := domain.Student{UserID: s.User.ID, EnrollmentNumber: "Testing"}
	createStudentTest(t, &student)
	tests := []struct {
		name    string // description of this test case
		student domain.Student
		wantErr bool
	}{
		{
			name:    "shouldn't edit student if all params are empty",
			student: domain.Student{UserID: s.User.ID, EnrollmentNumber: ""},
			wantErr: false,
		},
		{
			name:    "should edit only enrollment number",
			student: domain.Student{UserID: s.User.ID, EnrollmentNumber: "UPDATED"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := postgres.NewStudentRepo(testDB)
			gotErr := r.Update(t.Context(), tt.student)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Update() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Update() succeeded unexpectedly")
			}
			updatedStudent, err := r.GetByID(t.Context(), s.User.ID)
			require.NoError(t, err)
			if tt.student.EnrollmentNumber != "" {
				require.NotEqual(t, student.EnrollmentNumber, updatedStudent.EnrollmentNumber)
			}
		})
	}
	cleanupStudent(t, student.UserID, s.User.ID)
}

func TestStudentRepo_Delete(t *testing.T) {
	s := newStudentScenario(t)
	student1 := &domain.Student{UserID: s.User.ID, EnrollmentNumber: "Testing"}
	createStudentTest(t, student1)
	tests := []struct {
		name      string // description of this test case
		studentID int64
		wantErr   bool
	}{
		{
			name:      "should delete with valid userID",
			studentID: student1.UserID,
			wantErr:   false,
		}, {
			name:      "shouldn't delete with invalid ID",
			studentID: student1.UserID + 100,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := postgres.NewStudentRepo(testDB)
			gotErr := r.Delete(t.Context(), tt.studentID)
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
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM users WHERE id=$1", s.User.ID)
		require.NoError(t, err)
	})
}
