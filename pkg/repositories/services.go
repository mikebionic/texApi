package repositories

import (
	"context"
	"encoding/json"
	db "texApi/database"
	"texApi/pkg/queries"
	"texApi/pkg/schemas/request"
	"texApi/pkg/schemas/response"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func GetServices() ([]response.Service, error) {
	var services []response.Service

	err := pgxscan.Select(
		context.Background(), db.DB,
		&services, queries.GetServices,
	)

	if err != nil {
		return nil, err
	}

	serviceMap := make(map[int]*response.Service)

	for i := range services {
		serviceMap[services[i].ID] = &services[i]
	}

	var rootServices []*response.Service
	for i := range services {
		service := &services[i]
		if service.ParentID == nil {
			rootServices = append(rootServices, service)
		} else {
			parent := serviceMap[*service.ParentID]
			if parent != nil {
				parent.Children = append(parent.Children, service)
			}
		}
	}

	jsonData, err := json.MarshalIndent(rootServices, "", "  ")

	if err != nil {
		return nil, err
	}

	var result []response.Service
	json.Unmarshal(jsonData, &result)

	return result, nil
}

func GetServiceList() ([]response.ServiceList, error) {
	var services []response.ServiceList

	err := pgxscan.Select(
		context.Background(), db.DB,
		&services, queries.GetServiceList,
	)

	if err != nil {
		return nil, err
	}

	return services, nil
}

func GetService(id int) response.Service {
	var service response.Service

	db.DB.QueryRow(
		context.Background(), queries.GetService, id,
	).Scan(
		&service.ID,
		&service.Title,
		&service.Image,
		&service.ParentID,
	)

	return service
}

func CreateService(service request.CreateService) (int, error) {
	var id int

	db.DB.QueryRow(
		context.Background(), queries.CreateService, service.ParentID,
	).Scan(&id)

	_, err := db.DB.Exec(
		context.Background(), queries.CreateServiceTranslates,
		service.Title.TK, 1,
		service.Title.RU, 2,
		service.Title.EN, 3,
		id,
	)

	if err != nil {
		return id, err
	}

	return id, nil
}

func UpdateService(service request.UpdateService) error {
	if service.ParentID != 0 {
		_, err := db.DB.Exec(
			context.Background(), queries.UpdateService,
			service.ParentID, service.ID,
		)

		if err != nil {
			return err
		}
	}

	_, errTK := db.DB.Exec(
		context.Background(), queries.UpdateServiceTranslates,
		service.Title.TK, 1, service.ID,
	)

	if errTK != nil {
		return errTK
	}

	_, errRU := db.DB.Exec(
		context.Background(), queries.UpdateServiceTranslates,
		service.Title.RU, 2, service.ID,
	)

	if errRU != nil {
		return errTK
	}

	_, errEN := db.DB.Exec(
		context.Background(), queries.UpdateServiceTranslates,
		service.Title.EN, 3, service.ID,
	)

	if errEN != nil {
		return errTK
	}

	return nil
}

func DeleteService(id int) error {
	_, err := db.DB.Exec(context.Background(), queries.DeleteService, id)

	if err != nil {
		return err
	}

	return nil
}

func SetServiceImage(id int, image string) error {
	imageUrl := "/uploads/services/" + image

	_, err := db.DB.Exec(
		context.Background(), queries.SetServiceImage, imageUrl, id,
	)

	if err != nil {
		return err
	}

	return nil
}
