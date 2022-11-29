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

type MultiPolygon []Polygon

func (mp MultiPolygon) GormDataType() string {
	return fmt.Sprintf("geometry(MultiPolygon, %d)", srid)
}

func (mp *MultiPolygon) Scan(val interface{}) (err error) {
	data, err := hex.DecodeString(val.(string))
	if err != nil {
		return errors.InternalErr(err, "")
	}

	var emp ewkb.MultiPolygon
	err = emp.Scan(data)
	if err != nil {
		return errors.InternalErr(err, "")
	}

	for i := 0; i < emp.NumPolygons(); i++ {
		p := emp.Polygon(i)
		var polygon Polygon
		for i := 0; i < p.NumLinearRings(); i++ {
			l := p.LinearRing(i)
			var lr []Point
			for _, coord := range l.Coords() {
				lr = append(lr, []float64{coord.X(), coord.Y()})
			}
			polygon = append(polygon, lr)
		}
		*mp = append(*mp, polygon)
	}

	return nil
}

func (mp MultiPolygon) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	expr := clause.Expr{SQL: "ST_GeomFromText(?)"}
	text := ""
	for i, p := range mp {
		if i > 0 {
			text += ","
		}
		text += "("
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
		text += ")"
	}

	expr.Vars = []interface{}{fmt.Sprintf("MULTIPOLYGON(%s)", text)}
	return expr
}
