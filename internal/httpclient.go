package internal

import (
    _"fmt"
    "bytes"
    "io"
    "net/http"
)

//get logic for get requests
func Get(url string) (int, []byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        return 0, nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return resp.StatusCode, nil, err
    }

    return resp.StatusCode, body, nil
}

//post requests
func Post(url string, body []byte) (int, []byte, error) {
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
    if err != nil {
        return 0, nil, err
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return resp.StatusCode, nil, err
    }

    return resp.StatusCode, respBody, nil
}
