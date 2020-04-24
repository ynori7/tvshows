package email

import (
	"fmt"
)

func GetNewReleasesSubjectLine(startDate string, endDate string) string {
	return fmt.Sprintf("Newest premieres from %s through %s", startDate, endDate)
}
