# haste-ya-bin

For yabinning your hastebin.

_hasteyabin_ automates uploading your Hastebin file backup to Yabin.

## How To

### IMPORTANT - KEY LENGTH

YABin, by default, uses 5 character long paste keys. Where Haste uses 32 character ones, these will get stripped / won't work properly without a database modification on the YABin side.

To make this modification to the necessary column, run this against your YABin Postgres database:

```postgresql
ALTER TABLE Paste
ALTER COLUMN key TYPE VARCHAR(32)
```

>[!IMPORTANT]
> This _SHOULDN'T_ break anything in YABin, but you're still making an unsupported modification to its database.
> As such, you're on your own if stuff breaks, your house catches fire, or smoke starts billowing from your fingernails. You have been warned.

### Usage

Call like so:
```hasteyabin -api-url https://my.yabin.install ./path/to/my/hastebin/dump```

You can also call using the `-api-token` flag, where the token is the content of your YABin instance's `token` cookie (once you're logged in).
Fetching this is outside the scope of this readme - but feel free to raise an issue if you don't understand this.

>[!TIP]
> You'll **need** to use your token if you have the `PUBLIC_CUSTOM_PATHS_ENABLED` envvar sat on your YABin installation, otherwise you'll end up with
> a YABin full of keys that don't match your backup.

_hasteyabin_ will try to automatically determine the programming language of your uploaded files.
If you'd prefer that it _didn't_ do this, call with the `-no-autodetect` flag - all pastes will be set to plaintext.