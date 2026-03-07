package usecase_test

import (
	"encoding/json"
	"testing"

	"github.com/minguu42/harmattan/internal/api/handler"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthentication_SignUp(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Users{
			{
				ID:             testUserID,
				Email:          "user1@dummy.invalid",
				HashedPassword: "password",
			},
		},
	}))

	t.Run("ok", func(t *testing.T) {
		resp := doRequest(t, "POST", "/sign-up", `{"email": "newuser@dummy.invalid", "password": "Password123!"}`)
		defer resp.Body.Close()

		var body map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))

		assert.Equal(t, 200, resp.StatusCode)
		assert.NotEmpty(t, body["id_token"])
	})
	t.Run("duplicate_email", func(t *testing.T) {
		resp := doRequest(t, "POST", "/sign-up", `{"email": "user1@dummy.invalid", "password": "Password123!"}`)
		defer resp.Body.Close()

		assert.Equal(t, 409, resp.StatusCode)
		assertJSONEqual(t, resp, handler.ErrorResponse{Code: 409, Message: "そのメールアドレスは既に使用されています"})
	})
}

func TestAuthentication_SignIn(t *testing.T) {
	hashed, err := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
	require.NoError(t, err)

	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Users{
			{
				ID:             testUserID,
				Email:          "user1@dummy.invalid",
				HashedPassword: string(hashed),
			},
		},
	}))

	t.Run("ok", func(t *testing.T) {
		resp := doRequest(t, "POST", "/sign-in", `{"email": "user1@dummy.invalid", "password": "Password123!"}`)
		defer resp.Body.Close()

		var body map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))

		assert.Equal(t, 200, resp.StatusCode)
		assert.NotEmpty(t, body["id_token"])
	})
	t.Run("invalid_email", func(t *testing.T) {
		resp := doRequest(t, "POST", "/sign-in", `{"email": "unknown@dummy.invalid", "password": "Password123!"}`)
		defer resp.Body.Close()

		assert.Equal(t, 400, resp.StatusCode)
		assertJSONEqual(t, resp, handler.ErrorResponse{Code: 400, Message: "メールアドレスかパスワードに誤りがあります"})
	})
	t.Run("invalid_password", func(t *testing.T) {
		resp := doRequest(t, "POST", "/sign-in", `{"email": "user1@dummy.invalid", "password": "WrongPassword1!"}`)
		defer resp.Body.Close()

		assert.Equal(t, 400, resp.StatusCode)
		assertJSONEqual(t, resp, handler.ErrorResponse{Code: 400, Message: "メールアドレスかパスワードに誤りがあります"})
	})
}
