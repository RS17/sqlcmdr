# sqlcmdr


golang package for slightly easier sql commands

This is a working (but also very WIP) package intended for use with sqlite3 that makes it slightly easier to create basic sql select and insert commands for a simple sqlite database while minimizing putting long strings like "SELECT x FROM y" all over your code.  Over time I expect to add more features as I use this for side projects.  In the meantime use at your own risk.

This is also used extensively in my other project stratterm, which is developed concurrently with this project.  You can see several examples there.

## Examples:
### JustRunIt 
Runs command in one line, returns nothing.  Doesn't require initializing DB first.

```

sqlcmdr.JustRunIt( 	"CREATE TABLE IF NOT EXISTS developers (" + 
					"name VARCHAR( 50 ) NOT NULL PRIMARY KEY," +
					"description VARCHAR( 1000 ) )"  )
```

### InitDB 
Initializes DB and returns connection, necessary for select and insert  

 ```  conn := sqlcmdr.InitDB() ```

### InsertCmd/Insert 
Use this to do a simple insert. 

```
   // create command
   icmd := sqlcmdr.InsertCmd{ Tablename: "developers" }
   icmd.Add( "name", "Ravi" )
   icmd.Add( "description", "Me" )

   // run it
   conn := sqlcmdr.InitDB() 
   defer conn.Close()
   sqlcmdr.Insert( conn, icmd )
   
```

### SelectCmd/Select 
Do a simple select (there is also a very rudimentary ability to add a where clause, to be improved in future).  Returns array of arrays corresponding to 2D table.  Can do joins with AddJoin.

```
	// select all columns from developers table
	scmd := sqlcmdr.SelectCmd{ Tablename: "developers", Columns: "name, commitdate" } 
	scmd.AddJoin( "developers-name", "projects-developer", "LEFT" }				// jointable uses tablename-columnname to define the matching join columns, last value is type
	scmd.AddJoin( "projects-id", "commits-projectid", "RIGHT" }
	
	// run it
	conn := sqlcmdr.InitDB() 
   	defer conn.Close()
	result := sqlcmdr.Select( conn, scmd ) 							// result is [][]interface{}.  You can use built-in converters (see below) or make your own
```	

### ResultCSV
Converts [][]interface to a csv, output to a specified file

```	
	message = sqlcmdr.ResultCSV( result, "result.csv" );		// returns confirmation message
```	

### ResultString 
Converts[][]interface to a string

```	
	resultstr = sqlcmdr.ResultString( result );
```
