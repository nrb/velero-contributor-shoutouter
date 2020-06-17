# Velero Contributor Shoutouter

This is a small program that will look at Pull Requests on the various Velero GitHub repositories, and assemble their titles into Markdown for acknowledgement in the weekly community meeting.

To build it, run:

```shell
    go build -o velero-shoutouts
```

Then, run it with:

```shell
    ./velero-shoutouts
```

On macOS, you can also copy them directly to your clipboard.

```shell
    ./velero-shoutouts | pbcopy
```
