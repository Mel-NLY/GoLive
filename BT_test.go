package main

import (
	"testing"
	"GoLive/pkgs"
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

	tests := []struct{
		name string
		pass string
		valid bool
	}{
		{
			"NoCharacters",
			"",
			false,
		},
		{
			"JustEmptyStringAndWhitespace",
			" \n\t\r\v\f ",
			false,
		},
		{
			"MixtureOfEmptyStringAndWhitespace",
			"U u\n1\t?\r1\v2\f34",
			false,
		},
		{
			"MissingUpperCaseString",
			"uu1?1234",
			false,
		},
		{
			"MissingLowerCaseString",
			"UU1?1234",
			false,
		},
		{
			"MissingNumber",
			"Uua?aaaa",
			false,
		},
		{
			"MissingSymbol",
			"Uu101234",
			false,
		},
		{
			"LessThanRequiredMinimumLength",
			"Uu1?123",
			false,
		},
		{
			"ValidPassword",
			"Uu1?1234",
			true,
		},
	}
 
	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			if c.valid != pkgs.Password(c.pass) {
				t.Fatal("invalid password")
			}
		})
	}
}
