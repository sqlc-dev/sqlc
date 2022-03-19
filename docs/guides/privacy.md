# Privacy and data collection

These days, it feels like every piece of software is tracking you. From your
browser, to your phone, to your terminal, programs collect as much data about
you as possible and send it off to the cloud for analysis.

We believe the best way to keep data safe is to never collect it in the first
place.

## Our Privacy Pledge

The `sqlc` command line tool does not collect any information. It
does not send crash reports to a third-party. It does not gather anonymous
aggregate user behaviour analytics.

No analytics. 
No finger-printing.
No tracking.

Not now and not in the future.

### Distribution Channels

We distribute sqlc using popular package managers such as
[Homebrew](https://brew.sh/) and [Snapcraft](https://snapcraft.io/). These
package managers and their associated command-line tools do collect usage
metrics.

We use these services to make it easy to for users to install sqlc. There will
always be an option to download sqlc from a stable URL.

## Hosted Services

We provide a few hosted services in addition to the sqlc command line tool.

### sqlc.dev

* Hosted on [GitHub Pages](https://pages.github.com/)
* Analytics with [Plausible](https://plausible.io/privacy-focused-web-analytics)

### docs.sqlc.dev

* Hosted on [Read the Docs](https://readthedocs.org/)
* Analytics with [Plausible](https://plausible.io/privacy-focused-web-analytics)

### play.sqlc.dev

* Hosted on [Heroku](https://heroku.com)
* Playground data stored in [Google Cloud Storage](https://cloud.google.com/storage)
  * Automatically deleted after 30 days

### app.sqlc.dev / api.sqlc.dev

* Hosted on [Heroku](https://heroku.com)
* Error tracking and tracing with [Sentry](https://sentry.io)
