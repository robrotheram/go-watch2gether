package meta

type RoomSettings struct {
	Controls bool `json:"controls"`
	AutoSkip bool `json:"auto_skip"`
}

func NewRoomSettings() RoomSettings {
	return RoomSettings{
		Controls: true,
		AutoSkip: true,
	}
}
