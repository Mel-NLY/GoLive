package main

import (
	"GoLive/pkgs"
	"errors"
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

	tests := []struct {
		name  string
		pass  string
		valid bool
		err   error
	}{
		{
			"NoCharacters",
			"",
			false,
			errors.New("No whitespaces/empty inputs allowed"),
		},
		{
			"JustEmptyStringAndWhitespace",
			" \n\t\r\v\f ",
			false,
			errors.New("No whitespaces/empty inputs allowed"),
		},
		{
			"MixtureOfEmptyStringAndWhitespace",
			"U u\n1\t?\r1\v2\f34",
			false,
			errors.New("No whitespaces/empty inputs allowed"),
		},
		{
			"MissingUpperCaseString",
			"uu1?1234",
			false,
			errors.New("No Uppercase"),
		},
		{
			"MissingLowerCaseString",
			"UU1?1234",
			false,
			errors.New("No Lowercase"),
		},
		{
			"MissingNumber",
			"Uua?aaaa",
			false,
			errors.New("No Number"),
		},
		{
			"MissingSymbol",
			"Uu101234",
			false,
			errors.New("No Symbol"),
		},
		{
			"LessThanRequiredMinimumLength",
			"Uu1?123",
			false,
			errors.New("Lesser than 8 characters"),
		},
		{
			"ValidPassword",
			"Uu1?1234",
			true,
			nil,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			p, err := pkgs.Password(c.pass)
			if c.valid != p && err != nil {
				t.Fatal("invalid password")
			}
		})
	}
}
