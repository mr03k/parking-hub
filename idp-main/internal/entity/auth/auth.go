package authentities

import idpentities "application/internal/entity/idp"

type JWTData struct {
	User      idpentities.User   `json:"user"`
	Roles     []idpentities.Role `json:"roles"`
	VehicleID string             `json:"vehicle_id"`
}
