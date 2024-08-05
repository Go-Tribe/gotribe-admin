// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package jobs

import (
	"fmt"
	"time"
)

func exampleJob() {
	fmt.Printf("Every seconds, %s\n", time.Now().Format("15:04:05"))
}
