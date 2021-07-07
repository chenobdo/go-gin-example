package tag_service

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

type Tag struct {
	ID         int
	Name       string
	State      int
	CreatedBy  string
	ModifiedBy string

	PageNum  int
	PageSize int
}

func (t *Tag) Import(r io.Reader) error {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	rows := xlsx.GetRows("标签信息")
	for irow, row := range rows {
		if irow > 0 {
			var data []string
			for _, cell := range row {
				data = append(data, cell)
			}

			models.AddTag(data[1], 1, data[2])
		}
	}

	return nil
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMaps())
}

func (t *Tag) GetAll() ([]*models.Tag, error) {
	var (
		tags, cacheTags []*models.Tag
	)

	cache := cache_service.Tag{
		Name:     t.Name,
		State:    t.State,
		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}

	key := cache.GetTagsKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			err := json.Unmarshal(data, &cacheTags)
			if err != nil {
				return nil, err
			}
			return cacheTags, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	err = gredis.Set(key, tags, 3600)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("标签信息")
	if err != nil {
		return "", err
	}

	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	row := sheet.AddRow()

	var cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}

	for _, v := range tags {
		values := []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			strconv.Itoa(v.CreatedOn),
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
		}

		row = sheet.AddRow()
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}

	timestamp := strconv.Itoa(int(time.Now().Unix()))
	fileName := "tags-" + timestamp + ".xlsx"

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

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State != -1 {
		maps["state"] = t.State
	}

	return maps
}
