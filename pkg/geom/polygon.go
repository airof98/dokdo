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

type Polygon [][]Point

func (p Polygon) GormDataType() string {
	return fmt.Sprintf("geometry(Polygon, %d)", srid)
}

func (p *Polygon) Scan(val interface{}) (err error) {
	data, err := hex.DecodeString(val.(string))
	if err != nil {
		return errors.InternalErr(err, "")
	}

	var ep ewkb.Polygon
	err = ep.Scan(data)
	if err != nil {
		return errors.InternalErr(err, "")
	}

	for i := 0; i < ep.NumLinearRings(); i++ {
		l := ep.LinearRing(i)
		var lr []Point
		for _, coord := range l.Coords() {
			lr = append(lr, []float64{coord.X(), coord.Y()})
		}
		*p = append(*p, lr)
	}

	return nil
}

func (p Polygon) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	expr := clause.Expr{SQL: "ST_GeomFromText(?)"}
	text := ""
	for i, ls := range p {
		if i > 0 {
			text += ","
		}
		text += "("
		for j, p := range ls {
			if j > 0 {
				text += ","
			}
			text += fmt.Sprintf("%3.7f %3.7f", p[0], p[1])
		}
		text += ")"
	}

	expr.Vars = []interface{}{fmt.Sprintf("POLYGON(%s)", text)}
	return expr
}
