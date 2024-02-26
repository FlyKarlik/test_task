package handler

import (
	apiserver "mongo_db"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	headerToken         = "refresh_token"
	authorizationHeader = "access_token"
)

func (h *Handler) GetTokens(c *gin.Context) {
	var user apiserver.User
	userId := c.Query("id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid iser_id",
		})
		return
	}
	user.Id = userId
	refreshToken, err := h.service.GenerateRefreshToken()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user.RefreshToken = refreshToken
	if err := h.service.CreateUser(user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	accessToken, err := h.service.GenerateAccessToken(refreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})

}

func (h *Handler) RefreshTokens(c *gin.Context) {
	userId := c.Query("id")
	refreshTokenFromHeader := c.GetHeader(headerToken)
	accessTokenFromHeader := c.GetHeader(authorizationHeader)
	if refreshTokenFromHeader == "" && accessTokenFromHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "empty header",
		})
		return
	}

	parseToken, err := h.service.ParseAccessToken(accessTokenFromHeader)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if refreshTokenFromHeader != parseToken {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.service.GetUserById(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "user not found",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.RefreshToken), []byte(refreshTokenFromHeader)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	accessToken, err := h.service.GenerateAccessToken(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "could not generate access token",
		})
		return
	}

	refreshToken, err := h.service.GenerateRefreshToken()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "could not generate refresh token",
		})
		return
	}
	if err := h.service.UpdateRefreshToken(userId, refreshToken); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "could not save refresh token in database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"new_access_token":  accessToken,
		"new_refresh_token": refreshToken,
	})
}
