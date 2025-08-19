# 📱 Mobile Firebase Notification Guide

## 🔧 Setup Required

### 1. Register Your Device Token
**When to call:** App start, after login, when FCM token refreshes

**Send this:**
```http
POST /api/v1/auth/save-notification-token/
Authorization: Bearer your-access-token
```
```json
{
  "notification_token": "fGxqZ8vQR0mX...",
  "device_type": "android"
}
```

**You get this:**
```json
{
  "message": "Notification token saved successfully",
  "success": true
}
```

---

## 📩 Notification Types You'll Receive

### 💬 **Chat Message Notification**
**When:** Someone sends you a message and you're offline/disconnected

**Notification Display:**
- **Title:** "New message from John Doe"
- **Body:** "Hello, how are you?"

**Data Payload:**
```json
{
  "sender_name": "John Doe",
  "conversation_id": "123",
  "user_id": "456",
  "content": "Hello, how are you?",
  "type": "message",
  "app_name": "YourApp",
  "is_silent": "0",
  "created_at": "2024-01-15 14:30:25"
}
```

**What to do:** Open conversation with `conversation_id: 123`

---

### 🔔 **System Notification**
**When:** App sends important updates, announcements

**Notification Display:**
- **Title:** "System Notification"
- **Body:** "Your account has been updated"

**Data Payload:**
```json
{
  "sender_name": "system",
  "conversation_id": "0",
  "user_id": "0",
  "content": "Your account has been updated",
  "type": "notification",
  "app_name": "YourApp",
  "is_silent": "0",
  "created_at": "2024-01-15 14:30:25"
}
```

**What to do:** Show in-app notification or navigate to relevant screen

---

### 🔕 **Silent Notification**
**When:** Background data sync, status updates

**Notification Display:** *(No visual notification)*

**Data Payload:**
```json
{
  "sender_name": "system",
  "conversation_id": "456",
  "user_id": "789",
  "content": "User is typing...",
  "type": "typing",
  "app_name": "YourApp",
  "is_silent": "1",
  "created_at": "2024-01-15 14:30:25"
}
```

**What to do:** Update UI silently, no notification sound/banner

---

## 📋 Data Fields Explained

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `sender_name` | String | Who sent the message/notification | `"John Doe"`, `"system"` |
| `conversation_id` | String | Chat room to open | `"123"` (0 for system) |
| `user_id` | String | Sender's ID | `"456"` (0 for system) |
| `content` | String | Message text | `"Hello there!"` |
| `type` | String | Notification type | `"message"`, `"notification"`, `"typing"` |
| `is_silent` | String | Show notification? | `"0"` = show, `"1"` = silent |
| `created_at` | String | When sent | `"2024-01-15 14:30:25"` |
| `app_name` | String | Your app name | `"YourApp"` |

---

## 📱 Platform Differences

### 🤖 **Android**
- Channel ID: `messages`
- Click Action: `FLUTTER_NOTIFICATION_CLICK`
- High Priority notifications

### 🍎 **iOS**
- Badge count automatically managed
- Default notification sound
- Alert style notifications

---

## 🎯 Quick Implementation

### Flutter Example
```dart
FirebaseMessaging.onMessage.listen((RemoteMessage message) {
  final data = message.data;
  
  if (data['type'] == 'message') {
    // Navigate to chat: conversation_id
    Navigator.pushNamed(context, '/chat/${data['conversation_id']}');
  } else if (data['type'] == 'notification') {
    // Show system notification
    showDialog(...);
  }
});
```

### React Native Example
```javascript
messaging().onMessage(async remoteMessage => {
  const data = remoteMessage.data;
  
  if (data.type === 'message') {
    // Navigate to chat
    navigation.navigate('Chat', { conversationId: data.conversation_id });
  }
});
```

---

## ✅ Testing Checklist

- [ ] Token registration works on app start
- [ ] Notifications received when app is backgrounded
- [ ] Clicking notification opens correct chat
- [ ] Silent notifications don't show popup
- [ ] Token updates when FCM token refreshes
- [ ] Works on both Android and iOS



----

# Firebase Push Notification API Guide

## Overview
This API provides Firebase Cloud Messaging (FCM) integration for sending push notifications to mobile and web clients. It supports multiple device tokens per user and automatic cleanup of invalid tokens.

## 📱 Client Setup

### 1. Save Notification Token
Register a device token when the app starts or user logs in.

**Endpoint:** `POST {{host}}/{{prefix}}/auth/save-notification-token/`

**Headers:**
```
Authorization: Bearer {{access_token}}
Content-Type: application/json
```

**Request Body:**
```json
{
  "notification_token": "fGxqZ8vQR0mX...", // Required: FCM token from client
  "device_type": "android"                  // Optional: "android", "ios", "web"
}
```

**Response (Success - 200):**
```json
{
  "message": "Notification token saved successfully",
  "success": true,
  "data": {
    "notification_token": "fGxqZ8vQR0mX...",
    "device_type": "android"
  },
  "errorMsg": ""
}
```

**Response (Error - 400/500):**
```json
{
  "message": "Error message",
  "success": false,
  "data": null,
  "errorMsg": "Detailed error description"
}
```

## 🔧 Backend Integration

### 1. Automatic Notification Sending
Notifications are automatically sent in these scenarios:

- **WebSocket delivery fails**: When a user is online but message delivery fails
- **User offline**: When user is not connected via WebSocket
- **Conversation messages**: All participants except sender receive notifications

### 2. Manual Notification Sending

#### Send to Specific User
```go
import "texApi/internal/modules/firebasePush"

payload := firebasePush.NotificationPayload{
    SenderName:     "John Doe",
    ConversationID: 123,
    UserID:         456,
    Content:        "Hello, how are you?",
    Title:          &titleText,        // Optional: defaults to "New message from {SenderName}"
    CreatedAt:      time.Now().Format("2006-01-02 15:04:05"),
    Type:           "message",         // "message", "notification", "system"
    IsSilent:       0,                 // 0 = normal, 1 = silent
}

err := firebasePush.SendNotificationToUser(userID, payload)
```

#### Send to Conversation
```go
err := firebasePush.SendNotificationToConversation(
    conversationID,
    senderID,
    "John Doe",     // sender name
    "Hello everyone!" // message content
)
```

## 📊 Data Structures

### NotificationPayload
```go
type NotificationPayload struct {
    SenderName     string  `json:"sender_name"`     // Name of message sender
    ConversationID int     `json:"conversation_id"` // Chat/conversation ID
    UserID         int     `json:"user_id"`         // Sender's user ID
    Content        string  `json:"content"`         // Message content
    Title          *string `json:"title"`           // Optional: notification title
    CreatedAt      string  `json:"created_at"`      // Format: "2006-01-02 15:04:05"
    Type           string  `json:"type"`            // "message", "notification", "system"
    IsSilent       int     `json:"is_silent"`       // 0 = normal, 1 = silent
}
```

### Client Notification Data
When clients receive notifications, they get this data:

```json
{
  "data": {
    "sender_name": "John Doe",
    "conversation_id": "123",
    "user_id": "456",
    "content": "Hello, how are you?",
    "type": "message",
    "app_name": "YourAppName",
    "is_silent": "0",
    "created_at": "2024-01-15 14:30:25"
  },
  "notification": {
    "title": "New message from John Doe",
    "body": "Hello, how are you?"
  }
}
```

## 📱 Platform-Specific Features

### Android
- **Channel ID**: `messages`
- **Click Action**: `FLUTTER_NOTIFICATION_CLICK`
- **Priority**: `high`

### iOS (APNS)
- **Badge**: Auto-incremented
- **Sound**: `default`
- **Content Available**: Yes

## 🔄 Automatic Token Management

### Token Lifecycle
1. **Registration**: Tokens saved when client calls save-notification-token
2. **Validation**: Invalid tokens automatically detected during sending
3. **Cleanup**: Invalid tokens marked as `deleted = 1, active = 0`
4. **Multi-device**: Users can have multiple active tokens

### Database Schema
```sql
CREATE TABLE tbl_firebase_token (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES tbl_user(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    device_type TEXT,
    meta text,
    meta2 text,
    meta3 text,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    active INT NOT NULL DEFAULT 1,
    deleted INT NOT NULL DEFAULT 0
);
```

## 🚨 Error Handling

### Common Errors
- **Invalid Token**: Token automatically removed from database
- **User Not Found**: No tokens found for user (logged, not error)
- **Firebase Error**: Detailed logging for debugging

### Logging
All notification activities are logged:
```
Successfully sent message to user 123: projects/your-project/messages/0:1234567890
Failed to send to token: abc123, error: registration-token-not-registered
Removed invalid token: abc123
```

## 🔐 Security

- **Authentication**: Bearer token required for token registration
- **Token Validation**: Automatic cleanup of invalid tokens
- **Rate Limiting**: Handled by Firebase (default limits apply)
- **Data Privacy**: Only necessary message data sent in notifications

## 📝 Integration Checklist

- [ ] Firebase Admin SDK JSON file configured
- [ ] `FIREBASE_ADMINSDK_FILE` environment variable set
- [ ] Database table `tbl_firebase_token` exists
- [ ] Client apps integrated with FCM
- [ ] Token registration endpoint called on app start/login
- [ ] WebSocket failure handler updated to send notifications

## 🧪 Testing

### Test Token Registration
```bash
curl -X POST "{{host}}/{{prefix}}/auth/save-notification-token/" \
  -H "Authorization: Bearer your-access-token" \
  -H "Content-Type: application/json" \
  -d '{
    "notification_token": "test-token-123",
    "device_type": "android"
  }'
```

### Test Manual Notification
```go
// In your code
payload := firebasePush.NotificationPayload{
    SenderName:     "Test Sender",
    ConversationID: 1,
    UserID:         999,
    Content:        "Test notification",
    Type:           "message",
    CreatedAt:      time.Now().Format("2006-01-02 15:04:05"),
    IsSilent:       0,
}
err := firebasePush.SendNotificationToUser(targetUserID, payload)
```
