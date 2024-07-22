package handler

import (
	"log"
	"net/http"
	"pinjamtani_project/app/middlewares"
	"pinjamtani_project/internal/features/users"
	"pinjamtani_project/internal/utils/responses"
	"strings"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService users.ServiceUserInterface
}

func NewUserHandler(us users.ServiceUserInterface) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

// Register handles user registration
func (uh *UserHandler) Register(c echo.Context) error {
	newUser := UserRequest{}
	if errBind := c.Bind(&newUser); errBind != nil {
		log.Printf("Register: Error binding data: %v", errBind)
		return c.JSON(http.StatusBadRequest, responses.JSONWebResponse(http.StatusBadRequest, "error", "error binding data: "+errBind.Error(), nil))
	}

	file, err := c.FormFile("images")
	var imageURL string
	if err == nil {
		src, err := file.Open()
		if err != nil {
			log.Printf("Register: Error opening image file: %v", err)
			return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "error", "gagal membuka file gambar: "+err.Error(), nil))
		}
		defer src.Close()

		imageURL, err = newUser.uploadToCloudinary(src, file.Filename)
		if err != nil {
			log.Printf("Register: Error uploading image: %v", err)
			return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "error", "gagal mengunggah gambar: "+err.Error(), nil))
		}
	}

	dataUser := users.User{
		Images:      imageURL,
		UserName:    newUser.UserName,
		Email:       newUser.Email,
		Password:    newUser.Password,
		PhoneNumber: newUser.PhoneNumber,
		Address:     newUser.Address,
	}

	userID, errInsert := uh.userService.RegistrasiAccount(dataUser)
	if errInsert != nil {
		log.Printf("Register: Error registering user: %v", errInsert)
		if strings.Contains(errInsert.Error(), "validation") {
			return c.JSON(http.StatusBadRequest, responses.JSONWebResponse(http.StatusBadRequest, "failed", "user registration failed: "+errInsert.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "failed", "user registration failed: "+errInsert.Error(), nil))
	}

	userResponse := UserResponse{
		ID:          userID,
		Images:      imageURL,
		UserName:    newUser.UserName,
		Email:       newUser.Email,
		PhoneNumber: newUser.PhoneNumber,
		Address:     newUser.Address,
	}

	return c.JSON(http.StatusCreated, responses.JSONWebResponse(http.StatusCreated, "success", "user registration successful", userResponse))
}

// Login handles user login
func (uh *UserHandler) Login(c echo.Context) error {
	loginReq := LoginRequest{}
	if errBind := c.Bind(&loginReq); errBind != nil {
		log.Printf("Login: Error binding data: %v", errBind)
		return c.JSON(http.StatusBadRequest, responses.JSONWebResponse(http.StatusBadRequest, "failed", "error binding data: "+errBind.Error(), nil))
	}

	_, token, err := uh.userService.LoginAccount(loginReq.Email, loginReq.Password)
	if err != nil {
		log.Printf("Login: User login failed: %v", err)
		return c.JSON(http.StatusUnauthorized, responses.JSONWebResponse(http.StatusUnauthorized, "failed", "user login failed: "+err.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.JSONWebResponse(http.StatusOK, "success", "user login successful", echo.Map{"token": token}))
}

// Update handles user profile updates
func (uh *UserHandler) Update(c echo.Context) error {
	userID := middlewares.NewMiddlewares().ExtractTokenUserId(c)
	if userID == 0 {
		log.Println("Update: Unauthorized access")
		return c.JSON(http.StatusUnauthorized, responses.JSONWebResponse(http.StatusUnauthorized, "error", "unauthorized", nil))
	}

	newUser := UserRequest{}
	if errBind := c.Bind(&newUser); errBind != nil {
		log.Printf("Update: Error binding data: %v", errBind)
		return c.JSON(http.StatusBadRequest, responses.JSONWebResponse(http.StatusBadRequest, "error", "error binding data: "+errBind.Error(), nil))
	}

	file, err := c.FormFile("images")
	var imageURL string
	if err == nil {
		src, err := file.Open()
		if err != nil {
			log.Printf("Update: Error opening image file: %v", err)
			return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "error", "gagal membuka file gambar: "+err.Error(), nil))
		}
		defer src.Close()

		imageURL, err = newUser.uploadToCloudinary(src, file.Filename)
		if err != nil {
			log.Printf("Update: Error uploading image: %v", err)
			return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "error", "gagal mengunggah gambar: "+err.Error(), nil))
		}
	}

	dataUser := users.User{
		Images:      imageURL,
		UserName:    newUser.UserName,
		Email:       newUser.Email,
		Password:    newUser.Password,
		PhoneNumber: newUser.PhoneNumber,
		Address:     newUser.Address,
	}

	if errInsert := uh.userService.UpdateProfile(uint(userID), dataUser); errInsert != nil {
		log.Printf("Update: Error updating account: %v", errInsert)
		if strings.Contains(errInsert.Error(), "validation") {
			return c.JSON(http.StatusBadRequest, responses.JSONWebResponse(http.StatusBadRequest, "Failed", "failed to update account: "+errInsert.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "failed", "failed to update account: "+errInsert.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.JSONWebResponse(http.StatusOK, "success", "successfully updated account", nil))
}

// Delete handles account deletion
func (uh *UserHandler) Delete(c echo.Context) error {
	userID := middlewares.NewMiddlewares().ExtractTokenUserId(c)
	if userID == 0 {
		log.Println("Delete: Unauthorized access")
		return c.JSON(http.StatusUnauthorized, responses.JSONWebResponse(http.StatusUnauthorized, "error", "unauthorized", nil))
	}

	if errDelete := uh.userService.DeleteAccount(uint(userID)); errDelete != nil {
		log.Printf("Delete: Error deleting account: %v", errDelete)
		return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "failed", "failed to delete account: "+errDelete.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.JSONWebResponse(http.StatusOK, "success", "successfully deleted account", nil))
}

// GetProfile handles fetching user profile
func (uh *UserHandler) GetProfile(c echo.Context) error {
	userID := middlewares.NewMiddlewares().ExtractTokenUserId(c)
	if userID == 0 {
		log.Println("GetProfile: Unauthorized access")
		return c.JSON(http.StatusUnauthorized, responses.JSONWebResponse(http.StatusUnauthorized, "error", "unauthorized", nil))
	}

	profile, err := uh.userService.GetProfile(uint(userID))
	if err != nil {
		log.Printf("GetProfile: Error fetching user profile: %v", err)
		return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "failed", "get user profile failed: "+err.Error(), nil))
	}

	userResponse := UserResponse{
		ID:          profile.UserID,
		Images:      profile.Images,
		UserName:    profile.UserName,
		Email:       profile.Email,
		PhoneNumber: profile.PhoneNumber,
		Address:     profile.Address,
	}

	return c.JSON(http.StatusOK, responses.JSONWebResponse(http.StatusOK, "success", "get user profile successful", userResponse))
}

// Logout handles user logout
func (uh *UserHandler) Logout(c echo.Context) error {
	userToken := c.Request().Header.Get("authorization")
	if userToken == "" {
		log.Println("Logout: No token provided")
		return c.JSON(http.StatusBadRequest, responses.JSONWebResponse(http.StatusBadRequest, "error", "no token provided", nil))
	}

	userToken = strings.TrimPrefix(userToken, "bearer ")

	if err := uh.userService.LogoutAccount(userToken); err != nil {
		log.Printf("Logout: Error logging out: %v", err)
		return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "failed", "logout failed: "+err.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.JSONWebResponse(http.StatusOK, "success", "user logged out successfully", nil))
}
