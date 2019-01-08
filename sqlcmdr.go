package sqlcmdr

import (
	"database/sql"
	"strconv"
	"reflect"
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
	RowID bool
	Joins []JoinTable
}

type JoinTable struct{
	Type string
	LTablename string
	RTablename string
	Leftcol string
	Rightcol string
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
	checkErr( err, "SELECT PREPARE" )
	
	_, err = statement.Exec(icmd.values...)
	
	checkErr( err, "SELECT EXEC" )
}

func Select( conn *sql.DB, scmd SelectCmd ) [][]interface{}{
	// returns array of array 
	
	cols := ""
	if( scmd.RowID ){
		cols = cols+"rowid, "
	}
	sqlcmd := "SELECT "+ cols + scmd.Columns + " FROM " + scmd.Tablename
	sqlcmd += " " + joinstr( scmd.Joins )
	var selectall bool = len( scmd.Keycol ) == 0
	if( !selectall ){
		sqlcmd += " WHERE " + scmd.Keycol + scmd.Comparison + " ? "
	}
	
	var rows *sql.Rows
	var err error
	if( !selectall ){
		rows, err = conn.Query( sqlcmd, scmd.Keyval )
	}else{
		rows, err = conn.Query( sqlcmd )
	}
	fmt.Println( sqlcmd )
	checkErr( err, "SELECT" )
    
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

func ResultString( retvals [][]interface{} ) string{
	var result string
	for _, element := range( retvals ){
		for _, e2 := range( element ){
			if( e2 != nil ){
				switch e2.(type){
					case []uint8 : 
						result = result + " " + string(e2.([]uint8))
					case float64 : 
						result = result + " " + strconv.FormatFloat(e2.(float64), 'f', 6, 64)
					case int64 :
						result = result + " " + strconv.FormatInt( e2.(int64), 10 )
					default :
						result = result + " " + "ERR: unhandled type "+ reflect.TypeOf(e2).String()
					}
			}
		}
		result = result + "\n"
	}
	return result
}

/////////////////// privates ////////////////////////////////////

func checkErr(err error, origin string) {
	if err != nil {
		panic("INVALID " + origin + ": " + err.Error() )
	}
}

func joinstr( joins []JoinTable ) string{
	str := ""
	for _, join := range joins{
		str += join.Type + " JOIN " + join.RTablename + 
				" ON " + join.LTablename + "." + join.Leftcol + 
				" = " + join.RTablename + "." + join.Rightcol + " "
	}
	return str
}
