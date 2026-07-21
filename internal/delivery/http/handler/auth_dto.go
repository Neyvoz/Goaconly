package handler

// RegisterRequest — тело запроса на регистрацию нового пользователя.
type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	CompanyName string `json:"company_name"`
}

// RegisterResponse — публичное представление созданного пользователя.
// Намеренно не содержит PasswordHash и других внутренних полей —
// это отдельный тип от domain.User специально для того, чтобы
// исключить случайную утечку чувствительных данных наружу.
type RegisterResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	CompanyName string `json:"company_name"`
}

// LoginRequest — тело запроса на вход в систему.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse — ответ на успешный логин.
// RefreshToken сюда намеренно не входит — он передаётся через
// HttpOnly cookie, а не в теле JSON, чтобы исключить доступ
// к нему со стороны клиентского JavaScript (защита от XSS-кражи токена).
type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

// RefreshResponse — ответ на обновление токенов.
// Аналогично LoginResponse: новый refresh-токен уходит только в cookie.
type RefreshResponse struct {
	// refresh_token НЕ включаем сюда — он уйдёт через HttpOnly cookie, не JSON
	AccessToken string `json:"access_token"`
}
