package utils

import (
    `errors`
    `net`
    `regexp`

    `github.com/go-resty/resty/v2`
)

const (
    naverIpUrl = `https://search.naver.com/search.naver?where=nexearch&sm=top_hty&fbm=0&ie=utf8&query=%EB%82%B4%EC%95%84%EC%9D%B4%ED%94%BC`
    regexIp    = `(?m)<div.*?ip_chk_box">(.*?)</div>`
)

func GetIpForNaver(client *resty.Client) (ip net.IP, err error) {
    resp, err := client.
        R().
        SetHeader(`User-Agent`, `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36`).
        Get(naverIpUrl)
    if err != nil {
        return nil, err
    }

    var re = regexp.MustCompile(regexIp)
    match := re.FindSubmatch(resp.Body())

    if len(match) != 2 {
        return nil, errors.New("no match")
    }
    ip = net.ParseIP(string(match[1]))
    if ip == nil {
        return nil, errors.New("invalid ip")
    }

    return ip, nil
}
