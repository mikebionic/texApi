Вот универсальное решение таблицы котировок для логистической компании, подходящее для автоперевозок, авиа, ж/д и морского транспорта.

---

## 📘 **Документация по таблице котировок цен**

|Колонка|Назначение / Описание|Тип данных|Пример|
|---|---|---|---|
|`id`|Уникальный идентификатор котировки|Целое число|1|
|`transport_type`|Тип транспорта: `авто`, `авиа`, `ж/д`, `море`|Текст|авто|
|`route_from`|Город или страна отправления|Текст|Ашхабад|
|`route_to`|Город или страна назначения|Текст|Стамбул|
|`price_value`|Значение цены (без валюты)|Число|5000|
|`price_currency`|Валюта цены|Текст|USD|
|`price_unit`|За что считается цена: `рейс`, `тонна`, `куб.м`, `паллет`, `контейнер`, `погрузка` и т.д.|Текст|тонна|
|`min_weight_ton`|Минимальный тоннаж (если цена за вес/объем)|Число (с плавающей точкой)|1.5|
|`max_weight_ton`|Максимальный тоннаж|Число (с плавающей точкой)|20|
|`valid_until`|До какой даты действует цена|Дата|2025-07-15|
|`company_name`|Название компании, предоставившей котировку|Текст|"ABC Logistics"|
|`contact_info`|Контактные данные: email, телефон, Telegram|Текст|`logistics@abc.com`|
|`extra_costs`|Дополнительные расходы: погрузка, простои, брокерские и т.д. (коротко, описание)|Текст|Погрузка 50$, оформление 100$|
|`comment`|Любые комментарии или пояснения|Текст|Цена указана без учета НДС|
|`is_active`|Актуальна ли сейчас эта цена (удобно для фильтрации неактивных)|Булево (TRUE/FALSE)|TRUE|
|`created_at`|Дата добавления записи|Дата|2025-07-28|
|`updated_at`|Дата последнего обновления записи|Дата|2025-07-28|

---

## 🧾 Примечания

- `price_unit` поможет фильтровать котировки по назначению: например, если клиент спрашивает цену "за рейс" — исключаем "за тонну".
    
- Вес (`min_weight_ton`, `max_weight_ton`) помогает понять применимость тарифа.
    
- `valid_until` важен, если тарифы обновляются часто.
    
- `is_active` позволяет сохранять архивные записи, не удаляя их.
    
- `extra_costs` — колонка для указания условий, которые влияют на финальную стоимость.
    
- `comment` позволяет вставить любую текстовую ремарку — например, "не включает стоимость топлива".
    

---

## 📥 финальная таблица 

| id  | transport_type | route_from | route_to | price_value | price_currency | price_unit | min_weight_ton | max_weight_ton | valid_until | company_name  | contact_info      | extra_costs          | comment                      | is_active | created_at | updated_at |
| --- | -------------- | ---------- | -------- | ----------- | -------------- | ---------- | -------------- | -------------- | ----------- | ------------- | ----------------- | -------------------- | ---------------------------- | --------- | ---------- | ---------- |
| 1   | авто           | Ашхабад    | Стамбул  | 5000        | USD            | рейс       |                |                | 2025-07-15  | ABC Logistics | logistics@abc.com | оформление 100$      | без учета НДС                | TRUE      | 2025-07-28 | 2025-07-28 |
| 2   | авиа           | Ашхабад    | Москва   | 4000        | USD            | тонна      | 0.5            | 2              | 2025-08-20  | Sky Air       | sales@skyair.tj   | аэропортный сбор 50$ | зависит от расписания        | TRUE      | 2025-07-28 | 2025-07-28 |
| 3   | ж/д            | Тегеран    | Душанбе  |             |                |            |                |                |             | RailCo        | +998901234567     |                      | требуется уточнение маршрута | FALSE     | 2025-07-28 | 2025-07-28 |

Вот полный CSV-шаблон с:
1. Первая строка — заголовки (человеческие названия)
2. Вторая строка — SQL-названия колонок (snake_case)
3. Несколько строк с **примером данных**

```
ID,Тип транспорта,Маршрут (откуда-куда),Страна отправления,Город отправления,Страна назначения,Город назначения,Тип груза,Цена,Валюта,Тип расчета (за что цена),Вес (тонны),Объем (м³),Срок действия (до),Время в пути (дней),Частота рейсов,Наличие доп. услуг,Комментарий / Примечание
id,transport_type,route,country_from,city_from,country_to,city_to,cargo_type,price,currency,rate_type,weight_ton,volume_m3,valid_until,transit_days,frequency,extra_services,notes
1,Авто,Ашхабад - Стамбул,Туркменистан,Ашхабад,Турция,Стамбул,Строительные материалы,5000,USD,за рейс,20,60,2025-07-15,7,еженедельно,Страховка; Отслеживание,Нужна предоплата
2,Авиа,Ашхабад - Москва,Туркменистан,Ашхабад,Россия,Москва,Медикаменты,4000,USD,за тонну,1,2,2025-08-20,2,по запросу,Экспресс-доставка,Уточнить документы
3,Ж/Д,Тегеран - Душанбе,Иран,Тегеран,Таджикистан,Душанбе,Промышленное оборудование,9000,EUR,за м³,10,30,2025-09-01,10,2 раза в месяц,,Безопасный маршрут
```

---

📌 **Как использовать:**
- Строка 2 (snake_case) — это ключи для SQL-таблицы, JSON API или backend.  
- Строка 1 — то, что отображается пользователю.
- Данные можно импортировать в Excel, Google Sheets, PostgreSQL и др.

----
# Price Quote Analysis API

## Overview

The Price Quote Analysis API provides intelligent price suggestions by analyzing both real market offers and historical price quotes. It finds similar transportation requests based on multiple criteria and returns statistical pricing data with detailed analysis information.

## Endpoint

```
GET /texApp/price-quote/analyze
```

## Features

- **Smart Matching**: Finds similar offers based on geographic, transport, and logistic criteria
- **Dual Data Sources**: Uses real offer data as primary source, price quotes as fallback
- **Statistical Analysis**: Provides min/max/average pricing with sample size
- **Flexible Matching**: Supports strict and flexible matching modes
- **Market Intelligence**: Returns actual market data rather than static estimates

## Query Parameters

### Geographic Filters

| Parameter         | Type   | Description                            | Example   |
| ----------------- | ------ | -------------------------------------- | --------- |
| `from_country_id` | int    | Source country ID                      | `1`       |
| `from_city_id`    | int    | Source city ID                         | `101`     |
| `to_country_id`   | int    | Destination country ID                 | `2`       |
| `to_city_id`      | int    | Destination city ID                    | `201`     |
| `from_country`    | string | Source country name (alternative)      | `USA`     |
| `to_country`      | string | Destination country name (alternative) | `Canada`  |
| `from_region`     | string | Source region/state                    | `Texas`   |
| `to_region`       | string | Destination region/province            | `Ontario` |

### Transport Details

|Parameter|Type|Description|Example|
|---|---|---|---|
|`transport_type`|string|Transport mode|`auto`, `avia`, `rail`, `sea`|
|`sub_type`|string|Transport subtype|`FTL`, `LTL`, `refrigerated`, `express`|
|`vehicle_type_id`|int|Vehicle type identifier|`3`|
|`packaging_type_id`|int|Packaging type identifier|`2`|

### Pricing & Logistics

|Parameter|Type|Description|Example|
|---|---|---|---|
|`currency`|string|Currency code|`USD`, `EUR`, `CAD`|
|`payment_method`|string|Payment method|`cash`, `bank_transfer`, `credit_card`|
|`distance`|int|Distance in kilometers|`1250`|
|`distance_km`|int|Distance in kilometers (alternative)|`1250`|
|`min_volume`|float|Minimum volume/weight capacity|`500.5`|
|`max_volume`|float|Maximum volume/weight capacity|`2000.0`|

### Service Options

|Parameter|Type|Description|Example|
|---|---|---|---|
|`fuel_included`|boolean|Whether fuel costs are included|`true`, `false`|
|`customs_included`|boolean|Whether customs fees are included|`true`, `false`|
|`match_strict`|boolean|Matching mode (strict vs flexible)|`true`, `false`|

## Matching Logic

### Strict Mode (`match_strict=true`)

- Exact parameter matching
- Distance tolerance: ±10km
- Higher precision, fewer results

### Flexible Mode (`match_strict=false`)

- Allows partial matches
- Distance tolerance: ±50km
- Includes region-based matching
- More results, broader coverage

## Response Format

```json
{
  "success": true,
  "message": "Price analysis completed successfully",
  "data": {
    // Standard PriceQuote fields
    "transport_type": "auto",
    "currency": "USD",
    "from_country": "USA",
    "to_country": "Canada",
    "average_price": 2687.50,
    "min_price": 2100.00,
    "max_price": 3275.00,
    "cost_per_km": 2.15,
    "sample_size": 18,
    "data_source": "offer_based",
    "is_dynamic": true,
    "notes": "Found average among 18 offers, minimum price 2100.00 and maximum 3275.00",
    
    // Analysis-specific information
    "analysis_info": {
      "found_from_offers": true,
      "offer_count": 18,
      "price_quote_count": 0,
      "offer_min_price": 2100.00,
      "offer_max_price": 3275.00,
      "offer_avg_cost_per_km": 2.15,
      "matching_criteria": [
        "from_country",
        "to_country", 
        "vehicle_type_id",
        "currency",
        "distance"
      ]
    }
  }
}
```

## Response Fields

### Analysis Info Object

|Field|Type|Description|
|---|---|---|
|`found_from_offers`|boolean|Whether data came from offers (true) or price quotes (false)|
|`offer_count`|int|Number of matching offers found|
|`price_quote_count`|int|Number of matching price quotes found|
|`offer_min_price`|float|Minimum price from offers|
|`offer_max_price`|float|Maximum price from offers|
|`offer_avg_cost_per_km`|float|Average cost per kilometer from offers|
|`matching_criteria`|array|List of criteria used for matching|

### Data Source Types

|Value|Description|
|---|---|
|`offer_based`|Calculated from real market offers|
|`price_quote_based`|Based on historical price quotes|
|`no_match`|No matching data found|

## Error Responses

### 400 Bad Request

```json
{
  "success": false,
  "message": "Invalid query parameters",
  "error": "Invalid transport_type value"
}
```

### 500 Internal Server Error

```json
{
  "success": false,
  "message": "Failed to analyze price quotes with offers",
  "error": "Database connection error"
}
```

## Use Cases

1. **Price Estimation**: Get market-based price estimates for new transport requests
2. **Market Analysis**: Understand pricing trends in specific routes or transport types
3. **Competitive Intelligence**: Compare prices against market averages
4. **Dynamic Pricing**: Implement real-time pricing based on market conditions
5. **Quote Generation**: Generate accurate quotes for customers based on similar requests

## Best Practices

1. **Start Simple**: Use basic filters first, then add specific criteria
2. **Use Flexible Mode**: For broader market coverage, use `match_strict=false`
3. **Check Sample Size**: Higher `sample_size` indicates more reliable data
4. **Consider Data Source**: `offer_based` is more current than `price_quote_based`
5. **Handle No Matches**: Always check if matches were found and handle empty results
6. **Cache Results**: Consider caching for frequently requested routes

## Rate Limits

- No specific rate limits currently implemented
- Recommend reasonable usage patterns for production environments

## Authentication

- Optional: Include `Authorization: Bearer <token>` header if authentication is required
- Check with your system administrator for specific authentication requirements