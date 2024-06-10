package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/whitexwc/basic-go/webook/internal/domain"
	"github.com/whitexwc/basic-go/webook/internal/service"

	regexp "github.com/dlclark/regexp2"

	"github.com/gin-gonic/gin"
)

// UserHandler 定义用户相关的所有路由
type UserHandler struct {
	svc         *service.UserService
	EmailExp    *regexp.Regexp
	PasswordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		// 和上面比起来，用 ` 看起来就比较清爽
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)

	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		EmailExp:    emailExp,
		PasswordExp: passwordExp,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	{
		ug.POST("/signup", u.SignUp)

		ug.POST("/login", u.Login)

		ug.POST("/edit", u.Edit)

		ug.GET("/profile", u.Profile)
	}

}

func (u *UserHandler) SignUp(c *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq
	// bind 方法解析错了会直接返回 4xx 的错误
	if err := c.Bind(&req); err != nil {
		return
	}

	ok, err := u.EmailExp.MatchString(req.Email)
	if err != nil {
		// 记录日志
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		c.String(http.StatusBadRequest, "你的邮箱格式错误")
		return
	}

	if req.ConfirmPassword != req.Password {
		c.String(http.StatusOK, "两次输入的密码不一致")
	}

	ok, err = u.PasswordExp.MatchString(req.Password)
	if err != nil {
		// 记录日志
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		c.String(http.StatusBadRequest, "密码必须大于8位，包含特殊字符、数字")
		return
	}

	// 调用 svc 的方法
	err = u.svc.SignUp(c, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		c.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}

	c.String(http.StatusOK, "注册成功")
	fmt.Printf("%v", req)
}

func (u *UserHandler) Login(c *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return
	}
	du, err := u.svc.Login(c, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		c.String(http.StatusOK, "用户名或密码错误")
		return
	}
	if err == service.ErrUserNotFound {
		c.String(http.StatusOK, "用户不存在")
		return
	}
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	// 在这里登陆成功了
	// 设置 session
	sess := sessions.Default(c)
	// 可以随便设置值
	sess.Set("userId", du.Id)
	sess.Save()
	c.String(http.StatusOK, "登陆成功")
}

func (u *UserHandler) Edit(c *gin.Context) {
	type editReq struct {
		NickName string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutme"`
	}

	// 绑定参数
	var e editReq
	if err := c.Bind(&e); err != nil {
		return
	}
	userId := c.Query("id")
	if userId == "" {
		c.String(http.StatusOK, "用户id为空")
		return
	}
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.String(http.StatusOK, "用户id错误")
		return
	}

	// 校验参数
	if len(e.NickName) > 10 {
		c.String(http.StatusOK, "昵称不能超过10个字符")
		return
	}
	if len(e.AboutMe) > 1024 {
		c.String(http.StatusOK, "简介不超过1024个字符")
		return
	}
	birthday, err := time.Parse(time.DateOnly, e.Birthday)
	if err != nil {
		c.String(http.StatusOK, "生日格式不对")
		return
	}

	// 更新用户信息
	err = u.svc.UpdateUserProfile(c, domain.User{
		Id:       id,
		Birthday: birthday,
		NickName: e.NickName,
		AboutMe:  e.AboutMe,
	})

	if err != nil {
		if err == service.ErrUserNotFound {
			c.String(http.StatusOK, "用户不存在, 更新用户信息失败")
		} else {
			c.String(http.StatusOK, "更新用户信息失败")
		}
		return
	}

	// 返回更新成功的结果
	c.String(http.StatusOK, "更新成功")
}

func (u *UserHandler) Profile(c *gin.Context) {
	userId := c.Query("id")
	if userId == "" {
		c.String(http.StatusOK, "用户id为空")
		return
	}
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.String(http.StatusOK, "用户id错误")
		return
	}
	// 根据 id 获取用户profile
	du, err := u.svc.GetUserProfile(c, id)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.String(http.StatusOK, "用户不存在, 获取用户profile失败")
		} else {
			c.String(http.StatusOK, "获取用户Profile失败")
		}
		return
	}
	// 展示用户信息
	c.JSON(http.StatusOK, gin.H{
		"NickName": du.NickName,
		"Birthday": du.Birthday,
		"AboutMe":  du.AboutMe,
	})
}