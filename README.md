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

Nitter requires a real account due to Twitters recent API changes and removal of guest accounts.

Thanks to [983](https://github.com/zedeus/nitter/issues/983#issuecomment-1688495284) you can generate a suitable token by setting the username/password and then using the Python snippet below:

```python
import requests

username=""
password=""

authentication = None

session = requests.session()

authorization_bearer = 'Bearer AAAAAAAAAAAAAAAAAAAAAFXzAwAAAAAAMHCxpeSDG1gLNLghVe8d74hl6k4%3DRUMF4xAQLsbeBhTSRrCiQpJtxoGWeyHrDb5te2jpGskWDFW82F'
guest_token = requests.post("https://api.twitter.com/1.1/guest/activate.json", headers={'Authorization': authorization_bearer}).json()['guest_token']

headers = {
    'X-Guest-Token': guest_token,
    'Content-Type': 'application/json',
    'Authorization':  authorization_bearer
}
session.headers.update(headers)

task1 = session.post('https://api.twitter.com/1.1/onboarding/task.json',
    params={
        'flow_name': 'login',
        'api_version': '1',
        'known_device_token': '',
        'sim_country_code': 'us'
    },
    json={
        "flow_token": None,
        "input_flow_data": {
            "country_code": None,
            "flow_context": {
                "referrer_context": {
                    "referral_details": "utm_source=google-play&utm_medium=organic",
                    "referrer_url": ""
                },
                "start_location": {
                    "location": "deeplink"
                }
            },
            "requested_variant": None,
            "target_user_id": 0
        }
    }
)

session.headers['att'] = task1.headers.get('att')
task2 = session.post('https://api.twitter.com/1.1/onboarding/task.json', json={
    "flow_token": task1.json().get('flow_token'),
    "subtask_inputs": [{
            "enter_text": {
                "suggestion_id": None,
                "text": username,
                "link": "next_link"
            },
            "subtask_id": "LoginEnterUserIdentifier"
        }
    ]
})

task3 = session.post('https://api.twitter.com/1.1/onboarding/task.json', json={
    "flow_token": task2.json().get('flow_token'),
    "subtask_inputs": [{
            "enter_password": {
                "password": password,
                "link": "next_link"
            },
            "subtask_id": "LoginEnterPassword"
        }
    ],
})

task4 = session.post('https://api.twitter.com/1.1/onboarding/task.json', json={
    "flow_token": task3.json().get('flow_token'),
    "subtask_inputs": [{
            "check_logged_in_account": {
                "link": "AccountDuplicationCheck_false"
            },
            "subtask_id": "AccountDuplicationCheck"
        }
    ]
}).json()

for t4_subtask in task4.get('subtasks', []):
    if 'open_account' in t4_subtask:
        authentication = t4_subtask['open_account']
        break
    elif 'enter_text' in t4_subtask:
        response_text = t4_subtask['enter_text']['hint_text']
        code = input(f'Requesting {response_text}: ')
        task5 = session.post('https://api.twitter.com/1.1/onboarding/task.json', json={
            "flow_token": task4.get('flow_token'),
            "subtask_inputs": [{
                "enter_text": {
                    "suggestion_id": None,
                    "text": code,
                    "link": "next_link"
                },
                "subtask_id": "LoginAcid"
            }]
        }).json()
        for t5_subtask in task5.get('subtasks', []):
            if 'open_account' in t5_subtask:
                authentication = t5_subtask['open_account']

print(authentication)
{
    'attribution_event': 'login',
    'known_device_token': 'XXXXXXXXXXXXXXXXXXXXXX',
    'next_link': {
        'link_id': 'next_link',
        'link_type': 'subtask',
        'subtask_id': 'SuccessExit'
    },
    'oauth_token': 'XXXXXXXXXXXXXXXXXXXXXX',
    'oauth_token_secret': 'XXXXXXXXXXXXXXXXXXXXXX',
    'user': {
        'id': 'XXXXXXXXXXXXXXXXXXXXXX',
        'id_str': 'XXXXXXXXXXXXXXXXXXXXXX',
        'name': 'XXXXXXXXXXXXXXXXXXXXXX',
        'screen_name': 'XXXXXXXXXXXXXXXXXXXXXX'
    }
}
```

This then enables you to use nitter on the guest_accounts branch provided the 30 day expiry is disabled ().