package main

import (
	"testing"
)

func TestDB(t *testing.T) {
	//Connection timeout, cannot connect to mysql
	// var testUser User
	// mutex.Lock()
	// {
	// 	db := OpenDB()
	// 	defer db.Close()
	// 	testUser, _ = GetUser(db, "testUser")
	// }
	// mutex.Unlock()
	// if testUser.First != "TestFirstName" || testUser.Last != "TestLastName" || bcrypt.CompareHashAndPassword(testUser.Password, []byte("TestPassword")) != nil || testUser.Email != "testuser@test.com" {
	// 	t.Errorf("GetUser has obtained incorrect results, got: %v, want: %v.", testUser, User{"testUser", []byte("TestPassword"), "TestFirstName", "TestLastName", "testuser@test.com"})
	// }

	// g := Goblin(t)
	// g.Describe("Find Distance", func() {
	// 	g.It("Should obtain distance from A to B ", func() {
	// 		g.Assert(findDistance(1, 2)).Equal(3)
	// 		g.Assert(findDistance(1, 1)).Equal(2)
	// 	})
	// })
}
