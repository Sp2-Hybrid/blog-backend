package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go-gin-example/models"
	"go-gin-example/pkg/e"
	"go-gin-example/pkg/setting"
	"go-gin-example/pkg/util"
	"log"
	"net/http"
)

func GetArticle(c *gin.Context){
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	var data interface{}
	if !valid.HasErrors(){
		if models.ExistTagByID(id){
			code = e.SUCCESS
			data = models.GetArticle(id)
		}else{
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	}else{
		for _, err := range valid.Errors{
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : data,
	})
}

func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state");arg!=""{
		state = com.StrTo(arg).MustInt()
		maps["state"] = state

		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	var tagId int = -1
	if arg := c.Query("tag_id"); arg != ""{
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId

		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}
	code := e.INVALID_PARAMS
	if !valid.HasErrors(){
		code = e.SUCCESS
		data["lists"] = models.GetArticles(util.GetPage(c), setting.PageSize, maps)
		data["total"] = models.GetTagTotal(maps)
	}else{
		for _, err := range valid.Errors{
			log.Printf("err.Key: %s, err.Message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : data,
	})
}

func AddArticle(c *gin.Context)  {
	// com.Stro:转换成字符串
	// MustInt:转换成整数
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()

	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors(){
		if models.ExistTagByID(tagId){
			code = e.SUCCESS
			data := make(map[string]interface{})
			data["tag_id"] = tagId
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createdBy
			data["state"] = state

			models.AddArticle(data)
		}else{
			code = e.ERROR_NOT_EXIST_TAG
		}
	}else{
		for _, err := range valid.Errors{
			log.Printf("err.Key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]interface{}),
	})
}

func EditArticle(c *gin.Context)  {
	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("文章ID必须大于0")
	valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")

	var state int = -1
	if arg:=c.Query("state");arg!=""{
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("state取值只能为0或1")
	}

	code := e.INVALID_PARAMS
	if !valid.HasErrors(){
		if models.ExistArticleByID(id){
			if models.ExistTagByID(tagId){
				code = e.SUCCESS
				data := make(map[string]interface{})
				data["tag_id"] = tagId
				if title!=""{
					data["title"] = title
				}
				if desc != ""{
					data["desc"] = desc
				}
				if content != ""{
					data["content"] = content
				}

				data["modified_by"] = modifiedBy
				models.EditArticle(id, data)
			}else{
				code = e.ERROR_NOT_EXIST_TAG
			}
		}else{
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	}else{
		for _, err := range valid.Errors{
			log.Printf("err.Key: %s, err.Message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]interface{}),
	})
}

func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("文章ID必须大于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors(){
		if models.ExistArticleByID(id){
			code = e.SUCCESS
			models.DeleteArticle(id)
		}else{
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	}else{
		for _, err := range valid.Errors{
			log.Printf("err.Key: %s, err.Message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]interface{}),
	})
}
