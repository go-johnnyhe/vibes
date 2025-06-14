# 🌤️ Vibes

## Another weather app, seriously?!

> **Note:** Despite the name, this project is NOT vibe-coded

Instead of "15°C with 40% precipitation", you get:

```
$ vibes
Here's the current weather condition report for Seattle, WA:
Classic hoodie/light jacket zone
temp will be around the same in the next 4 hours
Current temperature: 13.2°C
Might rain in 2 hours - maybe keep a jacket handy.
```

Weather apps show numbers. Vibes tells you what to wear, not just the temperature.


Built with Go. Powered by Open-Meteo API. **Zero API keys required.**

## Cool Features
- 📍 **Auto-location** - Uses your IP to find you
- 🎨 **Color-coded temps** - Blue for freezing, red for hot, you get it
- 🌧️ **Rain alerts** - Tells you when to grab an umbrella
- 🌡️ **Temperature advice** - "T-shirt weather" vs "Bundle up"
- ⏰ **Custom forecasts** - Check the next hour or next week

## Installation


```bash
go install github.com/go-johnnyhe/vibes@latest
```

Or clone and build:
```bash
git clone https://github.com/go-johnnyhe/vibes.git
cd vibes
go build
```

## Usage

Quick check:
```bash
vibes
```

Planning your day:
```bash
vibes -d 8 # or "vibes --hours 8"
```

Metric system:
```bash
vibes --unit celsius
# or just
vibes -u c
```

## Commands

- `vibes` - Current weather + next 4 hours
- `vibes --hours 24` - Full day forecast  
- `vibes --unit fahrenheit` - For the Americans
- `vibes -u c -d 12` - Celsius, next 12 hours

## The Vibe

- **Freezing** (<5°C/41°F): "Bundle up"
- **Cold** (5-10°C/41-50°F): "Proper jacket weather"
- **Cool** (10-15°C/50-59°F): "Classic hoodie zone"
- **Mild** (15-20°C/59-68°F): "Maybe just a light layer"
- **Warm** (>20°C/68°F): "T-shirt weather!"

## Built With

- Go + [Cobra](https://github.com/spf13/cobra) for the CLI magic
- [Open-Meteo](https://open-meteo.com/) for weather data (no API key needed!)
- [ipinfo.io](https://ipinfo.io/) for location detection
- [fatih/color](https://github.com/fatih/color) for the pretty colors

## Contributing

Got ideas? Found a bug? PRs welcome. If you like this project or have any suggestions, I'd love to know!

## License

MIT - Do whatever

---

*Weather apps don't have to be boring* ✨
