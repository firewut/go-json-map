package gjm

import (
	"reflect"
)

func isKind(what interface{}, kind reflect.Kind) bool {
	return reflect.ValueOf(what).Kind() == kind
}
