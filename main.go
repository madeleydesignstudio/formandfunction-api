package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	app := fiber.New()

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-Requested-With",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Form & Function API",
			"version": "1.0.0",
			"endpoints": []string{
				"GET /beams",
				"GET /beams/:sectionDesignation",
				"POST /beams",
				"PUT /beams/:sectionDesignation",
				"DELETE /beams/:sectionDesignation",
			},
		})
	})

	app.Get("/beams", getBeams)
	app.Get("/beams/:sectionDesignation", getBeam)
	app.Post("/beams", createBeam)
	app.Put("/beams/:sectionDesignation", updateBeam)
	app.Delete("/beams/:sectionDesignation", deleteBeam)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Listen(":" + port)
}

func getBeams(c *fiber.Ctx) error {
	return c.JSON(beams)
}

func getBeam(c *fiber.Ctx) error {
	sectionDesignation := c.Params("sectionDesignation")
	for _, beam := range beams {
		if beam.SectionDesignation == sectionDesignation {
			return c.JSON(beam)
		}
	}
	return c.SendStatus(fiber.StatusNotFound)
}

func createBeam(c *fiber.Ctx) error {
	beam := new(SteelBeam)
	if err := c.BodyParser(beam); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	beams = append(beams, *beam)
	return c.JSON(beam)
}

func updateBeam(c *fiber.Ctx) error {
	sectionDesignation := c.Params("sectionDesignation")
	beamUpdate := new(SteelBeam)
	if err := c.BodyParser(beamUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	for i, beam := range beams {
		if beam.SectionDesignation == sectionDesignation {
			beams[i] = *beamUpdate
			return c.JSON(beamUpdate)
		}
	}
	return c.SendStatus(fiber.StatusNotFound)
}

func deleteBeam(c *fiber.Ctx) error {
	sectionDesignation := c.Params("sectionDesignation")
	for i, beam := range beams {
		if beam.SectionDesignation == sectionDesignation {
			beams = append(beams[:i], beams[i+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}
	return c.SendStatus(fiber.StatusNotFound)
}
