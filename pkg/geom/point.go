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

type Point []float64

func (p Point) GormDataType() string {
	return fmt.Sprintf("geometry(Point, %d)", srid)
}

func (p Point) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	expr := clause.Expr{SQL: "ST_GeomFromText(?)"}
	expr.Vars = []interface{}{fmt.Sprintf("POINT(%3.7f %3.7f)", p[0], p[1])}
	return expr
}

func (p *Point) Scan(val interface{}) (err error) {
	data, err := hex.DecodeString(val.(string))
	if err != nil {
		return errors.InternalErr(err, "")
	}

	var mp ewkb.Point
	err = mp.Scan(data)
	if err != nil {
		return errors.InternalErr(err, "")
	}

	coord := mp.Coords()
	*p = []float64{coord.X(), coord.Y()}
	return nil
}
