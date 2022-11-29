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

type MultiLineString []LineString

func (mls MultiLineString) GormDataType() string {
	return fmt.Sprintf("geometry(MultiLineString, %d)", srid)
}

func (mls *MultiLineString) Scan(val interface{}) (err error) {
	data, err := hex.DecodeString(val.(string))
	if err != nil {
		return errors.InternalErr(err, "")
	}

	var emls ewkb.MultiLineString
	err = emls.Scan(data)
	if err != nil {
		return errors.InternalErr(err, "")
	}

	for i := 0; i < emls.NumLineStrings(); i++ {
		l := emls.LineString(i)
		var ls LineString
		for _, coord := range l.Coords() {
			ls = append(ls, []float64{coord.X(), coord.Y()})
		}
		*mls = append(*mls, ls)
	}

	return nil
}

func (mls MultiLineString) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	expr := clause.Expr{SQL: "ST_GeomFromText(?)"}
	text := ""
	for i, ls := range mls {
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

	expr.Vars = []interface{}{fmt.Sprintf("MULTILINESTRING(%s)", text)}
	return expr
}
