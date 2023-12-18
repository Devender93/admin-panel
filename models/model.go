package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DBConfig struct {
	DB_User     string `json:"db_user"`
	DB_Pass     string `json:"db_pass"`
	DB_Host     string `json:"db_host"`
	DB_Port     string `json:"db_port"`
	DB_Database string `json:"db_database"`
	DB_Sslmode  string `json:"db_sslmode"`
}

type Product struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	ReferralLink string    `json:"referral_link"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Country struct {
	Code          int    `json:"code"`
	Name          string `json:"name"`
	ContinentName string `json:"continent_name"`
}

type Claims struct {
	UserID int     `json:"userID"`
	Email  string  `json:"email"`
	Role   *string `json:"role"`
	jwt.RegisteredClaims
}

type AdminUser struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Role     *string `json:"role"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
}

type User struct {
	ID                 int       `json:"id"`
	Username           string    `json:"username"`
	RoleID             *int      `json:"role_id"`
	ApiKey             string    `json:"api_key"`
	ClientID           string    `json:"client_id"`
	CountryCode        *int      `json:"country_code"`
	Email              string    `json:"email"`
	Password           string    `json:"password"`
	ValidationToken    string    `json:"validation_token"`
	Mobile             string    `json:"mobile"`
	ReferralCode       string    `json:"referral_code"`
	ProductID          int       `json:"product_id"`
	TotalInvitees      int       `json:"total_invitees"`
	SuccessfulReferral int       `json:"successful_referral"`
	IsActive           int       `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type PaginatedResponse struct {
	Status     bool        `json:"status"`
	Message    string      `json:"message"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

type JsonResponse struct {
	Status     bool        `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	StatusCode int         `json:"status_code"`
}

type ErrorResponse struct {
	Status     bool   `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}
