package domain

type ListItem struct {
	Id     int    `json:"id,omitempty"`
	UserId int    `json:"user_id,omitempty"`
	Item   string `json:"item"`
}

type ListDb interface {
	InsertListItem(bannerId string, userId int, item string) (int, error)
	GetListItems(bannerId string, userId int) ([]ListItem, error)
	UpdateListItem(bannerId string, userId int, itemId int, item string) error
	CheckVersion() (string, error)
}
