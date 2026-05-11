package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "reddit/internal/domain/auth"
	"reddit/internal/domain/user"
	"reddit/pkg/helpers"

	"go.uber.org/zap"
)

type AuthController struct {
	authService domain.AuthService
	logger      *zap.SugaredLogger
}

func New(authService domain.AuthService, logger *zap.SugaredLogger) *AuthController {
	return &AuthController{authService: authService, logger: logger}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var loginDTO UserLoginDTO
	if err := json.NewDecoder(r.Body).Decode(&loginDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if loginDTO.Username == "" || loginDTO.Password == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}
	token, err := c.authService.Login(loginDTO.Username, loginDTO.Password)
	if err != nil {
		if errors.Is(err, user.ErrWrongCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var registerDTO RegisterUserDTO
	if err := json.NewDecoder(r.Body).Decode(&registerDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if registerDTO.Username == "" || registerDTO.Password == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}
	token, err := c.authService.Register(registerDTO.Username, registerDTO.Password)
	if err != nil {
		if errors.Is(err, user.ErrUsernameTaken) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	helpers.WriteJSON(w, http.StatusCreated, map[string]string{"token": token})
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "token is required", http.StatusUnauthorized)
		return
	}
	err := c.authService.Logout(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}
