package models

import (
	"encoding/json"

	"winddies/manage-api/code"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var articleContentCollection = "content-article"
var articleTagCollection = "tag-article"

const articleIDInitialNumber = 100

// ArticleDetail 文章详情
type ArticleDetail struct {
	ID            bson.ObjectId `bson:"_id" json:"id,omitempty"`
	PostID        int           `bson:"post_id" json:"post_id"`
	RawContent    string        `bson:"rawContent" json:"rawContent,omitempty"`
	RenderContent string        `bson:"renderContent" json:"renderContent"`
	*ArticleSummary
}

// ArticleContent 存储在 redis 里的内容，同时也是文章详情获取的内容
type ArticleContent struct {
	RelatedID     bson.ObjectId `json:"related_id,omitempty"`
	PostID        int           `json:"post_id"`
	RenderContent string        `json:"renderContent"`
	*ArticleSummary
}

// Comments 文章评论
type Comments struct {
	ID      bson.ObjectId `bson:"_id" json:"id,omitempty"`
	PostID  string        `json:"post_id"`
	User    string        `json:"user"`
	Content string        `json:"content"`
}

// ArticleSummary 文章内容摘要
type ArticleSummary struct {
	ID         bson.ObjectId `bson:"_id" json:"id,omitempty"`
	RelatedID  bson.ObjectId `bson:"related_id" json:"related_id"`
	PostID     int           `bson:"post_id" json:"post_id"`
	Tag        []string      `bson:"tag" json:"tag"`
	Title      string        `bson:"title" json:"title"`
	Desc       string        `bson:"desc" json:"desc"`
	CreateTime int           `bson:"createTime" json:"createTime"`
	PutTime    int           `bson:"putTime" json:"putTime"`
}

// ArticleOperation 相应 interface 实现
type ArticleOperation interface {
	Insert() error
	Update() error
	Delete() error
	List() (summary []ArticleSummary, err error)
	Detail() (content ArticleContent, err error)
}

// NewArticleOperation ArticleOperation 的构造
func NewArticleOperation(params ...interface{}) ArticleOperation {
	var id bson.ObjectId
	if len(params) > 0 {
		id = params[0].(bson.ObjectId)
	} else {
		id = bson.NewObjectId()
	}

	return &ArticleDetail{
		ID: id,
	}
}

// Insert 插入操作
func (articleDetail *ArticleDetail) Insert() (err error) {
	summaryList, err := articleDetail.List()
	length := len(summaryList)
	if length == 0 {
		articleDetail.PostID = articleIDInitialNumber
	} else {
		articleDetail.PostID = length + 1
	}
	// 先将内容分别分到 redis 和 mongo 所需要的字段
	var mongoSummaryContent ArticleSummary
	var redisContent ArticleContent
	bytes, err := json.Marshal(articleDetail)
	err = json.Unmarshal(bytes, &mongoSummaryContent)
	err = json.Unmarshal(bytes, &redisContent)

	redisContent.RelatedID = articleDetail.ID
	mongoSummaryContent.RelatedID = articleDetail.ID

	redisBytes, _ := json.Marshal(redisContent)

	// markdown 存在 mongo
	articleDetail.Query(articleContentCollection, func(c *mgo.Collection) {
		err = c.Insert(articleDetail)
		if err != nil {
			return
		}
	})

	// 摘要、类别专门存一个 collection
	articleDetail.Query(articleTagCollection, func(c *mgo.Collection) {
		err = c.Insert(mongoSummaryContent)
		if err != nil {
			return
		}
	})

	// 渲染的内容放到 redis 缓存
	RedisDb.HSet(ctx, "article:"+string(articleDetail.PostID), string(redisBytes))
	return
}

// Update 更新存储内容
func (articleDetail *ArticleDetail) Update() (err error) {
	articleDetail.Query(articleContentCollection, func(c *mgo.Collection) {
		err = c.UpdateId(articleDetail.ID, articleDetail)
		if err != nil {
			return
		}
	})

	var redisContent ArticleContent
	var articleSummary ArticleSummary
	bytes, err := json.Marshal(articleDetail)
	json.Unmarshal(bytes, &redisContent)
	json.Unmarshal(bytes, &articleSummary)
	redisBytes, err := json.Marshal(redisContent)

	articleDetail.Query(articleTagCollection, func(c *mgo.Collection) {
		err = c.UpdateId(articleDetail.ID, articleSummary)
		if err != nil {
			return
		}
	})

	RedisDb.HSet(ctx, string(articleDetail.ID), string(redisBytes))
	return
}

// Delete 删除存储内容
func (articleDetail *ArticleDetail) Delete() (err error) {

	articleDetail.Query(articleContentCollection, func(c *mgo.Collection) {
		err = c.RemoveId(articleDetail.ID)
		if err != nil {
			return
		}
	})

	articleDetail.Query(articleTagCollection, func(c *mgo.Collection) {
		err = c.RemoveId(articleDetail.ID)
		if err != nil {
			return
		}
	})

	_, err = RedisDb.HDel(ctx, string(articleDetail.ID)).Result()

	return
}

// List 获取文章摘要列表
func (articleDetail *ArticleDetail) List() (summary []ArticleSummary, err error) {
	articleDetail.Query(articleTagCollection, func(c *mgo.Collection) {
		err = c.Find(nil).All(&summary)
	})

	return
}

// Detail 获取文章详情
func (articleDetail *ArticleDetail) Detail() (article ArticleContent, err error) {
	data, err := RedisDb.HMGet(ctx, "article"+string(articleDetail.PostID)).Result()
	bytes, err := json.Marshal(data)
	json.Unmarshal(bytes, &article)
	return
}

// Query mongo 查找
func (articleDetail *ArticleDetail) Query(name string, query func(c *mgo.Collection)) {
	mgoQuery.Query(name, query)
}

// Find 根据指定 id 查找文章详情
func (articleDetail *ArticleDetail) Find(id string) (result *ArticleDetail, err error) {
	if !bson.IsObjectIdHex(id) {
		return nil, code.ErrInvalidId
	}
	articleDetail.Query(articleContentCollection, func(c *mgo.Collection) {
		err = c.FindId(bson.ObjectIdHex(id)).One(&result)
	})
	return
}
