package user

import "testing"

func TestUserCreate(t *testing.T) {
	_, err := Create(
		"user_name",
		"user_surname",
		"+000000000000",
		"test_email@test.test",
		"password",
	)

	if err != nil {
		t.Error(err)
	}
}
