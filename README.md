# Tomato

Tomato is a command for running [pomodoro](https://en.wikipedia.org/wiki/Pomodoro_Technique) in background. It's designed mainly to stay in MacBook touchbar.

For example, my touchbar looks like this:

##### 1. Start

- **First button**: Tap to start.
- **Second button**: Tap to switch between **work** and **break** mode.

![Start](others/touchbar_1.png)

##### 2. Working (running)

- **First button**: Tap to pause/continue.
- **Second button**: Tap to skip the current interval.

![Working](others/touchbar_2.png)

##### 3. Short Break (running)

- **First button**: Tap to pause/continue.
- **Second button**: Tap to skip the current break interval.

![Short Break](others/touchbar_3.png)

### Modes:
- **Work**: The timer is running in working interval (25 mins).
- **Short Break**: The timer is running in short break interval (5 mins).
- **Long Break**: The timer is running in long break interval (15 mins). After a set of pomodoro (default to 4), a longer break is taken.

## Quick Start

1. Download the [prebuilt command](others/tomato) or [build from source](#build-from-source).

    You may need to run `chmod +x tomato` and put it in your `$PATH`.

2. Config BetterTouchTool as this screenshot (`port=12345`). Don't forget to click **Apply Changes**.

    ![](others/btt-settings.png)

3. Config touchbar buttons with these Apple scripts:
    - **Polling**: Display information on the first button.
    - **Start/Pause**: Run when the first button is tapped.
    - **Stop/Switch Mode**: Run when the second button is tapped.

    Or copy config for [the widget](https://raw.githubusercontent.com/ng-vu/tomato/master/others/btt-widget.json) and [the button](https://raw.githubusercontent.com/ng-vu/tomato/master/others/btt-button.json) then press `Command+V` inside BetterTouchTool.

4. Copy `UUID` of the widget

    ![](others/btt-uuid.png)

5. Start the `tomato` server. It will listen on `:12321` by default.

    ```
    tomato -uuid=[UUID] -port=12345
    ```
    
    If you have trouble running the command, see [this issue](https://github.com/ng-vu/tomato/issues/2).

## Usage

```
Tomato on TouchBar v1.2.0 (works with BetterTouchTool)

Default:
   tomato

With options:
   tomato -n=3 -colon=: -work=25m -short=300s -long=15m -listen=:12321

Send updates to BetterTouchTool:
   tomato -uuid=UUID -port=12345
   tomato -icon1=PATH_ICON1 -icon2=PATH_ICON2 -uuid=UUID -url=http://127.0.0.1:12345/update_touch_bar_widget/

Execute a command at the end of timer:
   tomato -command="terminal-notifier -title Pomodoro -message \"Hey, time is over\!\" -sound default"

Options:
  -async
    	Execute the command without waiting it to finish (use together with -command)
  -colon string
    	Custom separator (default ":")
  -colon-alt string
    	Alternative separator for break modes (default ":")
  -command string
    	Execute command at the end of timer
  -icon1 string
    	Icon for work (default red)
  -icon2 string
    	Icon for break session (default green)
  -listen string
    	Address to listen on (default ":12321")
  -long string
    	Long break interval (default "15m")
  -n int
    	Number of intervals between long break (default 4)
  -port string
    	BetterTouchTool port
  -short string
    	Short break interval (default "5m")
  -tick int
    	Duration in ms for sending updates (default 100) (default 100)
  -url string
    	URL to post update
  -uuid string
    	UUID of the widget
  -work string
    	Work interval (default "25m")
```

## Build from source

1. Install [Go](https://golang.org/doc/install)
2. `go build *.go`

## API

| API                                         | Sample Output               |Description
|---------------------------------------------|-----------------------------|-----------
| GET [/status](http://localhost:12321/status)| `[R] 17:43 1/3 work`        | Current status
| GET [/time](http://localhost:12321/time)    | `17:43`                     | Current timer
| POST /action/start                          | `17:43`        | Start/pause the current interval.
| POST /action/stop                           | `25:00` | Stop the current interval or switch mode.

### Output

1. State: `[S]` - stopped, `[R]` - running, `[P]` - paused.
2. Timer: `mm:ss` - work interval, `mmːss` - break interval.
3. Number of completed pomodoro in a set.
4. Mode: `work`, `short-break`, `long-break`.

### JSON

```bash
curl -H "Accept: application/json" http://localhost:12321/status
```

```
{"i":0,"mode":"work","n":4,"state":"[S]","timer":"25:00"}
```

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

## Notes

- [BetterTouchTool](https://www.boastr.net/) to customize the touchbar. It's an awesome app!
- [Automator](https://stackoverflow.com/questions/6442364/running-script-upon-login-mac) to start the script at login.

# License

- [MIT License](https://opensource.org/licenses/mit-license.php)
