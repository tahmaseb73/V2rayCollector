package collector

import (
    "net/http"
    "strconv"
    "strings"
    "github.com/PuerkitoBio/goquery"
    "github.com/projectdiscovery/gologger"
)

func HttpRequest(url string) *http.Response {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        gologger.Fatal().Msg("Error requesting: " + url + " Error: " + err.Error())
    }
    resp, err := client.Do(req)
    if err != nil {
        gologger.Fatal().Msg(err.Error())
    }
    return resp
}

func loadMore(link string) *goquery.Document {
    req, _ := http.NewRequest("GET", link, nil)
    resp, _ := client.Do(req)
    doc, _ := goquery.NewDocumentFromReader(resp.Body)
    return doc
}

func GetMessages(length int, doc *goquery.Document, number string, channel string) *goquery.Document {
    x := loadMore(channel + "?before=" + number)
    html2, _ := x.Html()
    reader2 := strings.NewReader(html2)
    doc2, _ := goquery.NewDocumentFromReader(reader2)
    doc.Find("body").AppendSelection(doc2.Find("body").Children())
    newDoc := goquery.NewDocumentFromNode(doc.Selection.Nodes[0])
    messages := newDoc.Find(".js-widget_message_wrap").Length()
    if messages > length {
        return newDoc
    }
    num, _ := strconv.Atoi(number)
    n := num - 21
    if n > 0 {
        return GetMessages(length, newDoc, strconv.Itoa(n), channel)
    }
    return newDoc
}
