package services

import (
	"log"
	"os"
	"strconv"
	"strings"
	"texApi/config"
	repo "texApi/internal/_other/repositories"
	"texApi/internal/_other/schemas/request"
	"texApi/pkg/utils"

	"github.com/gin-gonic/gin"
)

func GetOrders(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	pageStr := ctx.Query("page")
	countStr := ctx.Query("count")

	pageInt, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(countStr)

	// set default values if not presented
	if pageInt == 0 {
		pageInt = 1
	}

	if limit == 0 {
		limit = 20
	}

	offset := pageInt*limit - limit

	orders, err := repo.GetOrders(offset, limit)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(200, orders)
}

func GetNewOrders(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	newOrders, err := repo.GetNewOrders()

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(200, newOrders)
}

func GetOrdersByStatus(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	idStr := ctx.Param("id")
	id, error := strconv.Atoi(idStr)

	if error != nil || id == 0 {
		ctx.JSON(400, gin.H{"message": "ID must be positive integer number"})
		return
	}

	pageStr := ctx.Query("page")
	countStr := ctx.Query("count")

	pageInt, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(countStr)

	// set default values if not presented
	if pageInt == 0 {
		pageInt = 1
	}

	if limit == 0 {
		limit = 20
	}

	offset := pageInt*limit - limit

	orders, err := repo.GetOrdersByStatus(id, offset, limit)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(200, orders)
}

func GetOrdersByWorker(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "worker" {
		ctx.AbortWithStatus(403)
		return
	}

	workerID := ctx.MustGet("id").(int)

	orders, err := repo.GetOrdersByWorker(workerID)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(200, orders)
}

func GetOrdersByUser(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "user" {
		ctx.AbortWithStatus(403)
		return
	}

	userID := ctx.MustGet("id").(int)

	orders, err := repo.GetOrdersByUser(userID)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(200, orders)
}

func GetOrder(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "worker" && role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	idStr := ctx.Param("id")
	id, error := strconv.Atoi(idStr)

	if error != nil || id == 0 {
		ctx.JSON(400, gin.H{"message": "ID must be positive integer number"})
		return
	}

	order := repo.GetOrder(id)

	ctx.JSON(200, order)
}

// func CreateOrder(ctx *gin.Context) {
// 	role := ctx.MustGet("role").(string)

// 	if role != "user" {
// 		ctx.AbortWithStatus(403)
// 		return
// 	}

// 	userID := ctx.MustGet("id").(int)

// 	var order request.CreateOrder

// 	validationError := ctx.BindJSON(&order)

// 	if validationError != nil {
// 		ctx.JSON(400, validationError.Error())
// 		return
// 	}

// 	orderID, err := repo.CreateOrder(order, userID)

// 	if err != nil {
// 		log.Println(err.Error())
// 		ctx.JSON(500, gin.H{"message": err.Error()})
// 		return
// 	}

// 	utils.NotifyBySocket("admin1", strconv.Itoa(orderID))

// 	ctx.JSON(201, gin.H{"id": orderID})
// }

// func UpdateOrder(ctx *gin.Context) {
// 	role := ctx.MustGet("role").(string)

// 	if role != "worker" && role != "admin" {
// 		ctx.AbortWithStatus(403)
// 		return
// 	}

// 	var order request.UpdateOrder

// 	validationError := ctx.BindJSON(&order)

// 	if validationError != nil {
// 		ctx.JSON(400, validationError.Error())
// 		return
// 	}

// 	userID := repo.UpdateOrder(order)

// 	if userID == 0 {
// 		ctx.JSON(400, gin.H{"message": "Nothing updated (cause of order ID)"})
// 		return
// 	}

// 	utils.NotifyBySocket("worker"+strconv.Itoa(order.WorkerID), "New order")

// 	userNotificationToken := repo.GetUserNotificationToken(userID)
// 	userNotifyErr := utils.Notify("user", userNotificationToken)

// 	if userNotifyErr != nil {
// 		log.Println(userNotifyErr.Error())
// 	}

// 	ctx.JSON(200, gin.H{"message": "Successfully updated"})
// }

func UpdateOrderTimeDuration(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "worker" {
		ctx.AbortWithStatus(403)
		return
	}

	workerID := ctx.MustGet("id").(int)

	var order request.UpdateOrderTimeDuration

	validationError := ctx.BindJSON(&order)

	if validationError != nil {
		ctx.JSON(400, validationError.Error())
		return
	}

	err := repo.UpdateOrderTimeDuration(order, workerID)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Successfully updated"})
}

func UpdateOrderStatusStart(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "worker" {
		ctx.AbortWithStatus(403)
		return
	}

	workerID := ctx.MustGet("id").(int)

	var order request.UpdateOrderStatusStart

	validationError := ctx.BindJSON(&order)

	if validationError != nil {
		ctx.JSON(400, validationError.Error())
		return
	}

	err := repo.UpdateOrderStatusStart(order, workerID)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Successfully updated"})
}

func UpdateOrderRead(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	idStr := ctx.Param("id")
	orderID, error := strconv.Atoi(idStr)

	if error != nil || orderID == 0 {
		ctx.JSON(400, gin.H{"message": "ID must be positive integer number"})
		return
	}

	err := repo.UpdateOrderRead(orderID)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Successfully updated"})
}

func AbortOrder(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "user" {
		ctx.AbortWithStatus(403)
		return
	}

	order := ctx.Param("id")
	orderID, _ := strconv.Atoi(order)

	if orderID == 0 {
		ctx.JSON(400, gin.H{"message": "ID must be positive integer number"})
		return
	}

	userID := ctx.MustGet("id").(int)

	if repo.CheckOrderStatus(orderID) == 3 {
		ctx.JSON(400, gin.H{"message": "Can't abort"})
		return
	}

	err := repo.AbortOrder(orderID, userID)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Has been aborted"})
}

func DeleteOrder(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	idStr := ctx.Param("id")
	id, error := strconv.Atoi(idStr)

	if error != nil || id == 0 {
		ctx.JSON(400, gin.H{"message": "ID must be positive integer number"})
		return
	}

	paths := repo.DeleteOrder(id)

	for _, path := range paths {
		splitedPath := strings.Split(path, "/uploads/")[1]
		os.Remove(config.ENV.UPLOAD_PATH + splitedPath)
	}

	ctx.JSON(200, gin.H{"message": "Successfully deleted"})
}

func SaveOrderFiles(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "user" && role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	idStr := ctx.Param("id")
	orderID, _ := strconv.Atoi(idStr)

	if orderID == 0 {
		ctx.JSON(400, gin.H{"message": "ID must be positive integer number"})
		return
	}

	filePaths, err := utils.SaveFiles(ctx)

	if err != nil {
		log.Println(err.Error())
		return
	}

	orderExist := repo.CheckOrderExist(orderID)

	if orderExist == 0 {
		ctx.JSON(400, gin.H{"message": "This order doesn't exist"})
		return
	}

	saveFileErr := repo.SaveOrderFiles(filePaths, orderID)

	if saveFileErr != nil {
		log.Println(saveFileErr.Error())
		ctx.JSON(500, gin.H{"message": saveFileErr.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Has been saved"})
}
