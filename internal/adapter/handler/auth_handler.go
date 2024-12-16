package handler

import (
	"bwanews/internal/adapter/handler/request"
	"bwanews/internal/adapter/handler/response"
	"bwanews/internal/core/domain/entity"
	"bwanews/internal/core/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// buat variable global agar di handler lain bisa digunakan
var err error
var code string
var errorResp response.ErrorResponseDefault
var validate = validator.New()

type AuthHandler interface {
	Login(c *fiber.Ctx) error
}

type authHandler struct {
	authService service.AuthService
}

func (a *authHandler) Login(c *fiber.Ctx) error {
	req := request.LoginRequest{}
	resp := response.SuccessAuthResponse{}

	// JIKA GAGAL
	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] Login -1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	// Check Validate
	 if err = validate.Struct(req); err!= nil {
		code = "[HANDLER] Login -2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	 }

	 // IF SUCCESS
	 reqLogin := entity.LoginRequest{
		Email: req.Email,
		Password: req.Password,
	 }

	 result, err := a.authService.GetUserByEmail(c.Context(), reqLogin)
	 if err!= nil {
		code = "[HANDLER] Login - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Message = err.Error()

		if err.Error() == "invalid password" {
			 return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	 }

	 resp.Meta.Status = true
	 resp.Meta.Message = "Login SuccessFully"
	 resp.AccessToken = result.AccessToken
	 resp.ExpiresAt = string(result.ExpiresAt)

	 return c.JSON(resp)
}


func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{

	}
}

