# Travel App Backend API Documentation

Base URL: `/api/v1`

## Authentication

### POST `/auth/signup`
Register a new user.
**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123" // min 6 chars
}
```
**Response (201 Created)**:
```json
{
  "access_token": "jwt_token...",
  "refresh_token": "refresh_token..."
}
```

### POST `/auth/login`
Login with email and password.
**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```
**Response (200 OK)**:
```json
{
  "access_token": "jwt_token...",
  "refresh_token": "refresh_token..."
}
```

### POST `/auth/refresh`
Refresh access token using a refresh token.
**Request Body**:
```json
{
  "refresh_token": "refresh_token..."
}
```
**Response (200 OK)**:
```json
{
  "access_token": "new_jwt_token..."
}
```

### POST `/auth/google`
Login or Signup with Google OAuth token.
**Request Body**:
```json
{
  "token": "google_id_token..."
}
```
**Response (200 OK)**:
```json
{
  "access_token": "jwt_token...",
  "refresh_token": "refresh_token..."
}
```

## Trips

### POST `/trips`
Create a new trip.
**Request Body**:
```json
{
  "location": "Paris, France",
  "start_date": "2023-12-01T00:00:00Z",
  "end_date": "2023-12-10T00:00:00Z"
}
```
**Response (201 Created)**:
```json
{
  "id": "uuid...",
  "user_id": "uuid...",
  "location": "Paris, France",
  "start_date": "...",
  "end_date": "..."
}
```

### GET `/trips`
List authenticated user's trips.
**Response (200 OK)**:
```json
[
  {
    "id": "uuid...",
    "location": "Paris, France",
    "start_date": "...",
    "end_date": "..."
  }
]
```

### POST `/trips/:id/share`
Generate a share link for a trip.
**Response (200 OK)**:
```json
{
  "share_token": "random_string...",
  "url": "/preview/random_string..."
}
```

## Itineraries

### POST `/trips/:tripId/itineraries`
Create a day itinerary.
**Request Body**:
```json
{
  "slug": "day-1",
  "title": "Arrival & Check-in", // optional
  "date": "2023-12-01T00:00:00Z"
}
```
**Response (201 Created)**:
```json
{
  "id": "uuid...",
  "trip_id": "uuid...",
  "slug": "day-1",
  "date": "..."
}
```

## Activities

### POST `/itineraries/:itineraryId/activities`
Add an activity to an itinerary.
**Request Body**:
```json
{
  "name": "Visit Eiffel Tower",
  "description": "Tickets booked for 10 AM", // optional
  "location": "Champ de Mars, 5 Av. Anatole France", // optional
  "start_time": "2023-12-01T10:00:00Z", // optional
  "end_time": "2023-12-01T12:00:00Z", // optional
  "type": "sightseeing", // optional
  "status": "planned" // optional
}
```
**Response (201 Created)**:
```json
{
  "id": "uuid...",
  "itinerary_id": "uuid...",
  "name": "Visit Eiffel Tower",
  ...
}
```

## Locations

### GET `/locations/search`
Search for a place (Proxy to Google Places).
**Query Params**: `?query=Eiffel+Tower`
**Response (200 OK)**:
```json
[
  {
    "place_id": "...",
    "name": "Eiffel Tower",
    "formatted_address": "...",
    "geometry": {
      "location": { "lat": 48.8584, "lng": 2.2945 }
    }
  }
]
```

## Media

### POST `/media/upload`
Upload a media file (image/video).
**Content-Type**: `multipart/form-data`
**Form Fields**:
- `file`: (Binary file data)
- `activity_id`: (UUID of the activity)

**Response (200 OK)**:
```json
{
  "id": "uuid...",
  "url": "/uploads/1700000000_image.jpg",
  "type": "image",
  "activity_id": "uuid..."
}
```

## Users

### POST `/users/device-token`
Register a device for push notifications.
**Request Body**:
```json
{
  "token": "fcm_device_token..."
}
```
**Response (200 OK)**:
```json
{
  "message": "device registered"
}
```

## Public / Shared

### GET `/preview/:token`
Get trip details via share token (No Auth).
**Response (200 OK)**:
```json
{
  "id": "uuid...",
  "location": "Paris, France",
  "start_date": "...",
  "end_date": "...",
  "itineraries": [...]
}
```
