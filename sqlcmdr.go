package sqlcmdr

import (
	"database/sql"
	"strconv"
	"reflect"
	_ "github.com/mattn/go-sqlite3"
   	"encoding/csv"
   	"os"
   	"strings"
   	"fmt"
)

/////////////////////////// structs ////////////////////////////////////

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
	LSide string  			//e.g. "LeftTable-Column"
	RSide string 			//e.g. "RightTable-Column"
	
}

////////////////////////// public funcs ////////////////////////////////


func (icmd *InsertCmd) Add (column string, value interface{} ){
	icmd.columns = append( icmd.columns, column )
	icmd.values = append( icmd.values, value )
}

func (scmd *SelectCmd) AddJoin ( leftside string, rightside string, jointype string){
	jt := JoinTable{ Type: jointype, LSide: leftside, RSide: rightside }
	scmd.Joins = append( scmd.Joins, jt )
}

func InitDB() *sql.DB {
	// opens connection, needs to be closed with conn.Close later
	database, _ := sql.Open("sqlite3", "./data.db")
	return database
}
func JustRunIt( command string ){
	// just runs command with own connection.  No return, no parameterization
	conn := InitDB()
	defer conn.Close();
	statement, err := conn.Prepare( command )		
	checkErr( err, "JUST RUN" )	   
	statement.Exec()
}

func Insert( conn *sql.DB, icmd InsertCmd){
	// runs insert based on InsertCmd
	
	// convert icmd into string insert command
	colstring := ""
	valstring := ""
	for _, col :=  range icmd.columns 	{
		if( len(colstring) > 0 )	{
			colstring = colstring + ", "
			valstring = valstring + ", "
		}
		colstring = colstring + col
		valstring = valstring + "?"
	}
	sqlcmd := "INSERT INTO " + icmd.Tablename + " (" + colstring + " ) VALUES (" + valstring + ")"
	
	// prepare and run
	statement, err := conn.Prepare(sqlcmd)
	checkErr( err, "SELECT PREPARE" )
	
	_, err = statement.Exec(icmd.values...)
	
	checkErr( err, "SELECT EXEC" )
}

func Select( conn *sql.DB, scmd SelectCmd ) [][]interface{}{
	// runs select based on SelectCmd, returns array of array that can then be processed with resultstring/resultcsv
	
	// covert scmd object into string for query
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
	fmt.Println( sqlcmd )
	// do the query
	var rows *sql.Rows
	var err error
	if( !selectall ){
		rows, err = conn.Query( sqlcmd, scmd.Keyval )
	}else{
		rows, err = conn.Query( sqlcmd )
	}
	
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

func ResultCSV( records [][]interface{}, output string ) string{
	// turns 2d array of interfaces into output csv file
	file, err := os.Create(output)
	checkErr( err, "CSV to " + output )
	w := csv.NewWriter(file)
	defer file.Close()
	for _, record := range records {
		var strrecord []string
		for _, record2 := range record {
			strrecord = append( strrecord, cell2string( record2 ) )			
		}
		err := w.Write(strrecord)
		checkErr( err, "error writing record to csv:" )
	}

	w.Flush()  
	if err := w.Error(); err != nil {
		panic(err)
	}
	return "output to " + output
}

func ResultString( records [][]interface{} ) string{
	// turns 2d array of intefaces into single string
	var result string
	for _, element := range( records ){
		for _, e2 := range( element ){
			result += cell2string( e2 );
		}
		result = result + "\n"
	}
	return result
}

/////////////////// private funcs ////////////////////////////////////

func cell2string( e2 interface{} ) string{
	result := "";
	if( e2 != nil ){
			switch e2.(type){
				case []uint8 : 
					result =  " " + string(e2.([]uint8))
				case float64 : 
					result =  " " + strconv.FormatFloat(e2.(float64), 'f', 6, 64)
				case int64 :
					result =  " " + strconv.FormatInt( e2.(int64), 10 )
				default :
					result =  " " + "ERR: unhandled type "+ reflect.TypeOf(e2).String()
				}
		}
	return result
}

func checkErr(err error, origin string) {
	if err != nil {
		panic("INVALID " + origin + ": " + err.Error() )
	}
}

func joinstr( joins []JoinTable ) string{
	str := ""
	for _, join := range joins{
		lsides := strings.Split( join.LSide, "-" );
		rsides := strings.Split( join.RSide, "-" );
		str += join.Type + " JOIN " + rsides[0]+ 
				" ON " + lsides[0] + "." + lsides[1] + 
				" = " + rsides[0] + "." + rsides[1] + " "
	}
	return str
}
