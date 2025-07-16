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

func loadMore(link string) (*goquery.Document, error) {
    req, err := http.NewRequest("GET", link, nil)
    if err != nil {
        return nil, err
    }
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return nil, err
    }
    return doc, nil
}

func GetMessages(length int, doc *goquery.Document, number string, channel string) (*goquery.Document, error) {
    x, err := loadMore(channel + "?before=" + number)
    if err != nil {
        gologger.Error().Msg("Error loading more messages: " + err.Error())
        return doc, err
    }
    html2, err := x.Html()
    if err != nil {
        gologger.Error().Msg("Error getting HTML: " + err.Error())
        return doc, err
    }
    reader2 := strings.NewReader(html2)
    doc2, err := goquery.NewDocumentFromReader(reader2)
    if err != nil {
        gologger.Error().Msg("Error creating document: " + err.Error())
        return doc, err
    }
    doc.Find("body").AppendSelection(doc2.Find("body").Children())
    newDoc := goquery.NewDocumentFromNode(doc.Selection.Nodes[0])
    messages := newDoc.Find(".js-widget_message_wrap").Length()
    if messages > length {
        return newDoc, nil
    }
    num, err := strconv.Atoi(number)
    if err != nil {
        gologger.Error().Msg("Error converting number: " + err.Error())
        return newDoc, err
    }
    n := num - 21
    if n > 0 {
        return GetMessages(length, newDoc, strconv.Itoa(n), channel)
    }
    return newDoc, nil
}
