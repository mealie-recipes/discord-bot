package main

import "fmt"

func msg(msg string) string {
	return fmt.Sprintf("%s\n%s", msg, BotFooter)
}

var BotFooter = `Beep Boop, I'm a bot 🤖`

var DefaultCredentials = `
Having trouble logging in? Try the default username and password

Username: changeme@email.com
Password: MyPassword

Btw, this is in the docs 👇

[installation-checklist/#step-3-startup](https://nightly.mealie.io/documentation/getting-started/installation/installation-checklist/#step-3-startup)
`

var V1MigrationLinks = `
Looking for information on Migration to v1?

[Beta Release Discussion](https://github.com/hay-kot/mealie/discussions/1073)
[v0.5.x to v1.0.0 Migration Guide](https://nightly.mealie.io/documentation/getting-started/migrating-to-mealie-v1/)
`

var DockerTags = `
Not sure which tag to pull?

**frontend-<version>** - Latest release build of the frontend server
**api-<version>** - Latest release build of the backend API server

*See the [docker compose example](https://nightly.mealie.io/documentation/getting-started/installation/sqlite/) for most current tags for running Mealie*

**frontend-nightly** - Nightly build fresh off the mealie-next branch (if you're feeling brave)
**api-nightly** - Nightly build fresh off the mealie-next branch (if you're feeling brave)

❗Depreciated Tags❗

**latest** - Legacy tag used for the previously combined servers last updated (v0.5.5)

---

[See all available tags on dockerhub](https://hub.docker.com/r/hkotel/mealie/tags)
`

var DockerFAQ = `
Having trouble with docker?

1. Is your API_URL set correctly on the Frontend container?
2. Have you verified your volumes are configured correctly?
3. Have you reviewed the installation checklist?

**Links**
[Docker Compose Example](https://nightly.mealie.io/documentation/getting-started/installation/sqlite/)
[Checklist](https://nightly.mealie.io/documentation/getting-started/installation/installation-checklist/)
[Docker Tags and Diagrams](https://nightly.mealie.io/documentation/getting-started/installation/installation-checklist/#docker-diagram)
`

var TokenTime = `
'TOKEN_TIME' is the time in hours that an auth token is valid for, it's default is set to 48 hours.

You can change this value by setting the environment variable 'TOKEN_TIME' to a number of hours you want the token to be valid for.

**Example:**

environment:
  - TOKEN_TIME=999

This will set the token to be valid for 999 hours.
`
