package services

import (
    "fmt"
    "io"
    "net/http"
    "strings"
)

func FetchHTMLPage(url string) (string, error) {
    normalized := url
    if !strings.Contains(normalized, "://") {
        normalized = "https://" + url
    }

    resp, err := http.Get(normalized)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("HTTP error: %d", resp.StatusCode)
    }

    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(bodyBytes), nil
}