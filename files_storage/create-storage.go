package storage

var DB DataBase

func CreateAndMigrateDB() {
	DB = NewDB()
	DB.MigrateDBs()
}
