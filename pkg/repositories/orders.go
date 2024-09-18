package repositories

import (
	"context"
	db "texApi/database"
	"texApi/pkg/queries"
	"texApi/pkg/schemas/request"
	"texApi/pkg/schemas/response"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func GetOrders(offset, limit int) (response.Orders, error) {
	var orders response.Orders

	err := pgxscan.Select(
		context.Background(), db.DB, &orders.Orders, queries.GetOrders,
		offset, limit,
	)

	if err != nil {
		return response.Orders{}, err
	}

	db.DB.QueryRow(
		context.Background(), queries.GetOrdersTotal,
	).Scan(&orders.Total)

	return orders, nil
}

func GetNewOrders() ([]response.Order, error) {
	var newOrders []response.Order

	err := pgxscan.Select(
		context.Background(), db.DB, &newOrders, queries.GetNewOrders,
	)

	if err != nil {
		return nil, err
	}

	return newOrders, nil
}

func GetOrdersByStatus(id, offset, limit int) (response.Orders, error) {
	var orders response.Orders

	err := pgxscan.Select(
		context.Background(), db.DB, &orders.Orders, queries.GetOrdersByStatus,
		id, offset, limit,
	)

	if err != nil {
		return response.Orders{}, err
	}

	db.DB.QueryRow(
		context.Background(), queries.GetOrdersByStatusTotal, id,
	).Scan(&orders.Total)

	return orders, nil
}

func GetOrdersByWorker(workerID int) ([]response.OrderByWorker, error) {
	var orders []response.OrderByWorker

	err := pgxscan.Select(
		context.Background(), db.DB, &orders, queries.GetOrdersByWorker,
		workerID,
	)

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func GetOrdersByUser(userID int) ([]response.OrderByUser, error) {
	var orders []response.OrderByUser

	err := pgxscan.Select(
		context.Background(), db.DB, &orders, queries.GetOrdersByUser,
		userID,
	)

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func GetOrder(id int) response.Order {
	var order response.Order

	db.DB.QueryRow(
		context.Background(), queries.GetOrder, id,
	).Scan(
		&order.ID,
		&order.OrderNumber,
		&order.FilePaths,
		&order.UserID,
		&order.WorkerID,
		&order.Address,
		&order.Date,
		&order.Time,
		&order.TimeDuration,
		&order.Status,
		&order.Description,
		&order.Services,
	)

	return order
}

func CreateOrder(order request.CreateOrder, userID int) (int, error) {
	var orderID int

	db.DB.QueryRow(
		context.Background(), queries.CreateOrder, userID, order.Address,
		order.Date, order.Time, order.SecretWord, order.Description,
	).Scan(&orderID)

	for _, serviceID := range order.Services {
		_, err := db.DB.Exec(
			context.Background(), queries.CreateOrderServices,
			orderID, serviceID,
		)

		if err != nil {
			return orderID, err
		}
	}

	return orderID, nil
}

func UpdateOrder(order request.UpdateOrder) int {
	var userID int

	db.DB.QueryRow(
		context.Background(), queries.UpdateOrder, order.WorkerID,
		order.StatusID, order.ID,
	).Scan(&userID)

	return userID
}

func UpdateOrderTimeDuration(
	order request.UpdateOrderTimeDuration, workerID int,
) error {
	_, err := db.DB.Exec(
		context.Background(), queries.UpdateOrderTimeDuration,
		order.TimeDuration, order.ID, workerID,
	)

	if err != nil {
		return err
	}

	return nil
}

func UpdateOrderStatusStart(
	order request.UpdateOrderStatusStart, workerID int,
) error {
	_, err := db.DB.Exec(
		context.Background(), queries.UpdateOrderStatusStart,
		order.ID, workerID,
	)

	if err != nil {
		return err
	}

	return nil
}

func UpdateOrderRead(orderID int) error {
	_, err := db.DB.Exec(
		context.Background(), queries.UpdateOrderRead, orderID,
	)

	if err != nil {
		return err
	}

	return nil
}

func CheckOrderStatus(orderID int) int {
	var statusID int

	db.DB.QueryRow(
		context.Background(), queries.CheckOrderStatus, orderID,
	).Scan(&statusID)

	return statusID
}

func AbortOrder(orderID, userID int) error {
	_, err := db.DB.Exec(
		context.Background(), queries.AbortOrder, orderID, userID,
	)

	if err != nil {
		return err
	}

	return nil
}

func DeleteOrder(id int) []string {
	var paths []string

	db.DB.QueryRow(
		context.Background(), queries.DeleteOrder, id,
	).Scan(&paths)

	return paths
}

func CheckOrderExist(orderID int) int {
	var id int

	db.DB.QueryRow(
		context.Background(), queries.CheckOrderExist, orderID,
	).Scan(&id)

	return id
}

func SaveOrderFiles(filePaths []string, orderID int) error {
	for _, filePath := range filePaths {
		_, err := db.DB.Exec(
			context.Background(), queries.SaveOrderFile, orderID, filePath,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
