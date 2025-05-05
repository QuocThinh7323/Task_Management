package main

import (
    "log"
    "task-manager/config"
    "task-manager/controllers"
    "task-manager/middleware"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // Kết nối cơ sở dữ liệu
    config.ConnectDatabase()

    // Thiết lập middleware xác thực JWT
    authMiddleware, err := middleware.AuthMiddleware()
    if err != nil {
        log.Fatal("Lỗi khi khởi tạo JWT Middleware:", err)
    }

    // Route đăng nhập
    r.POST("/login", authMiddleware.LoginHandler)

    // Route đăng ký người dùng mới
    r.POST("/register", func(c *gin.Context) {
        // Xử lý đăng ký người dùng
    })

    // Nhóm các route yêu cầu xác thực
    auth := r.Group("/")
    auth.Use(authMiddleware.MiddlewareFunc())
    {
        auth.GET("/tasks", controllers.GetTasks)
        auth.POST("/tasks", controllers.CreateTask)
        auth.PUT("/tasks/:id", controllers.UpdateTask)
        auth.DELETE("/tasks/:id", controllers.DeleteTask)
    }

    // Chạy server trên cổng 8080
    if err := r.Run(":8080"); err != nil {
        log.Fatal("Không thể khởi động server:", err)
    }
}
