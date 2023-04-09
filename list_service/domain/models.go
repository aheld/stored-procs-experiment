package domain

type ListItem struct {
	Id     int    `json:"id,omitempty"`
	UserId int    `json:"user_id,omitempty"`
	Item   string `json:"item"`
}

type ListRepo interface {
	InsertListItem(userId int, item string) (int, error)
	GetListItems(userId int) ([]ListItem, error)
	UpdateListItem(userId int, itemId int, item string) error
	CheckVersion() (string, error)
}
