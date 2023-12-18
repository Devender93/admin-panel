package controllers

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kominkamen/rootds-admin/auth"
	_ "github.com/kominkamen/rootds-admin/docs"
	"github.com/kominkamen/rootds-admin/models"
)

type AdminLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(username, password string) (*models.AdminUser, string, time.Time, error) {

	query := `SELECT users.id,
	users.username,
	user_roles.name,
	users.email,
	users.password
	FROM users Left join user_roles on user_roles.id = users.role_id
	WHERE email = $1`

	var user models.AdminUser
	err := DB.QueryRow(context.Background(), query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Role,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", time.Time{}, errors.New("Invalid email or password")
		}
		return nil, "", time.Time{}, err
	}

	hashedPassword := auth.Sha256Hash(password)

	if user.Password != hashedPassword {
		return nil, "", time.Time{}, errors.New("Invalid email or password")
	}

	accessToken, expirationTime, err := auth.GenerateJWT(&user)
	if err != nil {
		return nil, "", time.Time{}, err
	}

	return &user, accessToken, expirationTime, nil
}

// HandlerAdminLogin godoc
//
//	@Summary        Login
//	@Tags           Auth
//	@Accept         json
//	@Produce        json
//
// @Param payload body AdminLogin true "UserModel"
// @Success        200 {object} models.JsonResponse "Successful response"
//
//	@Failure        400 {object} models.ErrorResponse "Bad Request"
//
// @Router /api/v1/login [post]
func HandlerAdminLogin(c *fiber.Ctx) error {
	var login AdminLogin
	if err := c.BodyParser(&login); err != nil {
		response := models.JsonResponse{
			Status:     false,
			Message:    "Invalid request body",
			StatusCode: fiber.StatusBadRequest,
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	user, accessToken, expirationTime, err := Login(login.Email, login.Password)
	if err != nil {
		response := models.JsonResponse{
			Status:     false,
			Message:    err.Error(),
			StatusCode: fiber.StatusUnauthorized,
		}
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	apiResponse := map[string]interface{}{
		"token":      accessToken,
		"userEmail":  user.Email,
		"expireDate": expirationTime.Format("2006-01-02"),
	}

	response := models.JsonResponse{
		Status:     true,
		Message:    "Authenticated successfully",
		Data:       apiResponse,
		StatusCode: fiber.StatusOK,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
