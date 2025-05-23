# app-server-oci

Uses packer and Ansible to create a custom image for use in Oracle Cloud Infrastructure (OCI) for my selfhosted applications.

Built on latest Ubuntu 24.04 OCI image. Builds every day at 3AM.

Includes the following applications:

- :newspaper_roll: [FreshRSS](https://github.com/FreshRSS/FreshRSS) for feed reading
- :mag_right: [FullTextRSS](https://github.com/heussd/fivefilters-full-text-rss-docker/tree/master) for converting partial feeds into full feeds for use in FreshRSS
- :bird: [Nitter](https://github.com/zedeus/nitter) for providing RSS feeds for Twitter
- :alien: [Redlib](https://github.com/redlib-org/redlib) for browsing Reddit
- :card_index_dividers: [Planka](https://github.com/plankanban/planka) for task management
- :globe_with_meridians: [OpenResty](https://github.com/openresty/openresty) to serve FreshRSS and reverse proxy to other services with support for OIDC auth.
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

Deployed via Terraform, see: [batinicaz/app-server](https://github.com/batinicaz/app-server)

### Nitter Sessions

Nitter requires a real account due to Twitters API changes and removal of guest accounts.

To get the session configuration from your account you can use the script provided in the [Nitter wiki](https://github.com/zedeus/nitter/wiki/Creating-session-tokens).
