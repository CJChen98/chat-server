package ws

import (
	"encoding/json"
	"gin/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

//客户端连接详情
type wsClient struct {
	Conn       *websocket.Conn `json:"conn"`
	RemoteAddr string          `json:"remote_addr"`
	Uid        float64         `json:"uid"`
	Username   string          `json:"username"`
	RoomId     string          `json:"room_id"`
	AvatarPath string          `json:"avatar_path"`
}

//client 和 serve 消息体
type msg struct {
	Status int             `json:"status"`
	Data   interface{}     `json:"data"`
	Conn   *websocket.Conn `json:"conn"`
}

var (
	wsUpgrader = websocket.Upgrader{}
	clientMsg  = msg{}
	mutex      = sync.Mutex{}
	rooms      = make(map[int][]wsClient)
	enterRooms = make(chan wsClient)
	sMsg       = make(chan msg)
	offline    = make(chan *websocket.Conn)
)

// 定义消息类型
const msgTypeOnline = 1        // 上线
const msgTypeOffline = 2       // 离线
const msgTypeSend = 3          // 消息发送
const msgTypeGetOnlineUser = 4 // 获取用户列表
const msgTypePrivateChat = 5   // 私聊

const roomCount = 6 // 房间总数
type GoServe struct {
	ServeInterface
}

func (goServe *GoServe) RunWs(c *gin.Context) {
	Run(c)
}
func (goServe *GoServe) GetOnlineUserCount() int {
	return GetOnlineUserCount()
}
func (goServe *GoServe) GetOnlineRoomUserCount(rId int) int {
	return GetOnlineRoomUserCount(rId)
}
func Run(c *gin.Context) {
	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, _ := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	defer conn.Close()
	log.Println("http update to ws!")
	go read(conn)
	go write()
	select {}
}

func read(c *websocket.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("read 发生错误:", err)
		}
	}()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			offline <- c
			log.Println("read message 发生错误", err)
			break
		}

		serveMsg := message
		if string(serveMsg) == `heartbeat` {
			_ = c.WriteMessage(websocket.TextMessage, []byte(`{"status":0,"data":"heartbeat ok"}`))
			continue
		}
		_ = json.Unmarshal(message, &clientMsg)
		if clientMsg.Data != nil {
			if clientMsg.Status == msgTypeOnline { //进入房间
				rid, _ := getRoomID()
				enterRooms <- wsClient{
					Conn:       c,
					RemoteAddr: c.RemoteAddr().String(),
					Uid:        clientMsg.Data.(map[string]interface{})["uid"].(float64),
					Username:   clientMsg.Data.(map[string]interface{})["username"].(string),
					RoomId:     rid,
					AvatarPath: clientMsg.Data.(map[string]interface{})["avatar_path"].(string),
				}
			}
			_, serveMsg := formatServeMsgStr(clientMsg.Status, c)
			sMsg <- serveMsg
		}
	}
}

func write() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("read 发生错误:", err)
		}
	}()

	for {
		select {
		case r := <-enterRooms:
			handleConnClients(r.Conn)
		case cl := <-sMsg:
			serveMsgStr, _ := json.Marshal(cl)
			switch cl.Status {
			case msgTypeOnline, msgTypeSend:
				notify(cl.Conn, string(serveMsgStr))
			case msgTypeGetOnlineUser:
				_ = cl.Conn.WriteMessage(websocket.TextMessage, serveMsgStr)
			case msgTypePrivateChat:
				toC := findToUserCoonClient()
				if toC != nil {
					_ = toC.(wsClient).Conn.WriteMessage(websocket.TextMessage, serveMsgStr)
				}
			}
		case o := <-offline:
			disconnect(o)
		}
	}
}

func handleConnClients(c *websocket.Conn) {
	rid, ridInt := getRoomID()
	for cKey, wcl := range rooms[ridInt] {
		if wcl.Uid == clientMsg.Data.(map[string]interface{})["uid"].(float64) {
			mutex.Lock()
			// 通知当前用户下线
			_ = wcl.Conn.WriteMessage(websocket.TextMessage, []byte(`{"status":-1,"data":[]}`))
			rooms[ridInt] = append(rooms[ridInt][:cKey], rooms[ridInt][cKey+1:]...)
			wcl.Conn.Close()
			mutex.Lock()

		}
	}
	mutex.Lock()
	rooms[ridInt] = append(rooms[ridInt], wsClient{
		Conn:       c,
		RemoteAddr: c.RemoteAddr().String(),
		Uid:        clientMsg.Data.(map[string]interface{})["uid"].(float64),
		Username:   clientMsg.Data.(map[string]interface{})["username"].(string),
		RoomId:     rid,
		AvatarPath: clientMsg.Data.(map[string]interface{})["avatar_path"].(string),
	})
	mutex.Unlock()
}

// 格式化传送给客户端的消息数据
func formatServeMsgStr(status int, conn *websocket.Conn) ([]byte, msg) {

	roomId, roomIdInt := getRoomID()

	data := map[string]interface{}{
		"username": clientMsg.Data.(map[string]interface{})["username"].(string),
		"uid":      clientMsg.Data.(map[string]interface{})["uid"].(float64),
		"room_id":  roomId,
		"time":     time.Now().UnixNano() / 1e6, // 13位  10位 => now.Unix()
	}

	if status == msgTypeSend || status == msgTypePrivateChat {
		data["avatar_id"] = clientMsg.Data.(map[string]interface{})["avatar_id"].(string)
		data["content"] = clientMsg.Data.(map[string]interface{})["content"].(string)

		toUidStr := clientMsg.Data.(map[string]interface{})["to_uid"].(string)
		toUid, _ := strconv.Atoi(toUidStr)

		// 保存消息
		stringUid := strconv.FormatFloat(data["uid"].(float64), 'f', -1, 64)
		intUid, _ := strconv.Atoi(stringUid)

		if _, ok := clientMsg.Data.(map[string]interface{})["image_url"]; ok {
			// 存在图片
			models.SaveContent(map[string]interface{}{
				"user_id":    intUid,
				"to_user_id": toUid,
				"content":    data["content"],
				"room_id":    data["room_id"],
				"image_url":  clientMsg.Data.(map[string]interface{})["image_url"].(string),
			})
		} else {
			models.SaveContent(map[string]interface{}{
				"user_id":    intUid,
				"to_user_id": toUid,
				"room_id":    data["room_id"],
				"content":    data["content"],
			})
		}

	}

	if status == msgTypeGetOnlineUser {
		ro := rooms[roomIdInt]
		data["count"] = len(ro)
		data["list"] = ro
	}

	jsonStrServeMsg := msg{
		Status: status,
		Data:   data,
		Conn:   conn,
	}
	serveMsgStr, _ := json.Marshal(jsonStrServeMsg)

	return serveMsgStr, jsonStrServeMsg
}

// 获取私聊的用户连接
func findToUserCoonClient() interface{} {
	_, roomIdInt := getRoomID()

	toUserUid := clientMsg.Data.(map[string]interface{})["to_uid"].(string)

	for _, c := range rooms[roomIdInt] {
		stringUid := strconv.FormatFloat(c.Uid, 'f', -1, 64)
		if stringUid == toUserUid {
			return c
		}
	}

	return nil
}

// 统一消息发放
func notify(conn *websocket.Conn, msg string) {
	_, roomIdInt := getRoomID()
	for _, con := range rooms[roomIdInt] {
		if con.RemoteAddr != conn.RemoteAddr().String() {
			con.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
}

// 离线通知
func disconnect(conn *websocket.Conn) {
	_, roomIdInt := getRoomID()
	for index, con := range rooms[roomIdInt] {
		if con.RemoteAddr == conn.RemoteAddr().String() {
			data := map[string]interface{}{
				"username": con.Username,
				"uid":      con.Uid,
				"time":     time.Now().UnixNano() / 1e6, // 13位  10位 => now.Unix()
			}

			jsonStrServeMsg := msg{
				Status: msgTypeOffline,
				Data:   data,
			}
			serveMsgStr, _ := json.Marshal(jsonStrServeMsg)

			disMsg := string(serveMsgStr)

			mutex.Lock()
			rooms[roomIdInt] = append(rooms[roomIdInt][:index], rooms[roomIdInt][index+1:]...)
			con.Conn.Close()
			mutex.Unlock()
			notify(conn, disMsg)
		}
	}
}
func getRoomID() (string, int) {
	rid := clientMsg.Data.(map[string]interface{})["room_id"].(string)
	ridInt, _ := strconv.Atoi(rid)
	return rid, ridInt
}

func GetOnlineUserCount() int {
	num := 0
	for i := 1; i <= roomCount; i++ {
		num = num + GetOnlineRoomUserCount(i)
	}
	return num
}

func GetOnlineRoomUserCount(roomId int) int {
	return len(rooms[roomId])
}
