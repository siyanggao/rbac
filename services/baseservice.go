package services

import "github.com/astaxie/beego/cache"

var bm cache.Cache

func init() {
	bm, _ = cache.NewCache("memory", `{"interval":60}`)
}

type baseService struct {
}

func (this *baseService) ParsePage(page int, length int) int {
	if page == 0 {
		page = 1
	}
	return (page - 1) * length
}
