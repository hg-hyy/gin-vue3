package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"pt/model"
	"pt/utils"
	"strconv"

	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Api account type
type Api struct {
	AdminAccounts []*Account            `json:"adminaccounts"`
	SeatsAccounts []*Account            `json:"seatsaccounts"`
	DevsAccounts  []*Account            `json:"devsaccounts"`
	AdminIDs      []string              `json:"adminids"`
	SeatsIDs      []string              `json:"seatsids"`
	DevsIDs       []string              `json:"devsids"`
	AdminInfos    []*AccountInfo        `json:"admininfos"`
	SeatsInfos    []*AccountInfo        `json:"seatsinfos"`
	DevsInfos     []*AccountInfo        `json:"devsinfos"`
	DevsOfSeat    []map[string][]string `json:"devsofseat"`
	Devs          []map[string]Dev      `json:"devs"`
	Netconf       []map[string]string   `json:"netconf"`
	Pthost        map[string]string     `json:"pthost"`
}

type Netinfo struct {
	Device  string `json:"device"`
	Ip      string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	Dns1    string `json:"dns1"`
}

type Dev struct {
	Finger string `json:"finger"`
	Name   string `json:"name"`
}

type Account struct {
	Key   string `json:"key"`
	Value Value  `json:"value"`
}

type Value struct {
	Finger      string   `json:"finger"`
	Name        string   `json:"name"`
	Users       []string `json:"users"`
	Active_time string   `json:"active_time"`
	End_time    string   `json:"end_time"`
}

// AccountInfo account need this
type AccountInfo struct {
	Nick   string   `json:"nick"`
	User   []string `json:"user"`
	Finger string   `json:"finger"`
	Token  string   `json:"token"`
}

// Routers init routers
func (a *Api) Routers(r *gin.Engine) {

	gin := r.Group("/gin")
	gin.GET("/sqpt", a.index)
	gin.GET("/page", a.page)
	r.GET("/n_node/v1.0/kv/seq", a.getAccounts)
	r.Any("/echo", a.wsocket)
	r.GET("/n_account/v1.0/device", a.accounts)
	r.GET("/d_sysop/v1.0/netConf", a.getnetconf)
	r.GET("/n_node/v1.0/sensor/seq", a.getseq)
	r.GET("/d_sysop/v1.0/plat_host", a.pthost)
	r.GET("/api/v1/settings/public", a.jumpserver)
	r.POST("/d_sysop/v1.0/plat_host", a.pthost)
	r.POST("/d_auth/v1.0/bindAccount", a.binding)
	r.POST("/d_auth/v1.0/devsToSeat", a.assignment)
	r.POST("/d_auth/v1.0/authInfo", a.authInfo)
	r.POST("/d_sysop/v1.0/login", a.signin)
	r.POST("/d_sysop/v1.0/netConf", a.setnetconf)
	r.POST("/postdata", a.postdata)
	r.DELETE("/d_auth/v1.0/unbindAccount", a.unbinding)
}

// var upgrader = websocket.Upgrader{} // use default options
var wg sync.WaitGroup

func (a *Api) getseq(c *gin.Context) {

	filter, err := os.OpenFile("./json/seq.json", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	var seq []map[string]interface{}
	decoder := json.NewDecoder(filter)
	err = decoder.Decode(&seq)
	if err != nil {
		log.Println(err)
	}
	// wg.Add(1)
	// go utils.Seq()
	c.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"msg":  "ok...",
		"data": seq,
	})

}

func (a *Api) wsocket(c *gin.Context) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()
	wg.Add(1)
	go a.Recvmsg(conn)
	wg.Add(1)
	go a.Sendmsg(conn, 1)
	wg.Wait()
}

func (a *Api) Recvmsg(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		log.Printf("recv: %s", message)
		if err != nil {
			log.Println("read:", err)
			wg.Done()
			break
		}
	}
}
func (a *Api) Sendmsg(conn *websocket.Conn, mt int) {
	for {
		var msg = []byte(time.Now().Format("2006-01-02 15：04：05") + ":CPU温度高报警")
		err := conn.WriteMessage(mt, msg)
		time.Sleep(time.Second * 10)
		if err != nil {
			log.Println("write:", err)
			wg.Done()
			break

		}
	}
}

// hello say hello to the world
func (a *Api) postdata(c *gin.Context) {
	var netData model.NetData
	c.ShouldBind(&netData)
	fmt.Println(netData.Data)
	c.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"msg":  "ok...",
		"data": netData,
	})

}

// hello say hello to the world
func (a *Api) pthost(c *gin.Context) {

	if a.Pthost == nil {
		a.Pthost = map[string]string{
			"plat_ip": "10.99.0.56",
		}
	}

	if c.Request.Method == "GET" {
		c.JSON(http.StatusOK, gin.H{
			"code": 1000,
			"msg":  "ok...",
			"data": a.Pthost,
		})
	} else {

		var netData model.PtHostData
		c.ShouldBind(&netData)
		a.Pthost["plat_ip"] = netData.Data.Plat_ip
		c.JSON(http.StatusOK, gin.H{
			"code": 1000,
			"msg":  "ok...",
			"data": "",
		})
	}

}

// hello say hello to the world
func (a *Api) jumpserver(c *gin.Context) {
	fmt.Println("helo,jmp")

	c.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"msg":  "ok...",
		"data": "a.Pthost",
	})

}

func (a *Api) getnetconf(c *gin.Context) {
	var gd model.GetData
	c.ShouldBind(&gd)
	if a.Netconf == nil {

		for i := 0; i < 2; i++ {
			data := map[string]string{
				"device":  "eth" + strconv.Itoa(i),
				"ip":      "192.168.20.120",
				"netmask": "255.255.255.0",
				"gateway": "192.168.20.254",
				"dns1":    "114.114.114.114",
			}
			a.Netconf = append(a.Netconf, data)
		}
	}

	var flag = false
	var index int
	for k, v := range a.Netconf {
		if c.Query("device") == v["device"] {
			flag = true
			index = k
		}
	}
	if flag {
		c.JSON(http.StatusOK, gin.H{
			"code": 1000,
			"msg":  "ok",
			"data": a.Netconf[index],
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1000,
			"msg":  "ok",
			"data": a.Netconf,
		})
	}

}
func (a *Api) setnetconf(c *gin.Context) {
	var netData model.NetData
	c.ShouldBind(&netData)
	for k, v := range a.Netconf {
		if netData.Data.Device == v["device"] {
			a.Netconf[k]["device"] = netData.Data.Device
			a.Netconf[k]["ip"] = netData.Data.Ip
			a.Netconf[k]["gateway"] = netData.Data.Gateway
			a.Netconf[k]["netmask"] = netData.Data.Netmask
			a.Netconf[k]["dns1"] = netData.Data.Dns1
			c.JSON(http.StatusOK, gin.H{
				"code": 1000,
				"msg":  "设置网卡成功",
				"data": "",
			})
		}
	}
}

func (a *Api) signin(c *gin.Context) {
	var signdata model.SignData
	c.ShouldBind(&signdata)
	fmt.Println(c.Request.Header.Get("authtoken"))
	if signdata.Data.Username == "admin" && signdata.Data.Password == "admin" {
		c.JSON(http.StatusOK, gin.H{
			"code": 1000,
			"msg":  "login seccuss",
			"data": map[string]interface{}{
				"isAuthenticated": 1,
				"username":        signdata.Data.Username,
				"authtoken":       "admin666",
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1000,
			"msg":  "login faild",
			"data": map[string]interface{}{
				"isAuthenticated": 0,
				"username":        "",
				"authtoken":       "",
			},
		})
	}

}
func (a *Api) index(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", gin.H{
		"code": 1000,
		"msg":  "this is hyy say hello to the world",
		"data": "",
	})

}
func (a *Api) page(c *gin.Context) {

	c.HTML(http.StatusOK, "page.html", gin.H{
		"code": 1000,
		"msg":  "this is hyy say hello to the world",
		"data": ""})
}

// getAccounts return configs
func (a *Api) getAccounts(c *gin.Context) {
	c.ShouldBindJSON(&model.TopicOrKeys{})
	topic := c.Query("topic")
	log.Printf("topic:%s", topic)
	if topic != "" {
		if topic == "enterprise1.seats" {
			// 这里需要注意否则解析json null
			if len(a.SeatsAccounts) == 0 {
				a.SeatsAccounts = []*Account{}
			}
			c.JSON(http.StatusOK, gin.H{
				"code": 1000,
				"msg":  "返回座席成功",
				"data": a.SeatsAccounts,
			})
		} else if topic == "enterprise1.devs" {
			c.JSON(http.StatusOK,
				gin.H{
					"code": 1000,
					"msg":  "返回终端成功",
					"data": a.DevsAccounts,
				})
		} else if topic == "enterprise1.admin" {
			c.JSON(http.StatusOK,
				gin.H{
					"code": 1000,
					"msg":  "返回管理员成功",
					"data": a.AdminAccounts,
				})
		} else {
			if len(a.DevsOfSeat) != 0 {
				for _, m := range a.DevsOfSeat {
					for k, v := range m {
						if k+".devs" == topic {
							res := a.devmap(v)
							if len(res) == 0 {
								res = []Dev{}
							}
							c.JSON(http.StatusOK, gin.H{
								"code": 1000,
								"msg":  "返回座席分配的终端成功",
								"data": res,
							})
							return
						}
					}
				}
				c.JSON(http.StatusOK, gin.H{
					"code": 1000,
					"msg":  "没有找到，请查看座席指纹",
					"data": "",
				})

			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": 1000,
					"msg":  "没有找到,请在座席下分配终端",
					"data": "",
				})
			}
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1000,
			"msg":  "请检查Topic",
			"data": "",
		})
	}
}

// binding bind configs
func (a *Api) binding(c *gin.Context) {
	var data model.DataOfBind
	c.ShouldBind(&data)
	for _, v := range a.SeatsAccounts {
		if v.Key == data.Data.AccountID {

			v.Value.Finger = data.Data.Finger
			v.Value.Name = data.Data.Name
			v.Value.Users = append(v.Value.Users, data.Data.User)
		}
	}
	for _, v := range a.AdminAccounts {
		if v.Key == data.Data.AccountID {

			v.Value.Finger = data.Data.Finger
			v.Value.Name = data.Data.Name
			v.Value.Users = append(v.Value.Users, data.Data.User)
		}
	}
	for _, v := range a.DevsAccounts {
		if v.Key == data.Data.AccountID {

			v.Value.Finger = data.Data.Finger
			v.Value.Name = data.Data.Name
			v.Value.Users = append(v.Value.Users, data.Data.User)

			dev := map[string]Dev{
				data.Data.Finger: {
					Finger: data.Data.Finger,
					Name:   data.Data.Name,
				},
			}
			a.Devs = append(a.Devs, dev)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"msg":  "绑定成功",
		"data": data})
}

// unbinding unbind configs
func (a *Api) unbinding(c *gin.Context) {
	var data model.DataOfUnbind
	c.ShouldBind(&data)
	for _, v := range a.SeatsAccounts {
		if v.Key == data.Data.AccountID {

			v.Value.Finger = ""
			v.Value.Name = ""
			v.Value.Users = []string{}
			for k := range a.DevsOfSeat {
				a.DevsOfSeat[k][v.Value.Finger] = []string{}
			}
		}
	}
	for _, v := range a.AdminAccounts {
		if v.Key == data.Data.AccountID {

			v.Value.Finger = ""
			v.Value.Name = ""
			v.Value.Users = []string{}
		}
	}
	var dev_finger = ""
	for _, v := range a.DevsAccounts {
		if v.Key == data.Data.AccountID {
			for index, devmap := range a.DevsOfSeat {
				for key, devs := range devmap {
					for _, dev := range devs {
						if dev == v.Value.Finger {
							a.DevsOfSeat[index][key] = append(a.DevsOfSeat[index][key][:index], a.DevsOfSeat[index][key][index+1:]...)
						}
					}
				}
			}
			v.Value.Finger = dev_finger
			v.Value.Finger = ""
			v.Value.Name = ""
			v.Value.Users = []string{}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"msg":  "解除绑定成功",
		"data": &model.UnBind{
			Enterprise_id: data.Data.Enterprise_id,
			AccountID:     data.Data.AccountID,
			Type:          data.Data.Type,
		}})
}

// assignment assign devs to seat
func (a *Api) assignment(c *gin.Context) {
	var data model.DataOfDev2Seat
	c.ShouldBind(&data)
	fmt.Println(data)
	if len(a.DevsOfSeat) != 0 {
		var flag bool = false
		index := 0
		for i, mp := range a.DevsOfSeat {
			for k := range mp {
				if k == data.Data.Seat_finger {
					flag, index = true, i
					fmt.Println("存在，追加。。。")
				} else {
					devsofseat := map[string][]string{
						data.Data.Seat_finger: data.Data.Dev_fingers,
					}
					a.DevsOfSeat = append(a.DevsOfSeat, devsofseat)
				}

			}

		}
		if flag {
			a.DevsOfSeat[index][data.Data.Seat_finger] = append(a.DevsOfSeat[index][data.Data.Seat_finger], data.Data.Dev_fingers...)
		} else {
			devsofseat := map[string][]string{
				data.Data.Seat_finger: data.Data.Dev_fingers,
			}
			a.DevsOfSeat = append(a.DevsOfSeat, devsofseat)
		}
	} else {
		devsofseat := map[string][]string{
			data.Data.Seat_finger: data.Data.Dev_fingers,
		}
		a.DevsOfSeat = append(a.DevsOfSeat, devsofseat)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"msg":  "分配成功",
		"data": a.DevsOfSeat,
	})
}

func (a *Api) authInfo(c *gin.Context) {
	var data model.AuthData
	c.ShouldBind(&data)
	a.add(data)

	// c.HTML(http.StatusOK, "index.html", gin.H{
	// 	"code": 1000,
	// 	"msg":  "账户授权成功",
	// 	"data": data,
	// })
	c.Redirect(http.StatusMovedPermanently, "http://127.0.0.1:8080/gin/sqpt")
}

// add  authdata to accounts
func (a *Api) add(addaccount model.AuthData) {
	active_time := time.Unix(addaccount.Active_time, 0).Format("2006-01-02 15:04:05")
	end_time := time.Unix(addaccount.End_time, 0).Format("2006-01-02 15:04:05")
	for i := 0; i < addaccount.Admin_num; i++ {
		res := &AccountInfo{
			Nick:   "admin" + strconv.Itoa(1+len(a.AdminAccounts)),
			User:   []string{"hyy", "fhh"},
			Finger: utils.Finger(4, "admin"),
			Token:  utils.Token(8),
		}

		admin_id := "admin.accountID" + strconv.Itoa(len(a.AdminIDs)+1)

		admin_account := &Account{
			Key: addaccount.Enterprise_id + "." + admin_id,
			Value: Value{
				Finger:      "",
				Name:        "",
				Users:       []string{},
				Active_time: active_time,
				End_time:    end_time,
			},
		}
		a.AdminIDs = append(a.AdminIDs, admin_id)
		a.AdminAccounts = append(a.AdminAccounts, admin_account)
		a.AdminInfos = append(a.AdminInfos, res)

	}
	for i := 0; i < addaccount.Seat_num; i++ {
		res := &AccountInfo{
			Nick:   "seat" + strconv.Itoa(1+len(a.SeatsAccounts)),
			User:   []string{"fhh", "mcc"},
			Finger: utils.Finger(4, "seat"),
			Token:  utils.Token(8),
		}
		seat_id := "seat.accountID" + strconv.Itoa(len(a.SeatsIDs)+1)

		seat_account := &Account{
			Key: addaccount.Enterprise_id + "." + seat_id,
			Value: Value{
				Finger:      "",
				Name:        "",
				Users:       []string{},
				Active_time: active_time,
				End_time:    end_time,
			},
		}
		a.SeatsIDs = append(a.SeatsIDs, seat_id)
		a.SeatsAccounts = append(a.SeatsAccounts, seat_account)
		a.SeatsInfos = append(a.SeatsInfos, res)

	}
	for i := 0; i < addaccount.Dev_num; i++ {
		res := &AccountInfo{
			Nick:   "dev" + strconv.Itoa(1+len(a.DevsAccounts)),
			User:   []string{"mcc", "hyy"},
			Finger: utils.Finger(4, "devs"),
			Token:  utils.Token(8),
		}
		dev_id := "dev.accountID" + strconv.Itoa(len(a.DevsIDs)+1)

		dev_account := &Account{
			Key: addaccount.Enterprise_id + "." + dev_id,
			Value: Value{
				Finger:      "",
				Name:        "",
				Users:       []string{},
				Active_time: active_time,
				End_time:    end_time,
			},
		}
		a.DevsIDs = append(a.DevsIDs, dev_id)
		a.DevsAccounts = append(a.DevsAccounts, dev_account)
		a.DevsInfos = append(a.DevsInfos, res)
	}
}

func (a *Api) accounts(c *gin.Context) {
	name := c.Query("accountname")
	if name == "admin" {
		c.JSON(http.StatusOK,
			gin.H{
				"code": 1000,
				"msg":  "",
				"data": a.AdminInfos,
			})
	} else if name == "seats" {
		c.JSON(http.StatusOK,
			gin.H{
				"code": 1000,
				"msg":  "",
				"data": a.SeatsInfos,
			})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1000,
			"msg":  "",
			"data": a.DevsInfos,
		})
	}
}

func (a Api) devmap(devs []string) (devss []Dev) {
	for _, v := range a.Devs {
		for _, dev := range devs {
			if v[dev].Finger == dev {
				d := Dev{
					Finger: dev,
					Name:   v[dev].Name,
				}
				devss = append(devss, d)
			}
		}
	}
	return

}
