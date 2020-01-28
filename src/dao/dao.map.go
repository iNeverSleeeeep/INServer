package dao

import (
	"INServer/src/common/dbobj"
	"INServer/src/common/logger"
	"INServer/src/proto/db"
)

func AllStaticMapQuery(DB *dbobj.DBObject) []*db.DBStaticMap {
	rows, err := DB.DB().Query("select ZoneID, MapID, UUID, SerializedData from staticmaps")
	if err != nil {
		logger.Fatal(err)
	}
	staticMaps := []*db.DBStaticMap{}
	for rows.Next() {
		staticMap := &db.DBStaticMap{}
		rows.Scan(&staticMap.ZoneID, &staticMap.MapID, &staticMap.UUID, &staticMap.SerializedData)
	}
	return staticMaps
}

func StaticMapInsert(DB *dbobj.DBObject, staticMap *db.DBStaticMap) error {
	_, err := DB.DB().Exec("insert INTO staticmaps(ZoneID,MapID,UUID,SerializedData) values(?,?,?,?)", staticMap.ZoneID, staticMap.MapID, staticMap.UUID, staticMap.SerializedData)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func StaticMapUpdate(DB *dbobj.DBObject, staticMap *db.DBStaticMap) error {
	_, err := DB.DB().Exec("UPDATE staticmaps set SerializedData=? where UUID=?", staticMap.SerializedData, staticMap.UUID)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func BulkStaticMapUpdate(DB *dbobj.DBObject, staticMaps []*db.DBStaticMap) error {
	tx, err := DB.DB().Begin()
	if err != nil {
		logger.Error(err)
		return err
	}
	stmt, err := tx.Prepare(`UPDATE staticmaps set SerializedData=? where UUID=?`)
	for _, staticMap := range staticMaps {
		_, err := stmt.Exec(staticMap.SerializedData, staticMap.UUID)
		if err != nil {
			logger.Error(err)
			return err
		}
	}
	tx.Commit()
	return nil
}

func DynamicMapQuery(DB *dbobj.DBObject, uuid string) (*db.DBDynamicMap, error) {
	dynamicMap := new(db.DBDynamicMap)
	row := DB.DB().QueryRow("select * from dynamicmaps where UUID=?", uuid)
	if err := row.Scan(&dynamicMap.UUID, &dynamicMap.SerializedData); err != nil {
		logger.Debug(err)
		return nil, err
	}
	return dynamicMap, nil
}

func DynamicMapInsert(DB *dbobj.DBObject, dynamicMap *db.DBDynamicMap) error {
	_, err := DB.DB().Exec("insert INTO dynamicmaps(UUID,SerializedData) values(?,?)", dynamicMap.UUID, dynamicMap.SerializedData)
	if err != nil {
		logger.Debug(err)
		return err
	}
	return nil
}

func DynamicMapUpdate(DB *dbobj.DBObject, dynamicMap *db.DBDynamicMap) error {
	_, err := DB.DB().Exec("UPDATE dynamicmaps set SerializedData=? where UUID=?", dynamicMap.SerializedData, dynamicMap.UUID)
	if err != nil {
		logger.Debug(err)
		return err
	}
	return nil
}
