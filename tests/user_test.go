package tests

import (
	"autosalon/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserModel_DefaultRole(t *testing.T) {
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	assert.Equal(t, "", user.Role, "по умолчанию поле Role должно быть пустым до установки")
}
 