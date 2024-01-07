package players

type Player struct {
	KitNumber   string `json:"kit_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Nationality string `json:"nationality"`
	Position    string `json:"position"`
}
