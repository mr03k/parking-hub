package oapi

import (
	"log/slog"

	"github.com/google/wire"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi3"
)

var OapiProviderSet = wire.NewSet(GetSpec, NewOAPI, NewHealthzOpenApi,
	NewUserOpenApi, NewRoleOpenApi, NewAuthOpenApi, NewVehicleOpenApi,
	NewContractorOpenApi, NewDeviceOpenApi, NewDriverOpenApi, NewRingOpenApi, NewWorkCalendarOpenApi,
	NewCityOpenApi,
	NewDistrictOpenApi,
)

type OAPI interface {
	GenerateOperationContext(method string, path string) openapi.OperationContext
	GetYamlData() ([]byte, error)
	Register(method string, path string, handler OPHandlerFunc, opt ...RegisterOPT)
}

type RegisterOPT func(o openapi.OperationContext)

// with tags for operation
func WithTags(tags ...string) RegisterOPT {
	return func(o openapi.OperationContext) {
		o.SetTags(tags...)
	}
}

type OAPIMP struct {
	reflector *openapi3.Reflector
	logger    *slog.Logger
}

func NewOAPI(reflector *openapi3.Reflector, logger *slog.Logger) OAPI {
	return &OAPIMP{
		reflector: reflector,
		logger:    logger,
	}
}

func (s *OAPIMP) GenerateOperationContext(method, path string) openapi.OperationContext {
	op, err := s.reflector.NewOperationContext(method, path)
	if err != nil {
		s.logger.Error("GenerateOperationContext", "error", err)
		panic(err)
	}
	return op
}

// GetYamlData return yaml data
func (s *OAPIMP) GetYamlData() ([]byte, error) {
	return s.reflector.Spec.MarshalYAML()
}

func (s *OAPIMP) Register(method, path string, handler OPHandlerFunc, opt ...RegisterOPT) {
	op := s.GenerateOperationContext(method, path)

	for _, o := range opt {
		o(op)
	}

	handler(op)
	err := s.reflector.AddOperation(op)
	if err != nil {
		s.logger.Error("Register", "error", err)
	}
}

func GetSpec(
	o OAPI,
	h *HealthzOpenApi,
	cs *UserOpenApi,
	rs *RoleOpenApi,
	as *AuthOpenApi,
	vapi *VehicleOpenApi,
	capi *ContractorOpenApi,
	ringAPISet *RingOpenApi,
	do *DriverOpenApi,
	devAPI *DeviceOpenApi,
	calAPI *WorkCalendarOpenApi,
	city *CityOpenApi,
	disAPI *DistrictOpenApi,
) *openapi3.Spec {
	cs.OpenApiSpec(o)
	rs.OpenApiSpec(o)
	as.OpenApiSpec(o)
	vapi.OpenApiSpec(o)
	capi.OpenApiSpec(o)
	ringAPISet.OpenApiSpec(o)
	do.OpenApiSpec(o)
	devAPI.OpenApiSpec(o)
	calAPI.OpenApiSpec(o)
	city.OpenApiSpec(o)
	disAPI.OpenApiSpec(o)

	return o.(*OAPIMP).reflector.Spec
}

type OPHandlerFunc func(openapi.OperationContext)
