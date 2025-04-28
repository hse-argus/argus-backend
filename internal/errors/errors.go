package customerrors

type AlreadyExistsError struct{}

func (aee AlreadyExistsError) Error() string {
	return "Уже существует"
}

type JWTTokenEmpty struct{}

func (jte JWTTokenEmpty) Error() string {
	return "JWT токен пустой"
}

type JWTTokenInvalidError struct{}

func (jte JWTTokenInvalidError) Error() string {
	return "JWT токен невалидный"
}
