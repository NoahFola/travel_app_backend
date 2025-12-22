# Travel App Backend API Documentation

Base URL: `/api/v1`

## Authentication
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/auth/signup` | Register a new user | No |
| `POST` | `/auth/login` | Login with email/password | No |
| `POST` | `/auth/refresh` | Refresh access token | No |
| `POST` | `/auth/logout` | Logout user | Yes |
| `POST` | `/auth/google` | Login/Signup with Google | No |

## Trips
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/trips` | Create a new trip | Yes |
| `GET` | `/trips` | List my trips | Yes |
| `GET` | `/trips/:id` | Get trip details | Yes |
| `PUT` | `/trips/:id` | Update trip details | Yes |
| `DELETE` | `/trips/:id` | Delete a trip | Yes |
| `POST` | `/trips/:id/share` | Generate a share link | Yes |

## Itineraries
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/trips/:tripId/itineraries` | Create an itinerary for a trip | Yes |
| `GET` | `/trips/:tripId/itineraries` | List itineraries for a trip | Yes |
| `GET` | `/itineraries/:id` | Get itinerary details | Yes |
| `PUT` | `/itineraries/:id` | Update itinerary | Yes |
| `DELETE` | `/itineraries/:id` | Delete itinerary | Yes |

## Activities
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/itineraries/:itineraryId/activities` | Add activity to itinerary | Yes |
| `GET` | `/itineraries/:itineraryId/activities` | List activities in itinerary | Yes |
| `GET` | `/activities/:id` | Get activity details | Yes |
| `PUT` | `/activities/:id` | Update activity | Yes |
| `DELETE` | `/activities/:id` | Delete activity | Yes |

## Locations
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `GET` | `/locations/search` | Search places (Google Proxy) | Yes |
| | | Params: `?query=Paris` | |

## Media
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/media/upload` | Upload media file | Yes |
| | | Form Data: `file`, `activity_id` | |

## Users & Notifications
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/users/device-token` | Register device for notifications | Yes |

## Public / Sharing
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `GET` | `/preview/:token` | View shared trip details | No |

## Operational
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | Health Check (DB Ping) | No |

## Static Files
| Endpoint | Description |
| :--- | :--- |
| `/uploads/*` | Access uploaded media files |
