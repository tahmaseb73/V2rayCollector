package collector

import (
    "os"
    "strings"
    "slices"
)

func ChangeUrlToTelegramWebUrl(input string) string {
    if !strings.Contains(input, "/s/") {
        index := strings.Index(input, "/t.me/")
        if index != -1 {
            return input[:index+len("/t.me/")] + "s/" + input[index+len("/t.me/"):]
        }
    }
    return input
}

func ReadFileContent(filePath string) (string, error) {
    content, err := os.ReadFile(filePath)
    if err != nil {
        return "", err
    }
    return string(content), nil
}

func Reverse(lines []string) []string {
    for i := 0; i < len(lines)/2; i++ {
        j := len(lines) - i - 1
        lines[i], lines[j] = lines[j], lines[i]
    }
    return lines
}

func RemoveDuplicate(config string) string {
    lines := strings.Split(config, "\n")
    slices.Sort(lines)
    lines = slices.Compact(lines)
    return strings.Join(lines, "\n")
}

func WriteToFile(fileContent string, filePath string) error {
    if _, err := os.Stat(filePath); err == nil {
        err = os.WriteFile(filePath, []byte{}, 0644)
        if err != nil {
            return err
        }
    } else if os.IsNotExist(err) {
        _, err = os.Create(filePath)
        if err != nil {
            return err
        }
    } else {
        return err
    }
    return os.WriteFile(filePath, []byte(fileContent), 0644)
}
