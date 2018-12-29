package sqlcmdr

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
   	"fmt"
)

type InsertCmd struct{
	Tablename string
	columns []string
	values []interface{}
}

type SelectCmd struct{
	Tablename string
	Columns string
	Keycol string
	Keyval string
	Comparison string
}

func InitDB() *sql.DB {
	database, _ := sql.Open("sqlite3", "./data.db")
	return database
}
func JustRunIt( command string ){
	conn := InitDB()
	statement, _ := conn.Prepare( command )
								   
	statement.Exec()
	conn.Close();
}

func (icmd *InsertCmd) Add(column string, value interface{} ){
	icmd.columns = append( icmd.columns, column )
	icmd.values = append( icmd.values, value )
}

func Insert( conn *sql.DB, icmd InsertCmd){
	colstring := ""
	valstring := ""
	for _, col :=  range icmd.columns 	{
		if( len(colstring) > 0 )		{
			colstring = colstring + ", "
			valstring = valstring + ", "
		}
		colstring = colstring + col
		valstring = valstring + "?"
	}
	sqlcmd := "INSERT INTO " + icmd.Tablename + " (" + colstring + " ) VALUES (" + valstring + ")"
	statement, err := conn.Prepare(sqlcmd)
	if( err != nil ){
		fmt.Errorf("INVALID INSERT%s\n", err)
	}
	statement.Exec(icmd.values...)
}

func Select( conn *sql.DB, scmd SelectCmd ) [][]interface{}{
	// returns array of array 
	sqlcmd := "SELECT "+ scmd.Columns + " FROM " + scmd.Tablename
	var selectall bool = len( scmd.Keycol ) == 0
	if( !selectall ){
		sqlcmd += " WHERE " + scmd.Keycol + scmd.Comparison + " ? "
	}
	
	var rows *sql.Rows
	var err error
	fmt.Println("selecting " + sqlcmd  )
	if( !selectall ){
		rows, err = conn.Query( sqlcmd, scmd.Keyval )
	}else{
		rows, err = conn.Query( sqlcmd )
	}
	
	if( err != nil ){
		panic( err)
	}
    
    // build an array of interface pointers because this is needed by scan    
	columns, _ := rows.Columns()
    count := len(columns)
    retface := make([]interface{}, count)
	for i, _ := range columns {
		var i4ptr interface{}
		retface[i] = &i4ptr
	}

	// convert rows to values
	var retvals [][]interface{}
	for rows.Next() {
		rows.Scan( retface... )
		
		retrow := make( []interface{},  count )
		for _, value := range retface {
			var retval = *(value.(*interface{}))
			retrow = append( retrow, retval )
		}
		retvals = append( retvals, retrow )
	}
	return retvals
}

