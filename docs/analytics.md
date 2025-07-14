# Analytics System Documentation

## Overview

The Analytics System is a comprehensive tracking and reporting solution designed to monitor key business metrics including user growth, offer activity, and revenue generation. It provides automated data collection, filtering capabilities, and administrative controls for managing analytics operations.

| Name                   | key                 | description                                                                                        |
|------------------------|---------------------|----------------------------------------------------------------------------------------------------|
| Общие пользователи     | "user_all"          | (all tbl_user COUNT)                                                                               |
| Грузоотправители       | "user_sender"       | (all tbl_user role: sender COUNT)                                                                  |
| Перевозчики            | "user_carrier"      | (all tbl_user role: carrier COUNT)                                                                 |
| Новые Грузоотправители | "user_sender_new"   | (tbl_user with role=sender, tbl_user where id is greater than last_user_id )                       |
| Новые Перевозчики      | "user_carrier_new"  | (tbl_user with role=carrier, tbl_user where id is greater than last_user_id )                      |
| Новые заявки           | "offer_new_sender"  | (tbl_offer where id is greater than last_offer_id & role: sender)                                  |
| новые перевозки        | "offer_new_carrier" | (tbl_offer where id is greater than last_offer_id & role: carrier)                                 |
| Кол-во сделок          | "offer_all"         | (all tbl_offer with offer_state != deleted, pending, disabled)                                     |
| Активные заявки        | "offer_active"      | ( offer-state : active, enabled, working)                                                          |
| заявки в ожидании      | "offer_pending"     | ( offer with state: pending)                                                                       |
| Завершенные сделки     | "offer_completed"   | (offer with state: completed, archived. tbl_offer.updated_at > last_completed_offer_id.updated_at) |
| Заявки без отклика     | "offer_no_response" | (ALL tbl_offer COUNT where tbl_offer_response count = 0 )                                          |

## Architecture

### Core Components

1. **Data Collection Service** (`GenerateAnalytics`) - Automated metrics calculation
2. **Query Service** (`GetAnalytics`) - Data retrieval with filtering and pagination
3. **Configuration Service** - System settings management
4. **Scheduler** - Automated periodic data generation

## Analytics Metrics

### User Metrics

#### Total Users (`user_all`)
- **Description**: Total count of all active users in the system
- **Query**: `SELECT COUNT(*) FROM tbl_user WHERE deleted = 0 AND active = 1`
- **Purpose**: Track overall user base size

#### Sender Users (`user_sender`)
- **Description**: Count of users with role 'sender' (cargo senders/shippers)
- **Query**: `SELECT COUNT(*) FROM tbl_user WHERE deleted = 0 AND active = 1 AND role = 'sender'`
- **Purpose**: Track supply side of the marketplace

#### Carrier Users (`user_carrier`)
- **Description**: Count of users with role 'carrier' (transportation providers)
- **Query**: `SELECT COUNT(*) FROM tbl_user WHERE deleted = 0 AND active = 1 AND role = 'carrier'`
- **Purpose**: Track demand side of the marketplace

#### New Sender Users (`user_sender_new`)
- **Description**: Count of new sender users since last analytics run
- **Query**: `SELECT COUNT(*) FROM tbl_user WHERE deleted = 0 AND active = 1 AND role = 'sender' AND id > {last_user_id}`
- **Purpose**: Track sender user growth between periods

#### New Carrier Users (`user_carrier_new`)
- **Description**: Count of new carrier users since last analytics run
- **Query**: `SELECT COUNT(*) FROM tbl_user WHERE deleted = 0 AND active = 1 AND role = 'carrier' AND id > {last_user_id}`
- **Purpose**: Track carrier user growth between periods

### Offer Metrics

#### Total Offers (`offer_all`)
- **Description**: Count of all active offers (excluding deleted, pending, disabled)
- **Query**: `SELECT COUNT(*) FROM tbl_offer WHERE deleted = 0 AND offer_state NOT IN ('deleted', 'pending', 'disabled')`
- **Purpose**: Track total marketplace activity

#### Active Offers (`offer_active`)
- **Description**: Count of offers in active working states
- **Query**: `SELECT COUNT(*) FROM tbl_offer WHERE deleted = 0 AND offer_state IN ('active', 'enabled', 'working')`
- **Purpose**: Track current marketplace liquidity

#### Pending Offers (`offer_pending`)
- **Description**: Count of offers awaiting approval or action
- **Query**: `SELECT COUNT(*) FROM tbl_offer WHERE deleted = 0 AND offer_state = 'pending'`
- **Purpose**: Track offer backlog

#### Completed Offers (`offer_completed`)
- **Description**: Count of successfully completed transactions
- **Query**: `SELECT COUNT(*) FROM tbl_offer WHERE deleted = 0 AND offer_state IN ('completed', 'archived')`
- **Purpose**: Track successful transaction volume

#### Offers Without Response (`offer_no_response`)
- **Description**: Count of offers that haven't received any responses
- **Query**: `SELECT COUNT(*) FROM tbl_offer o LEFT JOIN tbl_offer_response r ON o.id = r.offer_id WHERE o.deleted = 0 AND r.id IS NULL`
- **Purpose**: Track market efficiency and response rates

#### New Sender Offers (`offer_new_sender`)
- **Description**: Count of new offers from senders since last analytics run
- **Query**: `SELECT COUNT(*) FROM tbl_offer WHERE deleted = 0 AND offer_role = 'sender' AND id > {last_offer_id}`
- **Purpose**: Track sender activity growth

#### New Carrier Offers (`offer_new_carrier`)
- **Description**: Count of new offers from carriers since last analytics run
- **Query**: `SELECT COUNT(*) FROM tbl_offer WHERE deleted = 0 AND offer_role = 'carrier' AND id > {last_offer_id}`
- **Purpose**: Track carrier activity growth

### Financial Metrics

#### Total Revenue (`total_revenue`)
- **Description**: Sum of revenue from all completed offers
- **Query**: `SELECT SUM(cost_per_km * distance) FROM tbl_offer WHERE deleted = 0 AND offer_state = 'completed'`
- **Purpose**: Track total platform revenue

#### Average Cost Per Kilometer (`average_cost_per_km`)
- **Description**: Average pricing across all offers
- **Query**: `SELECT AVG(cost_per_km) FROM tbl_offer WHERE deleted = 0 AND cost_per_km > 0`
- **Purpose**: Track pricing trends

#### Total Distance (`total_distance`)
- **Description**: Sum of distances from completed offers
- **Query**: `SELECT SUM(distance) FROM tbl_offer WHERE deleted = 0 AND offer_state = 'completed'`
- **Purpose**: Track total service volume

#### Active Companies (`active_companies`)
- **Description**: Count of distinct companies with active offers
- **Query**: `SELECT COUNT(DISTINCT company_id) FROM tbl_offer WHERE deleted = 0 AND offer_state IN ('active', 'working')`
- **Purpose**: Track corporate engagement

## Data Generation Process

### Automated Scheduling

The system uses a configurable scheduler that:

1. **Checks Configuration**: Verifies if analytics generation is enabled
2. **Calculates Interval**: Determines if enough time has passed since last run
3. **Generates Metrics**: Calls `GenerateAnalytics()` function
4. **Updates Tracking**: Records last run time and next scheduled run

### Generation Workflow

1. **Baseline Retrieval**: Get last analytics record for incremental calculations
2. **Period Definition**: Set 24-hour analysis period (configurable)
3. **Metric Calculation**: Execute all helper functions to gather metrics
4. **Data Storage**: Insert new analytics record with timestamp
5. **Configuration Update**: Update last run time in configuration

### Key Features

#### Incremental Tracking
- Uses `last_user_id`, `last_offer_id`, `last_completed_offer_id` for delta calculations
- Enables tracking of "new" items since last analytics run
- Prevents double-counting and ensures accurate growth metrics

#### Period-based Analysis
- Each record represents a specific time period (default: 24 hours)
- `period_start` and `period_end` define the analysis window
- Enables time-series analysis and trend identification

#### Soft Delete Handling
- All queries respect `deleted = 0` flag
- Ensures accuracy by excluding logically deleted records
- Maintains data integrity across all metrics

## API Endpoints

### Data Retrieval

#### `GET /analytics/`
Retrieves analytics data with comprehensive filtering options.

**Query Parameters:**
- **Date Filters**: `date_from`, `date_to` - Filter by record creation date
- **Period Filters**: `period_start`, `period_end` - Filter by analysis period
- **User Filters**: `user_all_min/max`, `user_sender_min/max`, `user_carrier_min/max`
- **Offer Filters**: `offer_all_min/max`, `offer_active_min/max`, `offer_pending_min/max`, `offer_completed_min/max`
- **Revenue Filters**: `revenue_min`, `revenue_max`
- **Sorting**: `order_by`, `order_dir`
- **Pagination**: `page`, `per_page`

**Response Structure:**
```json
{
  "message": "Analytics",
  "success": true,
  "data": {
    "total": 1,
    "page": 1,
    "per_page": 10,
    "stats": {
      "total_records": 1,
      "avg_users_per_period": 19,
      "avg_offers_per_period": 14,
      "total_revenue": 0,
      "growth_rate": 0,
      "last_update": "2025-07-15T04:23:18.65493Z"
    },
    "data": [/* analytics records */]
  }
}
```

### Administrative Operations

#### `POST /analytics/admin/generate/`
Manually triggers analytics generation.

**Purpose**:
- Emergency data generation
- Testing analytics collection
- Immediate metric updates

#### `GET /analytics/admin/config/`
Retrieves current analytics configuration.

**Returns**:
- `enabled`: Whether analytics generation is active
- `log_interval_days`: Days between automatic generations
- `last_analytics_run`: Timestamp of last generation

#### `PUT /analytics/admin/config/`
Updates analytics configuration.

**Parameters**:
- `log_interval_days`: Integer (1-365) - Generation frequency
- `enabled`: Boolean - Enable/disable automatic generation

## Configuration Management

### Database Configuration (`tbl_analytics_config`)

The system uses database-stored configuration with the following keys:

- **`enabled`**: Controls whether automatic analytics generation is active
- **`log_interval_days`**: Defines the frequency of automatic generation (default: 1 day)
- **`last_analytics_run`**: Timestamp of the last successful generation

### Scheduler Integration

The scheduler monitors configuration changes and:
- Adjusts generation intervals dynamically
- Respects enable/disable settings
- Maintains accurate run timing

## Performance Considerations

### Query Optimization

1. **Indexed Fields**: Ensure proper indexing on frequently queried fields
2. **Soft Delete Performance**: Consider compound indexes on `deleted` + other fields
3. **Date Range Queries**: Index `created_at`, `period_start`, `period_end` fields

### Data Volume Management

1. **Pagination**: All queries support pagination to handle large datasets
2. **Filtering**: Comprehensive filtering reduces data transfer
3. **Archival Strategy**: Consider archiving old analytics records based on retention policies

## Error Handling

### Generation Failures
- Logs detailed error information
- Maintains system state consistency
- Allows manual retry via admin endpoint

### Query Failures
- Returns appropriate HTTP status codes
- Provides descriptive error messages
- Maintains API stability

## Security

### Authentication
- All endpoints require admin authentication
- Bearer token authentication implemented
- Role-based access control

### Data Protection
- No sensitive user data exposed in analytics
- Aggregated metrics only
- Audit trail for configuration changes

## Monitoring and Alerting

### Key Metrics to Monitor

1. **Generation Success Rate**: Track failed analytics generations
2. **Data Freshness**: Monitor time since last successful generation
3. **Performance**: Track query execution times
4. **Data Quality**: Validate metric consistency

### Recommended Alerts

1. **Generation Failures**: Alert on failed automatic generation
2. **Stale Data**: Alert if analytics haven't been generated within expected timeframe
3. **Performance Degradation**: Alert on slow query performance
4. **Configuration Changes**: Audit configuration modifications

## Usage Examples

### Basic Analytics Retrieval
```bash
GET /api/v1/analytics/
```

### Filtered Analytics (Last 30 Days)
```bash
GET /api/v1/analytics/?date_from=2025-06-15T00:00:00Z&date_to=2025-07-15T23:59:59Z&order_by=created_at&order_dir=desc
```

### Revenue Analysis
```bash
GET /api/v1/analytics/?revenue_min=1000&order_by=total_revenue&order_dir=desc
```

### Manual Generation
```bash
POST /api/v1/analytics/admin/generate/
```

### Configuration Update
```bash
PUT /api/v1/analytics/admin/config/
Content-Type: application/json

{
  "log_interval_days": 7,
  "enabled": true
}
```

## Best Practices

1. **Regular Monitoring**: Check analytics generation status regularly
2. **Backup Strategy**: Implement backup procedures for analytics data
3. **Performance Testing**: Test queries with production-like data volumes
4. **Configuration Management**: Use version control for configuration changes
5. **Documentation**: Keep metric definitions updated as business requirements change