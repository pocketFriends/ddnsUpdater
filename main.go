package main

import (
    `log`
    `time`

    `ddnsUpdater/internal/config`
    `ddnsUpdater/internal/utils`

    `github.com/go-resty/resty/v2`
)

func main() {

    for {
        time.Sleep(time.Minute)

        client := resty.New()

        ip, err := utils.GetIpForNaver(client)
        if err != nil {
            log.Fatal(err)
        }

        configInstance := config.ReadConfigs()
        domains := configInstance.GetDomains(ip)

        for _, domain := range domains {
            updateStatus := utils.UpdateDomainIfIPChanged(client, domain, ip)
            if updateStatus == utils.GoogleDomainAPIStatusGood {
                configInstance.Update(domain, ip)
            }
        }

        configInstance.Write(configInstance)
    }

}
