package auth

type TokenInfoGoogle struct {
	Name    string
	Email   string
	Picture string
	Aud     string
}

type ErrorGoogle struct {
	Error struct {
		Code    int
		Message string
		Status  string
	}
}

type UserAuth struct {
	ID    int
	Email string
	Role  int
}
