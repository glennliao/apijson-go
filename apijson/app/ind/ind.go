package ind

import (
	"github.com/glennliao/apijson-go/apijson/app/ind/internal/api"
)

func Routers() []any {
	return []any{
		api.Access,
		api.Request,
		api.Table,
	}
}
