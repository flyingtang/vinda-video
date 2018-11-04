package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"video/constant"
)

var schemas = []string{
	`
	create table if not exists tb_video (
		id int primary key auto_increment,
		origin_name varchar(255) not null,
		url_name varchar(255) not null,
		title varchar(255) not null unique,
		description text not null,
		status TINYINT default 1,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
		index index_url_name(url_name),
		index index_created_at(created_at desc)
	)
	`,
}
var globalDB *sqlx.DB

func init() {
	db, err := sqlx.Connect("mysql", constant.MysqlUrl)
	if err != nil {
		panic(err.Error())
	}

	tx := db.MustBegin()
	for i := 0; i < len(schemas); i++ {
		tx.MustExec(schemas[i])
	}
	err = tx.Commit()
	if err != nil {
		panic("initial database table error")
	}
	globalDB = db
}
