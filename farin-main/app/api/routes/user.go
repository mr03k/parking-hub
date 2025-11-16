package routes

import (
	controller "farin/app/api/controllers"
	"farin/app/api/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	userController             *controller.UserController
	am                         *middleware.AuthMiddleware
	contractController         *controller.ContractController
	driverController           *controller.DriverController
	contractorController       *controller.ContractorController
	adminMiddleware            *middleware.AdminMiddleware
	vehicleController          *controller.VehicleController
	deviceController           *controller.DeviceController
	roleController             *controller.RoleController
	calenderController         *controller.CalenderController
	ringController             *controller.RingController
	driverAssignmentController *controller.DriverAssignmentController
}

func NewUserRouter(userController *controller.UserController, am *middleware.AuthMiddleware,
	contractController *controller.ContractController, contractorController *controller.ContractorController,
	adminMiddleware *middleware.AdminMiddleware, driverController *controller.DriverController,
	vehicleController *controller.VehicleController, deviceController *controller.DeviceController, roleController *controller.RoleController,
	calenderController *controller.CalenderController, ringController *controller.RingController,
	driverAssignmentController *controller.DriverAssignmentController) *UserRouter {
	return &UserRouter{userController: userController, am: am, contractController: contractController,
		contractorController: contractorController, adminMiddleware: adminMiddleware,
		driverController: driverController, vehicleController: vehicleController, deviceController: deviceController,
		calenderController: calenderController, ringController: ringController, driverAssignmentController: driverAssignmentController,
		roleController: roleController}
}

func (rh *UserRouter) SetupRoutes(router *gin.Engine) {
	g := router.Group("/api/users")
	{
		g.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/", rh.userController.ListUsers)
		g.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/:id", rh.userController.GetUserDetail)
		g.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).POST("/", rh.userController.CreateUser)
		g.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).POST("/:id/picture", rh.userController.UploadUserPicture)
		g.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).PUT("/:id", rh.userController.UpdateUser)
		g.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).DELETE("/:id", rh.userController.DeleteUser)
	}

	c := router.Group("/api/contracts")
	{
		c.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/", rh.contractController.ListContracts)
		c.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/:id", rh.contractController.GetContractDetail)
		c.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).POST("/", rh.contractController.CreateContract)
		c.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).PUT("/:id", rh.contractController.UpdateContract)
		c.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).DELETE("/:id", rh.contractController.DeleteContract)
	}

	v := router.Group("/api/vehicles")
	{
		v.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/", rh.vehicleController.ListVehicles)
		v.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/:id", rh.vehicleController.GetVehicleDetail)
		v.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).POST("/", rh.vehicleController.CreateVehicle)
		v.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).PUT("/:id", rh.vehicleController.UpdateVehicle)
		v.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).DELETE("/:id", rh.vehicleController.DeleteVehicle)
	}

	co := router.Group("/api/contractors")
	{
		co.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/", rh.contractorController.ListContractors)
		co.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/:id", rh.contractorController.GetContractorDetail)
		co.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).POST("/", rh.contractorController.CreateContractor)
		co.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).PUT("/:id", rh.contractorController.UpdateContractor)
		co.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).DELETE("/:id", rh.contractorController.DeleteContractor)
	}

	d := router.Group("/api/drivers")
	{
		d.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/", rh.driverController.ListDrivers)
		d.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/:id", rh.driverController.GetDriverDetail)
		d.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).POST("/", rh.driverController.CreateDriver)
		d.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).PUT("/:id", rh.driverController.UpdateDriver)
		d.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).DELETE("/:id", rh.driverController.DeleteDriver)
	}

	dev := router.Group("/api/devices")
	{
		dev.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/", rh.deviceController.ListDevices)
		dev.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/:id", rh.deviceController.GetDeviceDetail)
		dev.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).POST("/", rh.deviceController.CreateDevice)
		dev.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).PUT("/:id", rh.deviceController.UpdateDevice)
		dev.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).DELETE("/:id", rh.deviceController.DeleteDevice)
	}

	role := router.Group("/api/roles")
	{
		role.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/", rh.roleController.ListRoles)
		role.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/:id", rh.roleController.GetRoleDetail)
		role.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).POST("/", rh.roleController.CreateRole)
		role.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).PUT("/:id", rh.roleController.UpdateRole)
		role.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).DELETE("/:id", rh.roleController.DeleteRole)
	}

	cal := router.Group("/api/calenders")
	{
		cal.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/", rh.calenderController.ListCalenders)
		cal.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/:id", rh.calenderController.GetCalenderDetail)
		cal.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).POST("/", rh.calenderController.CreateCalender)
		cal.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).PUT("/:id", rh.calenderController.UpdateCalender)
		cal.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).DELETE("/:id", rh.calenderController.DeleteCalender)
	}

	ring := router.Group("/api/rings")
	{
		ring.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/", rh.ringController.ListRings)
		ring.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/:id", rh.ringController.GetRingDetail)
		ring.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).POST("/", rh.ringController.CreateRing)
		ring.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).PUT("/:id", rh.ringController.UpdateRing)
		ring.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).DELETE("/:id", rh.ringController.DeleteRing)
	}

	da := router.Group("/api/driver-assignments")
	{
		da.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/", rh.driverAssignmentController.ListDriverAssignments)
		da.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).GET("/:id", rh.driverAssignmentController.GetDriverAssignmentDetail)
		da.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).POST("/", rh.driverAssignmentController.CreateDriverAssignment)
		da.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).PUT("/:id", rh.driverAssignmentController.UpdateDriverAssignment)
		da.Use(rh.am.Handle()).Use(rh.adminMiddleware.Handle()).DELETE("/:id", rh.driverAssignmentController.DeleteDriverAssignment)
	}

}
