// Copyright (c) 2015 HPE Software Inc. All rights reserved.
// Copyright (c) 2013 ActiveState Software Inc. All rights reserved.

package main

import (
	"github.com/hpcloud/tail/notify"
)

func main() {

	watcher := notify.NewNotifier()
	watcher.WatchDir()
	select {}
}
