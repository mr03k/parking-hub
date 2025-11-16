package controller

import (
	"strconv"
	"time"

	"git.abanppc.com/farin-project/vehicle-records/app/api/response"
	_ "git.abanppc.com/farin-project/vehicle-records/domain/entity"
	"git.abanppc.com/farin-project/vehicle-records/domain/service"
	"github.com/gin-gonic/gin"
)

type VehicleRecordController struct {
	service *service.VehicleRecordService
}

// GetVehicleRecordDetail godoc
// @Summary      Get vehicleRecord details
// @Description  Retrieve vehicleRecord details by ID
// @Tags         vehicleRecords
// @Param        id   path      string  true  "VehicleRecord ID"
// @Success      200   {object}  response.Response[entity.VehicleRecord]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /v1/vehicle-records/{id} [get]
func (v VehicleRecordController) Detail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.NotFound(c)
		return
	}
	record, err := v.service.Detail(c, id)
	if err != nil {
		response.InternalError(c)
		return
	}
	response.Ok(c, record, "")
}

func NewVehicleRecordController(service *service.VehicleRecordService) *VehicleRecordController {
	return &VehicleRecordController{service: service}
}

// Retry godoc
// @Summary      Retry vehicleRecord resend
// @Description  Retry vehicleRecord resend
// @Tags         vehicleRecords
// @Param        from   	query      string  true  "From time in RFC3339 format"
// @Param        limit  	query      int     true  "Limit"
// @Param        x-api-key  header     string  true  "Authorization token"
// @Success      200   {object}  response.Response[swagger.EmptyObject]
func (v VehicleRecordController) Retry(c *gin.Context) {

	XApiKey := c.GetHeader("x-api-key")
	if XApiKey == "" {
		response.BadRequest(c, "Missing x-api-key header")
		return
	}
	if XApiKey != "ettRJD9ai1W17EiZ" {
		response.BadRequest(c, "Invalid x-api-key")
		return
	}

	from := c.Query("from")
	parsedTime, err := time.Parse(time.RFC3339, from)
	if err != nil {
		response.BadRequest(c, "Invalid from time format")
		return
	}
	fromTime := parsedTime.Unix()

	limit := c.Query("limit")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		response.BadRequest(c, "Invalid limit parameter")
		return
	}

	err = v.service.FindResend(c, fromTime, limitInt)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Ok(c, struct{}{}, "")
}
