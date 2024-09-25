package service

import (
	"GinChat/models"
	"GinChat/utils"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Summary 获取用户列表
// @Tags 用户模块
// @Success 200 {string} json{"code", "message", "data"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "用户名已注册",
		"data":    data,
	})
}

// FindUserByNameAndPwd
// @Summary 名字密码查询用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code", "message", "data"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	// 用户是否存在
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "该用户不存在",
			"data":    data,
		})
		return
	}
	// 判断密码
	flag := utils.ValidPassword(password, user.Salt, user.Password)
	if !flag {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "密码错误",
			"data":    data,
		})
		return
	}
	// 查询
	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data":    data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code", "message", "data"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("Identity")
	fmt.Println(user.Name, password, repassword)
	// 用户名判断
	data := models.FindUserByName(user.Name)
	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "用户名或密码不能为空！",
			"data":    user,
		})
		return
	}
	if data.Name != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "该用户名已存在",
			"data":    user,
		})
		return
	}
	// 密码判断
	if password != repassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "两次输入的密码不一致",
			"data":    user,
		})
		return
	}
	user.LoginTime = time.Now()
	user.HeartbeatTime = time.Now()
	user.LoginOutTime = time.Now()
	// 密码加密
	salt := fmt.Sprintf("%06d", rand.Int31())
	user.Salt = salt
	user.Password = utils.MakePassword(password, salt)
	models.CreateUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建用户成功！",
		"data":    user,
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code", "message", "data"}
// @Router /user/deleteUser [post]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除用户成功！",
		"data":    user,
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code", "message", "data"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	// 从PUT请求中拿到id
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Avatar = c.PostForm("icon")
	user.Email = c.PostForm("email")

	// govalidator校验
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "修改参数不匹配!",
			"data":    user,
		})
	} else {
		models.UpdateUser(user)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "修改用户成功！",
			"data":    user,
		})
	}
}

// websocket的upgrate设定检查，防止跨域站点的伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// SendMsg
// @Summary 发送消息
// @Tags 聊天模块
// @Success 200 {string} json{"code", "message", "data"}
// @Router /user/sendMsg [get]
func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println("close:", err)
		}
	}(ws)
	MsgHandler(ws, c)
}
func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println("subscribe err:", err)
		}
		tm := time.Now().Format("2006-01-02T15:04:05")
		// 消息格式化
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println("WriteMessage err:", err)
		}
	}
}

// SendUserMsg
// @Summary 用户发送
// @Tags 聊天模块
// @Success 200 {string} json{"code", "message", "data"}
// @Router /user/sendUserMsg [get]
func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

// SearchFriends
// @Summary 加载好友
// @Tags 聊天模块
// @param userId formData string false "userId"
// @Success 200 {string} json{"code", "message", "data"}
// @Router /searchFriends [post]
func SearchFriends(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("userId"))
	users := models.SearchFriends(uint(id))
	// 查询OK列表
	utils.RespOKList(c.Writer, users, len(users))
}

// AddFriend
// @Summary 创建群聊
// @Tags 聊天模块
// @param userId formData string false "userId"
// @param userId formData string false "targetId"
// @Success 200 {string} json{"code", "message", "data"}
// @Router /contact/addfriend [post]
func AddFriend(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("userId"))
	targetid, _ := strconv.Atoi(c.Request.FormValue("targetId"))
	code, msg := models.AddFriend(uint(id), uint(targetid))
	// 查询OK列表
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// CreateCommunity
// @Summary 创建群聊
// @Tags 聊天模块
// @param ownerId formData string false "ownerId"
// @param Name formData string false "name"
// @param Img formData string false "icon"
// @param Desc formData string false "desc"
// @Success 200 {string} json{"code", "message"}
// @Router /contact/createCommunity [post]
func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	name := c.Request.FormValue("name")
	icon := c.Request.FormValue("icon")
	desc := c.Request.FormValue("desc")
	community := models.Community{}
	community.OwnerId = uint(ownerId)
	community.Name = name
	community.Img = icon
	community.Desc = desc
	code, msg := models.CreateCommunity(community)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// LoadCommunity
// @Summary 加载群列表
// @Tags 聊天模块
// @param ownerId formData string false "ownerId"
// @Success 200 {string} json{"code", "message"}
// @Router /contact/loadcommunity [post]
func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	//	name := c.Request.FormValue("name")
	data, msg := models.LoadCommunity(uint(ownerId))
	if len(data) != 0 {
		utils.RespList(c.Writer, 0, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// JoinGroups
// @Summary 加入群聊
// @Tags 聊天模块
// @param userId formData string false "userId"
// @param comId formData string false "comId"
// @Success 200 {string} json{"code", "message"}
// @Router /contact/joinGroup [post]
func JoinGroups(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	comId := c.Request.FormValue("comId")

	//	name := c.Request.FormValue("name")
	data, msg := models.JoinGroup(uint(userId), comId)
	if data == 0 {
		utils.RespOK(c.Writer, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

func FindByID(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	//	name := c.Request.FormValue("name")
	data := models.FindByID(uint(userId))
	utils.RespOK(c.Writer, data, "ok")
}

// RedisMsg
// @Summary 载入聊天缓存
// @Tags 聊天模块
// @param userIdA formData string false "userIdA"
// @param userIdB formData string false "userIdB"
// @param start formData string false "start"
// @param end formData string false "end"
// @param isRev formData bool false "isRev"
// @Success 200 {string} json{"code", "message", "data}
// @Router /user/redisMsg [post]
func RedisMsg(c *gin.Context) {
	userIdA, _ := strconv.Atoi(c.PostForm("userIdA"))
	userIdB, _ := strconv.Atoi(c.PostForm("userIdB"))
	start, _ := strconv.Atoi(c.PostForm("start"))
	end, _ := strconv.Atoi(c.PostForm("end"))
	isRev, _ := strconv.ParseBool(c.PostForm("isRev"))
	res := models.RedisMsg(int64(userIdA), int64(userIdB), int64(start), int64(end), isRev)
	utils.RespOKList(c.Writer, "ok", res)
}
