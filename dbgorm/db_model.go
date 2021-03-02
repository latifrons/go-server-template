package dbgorm

import (
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//type DbArticle struct {
//	gorm.Model
//	ProjectUniqueId       string `gorm:"size:120;index"`
//	PicUrl                string `gorm:"size:120"`
//	Html                  string `gorm:"type:LONGTEXT"`
//	Author                string `gorm:"size:120"`
//	AuthorTitle           string `gorm:"size:120"`
//	AvatarUrl             string `gorm:"size:1000"`
//	AuthorDescriptionHtml string `gorm:"type:LONGTEXT"`
//	Lang                  string `gorm:"size:120;index"`
//	Platform              string `gorm:"size:8;index"`
//  StartAt         time.Time
//}

type User struct {
	gorm.Model
	Email    string `gorm:"size:120;index"`
	Phone    string `gorm:"size:120;index"`
	Password string `gorm:"size:120;index"` // bcrypt
}
