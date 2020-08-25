package article

import (
	"winddies/manage-api/code"
	"winddies/manage-api/controllers/base"
	"winddies/manage-api/models"

	"github.com/gin-gonic/gin"
)

type manage struct {
	*base.Base
}

type Manage interface {
	CreateNewArticle(ctx *gin.Context)
	UpdateArticle(ctx *gin.Context)
	DeleteArticle(ctx *gin.Context)
	GetArticleSummaryList(ctx *gin.Context)
}

func Create() Manage {
	return &manage{}
}

// 新建文章
func (manage *manage) CreateNewArticle(ctx *gin.Context) {
	logger := manage.Logger(ctx)

	articleOperation := models.NewArticleOperation()
	err := ctx.BindJSON(&articleOperation)
	if err != nil {
		logger.Errorf("CreateNewArticle error:%s", err)
		return
	}
	err = articleOperation.Insert()
	if err != nil {
		logger.Errorf("articleOperation insert err: %s", err)
		return
	}

	manage.Send(ctx, code.OK, nil)
}

// 更新文章
func (manage *manage) UpdateArticle(ctx *gin.Context) {
	logger := manage.Logger(ctx)

	articleOperation := &models.ArticleDetail{}
	err := ctx.BindJSON(articleOperation)
	if err != nil {
		logger.Errorf("UpdateArticle error:%s", err)
		return
	}
	err = articleOperation.Update()
	if err != nil {
		logger.Errorf("articleOperation update err: %s", err)
		return
	}

	manage.Send(ctx, code.OK, nil)
}

// 删除文章
func (manage *manage) DeleteArticle(ctx *gin.Context) {
	logger := manage.Logger(ctx)
	id := ctx.Param("id")

	articleOperation := models.NewArticleOperation(id)

	err := articleOperation.Delete()
	if err != nil {
		logger.Errorf("DeleteArticle error:%s", err)
		return
	}

	manage.Send(ctx, code.OK, nil)
}

// 获取文章列表
func (manage *manage) GetArticleSummaryList(ctx *gin.Context) {
	logger := manage.Logger(ctx)

	articleOperation := models.ArticleDetail{}

	articleSummaryList, err := articleOperation.List()
	if err != nil {
		logger.Errorf("articleOperation List err: %s", err)
		return
	}

	manage.Send(ctx, code.OK, articleSummaryList)
}

// 获取文章详情
func (manage *manage) GetArticleDetail(ctx *gin.Context) {
	logger := manage.Logger(ctx)
	id := ctx.Param("id")
	articleOperation := models.NewArticleOperation(id)
	detail, err := articleOperation.Detail()
	if err != nil {
		logger.Errorf("articleOperation get detail err: %s", err)
	}

	manage.Send(ctx, code.OK, detail)
}
