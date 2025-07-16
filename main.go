package main

import (
    "flag"
    "time"
    "github.com/tahmaseb73/V2rayCollector/collector"
    "github.com/jszwec/csvutil"
    "github.com/projectdiscovery/gologger"
    "github.com/projectdiscovery/gologger/levels"
)

type ChannelsType struct {
    URL             string `csv:"URL"`
    AllMessagesFlag bool   `csv:"AllMessagesFlag"`
}

func main() {
    gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
    flag.Parse()

    fileData, err := collector.ReadFileContent("config/channels.csv")
    if err != nil {
        gologger.Fatal().Msg("error: " + err.Error())
    }

    var channels []ChannelsType
    if err = csvutil.Unmarshal([]byte(fileData), &channels); err != nil {
        gologger.Fatal().Msg("error: " + err.Error())
    }

    for _, channel := range channels {
        channel.URL = collector.ChangeUrlToTelegramWebUrl(channel.URL)
        resp := collector.HttpRequest(channel.URL)
        doc, err := goquery.NewDocumentFromReader(resp.Body)
        resp.Body.Close()
        if err != nil {
            gologger.Error().Msg(err.Error())
        }
        gologger.Info().Msg("Crawling " + channel.URL)
        collector.CrawlForV2ray(doc, channel.URL, channel.AllMessagesFlag)
        gologger.Info().Msg("Crawled " + channel.URL + " ! ")
    }

    gologger.Info().Msg("Creating output files !")
    // به‌روزرسانی README
    updateReadme()
    gologger.Info().Msg("All Done :D")
}

func updateReadme() {
    readmeContent := `# V2rayCollector

This project crawls V2Ray, Hysteria2, TUIC, and WireGuard configs from Telegram channels.

**Last Updated:** ` + time.Now().Format(time.RFC1123) + `

## Available Configs
- Shadowsocks: [config/ss_iran.txt](config/ss_iran.txt)
- VMess: [config/vmess_iran.txt](config/vmess_iran.txt)
- Trojan: [config/trojan_iran.txt](config/trojan_iran.txt)
- VLess: [config/vless_iran.txt](config/vless_iran.txt)
- Hysteria2: [config/hysteria2_iran.txt](config/hysteria2_iran.txt)
- TUIC: [config/tuic_iran.txt](config/tuic_iran.txt)
- WireGuard: [config/wireguard_iran.txt](config/wireguard_iran.txt)
- Mixed: [config/mixed_iran.txt](config/mixed_iran.txt)
`
    if err := collector.WriteToFile(readmeContent, "README.md"); err != nil {
        gologger.Error().Msg("Error updating README: " + err.Error())
    }
}
