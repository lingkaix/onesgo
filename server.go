package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	// ? all handlers need better error info output
	router.POST("/login", login)
	router.POST("/register", register)
	router.GET("/users/:id", authMiddleware, getUser)
	router.POST("/users/:id", authMiddleware, updateUser)
	router.DELETE("/users/:id", authMiddleware, deleteUser)

	return router
}

type auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func login(c *gin.Context) {
	var auth auth
	if err := c.BindJSON(&auth); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user, err := GetUserByEmail(db, auth.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(auth.Password))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	token, _ := GenerateToken(strconv.FormatUint(uint64(user.ID), 10), KEY)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func register(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if user.Password == "" || user.Email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Please enter email/password"})
		return
	}
	err := AddUser(db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ? filter out more sensitive data
	//
	// ! filter out password from json output
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func getUser(c *gin.Context) {
	jid, ok := c.Get("user-id-jwt")
	if !ok {
		// ? need log this rare error
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	// ? rarely occurs a panic if jid is not a string
	if c.Param("id") != jid.(string) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	user, err := GetUser(db, id)
	if err != nil {
		// ? do not expose internal error in production deployment
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func updateUser(c *gin.Context) {
	jid, _ := c.Get("user-id-jwt")
	if c.Param("id") != jid.(string) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
		return
	}

	var user User
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	user.ID, err = strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var keys []string
	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		if key, ok := token.(string); ok {
			keys = append(keys, key)
		}
	}

	// ? for deleted record, it will return 200 although it won't updadte the record
	// ? need a check if the record has been deleted
	err = UpdateUser(db, &user, keys)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
	jid, _ := c.Get("user-id-jwt")
	if c.Param("id") != jid.(string) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	err = DeleteUser(db, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
