package middleware

import (
    "github.com/appleboy/gin-jwt/v2"
    "github.com/gin-gonic/gin"
    "time"
    "task-manager/models" // Import the models package
)

var identityKey = "id"

func AuthMiddleware() (*jwt.GinJWTMiddleware, error) {
    return jwt.New(&jwt.GinJWTMiddleware{
        Realm:       "task manager",
        Key:         []byte("secret key"),
        Timeout:     time.Hour,
        MaxRefresh:  time.Hour,
        IdentityKey: identityKey,
        Authenticator: func(c *gin.Context) (interface{}, error) {
            var loginVals struct {
                Username string `json:"username"`
                Password string `json:"password"`
            }
            if err := c.ShouldBindJSON(&loginVals); err != nil {
                return "", jwt.ErrMissingLoginValues
            }
            username := loginVals.Username
            password := loginVals.Password

            // Example: Validate username and password (replace with actual DB check)
            if username == "admin" && password == "password" {
                return &models.User{Username: username, Role: "admin"}, nil
            }

            return nil, jwt.ErrFailedAuthentication
        },
        Authorizator: func(data interface{}, c *gin.Context) bool {
            if v, ok := data.(*models.User); ok && v.Role == "admin" {
                return true
            }
            return false
        },
        Unauthorized: func(c *gin.Context, code int, message string) {
            c.JSON(code, gin.H{"error": message})
        },
        TokenLookup: "header: Authorization, query: token, cookie: jwt",
        TokenHeadName: "Bearer",
        TimeFunc: time.Now,
    })
}