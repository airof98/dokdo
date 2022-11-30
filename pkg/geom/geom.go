package geom

import (
	"fmt"
	"reflect"
)

var srid = 4326

func Float64(obj interface{}) float64 {
	var f float64
	switch num := obj.(type) {
	case float64:
		f = float64(num)
	case int:
		f = float64(num)
	case float32:
		f = float64(num)
	case int64:
		f = float64(num)
	default:
		panic(fmt.Sprintf("Error: Cannot parse object: '%v' type: '%v' to float64!", obj, reflect.TypeOf(obj)))
	}
	return f
}

func SetStrid(newSrid int) {
	srid = newSrid
}
