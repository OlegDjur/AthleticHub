package controller

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"workout/internal/adapters/postgres"
	"workout/internal/dto"
	"workout/internal/service/auth"
)

func (h *Handler) Login(e echo.Context) error {
	var request dto.LoginRequest

	if err := e.Bind(&request); err != nil {
		return err
	}

	if err := validateLogin(request); err != nil {
		return err
	}

	token, err := h.auth.Login(e.Request().Context(), request)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			// todo: вынести ошибку в костанту
			return fmt.Errorf("%s", "invalid credentials")
		}
		return err
	}

	return e.JSON(http.StatusCreated, token)
}

func (h *Handler) Register(e echo.Context) error {
	var request dto.RegisterRequest

	if err := e.Bind(&request); err != nil {
		return err
	}

	if err := validateRegister(request); err != nil {
		return err
	}

	userID, err := h.auth.RegisterNewUser(e.Request().Context(), request)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return fmt.Errorf("%s", "user already exists")
		}

		return err
	}

	return e.JSON(http.StatusCreated, userID)
}

func (h *Handler) IsAdmin(e echo.Context) error {
	var request int64
	if err := e.Bind(&request); err != nil {
		return err
	}

	isAdmin, err := h.auth.IsAdmin(e.Request().Context(), request)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return fmt.Errorf("%s", "user not found")
		}
		return err
	}

	return e.JSON(http.StatusCreated, isAdmin)
}

func validateLogin(req dto.LoginRequest) error {
	if req.Login == "" || req.Password == "" {
		// todo: вынести ошибку в костанту
		return fmt.Errorf("%s", "login requires email and password")
	}

	return nil
}

func validateRegister(req dto.RegisterRequest) error {
	if req.Login == "" || req.Password == "" {
		// todo: вынести ошибку в костанту
		return fmt.Errorf("%s", "register requires email and password")
	}

	return nil
}
