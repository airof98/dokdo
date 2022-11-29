package geom

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn = "postgres://localhost:5432/dokdo?sslmode=disable"

var db *gorm.DB

type DokdoStop struct {
	gorm.Model
	Name   string
	LngLat Point
}

type DokdoMultiStop struct {
	gorm.Model
	Name    string
	LngLats MultiPoint
}

var dokdoEast = Point{131.8675129, 37.2390974}
var dokdoWest = Point{131.8644791, 37.2417323}
var ulleungdo = Point{130.861833523, 37.498403198}

func init() {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
}

func TestPoint(t *testing.T) {
	db.Migrator().DropTable(DokdoStop{}, DokdoMultiStop{})
	err := db.Migrator().AutoMigrate(DokdoStop{}, DokdoMultiStop{})
	assert.NoError(t, err)
	err = db.Create(&DokdoStop{Name: "독도는 대한믹국 영토", LngLat: dokdoEast}).Error
	assert.NoError(t, err)
	var stop DokdoStop
	err = db.First(&stop).Error
	assert.NoError(t, err)
	assert.InDelta(t, stop.LngLat[0], dokdoEast[0], 0.0000001)
	assert.InDelta(t, stop.LngLat[1], dokdoEast[1], 0.0000001)

	var lnglats []Point
	lnglats = append(lnglats, dokdoEast)
	lnglats = append(lnglats, dokdoWest)
	err = db.Create(&DokdoMultiStop{Name: "독도는 대한믹국 영토", LngLats: lnglats}).Error
	assert.NoError(t, err)
	var multiStop DokdoMultiStop
	err = db.First(&multiStop).Error
	assert.NoError(t, err)
	assert.InDelta(t, multiStop.LngLats[1][0], dokdoWest[0], 0.0000001)
	assert.InDelta(t, multiStop.LngLats[1][1], dokdoWest[1], 0.0000001)
}

type DokdoLine struct {
	Name string
	Line LineString
}

type DokdoMultiLine struct {
	Name      string
	MultiLine MultiLineString
}

func TestLine(t *testing.T) {
	db.Migrator().DropTable(DokdoLine{}, DokdoMultiLine{})
	err := db.Migrator().AutoMigrate(DokdoLine{}, DokdoMultiLine{})
	assert.NoError(t, err)
	var line []Point
	line = append(line, dokdoEast)
	line = append(line, dokdoWest)
	err = db.Create(&DokdoLine{Name: "독도는 대한믹국 영토", Line: line}).Error
	assert.NoError(t, err)
	var dokdoLine DokdoLine
	err = db.First(&dokdoLine).Error
	assert.NoError(t, err)
	assert.InDelta(t, dokdoLine.Line[1][0], dokdoWest[0], 0.0000001)
	assert.InDelta(t, dokdoLine.Line[1][1], dokdoWest[1], 0.0000001)

	var multiLine []LineString
	multiLine = append(multiLine, line)
	line = append(line, ulleungdo)
	multiLine = append(multiLine, line)
	err = db.Create(&DokdoMultiLine{Name: "독도는 대한믹국 영토", MultiLine: multiLine}).Error
	assert.NoError(t, err)
	var dokdoMultiLine DokdoMultiLine
	err = db.First(&dokdoMultiLine).Error
	assert.NoError(t, err)
	assert.InDelta(t, dokdoMultiLine.MultiLine[1][2][0], ulleungdo[0], 0.0000001)
	assert.InDelta(t, dokdoMultiLine.MultiLine[1][2][1], ulleungdo[1], 0.0000001)
}

type DokdoPolygon struct {
	Name    string
	Polygon Polygon
}

type DokdoMultiPolygon struct {
	Name         string
	MultiPolygon MultiPolygon
}

func TestPolygon(t *testing.T) {
	db.Migrator().DropTable(DokdoPolygon{}, DokdoMultiPolygon{})
	err := db.Migrator().AutoMigrate(DokdoPolygon{}, DokdoMultiPolygon{})
	assert.NoError(t, err)
	var polygon [][]Point
	var lring []Point
	lring = append(lring, dokdoEast)
	lring = append(lring, dokdoWest)
	lring = append(lring, ulleungdo)
	lring = append(lring, dokdoEast)
	polygon = append(polygon, lring)
	err = db.Create(&DokdoPolygon{Name: "독도는 대한믹국 영토", Polygon: polygon}).Error
	assert.NoError(t, err)

	var dokdoPolygon DokdoPolygon
	err = db.First(&dokdoPolygon).Error
	assert.NoError(t, err)
	assert.InDelta(t, dokdoPolygon.Polygon[0][2][0], ulleungdo[0], 0.0000001)
	assert.InDelta(t, dokdoPolygon.Polygon[0][2][1], ulleungdo[1], 0.0000001)

	var multiPolygon []Polygon
	multiPolygon = append(multiPolygon, polygon)
	multiPolygon = append(multiPolygon, polygon)
	err = db.Create(&DokdoMultiPolygon{Name: "독도는 대한믹국 영토", MultiPolygon: multiPolygon}).Error
	assert.NoError(t, err)

	var dokdoMultiPolygon DokdoMultiPolygon
	err = db.First(&dokdoMultiPolygon).Error
	assert.NoError(t, err)
	a := dokdoMultiPolygon.MultiPolygon[0]
	assert.InDelta(t, a[0][2][0], ulleungdo[0], 0.0000001)
	assert.InDelta(t, a[0][2][1], ulleungdo[1], 0.0000001)
}
