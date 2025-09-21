package internal

import (
    "fmt"
    "io/ioutil"
    "net/http"
)

func Get(url string) (int, []byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        return 0, nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return resp.StatusCode, nil, err
    }

    return resp.StatusCode, body, nil
}
