package handler

import (
	"net/http"

	"github.com/google/wire"
)

var HandlerProviderSet = wire.NewSet(
	NewOAPI,
	NewMuxHealthzHandler,
	NewIDPHandler,
	NewServiceList,
	NewAuthHandler,
	NewVehicleHandler,
	NewRingHandler,
	NewDriverHandler,
	NewDeviceHandler,
	NewConteractorHandler,
	NewWorkCalendarHandler,
	NewCityHandler,
	NewDistrictHandler,
)

// New ServiceList
func NewServiceList(healthzSvc *HealthzHandler, idpHDL *IDPHandler, authhdl *AuthHandler,
	vehicleHandler *VehicleHandler, rh *RingHandler, dh *DriverHandler, devH *DeviceHandler,
	ch *ConteractorHandler,
	calH *WorkCalendarHandler,
	city *CityHandler,
	disH *DistrictHandler,
) []Handler {
	return []Handler{
		healthzSvc,
		idpHDL,
		authhdl,
		vehicleHandler,
		rh,
		dh,
		devH,
		ch,
		calH,
		city,
		disH,
	}
}

// Service Interface
type Handler interface {
	RegisterMuxRouter(mux *http.ServeMux)
}

type OpenApiHandler interface {
	OpenApiSpec(o OAPI)
}

// NotImplemented
func NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}
