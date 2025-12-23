package api

import (
	"github.com/NoahFola/travel_app_backend/internal/handlers"
	"github.com/NoahFola/travel_app_backend/internal/middleware"
	"github.com/NoahFola/travel_app_backend/internal/repository"
	"github.com/NoahFola/travel_app_backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/time/rate"
)

// NewRouter initializes all dependencies and returns the configured Gin engine
func NewRouter(db *pgxpool.Pool) *gin.Engine {
	r := gin.New() // Use New() to skip default middlewares
	r.Use(gin.Recovery())

	// Rate Limiter: 20 requests per second with burst of 50
	limiter := middleware.NewIPRateLimiter(rate.Limit(20), 50)
	r.Use(middleware.RateLimit(limiter))
	r.Use(middleware.Logger())

	// --- 1. Initialize Repositories ---
	userRepo := repository.NewUserRepository(db)
	tripRepo := repository.NewTripRepository(db)
	itineraryRepo := repository.NewItineraryRepository(db)
	activityRepo := repository.NewActivityRepository(db)
	locationRepo := repository.NewLocationRepository(db)
	mediaRepo := repository.NewMediaRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	// --- 2. Initialize Services ---
	authService := &service.AuthService{Repo: userRepo}
	tripService := &service.TripService{Repo: tripRepo, ShareRepo: repository.NewShareRepository(db)}
	itineraryService := &service.ItineraryService{Repo: itineraryRepo, TripRepo: tripRepo}
	activityService := &service.ActivityService{Repo: activityRepo, ItineraryRepo: itineraryRepo}
	locationService := &service.LocationService{Repo: locationRepo}
	mediaService := &service.MediaService{Repo: mediaRepo}

	// --- 3. Initialize Handlers ---
	authHandler := &handlers.AuthHandler{Service: authService}
	tripHandler := &handlers.TripHandler{Service: tripService}
	itineraryHandler := &handlers.ItineraryHandler{Service: itineraryService}
	activityHandler := &handlers.ActivityHandler{Service: activityService}
	locationHandler := &handlers.LocationHandler{Service: locationService}
	mediaHandler := &handlers.MediaHandler{Service: mediaService}
	userHandler := &handlers.UserHandler{DeviceRepo: deviceRepo}
	healthHandler := &handlers.HealthHandler{DB: db}

	// Serve static files
	r.Static("/uploads", "./uploads")

	// Health Check
	r.GET("/health", healthHandler.HealthCheck)

	// --- 4. Register Routes ---

	// API Versioning Group (Good practice for future proofing)
	v1 := r.Group("/api/v1")
	{
		// Auth Routes
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", authHandler.Signup)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/google", authHandler.GoogleLogin)
		}

		// User Routes (Protected)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.POST("/device-token", userHandler.RegisterDevice)
		}

		// Trips Routes
		// Trips Routes
		trips := v1.Group("/trips")
		trips.Use(middleware.AuthMiddleware())
		{
			trips.POST("", tripHandler.CreateTrip)
			trips.GET("", tripHandler.ListMyTrips)

			trip := trips.Group("/:tripId")
			{
				trip.GET("", tripHandler.GetTrip)
				trip.PUT("", tripHandler.UpdateTrip)
				trip.DELETE("", tripHandler.DeleteTrip)

				// Sharing
				trip.POST("/share", tripHandler.ShareTrip)

				// Itineraries under a trip
				itineraries := trip.Group("/itineraries")
				{
					itineraries.POST("", itineraryHandler.CreateItinerary)
					itineraries.GET("", itineraryHandler.ListItineraries)
				}
			}
		}

		// Public Routes for Preview
		v1.GET("/preview/:token", tripHandler.GetSharedTrip)

		// Itineraries Routes (Direct access or strictly nested? User asked for /itineraries/{id}/activities)
		itineraries := v1.Group("/itineraries/:id")
		itineraries.Use(middleware.AuthMiddleware())
		{
			itineraries.GET("", itineraryHandler.GetItinerary)
			itineraries.PUT("", itineraryHandler.UpdateItinerary)
			itineraries.DELETE("", itineraryHandler.DeleteItinerary)

			// Nested Activities
			itineraries.POST("/activities", activityHandler.CreateActivity)
			itineraries.GET("/activities", activityHandler.ListActivities)
		}

		// Activities Routes
		activities := v1.Group("/activities")
		activities.Use(middleware.AuthMiddleware())
		{
			activities.GET("/:id", activityHandler.GetActivity)
			activities.PUT("/:id", activityHandler.UpdateActivity)
			activities.DELETE("/:id", activityHandler.DeleteActivity)
		}

		// Location Routes
		locations := v1.Group("/locations")
		locations.Use(middleware.AuthMiddleware())
		{
			locations.GET("/search", locationHandler.Search)
		}

		// Media Routes
		media := v1.Group("/media")
		media.Use(middleware.AuthMiddleware())
		{
			media.POST("/upload", mediaHandler.Upload)
		}
	}

	return r
}
