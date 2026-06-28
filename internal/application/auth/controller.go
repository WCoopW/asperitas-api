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
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if loginDTO.Username == "" || loginDTO.Password == "" {
		helpers.WriteError(w, http.StatusBadRequest, "username and password are required")
		return
	}
	token, err := c.authService.Login(loginDTO.Username, loginDTO.Password)
	if err != nil {
		if errors.Is(err, user.ErrWrongCredentials) {
			helpers.WriteError(w, http.StatusUnauthorized, "wrong credentials")
			return
		}
		if errors.Is(err, user.ErrNotFound) {
			helpers.WriteError(w, http.StatusNotFound, "user not found")
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var registerDTO RegisterUserDTO
	if err := json.NewDecoder(r.Body).Decode(&registerDTO); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if registerDTO.Username == "" || registerDTO.Password == "" {
		helpers.WriteError(w, http.StatusBadRequest, "username and password are required")
		return
	}
	token, err := c.authService.Register(registerDTO.Username, registerDTO.Password)
	if err != nil {
		if errors.Is(err, user.ErrUsernameTaken) {
			helpers.WriteError(w, http.StatusConflict, err.Error())
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.WriteJSON(w, http.StatusCreated, map[string]string{"token": token})
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		helpers.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	err := c.authService.Logout(token)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.WriteNoContent(w)
}
