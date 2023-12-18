package controllers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kominkamen/rootds-admin/auth"
	"github.com/kominkamen/rootds-admin/models"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	ID                 int     `json:"id"`
	Username           *string `json:"username"`
	RoleID             *int    `json:"role_id"`
	ApiKey             *string `json:"api_key"`
	ClientID           *string `json:"client_id"`
	CountryID          *int    `json:"country_code"`
	Email              *string `json:"email"`
	ValidationToken    *string `json:"validation_token"`
	Mobile             *string `json:"mobile"`
	ReferralCode       *string `json:"referral_code"`
	ProductID          *int    `json:"product_id"`
	TotalInvitees      int     `json:"total_invitees"`
	SuccessfulReferral int     `json:"successful_referral"`
	IsActive           int     `json:"is_active"`
}

type User struct{
	ID                 int     `json:"id"`
	Username           *string `json:"username"`
	Mobile             *string `json:"mobile"`
	Email              *string `json:"email"`
}

type CreateUser struct {
	ID                 int    `json:"id"`
	Username           string `json:"username"`
	RoleID             *int   `json:"role_id"`
	ApiKey             string `json:"api_key"`
	ClientID           string `json:"client_id"`
	CountryCode        *int   `json:"country_code"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	ValidationToken    string `json:"validation_token"`
	Mobile             string `json:"mobile"`
	ReferralCode       string `json:"referral_code"`
	ProductID          int    `json:"product_id"`
	TotalInvitees      int    `json:"total_invitees"`
	SuccessfulReferral int    `json:"successful_referral"`
	IsActive           int    `json:"is_active"`
}


// HandlerGetAllUser godoc
//
// @Summary     Get All Users
// @Tags        Users
// @Accept      json
// @Produce     json
// @Security    X-Access-Token
// @Success     200 {object} models.JsonResponse "Successful response"
// @Failure     400 {object} models.ErrorResponse "Bad Request"
// @Router      /api/v1/user [get]
func HandlerGetAllUser(c *fiber.Ctx) error {

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	query := `
	SELECT users.id, users.username, users.email, users.mobile
	FROM users
	ORDER BY users.created_at DESC LIMIT $1 OFFSET $2`

	countQuery := `SELECT COUNT(*) FROM users`

	var totalRows int
	if err := DB.QueryRow(context.Background(), countQuery).Scan(&totalRows); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error getting total count",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	rows, err := DB.Query(context.Background(), query, pageSize, offset)
	if err != nil {
		fmt.Println("Error executing the query:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error executing the query",
			StatusCode: fiber.StatusInternalServerError,
		})
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Mobile,
		); err != nil {
			fmt.Println("Error scanning rows:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
				Status:     false,
				Message:    "Error scanning rows",
				StatusCode: fiber.StatusInternalServerError,
			})
		}
		users = append(users, user)
	}

	totalPages := (totalRows + pageSize - 1) / pageSize

	if len(users) > 0 {
		paginatedResponse := models.PaginatedResponse{
			Status:     true,
			Message:    "Data found",
			Page:       page,
			PerPage:    pageSize,
			Total:      totalRows,
			TotalPages: totalPages,
			Data:       users,
		}
		return c.Status(fiber.StatusOK).JSON(paginatedResponse)
	} else {
		return c.Status(fiber.StatusNotFound).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Data not found",
			Data:       nil,
			StatusCode: fiber.StatusNotFound,
		})
	}
}

// HandlerCreateUser godoc
//
// @Summary Create A New User
// @Tags Users
// @Accept json
// @Produce json
// @Security X-Access-Token
// @Param user body CreateUser true "User object to be created"
// @Success        200 {object} models.JsonResponse "Successful response"
//
//	@Failure        400 {object} models.ErrorResponse "Bad Request"
//
// @Router /api/v1/user [post]
func HandlerCreateUser(c *fiber.Ctx) error {

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Invalid request body",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	hashedPassword := auth.Sha256Hash(user.Password)

	var roleID interface{}
	if user.RoleID == nil || *user.RoleID == 0 {
		roleID = nil
	} else {
		roleID = user.RoleID
	}

	var countryCode interface{}
	if user.CountryCode == nil || *user.CountryCode == 0 {
		countryCode = nil
	} else {
		countryCode = user.CountryCode
	}

	query := "INSERT INTO users(username, role_id, api_key, client_id, country_code, email, password, validation_token, mobile, referral_code, product_id, total_invitees, successful_referral, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)"

	_, err := DB.Exec(context.Background(), query,
		user.Username, roleID, user.ApiKey, user.ClientID, countryCode, user.Email, string(hashedPassword), user.ValidationToken, user.Mobile, user.ReferralCode, user.ProductID, user.TotalInvitees, user.SuccessfulReferral, user.IsActive)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error creating user",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	err = DB.QueryRow(context.Background(), "SELECT lastval()").Scan(&user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error retrieving user ID",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.JsonResponse{
		Status:     true,
		Message:    "User created successfully",
		Data:       user,
		StatusCode: fiber.StatusCreated,
	})
}

// HandlerGetOneUser godoc
//
// @Summary     Get One User
// @Tags        Users
// @Accept      json
// @Produce     json
// @Security    X-Access-Token
// @Param id path int true "User ID" format(int)
// @Success        200 {object} models.JsonResponse "Successful response"
//
//	@Failure        400 {object} models.ErrorResponse "Bad Request"
//
// @Router      /api/v1/user/{id} [get]

func HandlerGetOneUser(c *fiber.Ctx) error {

	id := c.Params("id")

	query := `SELECT
    users.id,
     users.username,
     users.role_id,
     users.api_key,
     users.client_id,
     users.country_code,
     users.email,
     users.validation_token,
     users.mobile,
     users.referral_code,
     users.product_id,
     users.total_invitees,
     users.successful_referral,
     users.is_active
	FROM users
	where users.id = $1`

	row := DB.QueryRow(context.Background(), query, id)

	var user Users
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.RoleID,
		&user.ApiKey,
		&user.ClientID,
		&user.CountryID,
		&user.Email,
		&user.ValidationToken,
		&user.Mobile,
		&user.ReferralCode,
		&user.ProductID,
		&user.TotalInvitees,
		&user.SuccessfulReferral,
		&user.IsActive,
	)

	if err != nil {
		fmt.Println("Error:", err)
		return c.Status(fiber.StatusNotFound).JSON(models.JsonResponse{
			Status:     false,
			Message:    "User not found",
			StatusCode: fiber.StatusNotFound,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.JsonResponse{
		Status:     true,
		Message:    "User retrieved successfully",
		Data:       user,
		StatusCode: fiber.StatusOK,
	})
}

// HandleUpdateUser godoc
//
// @Summary Update A User
// @Tags Users
// @Accept json
// @Produce json
// @Security X-Access-Token
// @Param id path int true "User ID" format(int)
// @Param user body CreateUser true "Updated user information"
// @Success        200 {object} models.JsonResponse "Successful response"
//
//	@Failure        400 {object} models.ErrorResponse "Bad Request"
//
// @Router /api/v1/user/{id} [put]
func HandleUpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var body models.User

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Failed to read body",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Unable to hash password",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	var roleID interface{}
	if body.RoleID == nil || *body.RoleID == 0 {
		roleID = nil
	} else {
		roleID = body.RoleID
	}

	var countryCode interface{}
	if body.CountryCode == nil || *body.CountryCode == 0 {
		countryCode = nil
	} else {
		countryCode = body.CountryCode
	}

	query := `
        UPDATE users
        SET username = $1, role_id = $2, api_key = $3, client_id = $4, country_code = $5,
        email = $6, password = $7, validation_token = $8, mobile = $9, referral_code = $10,
        product_id = $11, total_invitees = $12, successful_referral = $13, is_active = $14
        WHERE id = $15
    `

	_, err = DB.Exec(context.Background(), query,
		body.Username,
		roleID,
		body.ApiKey,
		body.ClientID,
		countryCode,
		body.Email,
		string(hash),
		body.ValidationToken,
		body.Mobile,
		body.ReferralCode,
		body.ProductID,
		body.TotalInvitees,
		body.SuccessfulReferral,
		body.IsActive,
		id,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error updating user",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(models.JsonResponse{
		Status:     true,
		Message:    "User updated successfully",
		StatusCode: fiber.StatusAccepted,
	})
}

// HandleDeleteUser godoc
//
// @Summary Delete A User
// @Tags Users
// @Accept json
// @Produce json
// @Security X-Access-Token
// @Param id path int true "User ID" format(int)
// @Success        200 {object} models.JsonResponse "Successful response"
//
//	@Failure        400 {object} models.ErrorResponse "Bad Request"
//
// @Router /api/v1/user/{id} [delete]
func HandleDeleteUser(c *fiber.Ctx) error {

	id := c.Params("id")

	query := `DELETE FROM users WHERE id = $1 `

	result, err := DB.Exec(context.Background(), query, id)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Unable to execute the query",
			StatusCode: fiber.StatusBadGateway,
		})
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.JsonResponse{
			Status:     false,
			Message:    "User not found",
			StatusCode: fiber.StatusNotFound,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(models.JsonResponse{
		Status:     true,
		Message:    "User deleted successfully",
		StatusCode: fiber.StatusAccepted,
	})
}

func HandleBulkDeleteUsers(c *fiber.Ctx) error {

	var request struct {
		UserIDs []string `json:"user_ids"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Failed to read request body",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	if len(request.UserIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.JsonResponse{
			Status:     false,
			Message:    "No user IDs provided for deletion",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	query := `DELETE FROM users WHERE id = ANY($1) RETURNING
    username,
    role_id,
    api_key,
    client_id,
    country_code,
    email,
    validation_token,
    mobile,
    referral_code,
    product_id,
    total_invitees,
    successful_referral,
    is_active`

	result, err := DB.Exec(context.Background(), query, request.UserIDs)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Unable to execute the query",
			StatusCode: fiber.StatusBadGateway,
		})
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.JsonResponse{
			Status:     false,
			Message:    "No users found for deletion",
			StatusCode: fiber.StatusNotFound,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(models.JsonResponse{
		Status:     true,
		Message:    "Users deleted successfully",
		StatusCode: fiber.StatusAccepted,
	})
}
