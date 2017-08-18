# TOMATO

Tomato is a command for running pomodoro in background. It's designed mainly to stay in MacBook touchbar.

For example, my touchbar looks like below:

##### 1. Start

- **First button**: Tap to start.
- **Second button**: Tap to switch to break mode.

![Start](others/touchbar_1.png)

##### 2. Working

- **First button**: Tap to pause.
- **Second button**: Tap to skip the current interval.

![Working](others/touchbar_2.png)

##### 3. Short Break

- **First button**: Tap to pause.
- **Second button**: Tap to skip the current break interval.

![Short Break](others/touchbar_3.png)

### Modes:
- **Work**: The timer is running in working interval (25 mins).
- **Short Break**: The timer is running in short break interval (5 mins).
- **Long Break**: The timer is running in long break internal (15 mins).


## Quick Start

1. Install [Go](https://golang.org/doc/install)
2. `go build tomato.go`
3. Put the command `tomato` in your $PATH
4. Start the `tomato` server. It will listen on `:12321`.
5. Config touchbar buttons with these Apple scripts:
    
    - **Polling**: Display information on the first button.
    - **Start/Pause**: Run when the first button is tapped.
    - **Stop/Switch Mode**: Run when the second button is tapped.

## API

| API                                         | Sample Output               |Description
|---------------------------------------------|-----------------------------|-----------
| GET [/status](http://localhost:12321/status)| `[R] 17:43 1/3 work`        | Current status
| GET [/time](http://localhost:12321/time)    | `17:43`                     | Current timer
| POST /action/start                          | `[R] 17:43 1/3 work`        | Start/pause the current interval.
| POST /action/stop                           | `[S] 25:00 1/3 short-break` | Stop the current interval or switch mode.

## AppleScript

### 1. Polling

```applescript
try
	set reqURL to "http://localhost:12321/time"
	do shell script "curl " & quoted form of reqURL
on error
	return "00:00"
end try
```

### 2. Start / Pause

```applescript
try
	set reqURL to "http://localhost:12321/action/start"
	do shell script "curl -X POST " & quoted form of reqURL
end try
```

### 3. Stop / Switch Mode

```applescript
try
	set reqURL to "http://localhost:12321/action/stop"
	do shell script "curl -X POST " & quoted form of reqURL
end try
```
