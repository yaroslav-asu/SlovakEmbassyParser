package init

import (
	"main/internal/logger"
	"main/internal/utils/random"
	"main/internal/utils/variables"
)

func Init() {
	random.InitRandom()
	logger.InitLogger()
	variables.InitEnv()
}
