# freshrss-oci

Uses packer and Ansible to create a custom image for use in Oracle Cloud Infrastructure (OCI).

Built on latest Ubuntu 22.04 OCI image. Builds every day at 3AM.

Includes the following applications:

- :newspaper_roll: [FreshRSS](https://github.com/FreshRSS/FreshRSS) for feed reading
- :mag_right: [FullTextRSS](https://github.com/heussd/fivefilters-full-text-rss-docker/tree/master) for converting partial feeds into full feeds for use in FreshRSS
- :bird: [Nitter](https://github.com/zedeus/nitter) for providing RSS feeds for Twitter
- :globe_with_meridians: Nginx to serve FreshRSS and reverse proxy to FullTextRSS and Nitter
- :closed_lock_with_key: [Tailscale](https://github.com/tailscale/tailscale) for management (SSH removed)

Along with the following config tweaks:

- :package: All packages updated
- :mechanical_arm: AppArmor configured with appropriate profiles
- :clock1: CRON jobs to automatically backup FreshRSS and update FreshRSS feeds
- :rewind: Custom restore script to support restoring the CRON created backups
- :no_entry_sign: iptables configured to limit allowed ingress
- :file_cabinet: Latest [site configs](https://github.com/fivefilters/ftr-site-config) mapped into FullTextRSS container
- :broom: Log rotate configured for all services
- :no_bell: MOTD advertisements/spam removed
- :toolbox: OCI CLI added
- :bar_chart: Oracle cloud agent added
- :no_good: Ubuntu Advantage removed
- :detective: Telemetry packages removed and telemetry domains blocked
- :robot: Nginx configured to prevent site being scraped by bots

## Deployment

Deployed via Terraform, see: [batinicaz/freshrss](https://github.com/batinicaz/freshrss)

### Nitter Guest Accounts

Nitter requires guest accounts due to Twitters recent API changes.

Until a service is integrated into the main nitter codebase to generate these (WIP see [#zedeus/nitter#983](https://github.com/zedeus/nitter/issues/983))

You can generate guest accounts using: https://gitlab.com/yawning/twitter-guest-account like so:

```bash
go run main.go -fetch-bearer-token -num-accounts 3
ansible-vault encrypt guest_accounts.json --output=guest_accounts.json.encrypted
```