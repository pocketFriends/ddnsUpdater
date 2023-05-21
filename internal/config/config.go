package config

import (
    `fmt`
    `log`
    `net`
    `os`

    `github.com/BurntSushi/toml`
)

type Domain struct {
    ID       string `toml:"id"`
    Password string `toml:"password"`
    Hostname string `toml:"hostname"`
    OldIP    string `toml:"old_ip"`
}

type Config struct {
    Domains []Domain `toml:"domain"`
}

func ReadConfigs() *Config {
    var config Config
    if _, err := toml.DecodeFile("config.toml", &config); err != nil {
        log.Fatal(err)
    }

    return &config
}

func (c *Config) GetDomains(ip net.IP) []string {
    var domains []string
    for _, domain := range c.Domains {
        domains = append(domains, getGoogleDomainUrl(domain.ID, domain.Password, domain.Hostname, ip))
    }
    return domains
}

func (c *Config) Update(hostname string, ip net.IP) {
    for i, domain := range c.Domains {
        if domain.Hostname == hostname {
            c.Domains[i].OldIP = ip.String()
        }
    }

}

func (c *Config) Write(config *Config) {
    file, err := os.Create("config.toml")
    if err != nil {
        log.Println("Config Write Error: ", err)
        return
    }
    if err := toml.NewEncoder(file).Encode(config); err != nil {
        log.Println("Config Write Error encode: ", err)
    } else {
        log.Println("Config Write Success")
    }
}

func getGoogleDomainUrl(userId, userPw, hostname string, ip net.IP) string {
    return fmt.Sprintf(
        "https://%s:%s@domains.google.com/nic/update?hostname=%s&myip=%s", userId, userPw, hostname, ip.String())
}
