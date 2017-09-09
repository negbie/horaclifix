package main

var createTableStatements = []string{
	`CREATE DATABASE IF NOT EXISTS horaclifix DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci';`,
	`USE horaclifix;`,
	`CREATE TABLE IF NOT EXISTS qos_report (
		id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
		callerIncSrcIP VARCHAR(255) NULL,
		PRIMARY KEY (id,date),
		KEY date (date),
		KEY callid (callid),

	)`,
}

func (conn *Connections) SendMySQL(i *IPFIX, s string) {

}

func createTable(conn *Connections) error {
	for _, stmt := range createTableStatements {
		_, err := conn.MySQL.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}
