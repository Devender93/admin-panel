package controllers

import (
	"context"
	"fmt"
	"strconv"
	"github.com/gofiber/fiber/v2"
	_ "github.com/kominkamen/rootds-admin/docs"
	"github.com/kominkamen/rootds-admin/models"
)

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Products struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ReferralLink string `json:"referral_link"`
	IsActive     bool   `json:"is_active"`
}

// HandlerGetAllProduct godoc
//
// @Summary     Get All Products
// @Tags        Products
// @Accept      json
// @Produce     json
// @Security    X-Access-Token
// @Success     200 {object} models.JsonResponse "Successful response"
// @Failure     400 {object} models.ErrorResponse "Bad Request"
// @Router      /api/v1/product [get]
func HandlerGetAllProduct(c *fiber.Ctx) error {

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	query := `select id, name,referral_link,is_active from products ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	countQuery := `SELECT COUNT(*) FROM products`

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

	var products []Products
	for rows.Next() {
		var product Products
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.ReferralLink,
			&product.IsActive,
		); err != nil {
			fmt.Println("Error scanning rows:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
				Status:     false,
				Message:    "Error scanning rows",
				StatusCode: fiber.StatusInternalServerError,
			})
		}
		products = append(products, product)
	}
	totalPages := (totalRows + pageSize - 1) / pageSize

	if len(products) > 0 {
		paginatedResponse := models.PaginatedResponse{
			Status:     true,
			Message:    "Data found",
			Page:       page,
			PerPage:    pageSize,
			Total:      totalRows,
			TotalPages: totalPages,
			Data:       products,
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

// HandlerCreateProduct godoc
//
// @Summary Create A New Product
// @Tags Products
// @Accept json
// @Produce json
// @Security X-Access-Token
// @Param product body Products true "Product object to be created"
// @Success        200 {object} models.JsonResponse "Successful response"
//
//	@Failure        400 {object} models.ErrorResponse "Bad Request"
//
// @Router /api/v1/product [post]
func HandlerCreateProduct(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Invalid request body",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	query := "INSERT INTO products(name,referral_link,is_active) VALUES ($1,$2,$3)"
	_, err := DB.Exec(context.Background(), query, product.Name, product.ReferralLink, product.IsActive)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error creating product",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	err = DB.QueryRow(context.Background(), "SELECT lastval()").Scan(&product.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error retrieving product id",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.JsonResponse{
		Status:     true,
		Message:    "Product created successfully",
		Data:       product,
		StatusCode: fiber.StatusCreated,
	})
}

// HandlerGetOneProduct godoc
//
// @Summary     Get One Product
// @Tags        Products
// @Accept      json
// @Produce     json
// @Security    X-Access-Token
// @Param id path int true "Product ID" format(int)
// @Success        200 {object} models.JsonResponse "Successful response"
//
//	@Failure        400 {object} models.ErrorResponse "Bad Request"
//
// @Router      /api/v1/product/{id} [get]
func HandlerGetOneProduct(c *fiber.Ctx) error {

	id := c.Params("id")

	query := "SELECT * from products where id = $1"

	row := DB.QueryRow(context.Background(), query, id)

	var product models.Product
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.ReferralLink,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Product not found",
			StatusCode: fiber.StatusNotFound,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.JsonResponse{
		Status:     true,
		Message:    "Product retrieved successfully",
		Data:       product,
		StatusCode: fiber.StatusOK,
	})
}

// HandleUpdateProduct godoc
//
// @Summary Update A Product
// @Tags Products
// @Accept json
// @Produce json
// @Security X-Access-Token
// @Param id path int true "Products ID" format(int)
// @Param product body Products true "Updated product information"
// @Success        200 {object} models.JsonResponse "Successful response"
//
//	@Failure        400 {object} models.ErrorResponse "Bad Request"
//
// @Router /api/v1/product/{id} [put]
func HandleUpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Failed to read body",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	query := "UPDATE products SET name = $1, referral_link = $2, is_active = $3 WHERE id = $4"

	_, err := DB.Exec(context.Background(), query, product.Name, product.ReferralLink, product.IsActive, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error updating product",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(models.JsonResponse{
		Status:     true,
		Message:    "Product updated successfully",
		StatusCode: fiber.StatusAccepted,
	})
}

// HandleDeleteProduct godoc
//
// @Summary Delete A Product
// @Tags Products
// @Accept json
// @Produce json
// @Security X-Access-Token
// @Param id path int true "Product ID" format(int)
// @Success        200 {object} models.JsonResponse "Successful response"
//
//	@Failure        400 {object} models.ErrorResponse "Bad Request"
//
// @Router /api/v1/product/{id} [delete]
func HandleDeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	query := `DELETE FROM products WHERE id = $1 `

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
			Message:    "Product not found",
			StatusCode: fiber.StatusNotFound,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(models.JsonResponse{
		Status:     true,
		Message:    "Product deleted successfully",
		StatusCode: fiber.StatusAccepted,
	})
}
