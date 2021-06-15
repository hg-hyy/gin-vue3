package model

type TopicOrKeys struct {
	Topic string `json:"topic"`
	Keys  string `json:"keys"`
}

type DataOfBind struct {
	Data Bind `json:"data"`
}
type DataOfUnbind struct {
	Data Bind `json:"data"`
}
type Bind struct {
	Enterprise_id string `json:"enterprise_id"`
	Finger        string `json:"finger"`
	Name          string `json:"name"`
	User          string `json:"user"`
	AccountID     string `json:"accountID"`
	Type          int    `json:"type"`
}
type UnBind struct {
	Enterprise_id string `json:"enterprise_id"`
	AccountID     string `json:"accountID"`
	Type          int    `json:"type"`
}

type DataOfDev2Seat struct {
	Data Dev2Seat
}
type Dev2Seat struct {
	Enterprise_id string   `json:"enterprise_id"`
	Seat_finger   string   `json:"seat_finger"`
	Dev_fingers   []string `json:"dev_fingers"`
}

type AuthData struct {
	Enterprise_id string `json:"enterprise_id" form:"enterprise_id"`
	Enterprise    string `json:"enterprise" form:"enterprise"`
	Admin_num     int    `json:"admin_num" form:"admin_num"`
	Seat_num      int    `json:"seat_num" form:"seat_num"`
	Dev_num       int    `json:"dev_num" form:"dev_num"`
	Active_time   int64  `json:"active_time" form:"active_time"`
	End_time      int64  `json:"end_time" form:"end_time"`
}

type SignData struct {
	Data AccountData `json:"data"`
}

type AccountData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NetData struct {
	Data PostData `json:"data"`
}

type PostData struct {
	Device  string `json:"device"`
	Ip      string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	Dns1    string `json:"dns1"`
}

type GetData struct {
	Device string `json:"device"`
}

type PtHostData struct {
	Data PtIp `json:"data"`
}

type PtIp struct {
	Plat_ip string `json:"plat_ip"`
}
