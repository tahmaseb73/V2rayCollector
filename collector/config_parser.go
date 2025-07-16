package collector

import (
    "encoding/base64"
    "encoding/json"
    "regexp"
    "strconv"
    "strings"
)

var (
    ConfigsNames = "@Vip_Security join us"
    ConfigFileIds = map[string]int32{
        "ss":         0,
        "vmess":      0,
        "trojan":     0,
        "vless":      0,
        "hysteria2":  0,
        "tuic":       0,
        "wireguard":  0,
        "mixed":      0,
    }
    myregex = map[string]string{
        "ss":         `(?m)(...ss:|^ss:)\/\/.+?(%3A%40|#)`,
        "vmess":      `(?m)vmess:\/\/.+`,
        "trojan":     `(?m)trojan:\/\/.+?(%3A%40|#)`,
        "vless":      `(?m)vless:\/\/.+?(%3A%40|#)`,
        "hysteria2":  `(?m)hy2:\/\/[a-f0-9-]{36}@[\w\.\:]+:\d+\?.*?#.+`,
        "tuic":       `(?m)tuic:\/\/[a-f0-9-]{36}:[^@]+@[\w\.\[\]\:]+:\d+\?.*?#.+`,
        "wireguard":  `(?m)wireguard:\/\/[^@]+@[\w\.\:]+:\d+\?.*?#.+`,
    }
)

func AddConfigNames(config string, configtype string) string {
    configs := strings.Split(config, "\n")
    newConfigs := ""
    for protoRegex, regexValue := range myregex {
        for _, extractedConfig := range configs {
            re := regexp.MustCompile(regexValue)
            matches := re.FindStringSubmatch(extractedConfig)
            if len(matches) > 0 {
                extractedConfig = strings.ReplaceAll(extractedConfig, " ", "")
                if extractedConfig != "" {
                    if protoRegex == "vmess" {
                        extractedConfig = EditVmessPs(extractedConfig, configtype, true)
                        if extractedConfig != "" {
                            newConfigs += extractedConfig + "\n"
                        }
                    } else if protoRegex == "ss" {
                        Prefix := strings.Split(matches[0], "ss://")[0]
                        if Prefix == "" {
                            ConfigFileIds[configtype]++
                            newConfigs += extractedConfig + ConfigsNames + " - " + strconv.Itoa(int(ConfigFileIds[configtype])) + "\n"
                        }
                    } else {
                        ConfigFileIds[configtype]++
                        newConfigs += extractedConfig + ConfigsNames + " - " + strconv.Itoa(int(ConfigFileIds[configtype])) + "\n"
                    }
                }
            }
        }
    }
    return newConfigs
}

func ExtractConfig(Txt string, Tempconfigs []string) string {
    for protoRegex, regexValue := range myregex {
        re := regexp.MustCompile(regexValue)
        matches := re.FindStringSubmatch(Txt)
        extractedConfig := ""
        if len(matches) > 0 {
            if protoRegex == "ss" {
                Prefix := strings.Split(matches[0], "ss://")[0]
                if Prefix == "" {
                    extractedConfig = "\n" + matches[0]
                } else if Prefix != "vle" {
                    d := strings.Split(matches[0], "ss://")
                    extractedConfig = "\n" + "ss://" + d[1]
                }
            } else {
                extractedConfig = "\n" + matches[0]
            }
            Tempconfigs = append(Tempconfigs, extractedConfig)
            Txt = strings.ReplaceAll(Txt, matches[0], "")
            ExtractConfig(Txt, Tempconfigs)
        }
    }
    return strings.Join(Tempconfigs, "\n")
}

func EditVmessPs(config string, fileName string, AddConfigName bool) string {
    if config == "" {
        return ""
    }
    slice := strings.Split(config, "vmess://")
    if len(slice) > 1 {
        decodedBytes, err := base64.StdEncoding.DecodeString(slice[1])
        if err == nil {
            var data map[string]interface{}
            if err = json.Unmarshal(decodedBytes, &data); err == nil {
                if AddConfigName {
                    ConfigFileIds[fileName]++
                    data["ps"] = ConfigsNames + " - " + strconv.Itoa(int(ConfigFileIds[fileName]))
                } else {
                    data["ps"] = ""
                }
                jsonData, _ := json.Marshal(data)
                return "vmess://" + base64.StdEncoding.EncodeToString(jsonData)
            }
        }
    }
    return ""
}
