package collector

import (
    "net/http"
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
        doc = GetMessages(maxMessages, doc, number, channelLink)
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
                for _, extractedConfig := range extractedConfigs {
                    extractedConfig = strings.ReplaceAll(extractedConfig, " ", "")
                    if extractedConfig != "" {
                        for protoRegex, regexValue := range myregex {
                            re := regexp.MustCompile(regexValue)
                            if re.MatchString(extractedConfig) {
                                if protoRegex == "vmess" {
                                    extractedConfig = EditVmessPs(extractedConfig, "mixed", false)
                                }
                                configs[protoRegex] += extractedConfig + "\n"
                            }
                        }
                        configs["mixed"] += extractedConfig + "\n"
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
                    re := regexp.MustCompile(regexValue)
                    matches := re.FindStringSubmatch(extractedConfig)
                    if len(matches) > 0 {
                        extractedConfig = strings.ReplaceAll(extractedConfig, " ", "")
                        if extractedConfig != "" {
                            if protoRegex == "vmess" {
                                extractedConfig = EditVmessPs(extractedConfig, protoRegex, false)
                            }
                            configs[protoRegex] += extractedConfig + "\n"
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
