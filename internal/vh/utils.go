package vh

import (
	"fmt"
	"path/filepath"
	"time"
)

func GenerateObjectName(billingPn string, origName string) string {
	t := time.Now()
	formatted := fmt.Sprintf(
		"%d-%02d-%02dT%02d:%02d:%02d:%04d",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond(),
	)

	return fmt.Sprintf(
		"%s/%s%s",
		billingPn,
		formatted,
		filepath.Ext(origName),
	)
}
