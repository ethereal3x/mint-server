package service

import (
	"context"

	agentpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/agent"
	"github.com/ethereal3x/mint-server/internal/logic"
	"github.com/ethereal3x/mint-server/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// MappingServer MappingService gRPC 处理器
type MappingServer struct {
	agentpb.UnimplementedMappingServiceServer
	logic *logic.Mapping
}

// NewMappingServer 创建 Mapping gRPC 处理器
func NewMappingServer(mappingLogic *logic.Mapping) *MappingServer {
	return &MappingServer{logic: mappingLogic}
}

func (s *MappingServer) ListMappings(ctx context.Context, req *agentpb.ListMappingsRequest) (*agentpb.ListMappingsResponse, error) {
	var (
		mappings []*model.ModelMapping
		err      error
	)
	if req.Manufacturer != "" {
		mappings, err = s.logic.ListByManufacturer(ctx, req.Manufacturer)
	} else {
		mappings, err = s.logic.ListAll(ctx)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list mappings: %v", err)
	}
	pbMappings := make([]*agentpb.ModelMapping, 0, len(mappings))
	for _, m := range mappings {
		pbMappings = append(pbMappings, mappingFromModel(m))
	}
	return &agentpb.ListMappingsResponse{Mappings: pbMappings}, nil
}

func (s *MappingServer) GetMapping(ctx context.Context, req *agentpb.GetMappingRequest) (*agentpb.ModelMapping, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "get mapping: invalid id")
	}
	mapping, err := s.logic.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get mapping: %v", err)
	}
	if mapping == nil {
		return nil, status.Errorf(codes.NotFound, "get mapping: not found")
	}
	return mappingFromModel(mapping), nil
}

func (s *MappingServer) CreateMapping(ctx context.Context, req *agentpb.CreateMappingRequest) (*agentpb.ModelMapping, error) {
	if req.ModelType == "" || req.Manufacturer == "" {
		return nil, status.Errorf(codes.InvalidArgument, "create mapping: missing model_type or manufacturer")
	}
	mapping := &model.ModelMapping{
		ModelType:    req.ModelType,
		Manufacturer: req.Manufacturer,
		Description:  req.Description,
	}
	if err := s.logic.Create(ctx, mapping); err != nil {
		return nil, status.Errorf(codes.Internal, "create mapping: %v", err)
	}
	return mappingFromModel(mapping), nil
}

func (s *MappingServer) UpdateMapping(ctx context.Context, req *agentpb.UpdateMappingRequest) (*agentpb.ModelMapping, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "update mapping: invalid id")
	}
	if err := s.logic.Update(ctx, &model.ModelMapping{
		ID:          req.Id,
		ModelType:   req.ModelType,
		Manufacturer: req.Manufacturer,
		Description: req.Description,
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "update mapping: %v", err)
	}
	updated, err := s.logic.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get updated mapping: %v", err)
	}
	return mappingFromModel(updated), nil
}

func (s *MappingServer) DeleteMapping(ctx context.Context, req *agentpb.DeleteMappingRequest) (*emptypb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "delete mapping: invalid id")
	}
	if err := s.logic.Delete(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "delete mapping: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func mappingFromModel(m *model.ModelMapping) *agentpb.ModelMapping {
	if m == nil {
		return nil
	}
	return &agentpb.ModelMapping{
		Id:           m.ID,
		ModelType:    m.ModelType,
		Manufacturer: m.Manufacturer,
		Description:  m.Description,
	}
}
