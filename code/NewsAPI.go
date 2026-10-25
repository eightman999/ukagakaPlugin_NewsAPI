package main

import (

    "fmt"
    "os"
    "path/filepath"

    "io"
    "net/http"
    "encoding/json"
    "time"


)
type newsConfigStruct struct {
    API             string
    SearchWord      string 
    ENG_SearchWord  string
    ENG_SOURCE      string
    JP_SearchWord   string
    JP_DOMAIN       string
    Max             string
}
var newsConfig newsConfigStruct

func LoadJson(){
	JsonNewsAPI, err := os.Open( filepath.Join( DIRECTORY , "Config.json" ))
	if err != nil {
        fmt.Println( err )
	}
	defer JsonNewsAPI.Close()
    decoder := json.NewDecoder( JsonNewsAPI )
    err     = decoder.Decode( &newsConfig )
	if err != nil {
        fmt.Println(  err  )
    }
}



var JP_news     string 
var ENG_news    string 
var News_Title   []string
var News_URL     []string

var News_Result_JP  string
var News_Result_ENG string
// {{{
type searchWord struct {
	WORD string
}
// }}}
// {{{
type source struct{
    Id      string      `json:"id"`
    Name    string      `json:"name"` 
}
// }}}
// {{{
type article struct{
   Source       source  `json:"source"`
   Author       string  `json:"author"`
   Title        string  `json:"title"`
   Description  string  `json:"description"`
   Url          string  `json:"url"`
   UrlToImage   string  `json:"urlToImage"`
   PublishedAt  string  `json:"publishedAt"`
   Content      string  `json:"content"`
}

// }}}
// {{{
type newsStruct struct{
    Status          string  `json:"status"`
    TotalResults    int     `json:"totalResults"`
    //Articles        string  `-`
    Articles        []article `json:"articles"`
}
// }}}


// {{{
func NewsInit(){
    JP_news     = ""
    ENG_news    = ""

	JsonNewsAPI, err := os.Open( filepath.Join( DIRECTORY , "News.json" ))
	if err != nil {
        fmt.Println( err )
	}
	defer JsonNewsAPI.Close()
    decoder := json.NewDecoder( JsonNewsAPI )
    err     = decoder.Decode( &newsConfig )
	if err != nil { fmt.Println(  err  ) }

    // 日本
    // JP_news  , _ = GetNews( newsConfig.JP_SearchWord )
    // 外国
    // ENG_news , _ = GetNews( newsConfig.ENG_SearchWord )

    JP_news  , _  = GetNews( "{\"WORD\":\"" + newsConfig.JP_SearchWord + "\"}" )
    // type指定してENG仕様にする必要がある。
    // ENG_news , _ = GetNews( "{\"WORD\":\"" + newsConfig.ENG_SearchWord + "\"}" )


}
// }}}
//// {{{
//func GetNews() {
//    apiKey  := Config.API 
//    if apiKey == "" {
//        return
//    }
//    url     := "https://newsapi.org/v2/everything"

//    req, _ := http.NewRequest("GET", url , nil)
//    req.Header.Add("Content-Type", "application/json")
//    req.Header.Add("X-Api-Key", apiKey )

//    q := req.URL.Query()
//    q.Add("sortBy"  , "publishedAt" )

//    //ロイター
//    //q.Add("sources"  , "reuters"         )
//    //q.Add("sources"  , "fox-news"         )
//    //q.Add("sources"  , "bbc-news"         )
//    // q.Add("sources"  , "bbc-news,fox-news"         )
//    // q.Add("sources"  , "bbc-news,fox-news"         )

//    //検索ワード
//    q.Add("q"       , "中国"    )
//    // q.Add("q"       , Config.SearchWord    )

//    nowUTC := time.Now().UTC().Format(time.RFC3339)
//    q.Add("to"  , nowUTC )

//    //一昨日の今の時間から昨日の今の時間までを検索対象に。
//    hourUTC := time.Now().UTC().AddDate( 0, 0, -2).Format(time.RFC3339)
//    q.Add("from"  , hourUTC )

//    //q.Add("domains"  , "bbc.co.uk" )
//    q.Add("domains"  , "nhk.or.jp" )
//    // q.Add("excludeDomains"  , "asahi.com,nhk.or.jp" )
//    //q.Add("pageSize", "1"         )
//    q.Add("pageSize", "100"         )

//    req.URL.RawQuery = q.Encode()
//    var client *http.Client = &http.Client{}
//    resp, _ := client.Do( req )
//    defer resp.Body.Close()

//    var news newsStruct
//    body, _ := io.ReadAll(resp.Body)
//    //全文
//    // fmt.Println( string( body ) )

//    json.Unmarshal( body , &news)
//    //fmt.Println( news.Status )
//    fmt.Println( len( news.Articles))
//    // fmt.Println( news.Articles[0] )

//    _i := 0
//    NewsPaper = ""
//    for( _i < len( news.Articles ) ){
//        NewsPaper = NewsPaper + "\\_a[OnNewsPaperOpenLink," + news.Articles[_i].Url + "]" + news.Articles[_i].Title + "\\_a\\n\\n"

//        fmt.Println( news.Articles[_i].Title )
//        fmt.Print( news.Articles[_i].Source.Name )
//        fmt.Print( news.Articles[_i].Description)
//        // fmt.Print( news.Articles[_i].PublishedAt )

//        _i++
//    }
//    if NewsPaper != "" {
//        NewsPaper = "\\0\\b[2]" + NewsPaper 
//    }
//}
//// }}}
// {{{
// func GetNews( TYPE int , WORD string , MAX string ) {
// この関数はAIに呼ばせる必要はないな。
func GetNews( search_json string ) ( string , error ) {
    fmt.Println( search_json )

    var search_word searchWord
    // {{{
    if err := json.Unmarshal([]byte( search_json ), &search_word); err != nil {
        fmt.Println(err)
        return "" , err
    }
    // }}}
    fmt.Println( "SEARCH NEWS : " + search_word.WORD )

    TYPE := 0
    MAX  := "100"

    url     := "https://newsapi.org/v2/everything"
    apiKey  := newsConfig.API 

    // {{{
    if apiKey == "" { 
        fmt.Println( "API キーがありません。" )
        return "ニュースはありません。" , nil
    }
// }}}

    req, _ := http.NewRequest("GET", url , nil)
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("X-Api-Key", apiKey )

    // クエリ作成
    // {{{
    q := req.URL.Query()
    // 時間でソート
    q.Add("sortBy"  , "publishedAt" )
    //一昨日の今の時間から
    hourUTC := time.Now().UTC().AddDate( 0, 0, -2).Format(time.RFC3339)
    q.Add("from"  , hourUTC )
    // 現在時刻までを検索対象に。
    nowUTC := time.Now().UTC().Format(time.RFC3339)
    q.Add("to"  , nowUTC )

    // 検索最大数
    q.Add( "pageSize", MAX         )
    // q.Add( "pageSize", newsConfig.MAX         )

    // {{{
    // // 使わないメモ
    // 検索先を指定する。
    // ロイター
    // q.Add("sources"  , "reuters"          )
    // q.Add("sources"  , "fox-news"         )
    // q.Add("sources"  , "bbc-news"         )
    // q.Add("sources"  , "bbc-news,fox-news")
    // q.Add("sources"  , "bbc-news,fox-news")
    // q.Add("domains"  , "bbc.co.uk" )
    // q.Add("excludeDomains"  , "asahi.com,nhk.or.jp" )
    // }}}

    // 検索ワード
    q.Add("q"       , search_word.WORD )
    // q.Add("q"       , "")

    // 検索圏を絞る。
    if TYPE == 0 {
        // 日本の検索は NHK だけにしておく。
        // Source.Nameが指定されていない為Domainで指定する。
        // q.Add("q"       , newsConfig.JP_SearchWord  )
        q.Add("domains" , newsConfig.JP_DOMAIN      )
    } else {
        // q.Add("q"       , newsConfig.ENG_SearchWord )
        q.Add("sources" , newsConfig.ENG_SOURCE     )
    }
    // }}}

    req.URL.RawQuery = q.Encode()
    var client *http.Client = &http.Client{}
    resp, err := client.Do( req )
    // {{{
    if err != nil {
        fmt.Println( err )
        return "ニュースはありません。" , nil
    }
    // }}}
    defer resp.Body.Close()

    body, err := io.ReadAll( resp.Body )
    // {{{
    if err != nil {
        fmt.Println( err )
        return "ニュースはありません。" , nil
    }
    // }}}
    // fmt.Println( string( body ) )

    var news newsStruct
    json.Unmarshal( body , &news)
    // 接続可否
    // fmt.Println( news.Status )
    // 取得記事数
    // fmt.Println( news.TotalResults )

    // fmt.Println( news.Articles[0] )

    if news.TotalResults == 0 { 
        fmt.Println( "検索結果無し。" )
        return "ニュースはありません。" , nil 
    }
    // res := ""
    _i  := 0
    // {{{
    for( _i < len( news.Articles ) ){
        // どちらかが無ければスキップ
        if news.Articles[_i].Title == "" || news.Articles[_i].Description == "" {
            continue
        }

        News_Title = append( News_Title , news.Articles[_i].Title )
        News_URL   = append( News_URL   , news.Articles[_i].Url   )
        // fmt.Println( news.Articles[_i].Title )
        // fmt.Println( news.Articles[_i].Url   )
        // fmt.Println( news.Articles[_i].Source.Name )
        // fmt.Println( news.Articles[_i].Description)
        // fmt.Println( news.Articles[_i].PublishedAt )
        // res = res + news.Articles[_i].Title + "\n"
        // res = res + "### " + news.Articles[_i].Title + "\n" + news.Articles[_i].Description + "\n\n"

        _i++
    }
    // }}}

    // fmt.Println( res )

    return "" , nil


}
// }}}










