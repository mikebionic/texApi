# Call Room API Documentation

## Overview
git 
The Call Room API provides secure integration with Jitsi for creating and managing video/audio calls within conversations. It generates unique call rooms, validates user access, and handles call lifecycle management.

**Creating a Call:**
Create/have conversation → 2. Authenticate → 3. POST to create call room → 4. Auto WebSocket notification → 5. Use Jitsi URL

**Receiving a Call:**
Get WebSocket notification → 2. Validate access via join endpoint (can be omitted) → 3. Open Jitsi URL → 4. Join video session


## Database Schema

```sql
CREATE TABLE tbl_call_room
(
    id              SERIAL PRIMARY KEY,
    uuid            UUID         NOT NULL DEFAULT gen_random_uuid(),
    conversation_id INT          NOT NULL DEFAULT 0,
    max_user        INT          NOT NULL DEFAULT 2,
    user_ids        TEXT[]       NOT NULL DEFAULT '{}',
    profile_ids     TEXT[]       NOT NULL DEFAULT '{}',
    title           VARCHAR(200) NOT NULL DEFAULT '',
    hex             VARCHAR(60)  NOT NULL DEFAULT '',
    duration        INT          NOT NULL DEFAULT 60, -- 1h
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active          INT          NOT NULL DEFAULT 1,
    deleted         INT          NOT NULL DEFAULT 0
);
```

## Data Structures

```go
type CallRoom struct {
    ID             int       `json:"id"`
    UUID           string    `json:"uuid"`
    ConversationID int       `json:"conversation_id"`
    MaxUser        int       `json:"max_user"`
    UserIDs        []string  `json:"user_ids"`
    ProfileIDs     []string  `json:"profile_ids"`
    Title          string    `json:"title"`
    Hex            string    `json:"hex"`
    Duration       int       `json:"duration"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
    Active         int       `json:"active"`
    Deleted        int       `json:"deleted"`
    JoinURL        string    `json:"join_url,omitempty"`
    JitsiURL       string    `json:"jitsi_url,omitempty"`
}

type CreateCallRoomRequest struct {
    UserIDs        []string `json:"user_ids" binding:"required"`
    ConversationID int      `json:"conversation_id" binding:"required"`
    ProfileIDs     []string `json:"profile_ids"`
    Title          string   `json:"title"`
    Duration       int      `json:"duration"`
    MaxUser        int      `json:"max_user"`
}
```

## API Endpoints

### 1. Create Call Room
**POST** `/api/v1/call-room/create/`

Creates a new call room and sends invitation via WebSocket.

**Headers:**
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

**Request Body:**
```json
{
    "user_ids": ["3"],
    "conversation_id": 11,
    "profile_ids": ["3"],
    "title": "Voice call",
    "duration": 60,
    "max_user": 2
}
```

**Response (201):**
```json
{
    "message": "Call room created",
    "data": {
        "id": 1,
        "uuid": "550e8400-e29b-41d4-a716-446655440000",
        "conversation_id": 11,
        "max_user": 2,
        "user_ids": ["3", "1"],
        "profile_ids": ["3", "1"],
        "title": "Voice call",
        "hex": "a1b2c3d4e5f6789012345678901234",
        "duration": 60,
        "created_at": "2025-01-15T10:30:00Z",
        "updated_at": "2025-01-15T10:30:00Z",
        "active": 1,
        "deleted": 0,
        "join_url": "https://api.example.com/api/v1/call-room/join/550e8400-e29b-41d4-a716-446655440000",
        "jitsi_url": "https://jitsi.example.com/a1b2c3d4e5f6789012345678901234"
    }
}
```

**Automatic Behavior:**
- Requesting user's ID and profile ID are automatically added to lists if not present
- Default duration: 60 minutes if not specified
- Default max_user: number of users in user_ids if not specified
- WebSocket notification automatically sent to conversation participants

### 2. Join Call Room
**GET** `/api/v1/call-room/join/{uuid}`

Validates user access and returns call room details.

**Headers:**
```
Authorization: Bearer {access_token}
```

**Path Parameters:**
- `uuid`: Call room UUID

**Response (200):**
```json
{
    "message": "Call room access granted",
    "data": {
        "id": 1,
        "uuid": "550e8400-e29b-41d4-a716-446655440000",
        "conversation_id": 11,
        "max_user": 2,
        "user_ids": ["3", "1"],
        "profile_ids": ["3", "1"],
        "title": "Voice call",
        "hex": "a1b2c3d4e5f6789012345678901234",
        "duration": 60,
        "created_at": "2025-01-15T10:30:00Z",
        "updated_at": "2025-01-15T10:30:00Z",
        "active": 1,
        "deleted": 0,
        "join_url": "https://api.example.com/api/v1/call-room/join/550e8400-e29b-41d4-a716-446655440000",
        "jitsi_url": "https://jitsi.example.com/a1b2c3d4e5f6789012345678901234"
    }
}
```

### 3. End Call Room
**POST** `/api/v1/call-room/end/{uuid}`

Ends the call room (sets active=0, deleted=1).

**Headers:**
```
Authorization: Bearer {access_token}
```

**Path Parameters:**
- `uuid`: Call room UUID

**Response (200):**
```json
{
    "message": "Call room ended",
    "data": null
}
```

## Complete Workflow

### 1. Creating a Call

**Prerequisites:**
- User must be authenticated with valid access token
- User must have access to the conversation

**Steps:**
1. **Create Call Room**:
   ```bash
   curl -X POST /api/v1/call-room/create/ \
     -H "Authorization: Bearer {token}" \
     -H "Content-Type: application/json" \
     -d '{
       "user_ids": ["3"],
       "conversation_id": 11,
       "profile_ids": ["3"],
       "title": "Voice call",
       "duration": 60,
       "max_user": 2
     }'
   ```

2. **System Actions**:
    - Validates user has access to conversation_id 11
    - Adds requesting user to user_ids and profile_ids automatically
    - Generates unique 30-character hex identifier
    - Creates database record
    - Sends WebSocket notification to conversation participants
    - Returns call room data with join_url and jitsi_url

3. **WebSocket Notification Sent**:
   ```json
   {
     "type": "notification",
     "content": "https://jitsi.example.com/a1b2c3d4e5f6789012345678901234",
     "conversation_id": 11,
     "sender_id": 1,
     "created_at": "2025-01-15T10:30:00Z"
   }
   ```

### 2. Receiving a Call

**Steps:**
1. **Receive WebSocket Notification** - All conversation participants get call invitation
2. **Extract Call Information** - Parse the notification message
3. **Validate Access** (Optional but recommended):
   ```bash
   curl -X GET /api/v1/call-room/join/550e8400-e29b-41d4-a716-446655440000 \
     -H "Authorization: Bearer {token}"
   ```
4. **Access Validation**:
    - Checks if call room is active (active=1, deleted=0)
    - Validates user_id is in user_ids OR profile_id is in profile_ids
    - Returns full call room data if authorized

### 3. Joining the Session

**Option A: Direct Jitsi URL**
- Use the `jitsi_url` from notification: `https://jitsi.example.com/a1b2c3d4e5f6789012345678901234`
- Open directly in browser or webview

**Option B: Validated Join (Recommended)**
1. Call join endpoint first to ensure access
2. Use returned `jitsi_url`
3. Provides additional security validation

**Option C: Embedded Integration**
```javascript
// Embed Jitsi in your application
const domain = 'jitsi.example.com';
const options = {
    roomName: 'a1b2c3d4e5f6789012345678901234', // The hex from call room
    width: '100%',
    height: 700,
    parentNode: document.querySelector('#jitsi-container')
};
const api = new JitsiMeetExternalAPI(domain, options);
```

## Security Features

- **Conversation Access Validation** - User must have access to conversation before creating call
- **UUID-based Room Access** - Non-sequential, non-guessable identifiers
- **Dual Access Control** - Users authorized by either user_id OR profile_id
- **Active Status Check** - Only active, non-deleted rooms can be joined
- **Unique Hex Generation** - 30-character collision-resistant room identifiers
- **Automatic User Addition** - Creator automatically added to authorized lists

## Error Responses

| Status | Error | Description |
|--------|-------|-------------|
| 400 | Invalid request data | JSON validation failed |
| 403 | Access denied to this conversation | User cannot access conversation |
| 403 | Call room is not active | Room ended or deleted |
| 403 | Access denied | User/profile not in authorized lists |
| 404 | Call room not found | Invalid UUID |
| 500 | Failed to generate unique hex | Hex generation failed |
| 500 | Failed to create call room | Database error |

## Environment Variables

```bash
API_SERVER_URL=https://api.example.com
API_PREFIX=/api/v1  
JITSI_URL=https://jitsi.example.com
```

## Implementation Notes

- **Hex Generation**: 30-character hex (15 bytes) with collision checking
- **WebSocket Integration**: Automatic notification via existing message routing system
- **Database Arrays**: PostgreSQL TEXT[] arrays for user_ids and profile_ids
- **UUID Generation**: Database-generated UUID for secure room identification
- **Access Control**: Dual-layer security (conversation access + user/profile lists)