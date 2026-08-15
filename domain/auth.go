package domain

type RegisterUserReq struct {
	Name     string
	Email    string
	Password string
}

type LoginUserReq struct {
	Email    string
	Password string
}

type LoginUserRes struct {
	AccessToken string
}
