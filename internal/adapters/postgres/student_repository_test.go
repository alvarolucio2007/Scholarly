package postgres

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
)

func createStudentTest(t *testing.T, student *domain.Student) {
	r := NewStudentRepo(testDB)
	err := r.Create(t.Context(), student)
	require.NoError(t, err)
}

func TestStudentRepo_Create(t *testing.T) {
	user1 := &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, user1)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM students WHERE user_id=$1", user1.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user1.ID)
		require.NoError(t, err)
	})
	tests := []struct {
		name    string // description of this test case
		student *domain.Student
		wantErr bool
	}{
		{
			name:    "should accept valid userIDs",
			student: &domain.Student{UserID: user1.ID, EnrollmentNumber: "TestingCreation"},
			wantErr: false,
		}, {
			name:    "should not accept invalid userIDs",
			student: &domain.Student{UserID: 0, EnrollmentNumber: "Invalid"},
			wantErr: true,
		}, {
			name:    "should not accept duplicated userIDs",
			student: &domain.Student{UserID: user1.ID, EnrollmentNumber: "DuplicatedCreation"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewStudentRepo(testDB)
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
}

func TestStudentRepo_GetByID(t *testing.T) {
	user1 := &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, user1)
	student1 := &domain.Student{UserID: user1.ID, EnrollmentNumber: "Testing"}
	createStudentTest(t, student1)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM students WHERE user_id=$1", user1.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user1.ID)
		require.NoError(t, err)
	})
	tests := []struct {
		name      string // description of this test case
		studentID int64
		want      *domain.Student
		wantErr   bool
	}{
		{
			name:      "should read valid userIDs",
			studentID: user1.ID,
			want:      student1,
			wantErr:   false,
		},
		{
			name:      "should not return anything if userID is invalid",
			studentID: -1,
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "should not return anything if userID doesn't exist",
			studentID: 100,
			want:      nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewStudentRepo(testDB)
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
			require.EqualValues(t, got.UserID, tt.want.UserID)
			require.EqualValues(t, got.EnrollmentNumber, tt.want.EnrollmentNumber)
		})
	}
}

func TestStudentRepo_GetByEnrollment(t *testing.T) {
	user1 := &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, user1)
	student1 := &domain.Student{UserID: user1.ID, EnrollmentNumber: "Testing"}
	createStudentTest(t, student1)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM students WHERE user_id=$1", user1.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user1.ID)
		require.NoError(t, err)
	})
	tests := []struct {
		name       string // description of this test case
		enrollment string
		want       *domain.Student
		wantErr    bool
	}{
		{
			name:       "should return valid user with valid enrollment number",
			enrollment: "Testing",
			want:       student1,
			wantErr:    false,
		}, {
			name:       "should not return anything if userID doesn't exist",
			enrollment: "DOESN'T EXIST",
			want:       nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewStudentRepo(testDB)
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
			require.EqualValues(t, got.UserID, tt.want.UserID)
			require.EqualValues(t, got.EnrollmentNumber, tt.want.EnrollmentNumber)
		})
	}
}

func TestStudentRepo_Update(t *testing.T) {
	user1 := &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, user1)
	student1 := &domain.Student{UserID: user1.ID, EnrollmentNumber: "Testing"}
	createStudentTest(t, student1)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM students WHERE user_id=$1", user1.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user1.ID)
		require.NoError(t, err)
	})
	tests := []struct {
		name    string // description of this test case
		student *domain.Student
		wantErr bool
	}{
		{
			name:    "shouldn't edit student if all params are empty",
			student: &domain.Student{UserID: user1.ID, EnrollmentNumber: ""},
			wantErr: false,
		},
		{
			name:    "should edit only enrollment number",
			student: &domain.Student{UserID: user1.ID, EnrollmentNumber: "UPDATED"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewStudentRepo(testDB)
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
			updatedStudent, err := r.GetByID(t.Context(), user1.ID)
			require.NoError(t, err)
			if tt.student.EnrollmentNumber != "" {
				require.NotEqual(t, student1.EnrollmentNumber, updatedStudent.EnrollmentNumber)
			}
		})
	}
}
