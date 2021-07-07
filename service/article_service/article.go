package article_service

import (
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/chenobdo/go-gin-example/models"
	"github.com/chenobdo/go-gin-example/pkg/export"
	"github.com/chenobdo/go-gin-example/pkg/gredis"
	"github.com/chenobdo/go-gin-example/pkg/logging"
	"github.com/chenobdo/go-gin-example/service/cache_service"
	"github.com/tealeg/xlsx"
	"io"
	"strconv"
	"time"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

func (a *Article) Import(r io.Reader) error {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	rows := xlsx.GetRows("文章信息")
	for irow, row := range rows {
		if irow > 0 {
			data := make(map[string]interface{})
			data["tag_id"] = 0
			data["title"] = row[1]
			data["desc"] = row[2]
			data["content"] = row[3]
			data["created_by"] = row[5]
			state, _ := strconv.Atoi(row[8])
			data["state"] = state

			models.AddArticle(data)
		}
	}

	return nil
}

func (a *Article) Export() (string, error) {
	articles, err := a.GetAll()
	if err != nil {
		return "", err
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("文章信息")
	if err != nil {
		return "", err
	}

	titles := []string{"ID", "文章标题", "简述", "内容", "创建时间", "创建人", "修改时间", "修改人", "状态"}
	row := sheet.AddRow()

	var cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}

	for _, v := range articles {
		values := []string{
			strconv.Itoa(v.ID),
			v.Title,
			v.Desc,
			v.Content,
			strconv.Itoa(v.CreatedOn),
			v.CreatedBy,
			strconv.Itoa(v.ModifiedOn),
			v.ModifiedBy,
			strconv.Itoa(v.State),
		}

		row = sheet.AddRow()
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}

	timestamp := strconv.Itoa(int(time.Now().Unix()))
	fileName := "articles-" + timestamp + ".xlsx"

	fullPath := export.GetExcelFullPath() + fileName

	err = export.CheckDir(export.GetExcelFullPath())
	if err != nil {
		logging.Warn(err)
		return "", err
	}

	err = file.Save(fullPath)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}

func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}

	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	gredis.Set(key, article, 3600)
	return article, nil
}

func (a *Article) GetAll() ([]*models.Article, error) {
	var (
		articles, cacheArticles []*models.Article
	)

	cache := cache_service.Article{
		TagID: a.TagID,
		State: a.State,

		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}

	key := cache.GetArticlesKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}

	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(key, articles, 3600)
	return articles, nil
}

func (a *Article) Count() (int, error) {
	return models.GetArticleTotal(a.getMaps())
}

func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}

	return maps
}
