package main

import (
	"flag"
	"fmt"
)

var (
	help   = flag.Bool("help", false, "this help")
	fields = flag.String("fields", "ip,port", "fields which you want to select")
	query  = flag.String("query", "", "query string")
	email  = flag.String("email", "", "an email which you login to fofa.so")
	key    = flag.String("key", "", "md5 string which you can find on userinfo page")
	page   = flag.Int("page", 1, "page number you want to query")
	out    = flag.String("out", "", "output file path")
	format = flag.String("format", "\t", "output format")
	count  = flag.Bool("count", false, "only count the total number of matches,true or false")
)

func usage() {
	fmt.Println(`
    Fofa is a tool for discovering assets.

    Usage:

            fofa init|info|search|version option argument ...

    The options are:
		
	    init:
		    email           the email which you login to fofa.so

		    key             the md5 string which you can find on userinfo page
		
	    search:
		    fields          fields which you want to select
		                    Use ip,port as default.

		    format          output format
		                    Default is /t splice, you can choose other.

		    query           query statement which is similar to the statement used in the fofa.so

		    page            page number you want to query, 100 records per page
		                    If page is not set or page is less than 1, page will be set to 1.

		    out             output file path
		                    Print to the terminal as default.

		    count           only count the total number of matches,true or false	
		                    False as default.
    `)
}
