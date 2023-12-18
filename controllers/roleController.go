package controllers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/kominkamen/rootds-admin/docs"
	"github.com/kominkamen/rootds-admin/models"
)

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// HandlerGetAllRole godoc
//
// @Summary     Get All Roles
// @Tags        Roles
// @Accept      json
// @Produce     json
// @Security    X-Access-Token
// @Success     200 {object} models.JsonResponse "Successful response"
// @Failure     400 {object} models.ErrorResponse "Bad Request"
// @Router      /api/v1/role [get]
func HandlerGetAllRole(c *fiber.Ctx) error {

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	query := `select id, name from user_roles LIMIT $1 OFFSET $2`

	countQuery := `SELECT COUNT(*) FROM user_roles`

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

	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(
			&role.ID,
			&role.Name,
		); err != nil {
			fmt.Println("Error scanning rows:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
				Status:     false,
				Message:    "Error scanning rows",
				StatusCode: fiber.StatusInternalServerError,
			})
		}
		roles = append(roles, role)
	}

	totalPages := (totalRows + pageSize - 1) / pageSize

	if len(roles) > 0 {
		paginatedResponse := models.PaginatedResponse{
			Status:     true,
			Message:    "Data found",
			Page:       page,
			PerPage:    pageSize,
			Total:      totalRows,
			TotalPages: totalPages,
			Data:       roles,
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

// HandlerCreateRole godoc
//
// @Summary Create A New Role
// @Tags Roles
// @Accept json
// @Produce json
// @Security X-Access-Token
// @Param role body models.Role true "Role object to be created"
// @Success        200 {object} models.JsonResponse "Successful response"
//
//	@Failure        400 {object} models.ErrorResponse "Bad Request"
//
// @Router /api/v1/role [post]
func HandlerCreateRole(c *fiber.Ctx) error {

	var role models.Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Invalid request body",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	query := "INSERT INTO user_roles(name) VALUES ($1)"
	_, err := DB.Exec(context.Background(), query, role.Name)
	if err != nil {
		fmt.Println("roleee", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error creating role",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	err = DB.QueryRow(context.Background(), "SELECT lastval()").Scan(&role.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error retrieving role ID",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.JsonResponse{
		Status:     true,
		Message:    "Role created successfully",
		Data:       role,
		StatusCode: fiber.StatusCreated,
	})
}

// HandlerGetOneRole godoc
//
// @Summary     Get One Role
// @Tags        Roles
// @Accept      json
// @Produce     json
// @Security    X-Access-Token
// @Param id path int true "Role ID" format(int)
// @Success        200 {object} models.JsonResponse "Successful response"
//
//	@Failure        400 {object} models.ErrorResponse "Bad Request"
//
// @Router      /api/v1/role/{id} [get]
func HandlerGetOneRole(c *fiber.Ctx) error {

	id := c.Params("id")

	query := "SELECT * from user_roles where id = $1"

	row := DB.QueryRow(context.Background(), query, id)

	var role models.Role
	err := row.Scan(
		&role.ID,
		&role.Name,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Role not found",
			StatusCode: fiber.StatusNotFound,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.JsonResponse{
		Status:     true,
		Message:    "Role retrieved successfully",
		Data:       role,
		StatusCode: fiber.StatusOK,
	})
}

// HandleUpdateRole godoc
//
// @Summary Update A Role
// @Tags Roles
// @Accept json
// @Produce json
// @Security X-Access-Token
// @Param id path int true "Role ID" format(int)
// @Param role body models.Role true "Updated role information"
// @Success        200 {object} models.JsonResponse "Successful response"
//
//	@Failure        400 {object} models.ErrorResponse "Bad Request"
//
// @Router /api/v1/role/{id} [put]
func HandleUpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")

	var role models.Role

	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Failed to read body",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	query := "UPDATE user_roles SET name = $1 WHERE id = $2"

	_, err := DB.Exec(context.Background(), query, role.Name, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error updating role",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(models.JsonResponse{
		Status:     true,
		Message:    "Role updated successfully",
		StatusCode: fiber.StatusAccepted,
	})
}

// HandleDeleteRole godoc
//
// @Summary Delete A Role
// @Tags Roles
// @Accept json
// @Produce json
// @Security X-Access-Token
// @Param id path int true "Role ID" format(int)
// @Success        200 {object} models.JsonResponse "Successful response"
//
//	@Failure        400 {object} models.ErrorResponse "Bad Request"
//
// @Router /api/v1/role/{id} [delete]
func HandleDeleteRole(c *fiber.Ctx) error {

	id := c.Params("id")

	query := `DELETE FROM user_roles WHERE id = $1 `

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
			Message:    "Role not found",
			StatusCode: fiber.StatusNotFound,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(models.JsonResponse{
		Status:     true,
		Message:    "Role deleted successfully",
		StatusCode: fiber.StatusAccepted,
	})
}
