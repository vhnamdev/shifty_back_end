package dto

type InviteData struct {
	InviteCode   string `json:"invite_code"`
	PositionID   string `json:"position_id"`
	RestaurantID string `json:"restaurant_id"`
}
