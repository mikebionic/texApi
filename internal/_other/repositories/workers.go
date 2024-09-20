package repositories

import (
	"context"
	db "texApi/database"
	"texApi/internal/_other/queries"
	"texApi/internal/_other/schemas/request"
	"texApi/internal/_other/schemas/response"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func GetWorkers(offset, limit int) ([]response.Worker, error) {
	var workers []response.Worker

	err := pgxscan.Select(
		context.Background(), db.DB,
		&workers, queries.GetWorkers, offset, limit,
	)

	if err != nil {
		return nil, err
	}

	return workers, nil
}

func GetWorker(id int) response.Worker {
	var worker response.Worker

	db.DB.QueryRow(
		context.Background(), queries.GetWorker, id,
	).Scan(
		&worker.ID,
		&worker.Fullname,
		&worker.Phone,
		&worker.Address,
		&worker.Photo,
		&worker.AboutSelf,
		&worker.Services,
		&worker.CreatedAt,
		&worker.UpdatedAt,
	)

	return worker
}

func CreateWorker(worker request.CreateWorker) (int, error) {
	var workerID int

	db.DB.QueryRow(
		context.Background(), queries.CreateWorker,
		worker.Fullname, worker.Phone, worker.Address, worker.AboutSelf,
		worker.Password, time.Now(), time.Now(),
	).Scan(&workerID)

	for _, serviceID := range worker.Services {
		_, err := db.DB.Exec(
			context.Background(), queries.CreateWorkerService,
			workerID, serviceID,
		)

		if err != nil {
			return workerID, err
		}
	}

	return workerID, nil
}

func CheckWorkerExist(phone string) string {
	var existPhone string

	db.DB.QueryRow(
		context.Background(), queries.CheckWorkerExist, phone,
	).Scan(&existPhone)

	return existPhone
}

func UpdateWorker(worker request.UpdateWorker) error {
	if worker.Password == "" {
		_, err := db.DB.Exec(
			context.Background(), queries.UpdateWorkerWithoutPassword,
			worker.Fullname, worker.Phone, worker.Address, worker.AboutSelf,
			time.Now(), worker.ID,
		)

		if err != nil {
			return err
		}
	}

	_, err := db.DB.Exec(
		context.Background(), queries.UpdateWorker,
		worker.Fullname, worker.Phone, worker.Address, worker.AboutSelf,
		worker.Password, time.Now(), worker.ID,
	)

	if err != nil {
		return err
	}

	for _, service := range worker.Services {
		_, err := db.DB.Exec(
			context.Background(), queries.UpdateWorkerService,
			service.NextID, worker.ID, service.PrevID,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteWorker(id int) error {
	_, err := db.DB.Exec(context.Background(), queries.DeleteWorker, id)

	if err != nil {
		return err
	}

	return nil
}

func SetWorkerImage(id int, photo string) error {
	photoUrl := "/uploads/workers/" + photo

	_, err := db.DB.Exec(
		context.Background(), queries.SetWorkerImage, photoUrl, id,
	)

	if err != nil {
		return err
	}

	return nil
}
