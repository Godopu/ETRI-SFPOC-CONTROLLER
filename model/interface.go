package model

// func newDBHandler(dbtype, path string) (*gorm.DB, error) {
// 	if dbtype == "sqlite" {
// 		return gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
// 	} else {
// 		dsn := "host=localhost user=user password=user_password dbname=godopudb port=5432 sslmode=disable TimeZone=Asia/Seoul"
// 		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	}
// }

type DBHandler interface {
	GetDevices() ([]*Device, int, error)
	AddDevice(device *Device) error
	GetSID(sname string) (string, error)
	GetServiceForDevice(did string) (string, error)
	GetDeviceID(dname string) (*Device, error)
	StatusCheck(did string, new map[string]interface{}) bool
}

var db DBHandler

func GetDBHandler(dbtype, path string) (DBHandler, error) {
	if db == nil {
		db, _ = newSqliteHandler(path)
		return db, nil
	}
	return db, nil
}
