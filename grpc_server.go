package main

import (
	"context"
	"log"
	"net"

	pb "formandfunction-api/proto"

	"google.golang.org/grpc"
)

// server is used to implement steelbeam.SteelBeamServiceServer
type server struct {
	pb.UnimplementedSteelBeamServiceServer
}

// Helper function to convert Go SteelBeam to protobuf SteelBeam
func steelBeamToProto(beam SteelBeam) *pb.SteelBeam {
	return &pb.SteelBeam{
		SectionDesignation:           beam.SectionDesignation,
		MassPerMetre:                 beam.MassPerMetre,
		DepthOfSection:               beam.DepthOfSection,
		WidthOfSection:               beam.WidthOfSection,
		ThicknessWeb:                 beam.ThicknessWeb,
		ThicknessFlange:              beam.ThicknessFlange,
		RootRadius:                   beam.RootRadius,
		DepthBetweenFillets:          beam.DepthBetweenFillets,
		RatiosForLocalBucklingWeb:    beam.RatiosForLocalBucklingWeb,
		RatiosForLocalBucklingFlange: beam.RatiosForLocalBucklingFlange,
		EndClearance:                 beam.EndClearance,
		Notch:                        beam.Notch,
		DimensionsForDetailingN:      beam.DimensionsForDetailingN,
		SurfaceAreaPerMetre:          beam.SurfaceAreaPerMetre,
		SurfaceAreaPerTonne:          beam.SurfaceAreaPerTonne,
		SecondMomentOfAreaAxisY:      beam.SecondMomentOfAreaAxisY,
		SecondMomentOfAreaAxisZ:      beam.SecondMomentOfAreaAxisZ,
		RadiusOfGyrationAxisY:        beam.RadiusOfGyrationAxisY,
		RadiusOfGyrationAxisZ:        beam.RadiusOfGyrationAxisZ,
		ElasticModulusAxisY:          beam.ElasticModulusAxisY,
		ElasticModulusAxisZ:          beam.ElasticModulusAxisZ,
		PlasticModulusAxisY:          beam.PlasticModulusAxisY,
		PlasticModulusAxisZ:          beam.PlasticModulusAxisZ,
		BucklingParameter:            beam.BucklingParameter,
		TorsionalIndex:               beam.TorsionalIndex,
		WarpingConstant:              beam.WarpingConstant,
		TorsionalConstant:            beam.TorsionalConstant,
		AreaOfSection:                beam.AreaOfSection,
	}
}

// Helper function to convert protobuf SteelBeam to Go SteelBeam
func protoToSteelBeam(pbBeam *pb.SteelBeam) SteelBeam {
	return SteelBeam{
		SectionDesignation:           pbBeam.SectionDesignation,
		MassPerMetre:                 pbBeam.MassPerMetre,
		DepthOfSection:               pbBeam.DepthOfSection,
		WidthOfSection:               pbBeam.WidthOfSection,
		ThicknessWeb:                 pbBeam.ThicknessWeb,
		ThicknessFlange:              pbBeam.ThicknessFlange,
		RootRadius:                   pbBeam.RootRadius,
		DepthBetweenFillets:          pbBeam.DepthBetweenFillets,
		RatiosForLocalBucklingWeb:    pbBeam.RatiosForLocalBucklingWeb,
		RatiosForLocalBucklingFlange: pbBeam.RatiosForLocalBucklingFlange,
		EndClearance:                 pbBeam.EndClearance,
		Notch:                        pbBeam.Notch,
		DimensionsForDetailingN:      pbBeam.DimensionsForDetailingN,
		SurfaceAreaPerMetre:          pbBeam.SurfaceAreaPerMetre,
		SurfaceAreaPerTonne:          pbBeam.SurfaceAreaPerTonne,
		SecondMomentOfAreaAxisY:      pbBeam.SecondMomentOfAreaAxisY,
		SecondMomentOfAreaAxisZ:      pbBeam.SecondMomentOfAreaAxisZ,
		RadiusOfGyrationAxisY:        pbBeam.RadiusOfGyrationAxisY,
		RadiusOfGyrationAxisZ:        pbBeam.RadiusOfGyrationAxisZ,
		ElasticModulusAxisY:          pbBeam.ElasticModulusAxisY,
		ElasticModulusAxisZ:          pbBeam.ElasticModulusAxisZ,
		PlasticModulusAxisY:          pbBeam.PlasticModulusAxisY,
		PlasticModulusAxisZ:          pbBeam.PlasticModulusAxisZ,
		BucklingParameter:            pbBeam.BucklingParameter,
		TorsionalIndex:               pbBeam.TorsionalIndex,
		WarpingConstant:              pbBeam.WarpingConstant,
		TorsionalConstant:            pbBeam.TorsionalConstant,
		AreaOfSection:                pbBeam.AreaOfSection,
	}
}

// GetBeams returns all steel beams
func (s *server) GetBeams(ctx context.Context, req *pb.GetBeamsRequest) (*pb.GetBeamsResponse, error) {
	log.Printf("gRPC GetBeams called")

	var protoBeams []*pb.SteelBeam
	for _, beam := range beams {
		protoBeams = append(protoBeams, steelBeamToProto(beam))
	}

	return &pb.GetBeamsResponse{
		Beams: protoBeams,
	}, nil
}

// GetBeam returns a specific steel beam by section designation
func (s *server) GetBeam(ctx context.Context, req *pb.GetBeamRequest) (*pb.GetBeamResponse, error) {
	log.Printf("gRPC GetBeam called with section: %s", req.SectionDesignation)

	for _, beam := range beams {
		if beam.SectionDesignation == req.SectionDesignation {
			return &pb.GetBeamResponse{
				Beam:  steelBeamToProto(beam),
				Found: true,
			}, nil
		}
	}

	return &pb.GetBeamResponse{
		Found: false,
	}, nil
}

// CreateBeam creates a new steel beam
func (s *server) CreateBeam(ctx context.Context, req *pb.CreateBeamRequest) (*pb.CreateBeamResponse, error) {
	log.Printf("gRPC CreateBeam called for section: %s", req.Beam.SectionDesignation)

	newBeam := protoToSteelBeam(req.Beam)
	beams = append(beams, newBeam)

	return &pb.CreateBeamResponse{
		Beam:    steelBeamToProto(newBeam),
		Success: true,
		Message: "Beam created successfully",
	}, nil
}

// GetStockStatus returns stock status for a product
func (s *server) GetStockStatus(ctx context.Context, req *pb.GetStockStatusRequest) (*pb.GetStockStatusResponse, error) {
	log.Printf("gRPC GetStockStatus called for product: %s, postcode: %s", req.ProductId, req.Postcode)

	status, err := GetStockStatus(req.ProductId, req.Postcode)
	if err != nil {
		return &pb.GetStockStatusResponse{
			ProductId: req.ProductId,
			Postcode:  req.Postcode,
			Success:   false,
			Message:   err.Error(),
		}, nil
	}

	return &pb.GetStockStatusResponse{
		ProductId: req.ProductId,
		Postcode:  req.Postcode,
		Status:    status,
		Success:   true,
		Message:   "Stock status retrieved successfully",
	}, nil
}

// StartGRPCServer starts the gRPC server on the specified port
func StartGRPCServer(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSteelBeamServiceServer(grpcServer, &server{})

	log.Printf("gRPC server starting on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}
