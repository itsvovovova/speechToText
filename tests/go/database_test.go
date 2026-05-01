package main

import (
	"testing"
)

func TestGetTasksWithPagination(t *testing.T) {
	if testStore == nil {
		t.Skip("DB not available")
	}
	tests := []struct {
		name      string
		username  string
		page      int
		pageSize  int
		expectErr bool
	}{
		{name: "Valid pagination", username: "testuser", page: 1, pageSize: 10},
		{name: "Invalid page (zero)", username: "testuser", page: 0, pageSize: 10},
		{name: "Invalid page size (zero)", username: "testuser", page: 1, pageSize: 0},
		{name: "Empty username", username: "", page: 1, pageSize: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, total, err := testStore.GetTasksWithPagination(tt.username, tt.page, tt.pageSize)
			if tt.expectErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tasks == nil {
				t.Errorf("Tasks should not be nil")
			}
			if total < 0 {
				t.Errorf("Total should not be negative")
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	if testStore == nil {
		t.Skip("DB not available")
	}
	tests := []struct {
		name      string
		taskID    string
		username  string
		expectErr bool
	}{
		{name: "Non-existent task", taskID: "non-existent-id", username: "testuser", expectErr: true},
		{name: "Empty task ID", taskID: "", username: "testuser", expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testStore.DeleteTask(tt.taskID, tt.username)
			if tt.expectErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestUserExists(t *testing.T) {
	if testStore == nil {
		t.Skip("DB not available")
	}
	tests := []struct {
		name      string
		username  string
		expectErr bool
	}{
		{name: "Existing user", username: "testuser"},
		{name: "Non-existing user", username: "nonexistentuser"},
		{name: "Empty username", username: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := testStore.ExistUsername(tt.username)
			if tt.expectErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if exists && tt.username == "" {
				t.Errorf("Empty username should not exist")
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	if testStore == nil {
		t.Skip("DB not available")
	}
	tests := []struct {
		name      string
		username  string
		password  string
		expectErr bool
	}{
		{name: "Valid user creation", username: "newuser", password: "newpass"},
		{name: "Empty username", username: "", password: "newpass", expectErr: true},
		{name: "Empty password", username: "newuser", password: "", expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testStore.AddAuthData(tt.username, tt.password)
			if tt.expectErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidateUser(t *testing.T) {
	if testStore == nil {
		t.Skip("DB not available")
	}
	tests := []struct {
		name      string
		username  string
		password  string
		expectErr bool
	}{
		{name: "Valid user validation", username: "testuser", password: "testpass"},
		{name: "Invalid username", username: "wronguser", password: "testpass"},
		{name: "Invalid password", username: "testuser", password: "wrongpass"},
		{name: "Empty username", username: "", password: "testpass", expectErr: true},
		{name: "Empty password", username: "testuser", password: "", expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := testStore.CheckAuthData(tt.username, tt.password)
			if tt.expectErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if valid && (tt.username == "" || tt.password == "") {
				t.Errorf("Empty credentials should not be valid")
			}
		})
	}
}
