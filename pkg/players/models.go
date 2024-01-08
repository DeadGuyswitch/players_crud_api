package players

type Player struct {
	ID          int64  `json:"id"`
	KitNumber   string `json:"kit_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Nationality string `json:"nationality"`
	Position    string `json:"position"`
}
