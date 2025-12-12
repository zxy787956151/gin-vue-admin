package localai

import "github.com/flipped-aurora/gin-vue-admin/server/service/localai"

type ApiGroup struct {
	LocalAIApi
}

var (
	localaiService = localai.GetRAGService()
)

