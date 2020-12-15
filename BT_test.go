package main

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestDB(t *testing.T) {
	var testUser User
	mutex.Lock()
	{
		db := OpenDB()
		defer db.Close()
		testUser, _ = GetUser(db, "testUser")
	}
	mutex.Unlock()
	if testUser.First != "TestFirstName" || testUser.Last != "TestLastName" || bcrypt.CompareHashAndPassword(testUser.Password, []byte("TestPassword")) != nil || testUser.Email != "testuser@test.com" {
		t.Errorf("GetUser has obtained incorrect results, got: %v, want: %v.", testUser, User{"testUser", []byte("TestPassword"),"TestFirstName", "TestLastName", "testuser@test.com"})
	}
}
