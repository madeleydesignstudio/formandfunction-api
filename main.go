package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SteelBeam represents the properties of a steel beam
type SteelBeam struct {
	SectionDesignation           string  `json:"section_designation"`
	MassPerMetre                 float64 `json:"mass_per_metre"`
	DepthOfSection               float64 `json:"depth_of_section"`
	WidthOfSection               float64 `json:"width_of_section"`
	ThicknessWeb                 float64 `json:"thickness_web"`
	ThicknessFlange              float64 `json:"thickness_flange"`
	RootRadius                   float64 `json:"root_radius"`
	DepthBetweenFillets          float64 `json:"depth_between_fillets"`
	RatiosForLocalBucklingWeb    float64 `json:"ratios_for_local_buckling_web"`
	RatiosForLocalBucklingFlange float64 `json:"ratios_for_local_buckling_flange"`
	EndClearance                 float64 `json:"end_clearance"`
	Notch                        float64 `json:"notch"`
	DimensionsForDetailingN      float64 `json:"dimensions_for_detailing_n"`
	SurfaceAreaPerMetre          float64 `json:"surface_area_per_metre"`
	SurfaceAreaPerTonne          float64 `json:"surface_area_per_tonne"`
	SecondMomentOfAreaAxisY      float64 `json:"second_moment_of_area_axis_y"`
	SecondMomentOfAreaAxisZ      float64 `json:"second_moment_of_area_axis_z"`
	RadiusOfGyrationAxisY        float64 `json:"radius_of_gyration_axis_y"`
	RadiusOfGyrationAxisZ        float64 `json:"radius_of_gyration_axis_z"`
	ElasticModulusAxisY          float64 `json:"elastic_modulus_axis_y"`
	ElasticModulusAxisZ          float64 `json:"elastic_modulus_axis_z"`
	PlasticModulusAxisY          float64 `json:"plastic_modulus_axis_y"`
	PlasticModulusAxisZ          float64 `json:"plastic_modulus_axis_z"`
	BucklingParameter            float64 `json:"buckling_parameter"`
	TorsionalIndex               float64 `json:"torsional_index"`
	WarpingConstant              float64 `json:"warping_constant"`
	TorsionalConstant            float64 `json:"torsional_constant"`
	AreaOfSection                float64 `json:"area_of_section"`
}

var beams = []SteelBeam{
	{
		SectionDesignation:           "UB406x178x74",
		MassPerMetre:                 74.6,
		DepthOfSection:               412.8,
		WidthOfSection:               179.5,
		ThicknessWeb:                 9.3,
		ThicknessFlange:              16.0,
		RootRadius:                   10.2,
		DepthBetweenFillets:          360.8,
		RatiosForLocalBucklingWeb:    38.8,
		RatiosForLocalBucklingFlange: 5.61,
		EndClearance:                 369.0,
		Notch:                        360.8,
		DimensionsForDetailingN:      45.0,
		SurfaceAreaPerMetre:          1.17,
		SurfaceAreaPerTonne:          15.7,
		SecondMomentOfAreaAxisY:      27400,
		SecondMomentOfAreaAxisZ:      1600,
		RadiusOfGyrationAxisY:        17.1,
		RadiusOfGyrationAxisZ:        4.22,
		ElasticModulusAxisY:          1330,
		ElasticModulusAxisZ:          178,
		PlasticModulusAxisY:          1500,
		PlasticModulusAxisZ:          275,
		BucklingParameter:            0.338,
		TorsionalIndex:               29.6,
		WarpingConstant:              0.581,
		TorsionalConstant:            53.8,
		AreaOfSection:                95.0,
	},
	{
		SectionDesignation:           "UB406x178x67",
		MassPerMetre:                 67.1,
		DepthOfSection:               406.4,
		WidthOfSection:               177.9,
		ThicknessWeb:                 8.6,
		ThicknessFlange:              12.8,
		RootRadius:                   10.2,
		DepthBetweenFillets:          360.8,
		RatiosForLocalBucklingWeb:    42.0,
		RatiosForLocalBucklingFlange: 6.95,
		EndClearance:                 362.6,
		Notch:                        360.8,
		DimensionsForDetailingN:      45.0,
		SurfaceAreaPerMetre:          1.15,
		SurfaceAreaPerTonne:          17.1,
		SecondMomentOfAreaAxisY:      23500,
		SecondMomentOfAreaAxisZ:      1350,
		RadiusOfGyrationAxisY:        16.6,
		RadiusOfGyrationAxisZ:        4.09,
		ElasticModulusAxisY:          1160,
		ElasticModulusAxisZ:          152,
		PlasticModulusAxisY:          1300,
		PlasticModulusAxisZ:          234,
		BucklingParameter:            0.364,
		TorsionalIndex:               25.4,
		WarpingConstant:              0.424,
		TorsionalConstant:            36.4,
		AreaOfSection:                85.5,
	},
}

func main() {
	// Configuration
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "9090"
	}

	// Create Fiber app for HTTP REST API (frontend consumption)
	app := fiber.New(fiber.Config{
		ReadBufferSize:  32768,            // Increase to 32KB to handle large headers
		WriteBufferSize: 32768,            // Increase write buffer size
		ReadTimeout:     time.Second * 30, // 30 second read timeout
		WriteTimeout:    time.Second * 30, // 30 second write timeout
		BodyLimit:       10 * 1024 * 1024, // 10MB body limit
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Printf("HTTP Request error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Add logging middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))

	// Add CORS middleware for frontend
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: false,
	}))

	// HTTP REST API Routes for Frontend
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":     "Form & Function API",
			"version":     "2.0.0",
			"description": "HTTP REST API for frontend + gRPC backend communication",
			"endpoints": []string{
				"GET /beams",
				"GET /beams/:sectionDesignation",
				"POST /beams",
				"PUT /beams/:sectionDesignation",
				"DELETE /beams/:sectionDesignation",
				"GET /stock?productId=<product_id>&postcode=<postcode>",
			},
			"grpc_port": grpcPort,
			"http_port": httpPort,
		})
	})

	app.Get("/beams", getBeams)
	app.Get("/beams/:sectionDesignation", getBeam)
	app.Post("/beams", createBeam)
	app.Put("/beams/:sectionDesignation", updateBeam)
	app.Delete("/beams/:sectionDesignation", deleteBeam)
	app.Get("/stock", getStockStatusHandler)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":       "healthy",
			"service":      "Form & Function API",
			"http_port":    httpPort,
			"grpc_port":    grpcPort,
			"endpoints":    "HTTP REST for frontend, gRPC for backend services",
			"beam_count":   len(beams),
			"architecture": "Hybrid HTTP/gRPC",
		})
	})

	// Set up graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start gRPC server in a goroutine (for backend services like Python calc engine)
	go func() {
		log.Printf("Starting gRPC server on port %s (for backend services)", grpcPort)
		StartGRPCServer(grpcPort)
	}()

	// Start HTTP REST API server in a goroutine (for frontend)
	go func() {
		log.Printf("Starting HTTP REST API server on port %s (for frontend)", httpPort)
		if err := app.Listen(":" + httpPort); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	log.Printf("🚀 Form & Function API Services Started:")
	log.Printf("   📱 Frontend HTTP REST API: http://localhost:%s", httpPort)
	log.Printf("   🔧 Backend gRPC Service:   localhost:%s", grpcPort)
	log.Printf("   💡 Architecture: HTTP for frontend, gRPC for backend services")

	// Wait for shutdown signal
	sig := <-sigCh
	log.Printf("Received signal %v, shutting down gracefully...", sig)

	// Graceful shutdown
	log.Println("Shutting down HTTP server...")
	if err := app.Shutdown(); err != nil {
		log.Printf("Error during HTTP server shutdown: %v", err)
	}

	log.Println("Services shut down successfully")
}

// HTTP REST API Handlers for Frontend

func getStockStatusHandler(c *fiber.Ctx) error {
	productID := c.Query("productId")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "productId query parameter is required"})
	}
	postcode := c.Query("postcode")
	if postcode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "postcode query parameter is required"})
	}

	status, err := GetStockStatus(productID, postcode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"productId": productID,
		"postcode":  postcode,
		"status":    status,
		"source":    "http_rest_api",
	})
}

func getBeams(c *fiber.Ctx) error {
	log.Printf("HTTP REST API: GET /beams called")
	return c.JSON(fiber.Map{
		"beams":  beams,
		"count":  len(beams),
		"source": "http_rest_api",
	})
}

func getBeam(c *fiber.Ctx) error {
	sectionDesignation := c.Params("sectionDesignation")
	log.Printf("HTTP REST API: GET /beams/%s called", sectionDesignation)

	for _, beam := range beams {
		if beam.SectionDesignation == sectionDesignation {
			return c.JSON(fiber.Map{
				"beam":   beam,
				"source": "http_rest_api",
			})
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error":  "Beam not found",
		"source": "http_rest_api",
	})
}

func createBeam(c *fiber.Ctx) error {
	log.Printf("HTTP REST API: POST /beams called")

	beam := new(SteelBeam)
	if err := c.BodyParser(beam); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  err.Error(),
			"source": "http_rest_api",
		})
	}

	beams = append(beams, *beam)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"beam":    beam,
		"message": "Beam created successfully",
		"source":  "http_rest_api",
	})
}

func updateBeam(c *fiber.Ctx) error {
	sectionDesignation := c.Params("sectionDesignation")
	log.Printf("HTTP REST API: PUT /beams/%s called", sectionDesignation)

	beamUpdate := new(SteelBeam)
	if err := c.BodyParser(beamUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  err.Error(),
			"source": "http_rest_api",
		})
	}

	for i, beam := range beams {
		if beam.SectionDesignation == sectionDesignation {
			beams[i] = *beamUpdate
			return c.JSON(fiber.Map{
				"beam":    beamUpdate,
				"message": "Beam updated successfully",
				"source":  "http_rest_api",
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error":  "Beam not found",
		"source": "http_rest_api",
	})
}

func deleteBeam(c *fiber.Ctx) error {
	sectionDesignation := c.Params("sectionDesignation")
	log.Printf("HTTP REST API: DELETE /beams/%s called", sectionDesignation)

	for i, beam := range beams {
		if beam.SectionDesignation == sectionDesignation {
			beams = append(beams[:i], beams[i+1:]...)
			return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"message": "Beam deleted successfully",
				"source":  "http_rest_api",
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error":  "Beam not found",
		"source": "http_rest_api",
	})
}
