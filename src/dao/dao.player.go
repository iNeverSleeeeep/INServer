package dao

import (
	"INServer/src/common/dbobj"
	"INServer/src/common/logger"
	"INServer/src/proto/db"
)

func AccountQuery(DB *dbobj.DBObject, username string) (*db.DBAccount, error) {
	account := new(db.DBAccount)
	row := DB.DB().QueryRow("select * from accounts where Name=?", username)
	if err := row.Scan(&account.Name, &account.PasswordHash, &account.PlayerUUID); err != nil {
		logger.Debug(err)
		return nil, err
	}
	return account, nil
}

func AccountInsert(DB *dbobj.DBObject, account *db.DBAccount) error {
	_, err := DB.DB().Exec("insert INTO accounts(Name,PasswordHash,PlayerUUID) values(?,?,?)", account.Name, account.PasswordHash, account.PlayerUUID)
	if err != nil {
		logger.Debug(err)
		return err
	}
	return nil
}

func AccountUpdate(DB *dbobj.DBObject, account *db.DBAccount) error {
	_, err := DB.DB().Exec("UPDATE accounts set PasswordHash=? PlayerUUID=? where Name=?", account.PasswordHash, account.PlayerUUID, account.Name)
	if err != nil {
		logger.Debug(err)
		return err
	}
	return nil
}

func PlayerQuery(DB *dbobj.DBObject, uuid string) (*db.DBPlayer, error) {
	player := new(db.DBPlayer)
	row := DB.DB().QueryRow("select * from players where UUID=?", uuid)
	if err := row.Scan(&player.UUID, &player.SerializedData); err != nil {
		logger.Debug(err)
		return nil, err
	}
	return player, nil
}

func PlayerInsert(DB *dbobj.DBObject, player *db.DBPlayer) error {
	_, err := DB.DB().Exec("insert INTO accounts(UUID,SerializedData) values(?,?)", player.UUID, player.SerializedData)
	if err != nil {
		logger.Debug(err)
		return err
	}
	return nil
}

func PlayerUpdate(DB *dbobj.DBObject, player *db.DBPlayer) error {
	_, err := DB.DB().Exec("UPDATE accounts set SerializedData=? where UUID=?", player.SerializedData, player.UUID)
	if err != nil {
		logger.Debug(err)
		return err
	}
	return nil
}
