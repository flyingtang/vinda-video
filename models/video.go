package models

import (
	"github.com/jmoiron/sqlx"
	"time"
	"video/constant"
)

type Video struct {
	Id          int    `json:"id" db:"id"`
	OriginName  string    `json:"originName" db:"origin_name"`
	UrlName     string    `json:"urlName" db:"url_name"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Status      int       `json:"status" db:"status"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

func CreateVideo(v *Video) error {
	const sql = "insert into tb_video (origin_name, url_name, title, description, status) value (?, ?, ?, ?, ?);"
	_, err := globalDB.Exec(sql, v.OriginName, v.UrlName, v.Title, v.Description, v.Status)
	return err
}

func FindVideoById(id string) (*Video, error) {
	var v Video
	const sql = "select origin_name, url_name from tb_video where id = ?;"
	err := globalDB.Get(&v, sql, id)
	return &v, err
}

func PatchVideo(id string, v *Video) error {
	const sql = "update tb_video set origin_name = ?, url_name= ?, title=?, description= ?, updated_at=?, status=?"
	_, err := globalDB.Exec(sql, v.OriginName, v.UrlName, v.Title, v.Description, time.Now(), v.Status)
	return err
}

func DeleteVideo(id string) error{
	const sql  = "update tb_video set status=0 where id = ?;"
	_,err := globalDB.Exec(sql, id)
	return err
}

func DeleteALLVideo(ids []int) error {

	const sql  = "update tb_video set status=0 where id in (?);"
	query, args, err := sqlx.In(sql, ids)
	query = globalDB.Rebind(query)
	_, err = globalDB.Query(query, args...)
	return err
}

// TODO 如果查询出错会直接崩掉，这个得解决
func QueryVideos(page int64) (vs []Video, total int64, err error) {
	pageSize := constant.PageSize
	offset := int(page -1) * pageSize
	const datasql = "select * from tb_video where status = 1 order by created_at desc limit ? offset ?;"
	//const datasql = "select * from tb_video where status = 1 order by created_at desc limit ? offset ?;"
	err = globalDB.Select(&vs, datasql, pageSize, offset)
	if err != nil {
		return
	}
	const totalSql = "select count(*) as total from  tb_video where status = 1"
	err = globalDB.Get(&total, totalSql)
	return
}
