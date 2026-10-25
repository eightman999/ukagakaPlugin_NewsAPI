package main

/*
   #include <stdlib.h>
   #include <string.h>
*/
import "C"

import (
    "fmt"
    "unsafe"
    "strings"
    "regexp"
	"strconv"
)


var DIRECTORY string
func main() {
    fmt.Println( "test" )
    DIRECTORY = ".."
    LoadJson()
}


var Count = 0

var References []string 
var CheckID         = regexp.MustCompile("^ID: ")
var CheckReference  = regexp.MustCompile("^Reference.+?: ")

type ResponseStruct struct {
    Shiori  string
    Sender  string
    Charset string
    Marker  string
    Value   string
}
func GetResponse( r *ResponseStruct ) string {
    V := ""
    if r.Value  != "" { V = "Value: "  + r.Value     + "\r\n" }
    res :=  r.Shiori    + "\r\n" + 
            r.Sender    + "\r\n" + 
            r.Charset   + "\r\n" + 
            V + "\r\n\r\n"
    return res
}

//export load
func load(h unsafe.Pointer, length C.long ) C.int {
    fmt.Println( "load NewsAPI" )
    DIRECTORY = C.GoStringN(( *C.char )( unsafe.Pointer( h )), ( C.int )( length ))
    fmt.Println( DIRECTORY  )

    //設定読み込み。
    LoadJson()

    //ニュース検索開始。
    go NewsInit()

	freeHGlobal( h )
	return 1
}


//export unload
func unload() C.int {
    fmt.Println( "unload NewsAPI" )
	return 1
}


//export request
func request( h unsafe.Pointer, length *C.long ) unsafe.Pointer {
	RequestText := C.GoStringN(( *C.char )( unsafe.Pointer( h )), ( C.int )( *length ))
	freeHGlobal( h )


    Value           := ""
    ID              := ""
    References      = []string{}
    //var NOTIFY bool = false

    Response := new( ResponseStruct )
    Response.Sender  = "Sender: NewsAPI"
    Response.Charset = "Charset: UTF-8"

    //IDとReference
    //必要な情報を分解する。
    RequestLines := strings.Split( RequestText , "\r\n" )
    for _ , line := range RequestLines {
        if( line == "NOTIFY PLUGIN/2.0" ){
            //"GET PLUGIN/2.0";
            //NOTIFY = true

        } else if CheckID.MatchString( line )  {
            //fmt.Println( line )
            ID = CheckID.ReplaceAllString( line , "" )

        } else if CheckReference.MatchString( line )  {
            //fmt.Println( line )
            ref := CheckReference.ReplaceAllString( line , "" )
            References = append( References , ref )

        } else {
            //fmt.Println( line )
        }
    }

    //実行関数
    if ID == "OnMenuExec" {
        _res := ""
        _i := 0
        for( _i < len( News_URL ) ){
            _res = _res + "\\_a[OnUrlSelectNewsAPI," + News_URL[_i] + "," + strconv.Itoa( _i ) + "]" + News_Title[_i] + "\\_a\\n" 
            _i++
        }
        Value = _res
    } else if ID == "OnUrlSelectNewsAPI" {
        _i , err := strconv.Atoi( References[ 1 ] )
        // 数字に変換できない場合
        if err != nil { return nil }

        _i++
        _res := ""
        // 開いたURL以降の情報を再表示
        for( _i < len( News_URL ) ){
            _res = _res + "\\_a[OnUrlSelectNewsAPI," + News_URL[_i] + "," + strconv.Itoa( _i ) + "]" + News_Title[_i] + "\\_a\\n" 
            _i++
        }
        Value = "\\j[" + References[ 0 ]+ "]" + _res
    } else {
        //fmt.Println( "no touch :" + ID )
        //fmt.Print( "NOTIFY : " )
        //fmt.Println( NOTIFY )
        //fmt.Print( "References : " )
        //fmt.Println( References )
        //fmt.Println( "" )
    }


    if Value == "" {
        Response.Shiori  = "PLUGIN/2.0 204 No Content"
    } else {
        Response.Shiori = "PLUGIN/2.0 200 OK"
        Response.Value  = Value
    }

    res_buf := C.CString( GetResponse( Response ))
    defer C.free( unsafe.Pointer( res_buf ))

	res_size := C.strlen( res_buf )
	ret      := allocHGlobal( uintptr( res_size ))
	C.memcpy(( unsafe.Pointer )( ret ) , ( unsafe.Pointer )( res_buf ) , res_size )
	*length = ( C.long )( res_size )
	return ret
}

















