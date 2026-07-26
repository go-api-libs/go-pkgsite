package gopkgsite

import "fmt"

func (e *PackagesBadRequestResponse) Error() string {
	return fmt.Sprintf("%d: %s - fixes: %v", e.Code, e.Message, e.Fixes)
}
