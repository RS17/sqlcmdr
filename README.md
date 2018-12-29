# sqlcmdr
go package for slightly easier sql commands

This is a working (but also very WIP) package intended for use with sqlite3 that makes it slightly easier to create basic sql select and insert commands for a simple sql database while minimizing putting long strings like "SELECT x FROM y" all over your code.  Over time I expect to add more features as I use this for side projects.  In the meantime use at your own risk.

Examples:
JustRunIt - runs command in one line, returns nothing.  Doesn't require initializing DB first.  
   sqlcmdr.JustRunIt( 	"CREATE TABLE IF NOT EXISTS developers (" +   
			"name VARCHAR( 50 ) NOT NULL PRIMARY KEY," +   
			"description VARCHAR( 1000 ) )"  )

InitDB - initializes DB and returns connection, necessary for select and insert  
   conn := sqlcmdr.InitDB() 

InsertCmd/Insert - use this to do a simple insert   
   // create command
   icmd := sqlcmdr.InsertCmd{ Tablename: "developers" }
   icmd.Add( "name", "Ravi" )
   icmd.Add( "description", "Me" )

   // run it
   conn := sqlcmdr.InitDB() 
   sqlcmdr.Insert( conn, icmd )
   conn.Close()

SelectCmd/Select - do a simple select (there is also ability to add a where clause but this is still untested so it's not recommonded).  Returns array of arrays corresponding to 2D table  
   
   // select all columns from developers table  
   scmd := sqlcmdr.SelectCmd{ Tablename: "developers", Columns: "*" }  
   retvals := sqlcmdr.Select( conn, scmd )  
  
   // retvals is [][]interface{}, so this is how you could convert it to a big long string or something  
   var result string  
   for _, element := range( retvals ){  
      for _, e2 := range( element ){  
         if( e2 != nil ){  
            result = result + string(e2.([]uint8))  
	 }  
      }  
   }  
   conn.Close()



It can also probably be used for other sql databases with minimal changes.
