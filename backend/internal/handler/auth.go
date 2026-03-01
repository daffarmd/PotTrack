package handler

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication
type AuthHandler struct {
	db  *sqlx.DB
	jwt string
}

func NewAuthHandler(db *sqlx.DB, jwtSecret string) *AuthHandler {
	return &AuthHandler{db: db, jwt: jwtSecret}
}

// Register creates a new pengguna
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Nama     string `json:"nama" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	_, err := h.db.Exec(`INSERT INTO pengguna (email, kata_sandi_hash, nama, zona_waktu, dibuat_pada) VALUES ($1,$2,$3,'UTC',now())`, req.Email, string(hash), req.Nama)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}
	c.Status(http.StatusCreated)
}

// Login authenticates and returns JWT
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var id int
	var hash string
	err := h.db.Get(&hash, "SELECT kata_sandi_hash FROM pengguna WHERE email=$1", req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	// fetch id
	h.db.Get(&id, "SELECT id FROM pengguna WHERE email=$1", req.Email)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokStr, _ := token.SignedString([]byte(h.jwt))
	c.JSON(http.StatusOK, gin.H{"token": tokStr})
}
