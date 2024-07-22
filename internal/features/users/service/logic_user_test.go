package service_test

import (
	"errors"
	users "projectBE23/internal/features/users"
	"projectBE23/internal/features/users/service"
	"projectBE23/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	mockUserData := new(mocks.DataUserInterface)
	mockHashService := new(mocks.HashInterface)
	mockMiddleware := new(mocks.MiddlewaresInterface)
	mockutils := new(mocks.AccountUtilityInterface)
	userService := service.New(mockUserData, mockHashService, mockMiddleware, mockutils)

	t.Run("success", func(t *testing.T) {
		user := users.User{
			UserName:    "testuser",
			Email:       "testuser@example.com",
			Password:    "password123",
			PhoneNumber: "1234567890",
			Address:     "123 Test St",
		}

		mockutils.On("EmailValidator", user.Email).Return(nil).Once()
		mockutils.On("PasswordValidator", user.Password).Return(nil).Once()
		mockutils.On("PhoneNumberValidator", user.PhoneNumber).Return(nil).Once()
		mockHashService.On("HashPassword", user.Password).Return("hashedpassword", nil).Once()
		mockUserData.On("CreateAccount", mock.AnythingOfType("users.User")).Return(uint(1), nil).Once()

		id, err := userService.RegistrasiAccount(user)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), id)

		mockutils.AssertExpectations(t)
		mockHashService.AssertExpectations(t)
		mockUserData.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		user := users.User{}

		id, err := userService.RegistrasiAccount(user)

		assert.Error(t, err)
		assert.Equal(t, uint(0), id)
		assert.Equal(t, "nama/email/password/phone/address tidak boleh kosong", err.Error())
	})

	t.Run("hash error", func(t *testing.T) {
		user := users.User{
			UserName:    "testuser",
			Email:       "test@example.com",
			Password:    "password123",
			PhoneNumber: "1234567890",
			Address:     "Test Address",
		}

		mockutils.On("EmailValidator", user.Email).Return(nil).Once()
		mockutils.On("PasswordValidator", user.Password).Return(nil).Once()
		mockutils.On("PhoneNumberValidator", user.PhoneNumber).Return(nil).Once()
		mockHashService.On("HashPassword", user.Password).Return("", errors.New("hash error")).Once()

		id, err := userService.RegistrasiAccount(user)

		assert.Error(t, err)
		assert.Equal(t, uint(0), id)
		assert.Equal(t, "hash error", err.Error())

		mockutils.AssertExpectations(t)
		mockHashService.AssertExpectations(t)
	})

	t.Run("email validation error", func(t *testing.T) {
		user := users.User{
			UserName:    "testuser",
			Email:       "invalid-email",
			Password:    "password123",
			PhoneNumber: "1234567890",
			Address:     "Test Address",
		}

		mockutils.On("EmailValidator", user.Email).Return(errors.New("invalid email")).Once()

		id, err := userService.RegistrasiAccount(user)

		assert.Error(t, err)
		assert.Equal(t, uint(0), id)
		assert.Equal(t, "invalid email", err.Error())

		mockutils.AssertExpectations(t)
	})

	t.Run("password validation error", func(t *testing.T) {
		user := users.User{
			UserName:    "testuser",
			Email:       "test@example.com",
			Password:    "weak",
			PhoneNumber: "1234567890",
			Address:     "Test Address",
		}

		mockutils.On("EmailValidator", user.Email).Return(nil).Once()
		mockutils.On("PasswordValidator", user.Password).Return(errors.New("weak password")).Once()

		id, err := userService.RegistrasiAccount(user)

		assert.Error(t, err)
		assert.Equal(t, uint(0), id)
		assert.Equal(t, "weak password", err.Error())

		mockutils.AssertExpectations(t)
	})

	t.Run("phone number validation error", func(t *testing.T) {
		user := users.User{
			UserName:    "testuser",
			Email:       "test@example.com",
			Password:    "password123",
			PhoneNumber: "invalid-phone",
			Address:     "Test Address",
		}

		mockutils.On("EmailValidator", user.Email).Return(nil).Once()
		mockutils.On("PasswordValidator", user.Password).Return(nil).Once()
		mockutils.On("PhoneNumberValidator", user.PhoneNumber).Return(errors.New("invalid phone number")).Once()

		id, err := userService.RegistrasiAccount(user)

		assert.Error(t, err)
		assert.Equal(t, uint(0), id)
		assert.Equal(t, "invalid phone number", err.Error())

		mockutils.AssertExpectations(t)
	})

	t.Run("regitration account error", func(t *testing.T) {
		user := users.User{
			UserName:    "testuser",
			Email:       "test@example.com",
			Password:    "password123",
			PhoneNumber: "1234567890",
			Address:     "Test Address",
		}

		mockutils.On("EmailValidator", user.Email).Return(nil).Once()
		mockutils.On("PasswordValidator", user.Password).Return(nil).Once()
		mockutils.On("PhoneNumberValidator", user.PhoneNumber).Return(nil).Once()
		mockHashService.On("HashPassword", user.Password).Return("hashedpassword", nil).Once()
		mockUserData.On("CreateAccount", mock.AnythingOfType("users.User")).Return(uint(0), errors.New("registration error")).Once()

		id, err := userService.RegistrasiAccount(user)

		assert.Error(t, err)
		assert.Equal(t, uint(0), id)
		assert.Equal(t, "registration error", err.Error())

		mockutils.AssertExpectations(t)
		mockHashService.AssertExpectations(t)
		mockUserData.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	mockUserData := new(mocks.DataUserInterface)
	mockHashService := new(mocks.HashInterface)
	mockMiddleware := new(mocks.MiddlewaresInterface)
	mockutils := new(mocks.AccountUtilityInterface)
	userService := service.New(mockUserData, mockHashService, mockMiddleware, mockutils)

	t.Run("success", func(t *testing.T) {
		email := "johndoe@example.com"
		password := "password123"
		hashedPassword := "hashedpassword"
		user := &users.User{
			UserID:      1,
			Images:      "http://images.example.com",
			UserName:    "John Doe",
			Email:       email,
			Password:    hashedPassword,
			PhoneNumber: "08123456789",
			Address:     "Some Address",
		}

		token := "sometoken"

		mockUserData.On("AccountByEmail", email).Return(user, nil).Once()
		mockHashService.On("CheckPasswordHash", hashedPassword, password).Return(true).Once()
		mockMiddleware.On("CreateToken", 1).Return("sometoken", nil).Once()

		returnedUser, returnedToken, err := userService.LoginAccount(email, password)

		assert.NoError(t, err)
		assert.Equal(t, user, returnedUser)
		assert.Equal(t, token, returnedToken)
		mockUserData.AssertExpectations(t)
		mockHashService.AssertExpectations(t)
		mockMiddleware.AssertExpectations(t)
	})

	t.Run("account not found", func(t *testing.T) {
		email := "johndoe@example.com"
		password := "password123"

		mockUserData.On("AccountByEmail", email).Return(nil, errors.New("account not found")).Once()

		returnedUser, returnedToken, err := userService.LoginAccount(email, password)

		assert.Error(t, err)
		assert.Equal(t, "account not found", err.Error())
		assert.Nil(t, returnedUser)
		assert.Empty(t, returnedToken)
		mockUserData.AssertExpectations(t)
	})

	t.Run("password verification error", func(t *testing.T) {
		email := "johndoe@example.com"
		password := "password123"
		hashedPassword := "hashedpassword"
		user := &users.User{
			UserName:    "John Doe",
			Email:       email,
			Password:    hashedPassword,
			PhoneNumber: "08123456789",
			Address:     "Some Address",
		}

		mockUserData.On("AccountByEmail", email).Return(user, nil).Once()
		mockHashService.On("CheckPasswordHash", hashedPassword, password).Return(false).Once()

		returnedUser, returnedToken, err := userService.LoginAccount(email, password)

		assert.Error(t, err)
		assert.Equal(t, "email atau password tidak sesuai", err.Error())
		assert.Nil(t, returnedUser)
		assert.Empty(t, returnedToken)
		mockUserData.AssertExpectations(t)
		mockHashService.AssertExpectations(t)
	})

	t.Run("jwt authentication error", func(t *testing.T) {
		email := "johndoe@example.com"
		password := "password123"
		hashedPassword := "hashedpassword"
		user := &users.User{
			UserID:      1,
			UserName:    "John Doe",
			Email:       email,
			Password:    hashedPassword,
			PhoneNumber: "08123456789",
			Address:     "Some Address",
		}

		mockUserData.On("AccountByEmail", email).Return(user, nil).Once()
		mockHashService.On("CheckPasswordHash", hashedPassword, password).Return(true).Once()
		mockMiddleware.On("CreateToken", int(user.UserID)).Return("", errors.New("jwt error")).Once()

		returnedUser, returnedToken, err := userService.LoginAccount(email, password)

		assert.Error(t, err)
		assert.Equal(t, "jwt error", err.Error())
		assert.Nil(t, returnedUser)
		assert.Empty(t, returnedToken)
		mockUserData.AssertExpectations(t)
		mockHashService.AssertExpectations(t)
		mockMiddleware.AssertExpectations(t)
	})
}

func TestUpdateProfile(t *testing.T) {
	mockUserData := new(mocks.DataUserInterface)
	mockHashService := new(mocks.HashInterface)
	mockMiddleware := new(mocks.MiddlewaresInterface)
	mockutils := new(mocks.AccountUtilityInterface)
	userService := service.New(mockUserData, mockHashService, mockMiddleware, mockutils)

	t.Run("success", func(t *testing.T) {
		userID := uint(1)
		user := users.User{
			UserName:    "updateduser",
			Email:       "updateduser@example.com",
			Password:    "newpassword123",
			PhoneNumber: "0987654321",
			Address:     "456 Test Ave",
		}

		mockutils.On("EmailValidator", user.Email).Return(nil).Once()
		mockutils.On("PasswordValidator", user.Password).Return(nil).Once()
		mockutils.On("PhoneNumberValidator", user.PhoneNumber).Return(nil).Once()
		mockHashService.On("HashPassword", user.Password).Return("newhashedpassword", nil).Once()
		mockUserData.On("UpdateAccount", userID, mock.AnythingOfType("users.User")).Return(nil).Once()

		err := userService.UpdateProfile(userID, user)
		assert.NoError(t, err)

		mockutils.AssertExpectations(t)
		mockHashService.AssertExpectations(t)
		mockUserData.AssertExpectations(t)
	})

	t.Run("update error", func(t *testing.T) {
		userID := uint(1)
		user := users.User{
			UserName:    "updateduser",
			Email:       "updateduser@example.com",
			Password:    "newpassword123",
			PhoneNumber: "0987654321",
			Address:     "456 Test Ave",
		}

		mockutils.On("EmailValidator", user.Email).Return(nil).Once()
		mockutils.On("PasswordValidator", user.Password).Return(nil).Once()
		mockutils.On("PhoneNumberValidator", user.PhoneNumber).Return(nil).Once()
		mockHashService.On("HashPassword", user.Password).Return("newhashedpassword", nil).Once()
		mockUserData.On("UpdateAccount", userID, mock.AnythingOfType("users.User")).Return(errors.New("update error")).Once()

		err := userService.UpdateProfile(userID, user)
		assert.Error(t, err)
		assert.Equal(t, "update error", err.Error())

		mockutils.AssertExpectations(t)
		mockHashService.AssertExpectations(t)
		mockUserData.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		userID := uint(1)
		user := users.User{}

		err := userService.UpdateProfile(userID, user)

		assert.Error(t, err)
		assert.Equal(t, "nama/email/password/phone/address tidak boleh kosong", err.Error())

		// Ensure that no methods were called on mocks since validation should prevent further processing
		mockUserData.AssertNotCalled(t, "HashPassword")
		mockUserData.AssertNotCalled(t, "UpdateAccount")
		mockHashService.AssertNotCalled(t, "HashPassword")
	})

	t.Run("hash error", func(t *testing.T) {
		userID := uint(1)
		user := users.User{
			UserName:    "updateduser",
			Email:       "updateduser@example.com",
			Password:    "newpassword123",
			PhoneNumber: "0987654321",
			Address:     "456 Test Ave",
		}

		mockutils.On("EmailValidator", user.Email).Return(nil).Once()
		mockutils.On("PasswordValidator", user.Password).Return(nil).Once()
		mockutils.On("PhoneNumberValidator", user.PhoneNumber).Return(nil).Once()
		mockHashService.On("HashPassword", user.Password).Return("", errors.New("hash error")).Once()

		err := userService.UpdateProfile(userID, user)
		assert.Error(t, err)
		assert.Equal(t, "hash error", err.Error())

		mockutils.AssertExpectations(t)
		mockHashService.AssertExpectations(t)
		mockUserData.AssertNotCalled(t, "UpdateAccount", mock.Anything)
	})

	t.Run("email validation error", func(t *testing.T) {
		userID := uint(1)
		user := users.User{
			UserName:    "updateduser",
			Email:       "updateduser@example.com",
			Password:    "newpassword123",
			PhoneNumber: "0987654321",
			Address:     "456 Test Ave",
		}

		mockutils.On("EmailValidator", user.Email).Return(errors.New("invalid email")).Once()

		err := userService.UpdateProfile(userID, user)
		assert.Error(t, err)
		assert.Equal(t, "invalid email", err.Error())

		mockutils.AssertExpectations(t)
		mockHashService.AssertNotCalled(t, "HashPassword")
		mockUserData.AssertNotCalled(t, "UpdateAccount")
	})

	t.Run("password validation error", func(t *testing.T) {
		userID := uint(1)
		user := users.User{
			UserName:    "updateduser",
			Email:       "updateduser@example.com",
			Password:    "newpassword123",
			PhoneNumber: "0987654321",
			Address:     "456 Test Ave",
		}

		mockutils.On("EmailValidator", user.Email).Return(nil).Once()
		mockutils.On("PasswordValidator", user.Password).Return(errors.New("invalid password")).Once()

		err := userService.UpdateProfile(userID, user)
		assert.Error(t, err)
		assert.Equal(t, "invalid password", err.Error())

		mockutils.AssertExpectations(t)
		mockHashService.AssertNotCalled(t, "HashPassword")
		mockUserData.AssertNotCalled(t, "UpdateAccount")
	})

	t.Run("phone number validation error", func(t *testing.T) {
		userID := uint(1)
		user := users.User{
			UserName:    "updateduser",
			Email:       "updateduser@example.com",
			Password:    "newpassword123",
			PhoneNumber: "0987654321",
			Address:     "456 Test Ave",
		}

		mockutils.On("EmailValidator", user.Email).Return(nil).Once()
		mockutils.On("PasswordValidator", user.Password).Return(nil).Once()
		mockutils.On("PhoneNumberValidator", user.PhoneNumber).Return(errors.New("invalid phone number")).Once()

		err := userService.UpdateProfile(userID, user)
		assert.Error(t, err)
		assert.Equal(t, "invalid phone number", err.Error())

		mockutils.AssertExpectations(t)
		mockHashService.AssertNotCalled(t, "HashPassword")
		mockUserData.AssertNotCalled(t, "UpdateAccount")
	})
}

func TestDeleteProfile(t *testing.T) {
	mockUserData := new(mocks.DataUserInterface)
	mockHashService := new(mocks.HashInterface)
	mockMiddleware := new(mocks.MiddlewaresInterface)
	mockutils := new(mocks.AccountUtilityInterface)
	userService := service.New(mockUserData, mockHashService, mockMiddleware, mockutils)

	t.Run("success", func(t *testing.T) {
		userID := uint(1)

		mockUserData.On("DeleteAccount", userID).Return(nil).Once()

		err := userService.DeleteAccount(userID)

		assert.NoError(t, err)
		mockUserData.AssertExpectations(t)
	})

	t.Run("delete error", func(t *testing.T) {
		userID := uint(1)

		mockUserData.On("DeleteAccount", userID).Return(errors.New("delete error")).Once()

		err := userService.DeleteAccount(userID)

		assert.Error(t, err)
		assert.Equal(t, "delete error", err.Error())
		mockUserData.AssertExpectations(t)
	})
}

func TestGetProfile(t *testing.T) {
	mockUserData := new(mocks.DataUserInterface)
	mockHashService := new(mocks.HashInterface)
	mockMiddleware := new(mocks.MiddlewaresInterface)
	mockutils := new(mocks.AccountUtilityInterface)
	userService := service.New(mockUserData, mockHashService, mockMiddleware, mockutils)

	t.Run("success", func(t *testing.T) {
		userID := uint(1)
		user := &users.User{
			Images:      "http:images.com",
			UserName:    "John Doe",
			Email:       "johndoe@example.com",
			PhoneNumber: "08123456789",
			Address:     "Some Address",
		}

		mockUserData.On("AccountById", userID).Return(user, nil).Once()

		returnedUser, err := userService.GetProfile(userID)

		assert.NoError(t, err)
		assert.Equal(t, user, returnedUser)
		mockUserData.AssertExpectations(t)
	})

	t.Run("profile not found", func(t *testing.T) {
		userID := uint(1)

		mockUserData.On("AccountById", userID).Return(nil, errors.New("profile not found")).Once()

		returnedUser, err := userService.GetProfile(userID)

		assert.Error(t, err)
		assert.Equal(t, "profile not found", err.Error())
		assert.Nil(t, returnedUser)
		mockUserData.AssertExpectations(t)
	})
}

func TestLogout(t *testing.T) {
	mockUserData := new(mocks.DataUserInterface)
	mockHashService := new(mocks.HashInterface)
	mockMiddleware := new(mocks.MiddlewaresInterface)
	mockutils := new(mocks.AccountUtilityInterface)
	userService := service.New(mockUserData, mockHashService, mockMiddleware, mockutils)

	t.Run("success", func(t *testing.T) {
		token := "validToken"

		mockMiddleware.On("InvalidateToken", token).Return(nil).Once()

		err := userService.LogoutAccount(token)
		assert.NoError(t, err)

		mockMiddleware.AssertExpectations(t)
	})

	t.Run("invalidate token error", func(t *testing.T) {
		token := "invalidToken"

		mockMiddleware.On("InvalidateToken", token).Return(errors.New("invalidate error")).Once()

		err := userService.LogoutAccount(token)
		assert.Error(t, err)
		assert.Equal(t, "invalidate error", err.Error())

		mockMiddleware.AssertExpectations(t)
	})
}
