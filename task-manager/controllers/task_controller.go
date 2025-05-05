package controllers

import (
    "net/http"
    "task-manager/config"
    "task-manager/models"

    "github.com/gin-gonic/gin"
)

// Tạo mới một task
func CreateTask(c *gin.Context) {
    var input models.Task
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    task := models.Task{Title: input.Title, Description: input.Description, Completed: input.Completed}
    config.DB.Create(&task)
    c.JSON(http.StatusOK, task)
}

// Lấy danh sách tất cả các task
func GetTasks(c *gin.Context) {
    var tasks []models.Task
    config.DB.Find(&tasks)
    c.JSON(http.StatusOK, tasks)
}

// Cập nhật một task
func UpdateTask(c *gin.Context) {
    var task models.Task
    if err := config.DB.Where("id = ?", c.Param("id")).First(&task).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task không tồn tại"})
        return
    }
    var input models.Task
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    config.DB.Model(&task).Updates(input)
    c.JSON(http.StatusOK, task)
}

// Xóa một task
func DeleteTask(c *gin.Context) {
    var task models.Task
    if err := config.DB.Where("id = ?", c.Param("id")).First(&task).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task không tồn tại"})
        return
    }
    config.DB.Delete(&task)
    c.JSON(http.StatusOK, gin.H{"data": true})
}
