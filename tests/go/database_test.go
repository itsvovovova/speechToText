package main

import (
	"speechToText/src/db"
	"testing"
)

func TestGetTasksWithPagination(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		page      int
		pageSize  int
		expectErr bool
	}{
		{
			name:      "Valid pagination",
			username:  "testuser",
			page:      1,
			pageSize:  10,
			expectErr: false,
		},
		{
			name:      "Invalid page (zero)",
			username:  "testuser",
			page:      0,
			pageSize:  10,
			expectErr: false,
		},
		{
			name:      "Invalid page size (zero)",
			username:  "testuser",
			page:      1,
			pageSize:  0,
			expectErr: false,
		},
		{
			name:      "Empty username",
			username:  "",
			page:      1,
			pageSize:  10,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, total, err := db.GetTasksWithPagination(tt.username, tt.page, tt.pageSize)

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

func TestGetAllTasksWithPagination(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		pageSize  int
		expectErr bool
	}{
		{
			name:      "Valid pagination",
			page:      1,
			pageSize:  10,
			expectErr: false,
		},
		{
			name:      "Invalid page (zero)",
			page:      0,
			pageSize:  10,
			expectErr: false,
		},
		{
			name:      "Invalid page size (zero)",
			page:      1,
			pageSize:  0,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, total, err := db.GetAllTasksWithPagination(tt.page, tt.pageSize)

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

func TestUserExists(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		expectErr bool
	}{
		{
			name:      "Existing user",
			username:  "testuser",
			expectErr: false,
		},
		{
			name:      "Non-existing user",
			username:  "nonexistentuser",
			expectErr: false,
		},
		{
			name:      "Empty username",
			username:  "",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := db.ExistUsername(tt.username)

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
	tests := []struct {
		name      string
		username  string
		password  string
		expectErr bool
	}{
		{
			name:      "Valid user creation",
			username:  "newuser",
			password:  "newpass",
			expectErr: false,
		},
		{
			name:      "Empty username",
			username:  "",
			password:  "newpass",
			expectErr: true,
		},
		{
			name:      "Empty password",
			username:  "newuser",
			password:  "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.AddAuthData(tt.username, tt.password)

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
	tests := []struct {
		name      string
		username  string
		password  string
		expectErr bool
	}{
		{
			name:      "Valid user validation",
			username:  "testuser",
			password:  "testpass",
			expectErr: false,
		},
		{
			name:      "Invalid username",
			username:  "wronguser",
			password:  "testpass",
			expectErr: false,
		},
		{
			name:      "Invalid password",
			username:  "testuser",
			password:  "wrongpass",
			expectErr: false,
		},
		{
			name:      "Empty username",
			username:  "",
			password:  "testpass",
			expectErr: true,
		},
		{
			name:      "Empty password",
			username:  "testuser",
			password:  "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := db.CheckAuthData(tt.username, tt.password)

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
