package routes

import (
	controller "git.abanppc.com/farin-project/vehicle-records/app/api/controllers"
	"github.com/gin-gonic/gin"
)

type VehicleRecordRouter struct {
	vehicleRecordController *controller.VehicleRecordController
}

func NewVehicleRecordRouter(vehicleRecordController *controller.VehicleRecordController) *VehicleRecordRouter {
	return &VehicleRecordRouter{vehicleRecordController: vehicleRecordController}
}

func (rh *VehicleRecordRouter) SetupRoutes(router *gin.Engine) {
	router.GET("api/v1/vehicle-records/:id", rh.vehicleRecordController.Detail)
	router.GET("api/v1/vehicle-records/retry", rh.vehicleRecordController.Retry)
}
