package device

type RegisterDeviceIdRequest struct {
	DeviceID   string `json:"device_id" binding:"required"`
	DeviceName string `json:"device_name" binding:"required"`
	Platform   string `json:"platform" binding:"required"`
}

type RegisterDeviceIdResponse struct {
	Success  bool   `json:"success"`
	DeviceID string `json:"device_id"`
	Message  string `json:"message"`
}
