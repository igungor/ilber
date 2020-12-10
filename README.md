# ilber

`ilber` is a multi-purpose Telegram bot.

`ilber` is currently running on Google Cloud Functions. The following
environment variables can be used to configure the application:

| KEY                           | TYPE          | DEFAULT | REQUIRED | DESCRIPTION                                                  |
|-------------------------------|---------------|---------|----------|--------------------------------------------------------------|
| ILBER_TOKEN                   | String        |         | true     | Telegram bot token received from Botfather                   |
| ILBER_DEBUG                   | TRUE or FALSE | FALSE   | false    | Enable debug logging                                         |
| ILBER_GOOGLE_API_KEY          | String        |         | false    | Google Custom Search API Key used for Google Search requests |
| ILBER_GOOGLE_SEARCH_ENGINE_ID | String        |         | false    | Google Search Engine ID, passed to Custom Search API         |
| ILBER_OPENWEATHERMAP_APP_ID   | String        |         | false    | OpenWeather ID used for weather requests                     |


## license

MIT. See LICENSE.
