package gc

import (
	"auditlog/database"
	"fmt"
	"time"
)

func Gc(et, tt time.Duration, db *database.Database) {
	sqlstr := "DELETE FROM auditlog WHERE timestamp<=?"
	ticker := time.NewTicker(time.Second * tt)
	for {
		select {
		case <-ticker.C:
			ets := time.Now().Unix() - int64(et.Seconds())
			args := []interface{}{ets}
			c, err := db.Delete(sqlstr, args...)
			fmt.Printf("Deleted %d expired log(s): %v\n", c, err)
		}
	}
}
