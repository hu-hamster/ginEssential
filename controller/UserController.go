package controller

import (
	"hu-hamster/ginEssential/common"
	"hu-hamster/ginEssential/dto"
	"hu-hamster/ginEssential/model"
	"hu-hamster/ginEssential/response"
	"hu-hamster/ginEssential/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(ctx *gin.Context) {
	//获取参数
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	// var requestMap = make(map[string]string)
	// json.NewDecoder(ctx.Request.Body).Decode(&requestMap)
	// var resquestUser = model.User{}
	// ctx.Bind(&resquestUser)
	// name := resquestUser.Name
	// telephone := resquestUser.Telephone
	// password := resquestUser.Password
	//数据验证
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能小于6位")

		return
	}
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	//判断电话号是否存在
	if isTelephoneExist(common.DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		return
	}

	//创建用户
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "密码加密失败")
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashPassword),
	}
	common.DB.Create(&newUser)

	//发放token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}
	//返回结果
	response.Response(ctx, http.StatusUnprocessableEntity, 200, gin.H{"token": token}, "登录成功")

	log.Println(name, telephone, password)
	response.Response(ctx, http.StatusUnprocessableEntity, 200, nil, "注册成功")
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone=?", telephone).First(&user)
	return user.ID != 0
}

func Login(ctx *gin.Context) {
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	//数据验证
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能小于6位")
		return
	}
	//判断手机号是否存在
	var user model.User
	common.DB.Where("telephone=?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "密码错误",
		})
		return
	}
	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}
	//返回结果
	response.Response(ctx, http.StatusUnprocessableEntity, 200, gin.H{"token": token}, "登录成功")
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	response.Response(ctx, http.StatusUnprocessableEntity, 200, gin.H{"user": dto.ToUserDto(user.(model.User))}, "")
}
