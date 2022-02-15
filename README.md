# Contributor Shoutouter

This is a small program that will look at Pull Requests on the various Velero GitHub repositories, and assemble their titles into Markdown for acknowledgement in the weekly community meeting.

To use it, run:

```shell
    
    go run . --config my_file.yaml --token <GITHUB API TOKEN>
```

Optionally, you can specify a number of days to find PRs within. To filter by PRs in the last month:

```shell
    go run . --config my_file.yaml --token <GITHUB API TOKEN> --days 30
```

On macOS, you can also copy them directly to your clipboard.

```shell
    ./velero-shoutouts | pbcopy
```

## Making a config file

You can make a config file with YAML that specifies the following values:


```yaml
---
# Name of the shoutouts you're doing
name: velero
# GitHub organization to look in. All teams and repos must be within this org
org: vmware-tanzu
# List of repos to check
repos:
  - velero
  - velero-plugin-for-aws
  - velero-plugin-for-gcp
  - velero-plugin-for-microsoft-azure
  - helm-charts
# Core devs, if they are not in a team (optional)
devs:
  - ashish-amarnath
  - carlisia
  - jonasrosland
  - michmike
  - nrb
  - zubron
  - dsu-igeek
# GitHub teams to filter by.
teams:
  - velero-authors

```
