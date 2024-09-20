package repositories

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	db "texApi/database"
	"texApi/internal/_other/queries"
	"texApi/internal/_other/schemas/request"
	"texApi/internal/_other/schemas/response"
)

func GetSubscriptions() ([]response.Subscription, error) {
	var subscriptions []response.Subscription

	err := pgxscan.Select(
		context.Background(), db.DB,
		&subscriptions, queries.GetSubscriptions,
	)

	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func GetSubscription(id int) response.Subscription {
	var subscription response.Subscription

	db.DB.QueryRow(
		context.Background(), queries.GetSubscription, id,
	).Scan(
		&subscription.ID,
		&subscription.Title,
		&subscription.Description,
		&subscription.StartAt,
		&subscription.EndAt,
		&subscription.Days,
		&subscription.Count,
		&subscription.Price,
	)

	return subscription
}

func CreateSubscription(subscription request.CreateSubscription) (int, error) {
	var id int

	db.DB.QueryRow(
		context.Background(), queries.CreateSubscription,
		subscription.StartAt, subscription.EndAt, subscription.Days,
		subscription.Count, subscription.Price,
	).Scan(&id)

	_, err := db.DB.Exec(
		context.Background(), queries.CreateSubscriptionTranslates,
		subscription.Title.TK, subscription.Desc.TK, 1,
		subscription.Title.RU, subscription.Desc.RU, 2,
		subscription.Title.EN, subscription.Desc.EN, 3,
		id,
	)

	if err != nil {
		return id, err
	}

	return id, nil
}

func UpdateSubscription(subscription request.UpdateSubscription) error {
	_, err := db.DB.Exec(
		context.Background(), queries.UpdateSubscription,
		subscription.StartAt, subscription.EndAt, subscription.Days,
		subscription.Count, subscription.Price, subscription.ID,
	)

	if err != nil {
		return err
	}

	_, errTK := db.DB.Exec(
		context.Background(), queries.UpdateSubscriptionTranslates,
		subscription.Title.TK, subscription.Desc.TK, 1, subscription.ID,
	)

	if errTK != nil {
		return errTK
	}

	_, errRU := db.DB.Exec(
		context.Background(), queries.UpdateSubscriptionTranslates,
		subscription.Title.RU, subscription.Desc.RU, 2, subscription.ID,
	)

	if errRU != nil {
		return errTK
	}

	_, errEN := db.DB.Exec(
		context.Background(), queries.UpdateSubscriptionTranslates,
		subscription.Title.EN, subscription.Desc.EN, 3, subscription.ID,
	)

	if errEN != nil {
		return errTK
	}

	return nil
}

func DeleteSubscription(id int) error {
	_, err := db.DB.Exec(context.Background(), queries.DeleteSubscription, id)

	if err != nil {
		return err
	}

	return nil
}

func GetSubscriptionPrice(id int) float32 {
	var price float32

	db.DB.QueryRow(
		context.Background(), queries.GetSubscriptionPrice, id,
	).Scan(&price)

	return price
}
