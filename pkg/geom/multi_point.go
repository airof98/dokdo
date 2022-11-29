package geom

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/bnkamalesh/errors"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MultiPoint []Point

func (mp MultiPoint) GormDataType() string {
	return fmt.Sprintf("geometry(MultiPoint, %d)", srid)
}

func (mp *MultiPoint) Scan(val interface{}) (err error) {
	data, err := hex.DecodeString(val.(string))
	if err != nil {
		return errors.InternalErr(err, "")
	}
	var emp ewkb.MultiPoint
	err = emp.Scan(data)
	if err != nil {
		return errors.InternalErr(err, "")
	}

	for _, coord := range emp.Coords() {
		*mp = append(*mp, []float64{coord.X(), coord.Y()})
	}
	return nil
}

func (mp MultiPoint) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	expr := clause.Expr{SQL: "ST_GeomFromText(?)"}
	text := ""
	for i, p := range mp {
		if i > 0 {
			text += ","
		}
		text += fmt.Sprintf("%3.7f %3.7f", p[0], p[1])
	}

	expr.Vars = []interface{}{fmt.Sprintf("MULTIPOINT(%s)", text)}
	return expr
}
