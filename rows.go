// Go driver for MySQL X Protocol
// Based heavily on Go MySQL Driver - A MySQL-Driver for Go's database/sql package
//
// Copyright 2012 The Go-MySQL-Driver Authors. All rights reserved.
// Copyright 2016 Simon J Mudd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package mysql

import (
	"database/sql/driver"
	"fmt"
	"io"
	"log"

	"github.com/golang/protobuf/proto"

	"github.com/sjmudd/go-mysqlx-driver/Mysqlx"
	"github.com/sjmudd/go-mysqlx-driver/Mysqlx_Resultset"

	"github.com/sjmudd/go-mysqlx-driver/debug"
)

type mysqlXRows struct {
	columns [](*Mysqlx_Resultset.ColumnMetaData) // holds column metadata (if present) for a row
	mc      *mysqlXConn
	state   queryState
	err     error // provides the error received from a query (if present)
}

// readMsgIfNecessary reads in a message only if we don't have one already
func (rows *mysqlXRows) readMsgIfNecessary() error {
	// safety checks (which maybe can removed later
	if rows == nil {
		return fmt.Errorf("mysqlXRows.readMsgIfNecessary: rows == nil")
	}
	if rows.mc == nil {
		return fmt.Errorf("mysqlXRows.readMsgIfNecessary: rows.mc == nil")
	}
	// if we already have a protobuf message then no need to read a new one
	if rows.mc.pb != nil {
		debug.Msg("mysqlXRows.readMsgIfNecessary: DO NOT read new message (pb != nil)")
		return nil
	}

	debug.Msg("mysqlXRows.readMsgIfNecessary: read NEW MESSAGE")

	var err error
	rows.mc.pb, err = rows.mc.readMsg()
	if err != nil {
		err = fmt.Errorf("mysqlXRows.readMsgIfNecessary rows.mc.readMsg failed: %v", err)

		rows.err = err
		rows.state = queryStateError
	}
	return err
}

// Columns returns the column meta data of a row and may need to
// read in some of the metadata messages from the network.
func (rows *mysqlXRows) Columns() []string {
	rows.collectColumnMetaData()

	columns := make([]string, len(rows.columns))
	for i := range rows.columns {
		// FIXME: handle: 	if rows.mc.cfg.columnsWithAlias { ....
		columns[i] = string(rows.columns[i].GetName())
	}
	if len(columns) == 0 {
		debug.Msg("mysqlXRows.Columns: return empty []string with %d entries (probably due to SQL error)", len(columns))
	} else {
		debug.Msg("mysqlXRows.Columns: return %+v", columns)
	}
	return columns
}

// we have finished with the iterator
// - given Close can be called at any time we may have pending
//   messages in the queue which need skipping so we really need
//   to keep the state of where we are.
func (rows *mysqlXRows) Close() error {
	debug.Msg("mysqlXRows.Close: entry")

	// safety checks
	if rows == nil {
		debug.Msg("mysqlXRows.Close: rows == nil, ignoring")
		return nil // to avoid breakage. Fix the calling code later
	}
	if rows.mc == nil {
		return nil // no connection information
	}
	if rows.mc.netConn == nil {
		return ErrInvalidConn
	}

	// We may have "query packets" which have not yet been
	// processed. If so just let them through but ignore them.
	for rows.state != queryStateDone && rows.state != queryStateError {
		if err := rows.readMsgIfNecessary(); err != nil {
			debug.Msg("mysqlXRows.Close: got an error trying to read rows: %v", err)
			break
		}

		// Finish if we get an error or if the mssage type is EXECUTE_OK or ERROR
		switch Mysqlx.ServerMessages_Type(rows.mc.pb.msgType) {
		case Mysqlx.ServerMessages_ERROR:
			rows.mc.processErrorMsg()
			rows.state = queryStateError
		case Mysqlx.ServerMessages_SQL_STMT_EXECUTE_OK:
			rows.state = queryStateDone
		case Mysqlx.ServerMessages_NOTICE:
			rows.mc.processNotice("mysqlXRows.Close")
		default:
			// do nothing
		}
		rows.mc.pb = nil
	}

	// clean up
	rows.columns = nil
	rows.mc.pb = nil
	rows.mc = nil
	rows.state = queryStateNotStarted

	debug.Msg("mysqlXRows.Close: exit")
	return nil
}

// add the column information to the row
func (rows *mysqlXRows) addColumn() error {
	if rows == nil {
		return fmt.Errorf("mysqlXrows.addColumn: rows == nil")
	}

	column := new(Mysqlx_Resultset.ColumnMetaData)
	if err := proto.Unmarshal(rows.mc.pb.payload, column); err != nil {
		return fmt.Errorf("error unmarshalling ColumnMetaData: %v", err)
	}

	debug.Msg("mysqlXRows.addColumn: %s", printableColumnMetaData(rows.mc.pb))

	rows.columns = append(rows.columns, column)
	rows.mc.pb = nil

	return nil
}

// process a single row (in rows.mc.pb) and return if there was an error
func processRow(rows *mysqlXRows, dest []driver.Value) error {
	var err error

	myRow := new(Mysqlx_Resultset.Row)
	if err = proto.Unmarshal(rows.mc.pb.payload, myRow); err != nil {
		return fmt.Errorf("error unmarshalling Row: %v", err)
	}
	rows.mc.pb = nil // consume the message

	debug.Msg("processRow: row has %d columns", len(myRow.GetField()))
	// copy over data converting each type to a dest type
	for i := range dest {
		if dest[i], err = convertColumnData(rows.columns[i], myRow.GetField()[i]); err != nil {
			return fmt.Errorf("processRow: failed to convert data for column %d: %v", i, err)
		}
	}

	return nil // no error
}

// Read a row of data from the connection until no more and then return io.EOF to indicate we have finished
func (rows *mysqlXRows) Next(dest []driver.Value) error {
	// safety checks
	if rows == nil {
		log.Fatal("mysqlXRows.Next: rows == nil")
	}
	if rows.mc == nil {
		log.Fatal("mysqlXRows.Next: rows.mc == nil")
	}

	debug.Msg("mysqlXrows.Next: entry state: %q", rows.state.String())

	// Finished? Don't continue
	if rows.state.Finished() {
		debug.Msg("mysqlXrows.Next: rows.state.Finished() is true")
		return io.EOF
	}

	// Have we read the column data yet? If not read it.
	if rows.state == queryStateWaitingColumnMetaData {
		if err := rows.collectColumnMetaData(); err != nil {
			return err
		}
	}

	debug.Msg("mysqlXrows.Next: rows.state: %v, dest has %d elements", rows.state.String(), len(dest))

	// clean this logic up into a smaller more readable loop
	done := false
	for !done {
		debug.Msg("mysqlXrows.Next: loop state: %v", rows.state.String())

		switch rows.state {
		case queryStateWaitingRow:
			{

				// pull in a message if needed
				if err := rows.readMsgIfNecessary(); err != nil {
					log.Fatalf("DEBUG: mysqlXRow.Next: failed to read data if necessary")
				}

				// check if it's a Row message!
				switch Mysqlx.ServerMessages_Type(rows.mc.pb.msgType) {
				case Mysqlx.ServerMessages_RESULTSET_ROW:
					{
						if err := processRow(rows, dest); err != nil {
							return err
						}
						done = true
					}
				case Mysqlx.ServerMessages_NOTICE:
					rows.mc.processNotice("mysqlXRows.Next")
				case Mysqlx.ServerMessages_RESULTSET_FETCH_DONE:
					{
						rows.state = queryStateWaitingExecuteOk
						done = true
						rows.mc.pb = nil
					}
				case Mysqlx.ServerMessages_ERROR:
					{
						// should treat each message
						rows.state = queryStateDone
						done = true
						rows.mc.pb = nil
					}
				default:
					{
						log.Fatalf("mysqlXRowx.Next received unexpected message type: %s", printableMsgTypeIn(Mysqlx.ServerMessages_Type(rows.mc.pb.msgType)))
					}
				}
			}
		case queryStateDone, queryStateWaitingExecuteOk:
			{
				return io.EOF
			}
		default:
			{
				log.Fatalf("mysqlXRows.Next: called in unexpected state: %v", rows.state.String())
				// otherwise assume everything is fine
			}
		}
	}

	return nil
}

// Expectation here is to receive one of
// - RESULTSET_COLUMN_META_DATA (expected)
// - NOTICE (may happen, not expected)
// - RESULTSET_ROW (expected, changes state)
func (rows *mysqlXRows) collectColumnMetaData() error {
	if rows == nil {
		return fmt.Errorf("BUG: mysqlXRows.collectColumnMetaData: rows == nil")
	}
	debug.Msg("mysqlXRows.collectColumnMetaData: entry, rows.state: %q", rows.state.String())

	for !rows.state.Finished() && rows.state != queryStateWaitingRow {
		// debug.Msg("mysqlXRows.collectColumnMetaData: loop")
		if err := rows.readMsgIfNecessary(); err != nil {
			return fmt.Errorf("DEBUG: mysqlXRows.collectColumnMetaData: failed to read data if necessary")
		}

		switch Mysqlx.ServerMessages_Type(rows.mc.pb.msgType) {
		case Mysqlx.ServerMessages_RESULTSET_COLUMN_META_DATA:
			{
				if err := rows.addColumn(); err != nil {
					return fmt.Errorf("DEBUG: mysqlXRows.collectColumnMetaData: failed to addColumn: %v", err)
				}
			}
		case Mysqlx.ServerMessages_RESULTSET_ROW:
			{
				rows.state = queryStateWaitingRow
				debug.Msg("mysqlXRows.collectColumnMetaData: got RESULTSET_ROW: change state to %q", rows.state.String())
			}
		case Mysqlx.ServerMessages_NOTICE:
			{
				// don't really expect a notice but process it
				debug.Msg("mysqlXRows.collectColumnMetaData: got NOTICE: processing it")
				rows.mc.processNotice("mysqlxRows.collectColumnMetaData")
			}
		case Mysqlx.ServerMessages_ERROR:
			{
				debug.Msg("mysqlXRows.collectColumnMetaData: got ERROR: process it and change state to queryStateDone")
				rows.mc.processErrorMsg()
				rows.state = queryStateError
			}
		default:
			{
				e := fmt.Errorf("mysqlXRows.collectColumnMetaData: received unexpected message type: %s",
					printableMsgTypeIn(Mysqlx.ServerMessages_Type(rows.mc.pb.msgType)))
				rows.state = queryStateError
				rows.mc.pb = nil
				return e
			}
		}
	}
	return nil
}
