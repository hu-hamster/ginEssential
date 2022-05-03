package controller

import (
	"hu-hamster/ginEssential/common"
	"hu-hamster/ginEssential/model"
	"hu-hamster/ginEssential/repository"
	"hu-hamster/ginEssential/response"
	"hu-hamster/ginEssential/vo"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ICategoryController interface {
	RestController
}

type CategoryController struct {
	Repository repository.CategoryRepository
}

func NewCategoryController() CategoryController {
	repository := repository.NewCategoryRepository()
	db := common.DB
	db.AutoMigrate(model.Category{})

	return CategoryController{Repository: repository}
}

func (c CategoryController) Create(ctx *gin.Context) {
	// name := ctx.PostForm("name")
	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}
	var category *model.Category
	var err error
	if category, err = c.Repository.Create(requestCategory.Name); err != nil {
		panic(err)
	}
	response.Success(ctx, gin.H{"category": category}, "")
}

func (c CategoryController) Update(ctx *gin.Context) {
	//绑定body中的参数
	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}

	//获取path中的参数
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	updateCategory, err := c.Repository.SelectById(categoryId)
	if err != nil {
		response.Fail(ctx, nil, "分类不存在")
		return
	}
	category, err2 := c.Repository.Update(*updateCategory, requestCategory.Name)
	if err2 != nil {
		panic(err)
	}

	response.Success(ctx, gin.H{"category": category}, "修改成功")

}

func (c CategoryController) Show(ctx *gin.Context) {
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	category, err := c.Repository.SelectById(categoryId)
	if err != nil {
		response.Fail(ctx, nil, "分类不存在")
		return
	}
	response.Success(ctx, gin.H{"category": category}, "")
}

func (c CategoryController) Delete(ctx *gin.Context) {
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))

	if err := c.Repository.DeleteById(categoryId); err != nil {
		response.Fail(ctx, nil, "删除失败，请重试")
		return
	}
	response.Success(ctx, nil, "")
}
