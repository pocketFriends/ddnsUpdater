package utils

import (
    `fmt`
    `log`
    `net`
    `strings`

    `github.com/go-resty/resty/v2`
)

type GoogleDomainAPIStatus string

const (
    GoogleDomainAPIStatusGood     GoogleDomainAPIStatus = "good"     // 업데이트에 성공했습니다. IP 주소가 변경될 때까지 업데이트를 더 시도해서는 안 됩니다.
    GoogleDomainAPIStatusNochg    GoogleDomainAPIStatus = "nochg"    // 지정한 IP 주소가 이 호스트에 이미 설정되어 있습니다. IP 주소가 변경될 때까지 업데이트를 더 시도해서는 안 됩니다.
    GoogleDomainAPIStatusBadauth  GoogleDomainAPIStatus = "badauth"  // 지정한 호스트에서 사용자 이름과 비밀번호 조합이 유효하지 않습니다.
    GoogleDomainAPIStatusNotfqdn  GoogleDomainAPIStatus = "notfqdn"  // 지정한 호스트 이름이 올바른 정규화된 도메인 이름이 아닙니다.
    GoogleDomainAPIStatusNohost   GoogleDomainAPIStatus = "nohost"   // 호스트 이름이 없거나 동적 DNS가 사용 설정되지 않았습니다.
    GoogleDomainAPIStatusAbuse    GoogleDomainAPIStatus = "abuse"    // 이전 응답을 올바르게 해석하지 못하여 호스트 이름에 대한 동적 DNS 액세스가 차단되었습니다.
    GoogleDomainAPIStatusBadagent GoogleDomainAPIStatus = "badagent" // 동적 DNS 클라이언트의 요청이 잘못되었습니다. 사용자 에이전트가 요청에 설정되어 있는지 확인하세요.
    GoogleDomainAPIStatus911      GoogleDomainAPIStatus = "911"      // Google 시스템에서 오류가 발생했습니다. 5분 후에 다시 시도하세요.
    NotDefined                    GoogleDomainAPIStatus = ""         // 정의되지 않은 상태
)

type GoogleDomainAPIUpdate struct {
    Hostname string
    Status   GoogleDomainAPIStatus
}

func UpdateDomainIfIPChanged(client0 *resty.Client, url string, ip net.IP) GoogleDomainAPIStatus {
    isNeedUpdate := shouldUpdateRegisteredIP(url, ip)
    log.Println("UpdateDomainIfIPChanged isNeedUpdate: ", isNeedUpdate)
    if !isNeedUpdate {
        return NotDefined
    }

    result := updateDomain(client0, url)
    log.Println("UpdateDomainIfIPChanged result: ", result)
    if result == NotDefined {
        log.Println("UpdateDomainIfIPChanged Error: ", result)
        return NotDefined
    }
    return GoogleDomainAPIStatusGood
}

func shouldUpdateRegisteredIP(url string, ip net.IP) bool {
    urlToHostname, err := urlToHostname(url)
    if err != nil {
        log.Println("ShouldUpdateRegisteredIP Error urlToHostname: ", err)
        return false
    }

    oldIp, err := net.LookupIP(urlToHostname)
    if err != nil {
        if !strings.Contains(err.Error(), "The requested name is valid") {
            log.Println("ShouldUpdateRegisteredIP Error LookupIP: ", err)
            return false
        } else {
            return true
        }
    }

    if oldIp[0].Equal(ip) {
        return false
    }
    return true
}
func updateDomain(client *resty.Client, url string) GoogleDomainAPIStatus {
    resp, err := client.R().
        SetHeader("Content-Type", "application/x-www-form-urlencoded").
        Get(url)
    if err != nil {
        log.Println("UpdateDomainIfIPChanged Error: ", err)
        return NotDefined
    }

    status, err := domainRequestStatus(resp.String())
    if err != nil {
        log.Println("UpdateDomainIfIPChanged Error: ", err)
        return NotDefined
    }

    return status
}

func domainRequestStatus(responseText string) (GoogleDomainAPIStatus, error) {
    log.Println("get domainRequestStatus: ", responseText)

    if strings.Contains(responseText, "good") {
        return GoogleDomainAPIStatusGood, nil
    } else if strings.Contains(responseText, "nochg") {
        return GoogleDomainAPIStatusNochg, nil
    } else if strings.Contains(responseText, "badauth") {
        return GoogleDomainAPIStatusBadauth, nil
    } else if strings.Contains(responseText, "notfqdn") {
        return GoogleDomainAPIStatusNotfqdn, nil
    } else if strings.Contains(responseText, "nohost") {
        return GoogleDomainAPIStatusNohost, nil
    } else if strings.Contains(responseText, "abuse") {
        return GoogleDomainAPIStatusAbuse, nil
    } else if strings.Contains(responseText, "badagent") {
        return GoogleDomainAPIStatusBadagent, nil
    } else if strings.Contains(responseText, "911") {
        return GoogleDomainAPIStatus911, nil
    } else {
        return "", fmt.Errorf("Unknown responseText: %s", responseText)
    }
}

func urlToHostname(url string) (string, error) {
    // https://WZ2jX47tV5DxFYMo:aT3B9r1kcj8X0lgg@@domains.google.com/nic/update?hostname=remote.pocketfriends.chat&myip=61.78.126.99
    splitURL := strings.Split(url, "hostname=")
    if len(splitURL) != 2 {
        return "", fmt.Errorf("urlToHostname Error: %s", url)
    }

    splitURL = strings.Split(splitURL[1], "&myip=")
    if len(splitURL) != 2 {
        return "", fmt.Errorf("urlToHostname Error: %s", url)
    }
    return splitURL[0], nil
}
