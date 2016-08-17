package dbAccess

import (
	"database/sql"
	"fmt"
	"xmlConfig"
	//Init MS SQL Server dirver, so we can open Database as "mssql"
	//Select the related SQL Connect Drivers from :https://github.com/golang/go/wiki/SQLDrivers
	//_ "github.com/denisenkom/go-mssqldb", use this to import the package from your download place, ect go\bin
	_ "github.com/denisenkom/go-mssqldb" //use this to import you installed package from /pkg/XX/*.a
)

var connectStr string

func getConnectString() string {

	if len(connectStr) < 1 {
		var configMap = xmlConfig.GetConfig()

		userName := xmlConfig.GetElement("userName", configMap)
		password := xmlConfig.GetElement("password", configMap)
		sqlHost := xmlConfig.GetElement("sqlHost", configMap)
		dbName := xmlConfig.GetElement("dbName", configMap)
		port := xmlConfig.GetElement("sqlPort", configMap)

		connectStr = "server=" + sqlHost + ";user id=" + userName + ";password=" + password + ";port=" + port + ";database=" + dbName
		fmt.Println(connectStr)
	}

	return connectStr
}

//Get the encryption Key from database
func GetMacEncryptKey(macID string) []byte {

	var encryptionKey []byte

	dbContext, err := sql.Open("mssql", getConnectString())
	if err != nil {
		fmt.Println(err)
		return encryptionKey
	}
	defer dbContext.Close()

	rows, err := dbContext.Query("select EncryptionKey from EncryptionKeys where MacID = ?", macID)
	if err != nil {
		fmt.Println(err)
		return encryptionKey
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&encryptionKey)
		fmt.Println(encryptionKey)
	}

	return encryptionKey
}
