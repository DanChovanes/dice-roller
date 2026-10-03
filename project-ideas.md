# Randomness Source Ideas

These are candidate ways to give a roll a legible story, not claims of cryptographic unpredictability. Public data can be stale, predictable, or unavailable; every network-backed source needs a timeout and a fallback.

| Idea | Measures and squash | Explain-flag story | Needs and fragility | Feel |
| --- | --- | --- | --- | --- |
| The Earth Shook While You Blinked | Use the USGS worldwide “all earthquakes in the past hour” feed. Hash the event snapshot and its timestamp, then map that to the die. | “USGS recorded 18 earthquakes worldwide in the past hour; the planet contributed a d20.” | Network required. The public GeoJSON feed is updated every minute; a quiet hour or outage may leave the data stale or unavailable. | 5/5: chunky and easy to explain |
| The Sky’s Pressure Cooker | Use current cloud cover and wind gusts from Open-Meteo for a configured location. Hash the returned observation and timestamp, then map that to the die. | “At [location], the sky was 76% clouded and gusts reached 31 km/h; the weather tipped your d20.” | Network and location required. The free API is non-commercial, has published call limits, and provides data under CC BY 4.0. Weather values are model data, not direct instrument readings. | 4/5: tactile, but location-dependent |
| Wikipedia Attention Tides | Use Wikimedia pageview totals or the most-viewed pages for a completed day. Hash the page counts and date, then map that to the die. | “Yesterday, Wikipedia’s top pages drew 9.2 million views; the crowd chose your d20.” | Network required. Requests need a descriptive `User-Agent`; observe the API’s rate limits. Daily data is delayed, not live. | 3/5: wonderfully odd, less immediate |
| The Font Renderer’s Revenge | Time a bounded render of a sample glyph using installed fonts; hash the timing summary and timestamp, then map that to the die. | “Your system took 43 milliseconds to render ‘A’ across its typefaces; the fonts have spoken.” | No network required. Rendering behavior and available fonts vary by operating system and installed fonts. | 4/5: weird local flavor; keep it |

## Public Source References

- [USGS GeoJSON feeds](https://earthquake.usgs.gov/earthquakes/feed/v1.0/geojson.php): public earthquake summary feeds, including a worldwide past-hour feed updated every minute.
- [Open-Meteo API documentation](https://open-meteo.com/en/docs) and [terms](https://open-meteo.com/en/terms): current weather variables and free non-commercial usage limits and conditions.
- [Wikimedia pageview API reference](https://doc.wikimedia.org/generated-data-platform/aqs/analytics-api/reference/page-views.html) and [access policy](https://doc.wikimedia.org/generated-data-platform/aqs/analytics-api/documentation/access-policy.html): pageview endpoints, client identification, and rate-limit guidance.