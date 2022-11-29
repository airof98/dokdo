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

type LineString []Point

func (ls LineString) GormDataType() string {
	return fmt.Sprintf("geometry(LineString, %d)", srid)
}

func (ls *LineString) Scan(val interface{}) (err error) {
	data, err := hex.DecodeString(val.(string))
	if err != nil {
		return errors.InternalErr(err, "")
	}

	var els ewkb.LineString
	err = els.Scan(data)
	if err != nil {
		return errors.InternalErr(err, "")
	}

	for _, coord := range els.Coords() {
		*ls = append(*ls, []float64{coord.X(), coord.Y()})
	}

	return nil
}

func (ls LineString) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	expr := clause.Expr{SQL: "ST_GeomFromText(?)"}
	text := ""
	for i, p := range ls {
		if i > 0 {
			text += ","
		}
		text += fmt.Sprintf("%3.7f %3.7f", p[0], p[1])
	}

	expr.Vars = []interface{}{fmt.Sprintf("LINESTRING(%s)", text)}
	return expr
}
