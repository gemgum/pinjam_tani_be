package service

import (
	"errors"
	"log"
	"pinjamtani_project/app/middlewares"
	"pinjamtani_project/internal/features/users"
	"pinjamtani_project/internal/utils"
	"pinjamtani_project/internal/utils/encrypts"
)

type userService struct {
	userData          users.QueryUserInterface
	hashService       encrypts.HashInterface
	middlewareservice middlewares.MiddlewaresInterface
	accountUtility    utils.AccountUtilityInterface
}

func NewUserService(ud users.QueryUserInterface, hash encrypts.HashInterface, mi middlewares.MiddlewaresInterface, au utils.AccountUtilityInterface) users.ServiceUserInterface {
	return &userService{
		userData:          ud,
		hashService:       hash,
		middlewareservice: mi,
		accountUtility:    au,
	}

}

// LoginAccount implements users.ServiceUserInterface.
func (u *userService) LoginAccount(email string, password string) (data *users.User, token string, err error) {
	data, err = u.userData.AccountByEmail(email)
	if err != nil {
		log.Println("Error logging in:", err)
		return nil, "", err
	}

	isLoginValid := u.hashService.CheckPasswordHash(data.Password, password)
	if !isLoginValid {
		return nil, "", errors.New("email atau password tidak sesuai")
	}

	token, errJWT := u.middlewareservice.CreateToken(int(data.UserID))
	if errJWT != nil {
		log.Println("Error creating token:", errJWT)
		return nil, "", errJWT
	}
	return data, token, nil
}

// RegistrasiAccount implements users.ServiceUserInterface.
func (u *userService) RegistrasiAccount(accounts users.User) (uint, error) {
	if accounts.UserName == "" || accounts.Email == "" || accounts.Password == "" || accounts.PhoneNumber == "" || accounts.Address == "" {
		return 0, errors.New("nama/email/password/phone/address tidak boleh kosong")
	}

	if err := u.accountUtility.EmailValidator(accounts.Email); err != nil {
		return 0, err
	}
	if err := u.accountUtility.PasswordValidator(accounts.Password); err != nil {
		return 0, err
	}
	if err := u.accountUtility.PhoneNumberValidator(accounts.PhoneNumber); err != nil {
		return 0, err
	}

	// Hash password
	var errHash error
	if accounts.Password, errHash = u.hashService.HashPassword(accounts.Password); errHash != nil {
		return 0, errHash
	}

	id, err := u.userData.CreateAccount(accounts)
	if err != nil {
		log.Println("Error registering account:", err)
		return 0, err
	}

	return id, nil
}

// UpdateProfile implements users.ServiceUserInterface.
func (u *userService) UpdateProfile(userid uint, accounts users.User) error {
	if accounts.UserName == "" || accounts.Email == "" || accounts.Password == "" || accounts.PhoneNumber == "" || accounts.Address == "" {
		return errors.New("nama/email/password/phone/address tidak boleh kosong")
	}
	if err := u.accountUtility.EmailValidator(accounts.Email); err != nil {
		return err
	}
	if err := u.accountUtility.PasswordValidator(accounts.Password); err != nil {
		return err
	}
	if err := u.accountUtility.PhoneNumberValidator(accounts.PhoneNumber); err != nil {
		return err
	}

	// Hash password
	var errHash error
	if accounts.Password, errHash = u.hashService.HashPassword(accounts.Password); errHash != nil {
		return errHash
	}

	err := u.userData.UpdateAccount(userid, accounts)
	if err != nil {
		log.Println("Error updating profile:", err)
		return err
	}

	return nil
}

func (us *userService) DeleteAccount(userid uint) error {
	err := us.userData.DeleteAccount(userid)
	if err != nil {
		log.Println("Error deleting account:", err)
		return err
	}
	return nil
}

func (us *userService) GetProfile(userid uint) (*users.User, error) {
	profile, err := us.userData.AccountById(userid)
	if err != nil {
		log.Println("Error fetching profile:", err)
		return nil, err
	}
	return profile, nil
}

func (us *userService) LogoutAccount(token string) error {
	err := us.middlewareservice.InvalidateToken(token)
	if err != nil {
		log.Println("Error logging out:", err)
		return err
	}
	return nil
}
