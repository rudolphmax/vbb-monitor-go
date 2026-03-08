# VBB Monitor

Monitoring VBB stops at home using the dedicated API.

Access to the VBB API is required for usage, refer to VBB for further information.

## Prerequisites

- [GO](https://go.dev/) >= 1.25.5
- [pre-commit](https://pre-commit.com)


## Building from source

To build the project from source, run the following command.

```bash
$ go build ./cmd/vbbmon
```

The built executable can then be found under `vbbmon`.


## Running

```bash
$ go run ./cmd/vbbmon
```

Or simply run the executable built previously:

```bash
$ ./vbbmon
```

Necessary information has to specified via a config file. Refer to the example in the repo (`config.json`) and the table below for guidance.

Command-line options can be viewed with `-help`.

### Config

"display": {
    "numLines": 8,
    "theme": {
      "font": {
        "sizeBase": 17,
        "sizeMedium": 15,
        "sizeSmall": 12
      },
      "globalForegroundColor": "#FFFFFF",
      "globalBackgroundColor": "#212121"
    }
  }

| Key                                   | Description                                                                    | Example                                                 |
| ------------------------------------- | ------------------------------------------------------------------------------ | --------------------------------------------------------|
| `departureFetchInterval`              | Time interval in which to re-fetch departures (in seconds).                    | `20`                                                    |
| `messageFetchInterval`                | Time interval in which to re-fetch (disruption) messages (in seconds).         | `120`                                                   |
| `api.base`                            | The base url of the HAFAS API (no trailing slash!)                             | `https://api.example.com/api`                           |
| `api.accessId`                        | Your HAFAS Access-ID                                                           | Refer to HAFAS documentation                            |
| `api.stops`                           | List of objects containing stop-information.                                   | `[{ ID: "HAFAS Stop ID" }]`                             |
| `api.stops[i].ID`                     | The ID of this stop.                                                           | Refer to HAFAS documentation                            |
| `api.stops[i].lines`                  | The lines to be fetched at this stop.                                          | Refer to HAFAS documentation                            |
| `api.stops[i].maxDepartures`          | How many departures to fetch for this stop.                                    | Refer to HAFAS documentation                            |
| `api.stops[i].timeOffset`             | Offset for the arrival of the departures at this stop in minutes from "now".   | `10` -> no arrivals fetched that arrive in < 10 minutes |
| `api.stops[i].direction`              | The ID of this stop.                                                           | Refer to HAFAS documentation                            |
| `display.numLines`                    | The number of lines displayed. Has no effect on fetching (see `maxDepartures`) | 9                                                       |
| `display.scrollSpeed`                 | The number of pixels to scroll the message bar by each frame                   | 5                                                       |
| `display.theme.font.sizeBase`         | Sets the size of large text.                                                   | 17                                                      |
| `display.theme.font.sizeMedium`       | Sets the size of medium text.                                                  | 15                                                      |
| `display.theme.font.sizeSmall`        | Sets the size of small text.                                                   | 12                                                      |
| `display.theme.globalForegroundColor` | The color of the foreground as hex string (text and light backgrounds)         | #FFFFFF                                                 |
| `display.theme.globalForegroundColor` | The color of the background as hex string (and dark text)                      | #000000                                                 |


## Development

Before beginning development, install the git hooks provided via `pre-commit`. Then, fill the config file and run the app.

```bash
$ pre-commit install
$ go run ./cmd/vbbmon
```
