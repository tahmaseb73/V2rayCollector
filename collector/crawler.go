package collector

import (
    "net/http"
    "regexp"
    "strings"
    "github.com/PuerkitoBio/goquery"
    "github.com/projectdiscovery/gologger"
)

var client = &http.Client{}
var maxMessages = 100

func CrawlForV2ray(doc *goquery.Document, channelLink string, HasAllMessagesFlag bool) {
    messages := doc.Find(".tgme_widget_message_wrap").Length()
    link, exist := doc.Find(".tgme_widget_message_wrap .js-widget_message").Last().Attr("data-post")

    if messages < maxMessages && exist {
        number := strings.Split(link, "/")[1]
        var err error
        doc, err = GetMessages(maxMessages, doc, number, channelLink) // مدیریت هر دو مقدار بازگشتی
        if err != nil {
            gologger.Error().Msg("Error getting messages: " + err.Error())
            return
        }
    }

    configs := map[string]string{
        "ss":         "",
        "vmess":      "",
        "trojan":     "",
        "vless":      "",
        "hysteria2":  "",
        "tuic":       "",
        "wireguard":  "",
        "mixed":      "",
    }

    if HasAllMessagesFlag {
        doc.Find(".tgme_widget_message_text").Each(func(j int, s *goquery.Selection) {
            messageText, _ := s.Html()
            str := strings.Replace(messageText, "<br/>", "\n", -1)
            doc, _ := goquery.NewDocumentFromReader(strings.NewReader(str))
            messageText = doc.Text()
            line := strings.TrimSpace(messageText)
            lines := strings.Split(line, "\n")
            for _, data := range lines {
                extractedConfigs := strings.Split(ExtractConfig(data, []string{}), "\n")
                for _, config := range extractedConfigs {
                    config = strings.ReplaceAll(config, " ", "")
                    if config != "" {
                        for protoRegex, regexValue := range myregex {
                            re := regexp.MustCompile(regexValue)
                            if re.MatchString(config) {
                                if protoRegex == "vmess" {
                                    config = EditVmessPs(config, "mixed", false)
                                }
                                configs[protoRegex] += config + "\n"
                            }
                        }
                        configs["mixed"] += config + "\n"
                    }
                }
            }
        })
    } else {
        doc.Find("code,pre").Each(func(j int, s *goquery.Selection) {
            messageText, _ := s.Html()
            str := strings.ReplaceAll(messageText, "<br/>", "\n")
            doc, _ := goquery.NewDocumentFromReader(strings.NewReader(str))
            messageText = doc.Text()
            line := strings.TrimSpace(messageText)
            lines := strings.Split(line, "\n")
            for _, data := range lines {
                extractedConfigs := strings.Split(ExtractConfig(data, []string{}), "\n")
                for protoRegex, regexValue := range myregex {
                    for _, config := range extractedConfigs {
                        re := regexp.MustCompile(regexValue)
                        matches := re.FindStringSubmatch(config)
                        if len(matches) > 0 {
                            config = strings.ReplaceAll(config, " ", "")
                            if config != "" {
                                if protoRegex == "vmess" {
                                    config = EditVmessPs(config, protoRegex, false)
                                }
                                configs[protoRegex] += config + "\n"
                            }
                        }
                    }
                }
            }
        })
    }

    for proto, configcontent := range configs {
        lines := RemoveDuplicate(configcontent)
        lines = AddConfigNames(lines, proto)
        lines = strings.TrimSpace(lines)
        if err := WriteToFile(lines, "config/"+proto+"_iran.txt"); err != nil {
            gologger.Error().Msg(err.Error())
        }
    }
}
