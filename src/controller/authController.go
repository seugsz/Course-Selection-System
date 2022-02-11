package controller

import (
	global "course_select/src/global"
	"course_select/src/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
)

var cookiesName string = "camp-session"

func Login(c *gin.Context) {
	loginRequest := global.LoginRequest{}
	if err := c.ShouldBind(&loginRequest); err != nil {
		c.JSON(http.StatusOK, global.LoginResponse{Code: global.UnknownError})
		return
	}
	log.Println(loginRequest)

	user, err := model.GetMemberByUsernameAndPassword(loginRequest.Username, loginRequest.Password)
	//用户不存在或者密码错误
	if err != nil {
		c.JSON(http.StatusOK, global.LoginResponse{Code: global.WrongPassword})
		return
	}
	//用户已删除
	if user.IsDeleted {
		c.JSON(http.StatusOK, global.LoginResponse{Code: global.UserHasDeleted})
		return
	}

	session := sessions.Default(c)
	var sessionId = getSessionId()

	log.Println(sessionId, user)
	v := global.TMember{
		UserID:   user.UserID,
		Nickname: user.Nickname,
		Username: user.Username,
		UserType: user.UserType,
	}
	session.Set(sessionId, v)
	session.Save()

	c.SetCookie(cookiesName, sessionId, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, global.LoginResponse{
		Code: global.OK,
		Data: struct {
			UserID string
		}{user.UserID},
	})
}

func Logout(c *gin.Context) {
	sessionId, err := c.Cookie(cookiesName)
	if err != nil {
		c.JSON(http.StatusOK, global.LogoutResponse{Code: global.LoginRequired})
		return
	}

	session := sessions.Default(c)

	session.Delete(sessionId)
	session.Save()

	c.SetCookie(cookiesName, sessionId, -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, global.LogoutResponse{Code: global.OK})
}

func WhoAmi(c *gin.Context) {
	sessionId, err := c.Cookie(cookiesName)
	if err != nil {
		c.JSON(http.StatusOK, global.WhoAmIResponse{Code: global.LoginRequired})
		return
	}

	session := sessions.Default(c)
	v := session.Get(sessionId)
	if v == nil {
		c.JSON(http.StatusOK, global.WhoAmIResponse{Code: global.UnknownError})
		return
	}
	log.Println(v)
	user := v.(global.TMember)
	c.JSON(http.StatusOK, global.WhoAmIResponse{Code: global.OK, Data: user})
}

func getSessionId() string {
	//b := make([]byte, 32)
	//if _, err := io.ReadFull(rand.Reader, b); err != nil {
	//	return ""
	//}
	//return base64.StdEncoding.EncodeToString(b)
	return uuid.NewV4().String()
}
