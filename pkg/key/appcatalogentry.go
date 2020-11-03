package key

import (
	"fmt"
)

func AppCatalogEntryManagedBy(projectName string) string {
	return fmt.Sprintf("%s-unique", projectName)
}
