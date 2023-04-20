package funcs

import (
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

func Sleep() {
	SleepTime(2, 5)
}

func SleepTime(from, to float64) {
	sleepingTime := rand.Float64()*(to-from) + from
	zap.L().Info("Started random sleeping with time: " + fmt.Sprintf("%.2f", sleepingTime))
	time.Sleep(time.Duration(sleepingTime) * time.Second)
	zap.L().Info("Finished random sleeping with time: " + fmt.Sprintf("%.2f", sleepingTime))
}
