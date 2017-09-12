package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

var createTableStatements = []string{
	`CREATE DATABASE IF NOT EXISTS horaclifix DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci';`,
	`USE horaclifix;`,
	`CREATE TABLE IF NOT EXISTS qos_report (
		id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
		start TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		end TIMESTAMP NULL DEFAULT NULL,
		duration INT NOT NULL DEFAULT '0',
		sbcName VARCHAR(127) NOT NULL DEFAULT '',
		mediaType TINYINT UNSIGNED NOT NULL DEFAULT '0',
		incomingCallid VARCHAR(255) NOT NULL DEFAULT '',
		incomingRealm VARCHAR(127) NOT NULL DEFAULT '',
		incomingCallerSrcIP VARCHAR(60) NOT NULL DEFAULT '',
		incomingCallerDstIP VARCHAR(60) NOT NULL DEFAULT '',
		incomingCallerSrcPort VARCHAR(10) NOT NULL DEFAULT '',
		incomingCallerDstPort VARCHAR(10) NOT NULL DEFAULT '',
		incomingCalleeSrcIP VARCHAR(60) NOT NULL DEFAULT '',
		incomingCalleeDstIP VARCHAR(60) NOT NULL DEFAULT '',
		incomingCalleeSrcPort VARCHAR(10) NOT NULL DEFAULT '',
		incomingCalleeDstPort VARCHAR(10) NOT NULL DEFAULT '',
		incomingRtpPackets BIGINT UNSIGNED NOT NULL DEFAULT '0',
		incomingRtpLostPackets INT UNSIGNED NOT NULL DEFAULT '0',
		incomingRtpAvgJitter INT UNSIGNED NOT NULL DEFAULT '0',
		incomingRtpMaxJitter INT UNSIGNED NOT NULL DEFAULT '0',
		incomingRtcpPackets BIGINT UNSIGNED NOT NULL DEFAULT '0',
		incomingRtcpLostPackets INT UNSIGNED NOT NULL DEFAULT '0',
		incomingRtcpAvgJitter INT UNSIGNED NOT NULL DEFAULT '0',
		incomingRtcpMaxJitter INT UNSIGNED NOT NULL DEFAULT '0',
		incomingRtcpAvgLat INT UNSIGNED NOT NULL DEFAULT '0',
		incomingRtcpMaxLat INT UNSIGNED NOT NULL DEFAULT '0',
		incomingRval FLOAT(4,2) NOT NULL DEFAULT '0.0',
		incomingMos FLOAT(4,2) NOT NULL DEFAULT '0.0',
		outgoingCallid VARCHAR(255) NOT NULL DEFAULT '',
		outgoingRealm VARCHAR(127) NOT NULL DEFAULT '',
		outgoingCallerSrcIP VARCHAR(60) NOT NULL DEFAULT '',
		outgoingCallerDstIP VARCHAR(60) NOT NULL DEFAULT '',
		outgoingCallerSrcPort VARCHAR(10) NOT NULL DEFAULT '',
		outgoingCallerDstPort VARCHAR(10) NOT NULL DEFAULT '',
		outgoingCalleeSrcIP VARCHAR(60) NOT NULL DEFAULT '',
		outgoingCalleeDstIP VARCHAR(60) NOT NULL DEFAULT '',
		outgoingCalleeSrcPort VARCHAR(10) NOT NULL DEFAULT '',
		outgoingCalleeDstPort VARCHAR(10) NOT NULL DEFAULT '',
		outgoingRtpPackets BIGINT UNSIGNED NOT NULL DEFAULT '0',
		outgoingRtpLostPackets INT UNSIGNED NOT NULL DEFAULT '0',
		outgoingRtpAvgJitter INT UNSIGNED NOT NULL DEFAULT '0',
		outgoingRtpMaxJitter INT UNSIGNED NOT NULL DEFAULT '0',
		outgoingRtcpPackets BIGINT UNSIGNED NOT NULL DEFAULT '0',
		outgoingRtcpLostPackets INT UNSIGNED NOT NULL DEFAULT '0',
		outgoingRtcpAvgJitter INT UNSIGNED NOT NULL DEFAULT '0',
		outgoingRtcpMaxJitter INT UNSIGNED NOT NULL DEFAULT '0',
		outgoingRtcpAvgLat INT UNSIGNED NOT NULL DEFAULT '0',
		outgoingRtcpMaxLat INT UNSIGNED NOT NULL DEFAULT '0',
		outgoingRval FLOAT(4,2) NOT NULL DEFAULT '0.0',
		outgoingMos FLOAT(4,2) NOT NULL DEFAULT '0.0',
		PRIMARY KEY (id),
		KEY start (start),
		KEY end (end),
		KEY sbcName (sbcName),
		KEY incomingCallid (incomingCallid),
		KEY incomingCallerSrcIP (incomingCallerSrcIP),
		KEY incomingCallerDstIP (incomingCallerDstIP),
		KEY incomingCalleeSrcIP (incomingCalleeSrcIP),
		KEY incomingCalleeDstIP (incomingCalleeDstIP),
		KEY outgoingCallid (outgoingCallid),
		KEY outgoingCallerSrcIP (outgoingCallerSrcIP),
		KEY outgoingCallerDstIP (outgoingCallerDstIP),
		KEY outgoingCalleeSrcIP (outgoingCalleeSrcIP),
		KEY outgoingCalleeDstIP (outgoingCalleeDstIP)
	)`,
}

const insertStatement = `
INSERT INTO qos_report (
  start,
  end,
  duration,
  sbcName,
  mediaType,
  incomingCallid,
  incomingRealm,
  incomingCallerSrcIP,
  incomingCallerDstIP,
  incomingCallerSrcPort,
  incomingCallerDstPort,
  incomingCalleeSrcIP,
  incomingCalleeDstIP,
  incomingCalleeSrcPort,
  incomingCalleeDstPort,
  incomingRtpPackets,
  incomingRtpLostPackets,
  incomingRtpAvgJitter,
  incomingRtpMaxJitter,
  incomingRtcpPackets,
  incomingRtcpLostPackets,
  incomingRtcpAvgJitter,
  incomingRtcpMaxJitter,
  incomingRtcpAvgLat,
  incomingRtcpMaxLat,
  incomingRval,
  incomingMos,
  outgoingCallid,
  outgoingRealm,
  outgoingCallerSrcIP,
  outgoingCallerDstIP,
  outgoingCallerSrcPort,
  outgoingCallerDstPort,
  outgoingCalleeSrcIP,
  outgoingCalleeDstIP,
  outgoingCalleeSrcPort,
  outgoingCalleeDstPort,
  outgoingRtpPackets ,
  outgoingRtpLostPackets,
  outgoingRtpAvgJitter,
  outgoingRtpMaxJitter,
  outgoingRtcpPackets,
  outgoingRtcpLostPackets,
  outgoingRtcpAvgJitter,
  outgoingRtcpMaxJitter,
  outgoingRtcpAvgLat,
  outgoingRtcpMaxLat,
  outgoingRval,
  outgoingMos
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

// SendMySQL inserts a given qos value into the qos_report table
func (conn *Connections) SendMySQL(i *IPFIX, s string) {
	start := time.Unix(int64(i.Data.QOS.BeginTimeSec), 0).Local()
	end := time.Unix(int64(i.Data.QOS.EndTimeSec), 0).Local()
	duration := int(i.Data.QOS.EndTimeSec - i.Data.QOS.BeginTimeSec)

	if i.Data.QOS.BeginTimeSec == 0 {
		start = time.Now().Local()
		end = start
		duration = 0
	} else if i.Data.QOS.EndTimeSec == 0 {
		end = time.Now().Local()
		duration = 0
	}

	_, err := execAffectingOneRow(
		conn.MySQL.insert,
		start,
		end,
		duration,
		*name,
		i.Data.QOS.Type,
		string(i.Data.QOS.IncCallID),
		string(i.Data.QOS.IncRealm),
		stringIPv4(i.Data.QOS.CallerIncSrcIP),
		stringIPv4(i.Data.QOS.CallerIncDstIP),
		i.Data.QOS.CallerIncSrcPort,
		i.Data.QOS.CallerIncDstPort,
		stringIPv4(i.Data.QOS.CalleeIncSrcIP),
		stringIPv4(i.Data.QOS.CalleeIncDstIP),
		i.Data.QOS.CalleeIncSrcPort,
		i.Data.QOS.CalleeIncDstPort,
		i.Data.QOS.IncRtpPackets,
		i.Data.QOS.IncRtpLostPackets,
		i.Data.QOS.IncRtpAvgJitter,
		i.Data.QOS.IncRtpMaxJitter,
		i.Data.QOS.IncRtcpPackets,
		i.Data.QOS.IncRtcpLostPackets,
		i.Data.QOS.IncRtcpAvgJitter,
		i.Data.QOS.IncRtcpMaxJitter,
		i.Data.QOS.IncRtcpAvgLat,
		i.Data.QOS.IncRtcpMaxLat,
		float32(i.Data.QOS.IncrVal)/100,
		float32(i.Data.QOS.IncMos)/100,
		string(i.Data.QOS.OutCallID),
		string(i.Data.QOS.OutRealm),
		stringIPv4(i.Data.QOS.CallerOutSrcIP),
		stringIPv4(i.Data.QOS.CallerOutDstIP),
		i.Data.QOS.CallerOutSrcPort,
		i.Data.QOS.CallerOutDstPort,
		stringIPv4(i.Data.QOS.CalleeOutSrcIP),
		stringIPv4(i.Data.QOS.CalleeOutDstIP),
		i.Data.QOS.CalleeOutSrcPort,
		i.Data.QOS.CalleeOutDstPort,
		i.Data.QOS.OutRtpPackets,
		i.Data.QOS.OutRtpLostPackets,
		i.Data.QOS.OutRtpAvgJitter,
		i.Data.QOS.OutRtpMaxJitter,
		i.Data.QOS.OutRtcpPackets,
		i.Data.QOS.OutRtcpLostPackets,
		i.Data.QOS.OutRtcpAvgJitter,
		i.Data.QOS.OutRtcpMaxJitter,
		i.Data.QOS.OutRtcpAvgLat,
		i.Data.QOS.OutRtcpMaxLat,
		float32(i.Data.QOS.OutrVal)/100,
		float32(i.Data.QOS.OutMos)/100)

	if err != nil {
		log.Printf("insert error: %s", err)
	}
}

// newMySQLDB creates a new horaclifix database backed by a given MySQL server.
func newMySQLDB() (*mysqlDB, error) {

	// Check if database and table exists. If not, create them.
	if err := ensureTableExists(); err != nil {
		return nil, err
	}

	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/horaclifix", *muser, *mpass, *maddr))
	if err != nil {
		return nil, fmt.Errorf("mysql: could not get a connection: %v", err)
	}
	// Check the connection.
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("mysql: could not establish a connection: %v", err)
	}

	db := &mysqlDB{
		conn: conn,
	}

	// Prepared sql statements.
	if db.insert, err = conn.Prepare(insertStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare insert: %v", err)
	}

	return db, nil
}

// ensureTableExists checks the table exists. If not, it creates it.
func ensureTableExists() error {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/", *muser, *mpass, *maddr))
	if err != nil {
		return fmt.Errorf("mysql: could not get a connection: %v", err)
	}
	defer conn.Close()

	// Check the connection.
	if conn.Ping() == driver.ErrBadConn {
		return fmt.Errorf("mysql: could not connect to the database. " +
			"could be bad address, or this address is not whitelisted for access.")
	}

	if _, err := conn.Exec("USE horaclifix"); err != nil {
		// MySQL error 1049 is "database does not exist"
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1049 {
			return createTable(conn)
		}
	}

	if _, err := conn.Exec("DESCRIBE qos_report"); err != nil {
		// MySQL error 1146 is "table does not exist"
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1146 {
			return createTable(conn)
		}
		// Unknown error.
		return fmt.Errorf("mysql: could not connect to the database: %v", err)
	}
	return nil
}

func createTable(conn *sql.DB) error {
	for _, stmt := range createTableStatements {
		_, err := conn.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// execAffectingOneRow executes a given statement, expecting one row to be affected.
func execAffectingOneRow(stmt *sql.Stmt, args ...interface{}) (sql.Result, error) {
	r, err := stmt.Exec(args...)
	if err != nil {
		return r, fmt.Errorf("mysql: could not execute statement: %v", err)
	}
	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return r, fmt.Errorf("mysql: could not get rows affected: %v", err)
	} else if rowsAffected != 1 {
		return r, fmt.Errorf("mysql: expected 1 row affected, got %d", rowsAffected)
	}
	return r, nil
}
