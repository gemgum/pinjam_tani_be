package users

type User struct {
	UserID      uint
	Images      string
	UserName    string
	Email       string
	Password    string
	PhoneNumber string
	Address     string
}

type DataUserInterface interface {
	CreateAccount(account User) (uint, error)
	AccountByEmail(email string) (*User, error)
	AccountById(userid uint) (*User, error)
	UpdateAccount(userid uint, account User) error
	DeleteAccount(userid uint) error
}

type ServiceUserInterface interface {
	RegistrasiAccount(accounts User) (uint, error)
	LoginAccount(email string, password string) (data *User, token string, err error)
	LogoutAccount(token string) error
	GetProfile(userid uint) (data *User, err error)
	UpdateProfile(userid uint, accounts User) error
	DeleteAccount(userid uint) error
}
