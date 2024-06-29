package jobs

import (
	"fmt"
	"time"
)

func exampleJob() {
	fmt.Printf("Every seconds, %s\n", time.Now().Format("15:04:05"))
}
