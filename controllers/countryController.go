package controllers

import (
	"context"
	"fmt"
	"strconv"
	"github.com/gofiber/fiber/v2"
	_ "github.com/kominkamen/rootds-admin/docs"
	"github.com/kominkamen/rootds-admin/models"
)

type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// HandlerGetAllCountry godoc
//
// @Summary     Get All Countries
// @Tags        Countries
// @Accept      json
// @Produce     json
// @Security    X-Access-Token
// @Success     200 {object} models.JsonResponse "Successful response"
// @Failure     400 {object} models.ErrorResponse "Bad Request"
// @Router      /api/v1/country [get]
func HandlerGetAllCountry(c *fiber.Ctx) error {

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	query := "select * from countries LIMIT $1 OFFSET $2"

	countQuery := `SELECT COUNT(*) FROM countries`

	var totalRows int
	if err := DB.QueryRow(context.Background(), countQuery).Scan(&totalRows); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error getting total count",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	rows, err := DB.Query(context.Background(), query,pageSize,offset)
	if err != nil {
		fmt.Println("Error executing the query:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error executing the query",
			StatusCode: fiber.StatusInternalServerError,
		})
	}
	defer rows.Close()

	var countries []models.Country
	for rows.Next() {
		var country models.Country
		if err := rows.Scan(
			&country.Code,
			&country.Name,
			&country.ContinentName,
		); err != nil {
			fmt.Println("Error scanning rows:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
				Status:     false,
				Message:    "Error scanning rows",
				StatusCode: fiber.StatusInternalServerError,
			})
		}
		countries = append(countries, country)
	}

	totalPages := (totalRows + pageSize - 1) / pageSize

	if len(countries) > 0 {
		paginatedResponse := models.PaginatedResponse{
			Status:     true,
			Message:    "Data found",
			Page:       page,
			PerPage:    pageSize,
			Total:      totalRows,
			TotalPages: totalPages,
			Data:       countries,
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

// HandlerCreateCountry godoc
//
// @Summary Create A New Country
// @Tags Countries
// @Accept json
// @Produce json
// @Security X-Access-Token
// @Param country body models.Country true "Country object to be created"
// @Success        200 {object} models.JsonResponse "Successful response"
//
//	@Failure        400 {object} models.ErrorResponse "Bad Request"
//
// @Router /api/v1/country [post]
func HandlerCreateCountry(c *fiber.Ctx) error {

	var country models.Country
	if err := c.BodyParser(&country); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Invalid request body",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	query := "INSERT INTO countries(name,continent_name) VALUES ($1,$2)"
	_, err := DB.Exec(context.Background(), query, country.Name, country.ContinentName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error creating country",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	err = DB.QueryRow(context.Background(), "SELECT lastval()").Scan(&country.Code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error retrieving country code",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.JsonResponse{
		Status:     true,
		Message:    "Country created successfully",
		Data:       country,
		StatusCode: fiber.StatusCreated,
	})
}

// HandlerGetOneCountry godoc
//
// @Summary     Get One Country
// @Tags        Countries
// @Accept      json
// @Produce     json
// @Security    X-Access-Token
// @Param id path int true "Country ID" format(int)
// @Success        200 {object} models.JsonResponse "Successful response"
//
//	@Failure        400 {object} models.ErrorResponse "Bad Request"
//
// @Router      /api/v1/country/{id} [get]
func HandlerGetOneCountry(c *fiber.Ctx) error {

	id := c.Params("id")

	query := "SELECT * from countries where code = $1"

	row := DB.QueryRow(context.Background(), query, id)

	var country models.Country
	err := row.Scan(
		&country.Code,
		&country.Name,
		&country.ContinentName,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Country not found",
			StatusCode: fiber.StatusNotFound,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.JsonResponse{
		Status:     true,
		Message:    "Country retrieved successfully",
		Data:       country,
		StatusCode: fiber.StatusOK,
	})
}

// HandleUpdateCountry godoc
//
// @Summary Update A Country
// @Tags Countries
// @Accept json
// @Produce json
// @Security X-Access-Token
// @Param id path int true "Country ID" format(int)
// @Param country body models.Country true "Updated country information"
// @Success        200 {object} models.JsonResponse "Successful response"
//
//	@Failure        400 {object} models.ErrorResponse "Bad Request"
//
// @Router /api/v1/country/{id} [put]
func HandleUpdateCountry(c *fiber.Ctx) error {

	id := c.Params("id")

	var country models.Country

	if err := c.BodyParser(&country); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Failed to read body",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	query := "UPDATE countries SET name = $1, continent_name = $2 WHERE code = $3"

	_, err := DB.Exec(context.Background(), query, country.Name, country.ContinentName, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error updating country",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(models.JsonResponse{
		Status:     true,
		Message:    "Country updated successfully",
		StatusCode: fiber.StatusAccepted,
	})
}

// HandleDeleteCountry godoc
//
// @Summary Delete A Country
// @Tags Countries
// @Accept json
// @Produce json
// @Security X-Access-Token
// @Param id path int true "Country ID" format(int)
// @Success        200 {object} models.JsonResponse "Successful response"
//
//	@Failure        400 {object} models.ErrorResponse "Bad Request"
//
// @Router /api/v1/country/{id} [delete]
func HandleDeleteCountry(c *fiber.Ctx) error {

	id := c.Params("id")

	query := `DELETE FROM countries WHERE code = $1 `

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
			Message:    "Country not found",
			StatusCode: fiber.StatusNotFound,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(models.JsonResponse{
		Status:     true,
		Message:    "Country deleted successfully",
		StatusCode: fiber.StatusAccepted,
	})
}
