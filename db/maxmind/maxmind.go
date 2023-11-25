package maxmind

import (
	"fmt"
	"net"
	"sync"

	"github.com/oschwald/maxminddb-golang"
)

// Info is a geographic information and other data associated with specific Internet protocol addresses
type Info struct {
	Country string
	Region  string
	City    string
}

// DB is a GeoIP database
type DB interface {
	Lookup(net.IP) (Info, error)
	Close() error
}

// MaxmindDB is a Maxmind GeoIP database
type MaxmindDB struct {
	r *maxminddb.Reader
}

var _ DB = (*MaxmindDB)(nil)

func OpenMaxmindDB(filePath string) (*MaxmindDB, error) {
	r, err := maxminddb.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &MaxmindDB{r}, nil
}

type maxmindDBData struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
	City struct {
		Names struct {
			En string `maxminddb:"en"`
		} `maxminddb:"names"`
	} `maxminddb:"city"`
	Region []struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"subdivisions"`
}

var mdPool = &sync.Pool{
	New: func() interface{} {
		return &maxmindDBData{}
	},
}

// getMaxmindDBData returns a maxmindDBData from the pool.
func getMaxmindDBData() *maxmindDBData {
	return mdPool.Get().(*maxmindDBData)
}

// putMaxmindDBData returns a maxmindDBData to the pool.
// The maxmindDBData is reset before it is put back into circulation.
func putMaxmindDBData(d *maxmindDBData) {
	d.Country.ISOCode = ""
	d.City.Names.En = ""
	for i := range d.Region {
		d.Region[i].ISOCode = ""
	}
	d.Region = d.Region[:0]
	mdPool.Put(d)
}

func (db *MaxmindDB) Lookup(ip net.IP) (info Info, err error) {
	res := getMaxmindDBData()
	defer putMaxmindDBData(res)

	err = preventPanic(func() error {
		return db.r.Lookup(ip, res)
	})
	if err != nil {
		return
	}

	info.Country = res.Country.ISOCode
	info.City = res.City.Names.En
	if len(res.Region) > 0 {
		info.Region = res.Region[0].ISOCode
	}

	return
}

type maxmindDBNetworkData struct {
	ConnectionType string `maxminddb:"connection_type"`
}

func (db *MaxmindDB) LookupNetwork(ip net.IP) (connectionType string, err error) {
	res := new(maxmindDBNetworkData)
	err = preventPanic(func() error {
		_, _, err := db.r.LookupNetwork(ip, res)
		return err
	})
	connectionType = res.ConnectionType
	return connectionType, err
}

func preventPanic(f func() error) (err error) {
	defer func() {
		fErr := recover()
		if fErr != nil {
			err = fmt.Errorf("%v", fErr)
		}
	}()
	return f()
}

func (db *MaxmindDB) Close() error {
	return db.r.Close()
}
