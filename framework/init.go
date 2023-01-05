package framework

import (
	"github.com/glennliao/apijson-go/config/db"
	_ "github.com/glennliao/apijson-go/framework/gf_orm"
)

func Init() {
	db.Init()
}
