
![logo](logo.jpg)

Bot for transferring chat messages between messengers (Utopia, Telegram)

The finished compiled version can be found on the :arrow_down: [releases](releases) page.

# :robot: use cases

1. To unite different communities.
2. To work together from different devices and messengers with your team.

# :whale2: run with docker

Get the Utopia docker image:

```bash
docker pull uto9234/utopia-api
```

start new Utopia container:

```bash
docker container run -d uto9234/utopia-api
```

*Hint: `-d` is optional parameter means to run the container in the background.*

To find the IP of Utopia client host from the docker container, do the following:

```bash
docker ps
```

copy `container ID` & paste to the following:
(example for container ID 5d2df19066ac)

```bash
docker container inspect 5d2df19066ac
```

Find the IP in the `IPAddress` field.

Next, go to the "manual build" step or download the finished [release](releases) of the bot.


# :large_blue_circle: manual build

```bash
git clone https://github.com/Sagleft/utopia-telegram-bridge bridge-bot && cd bridge-bot
go build
cp config.example.json config.json
```

Then update the parameters in `config.json` file.

Use the following parameters to connect:

* API host: host or IP.
* API port: `22825`
* API protocol: `http`
* API token: `FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF`

To start the bot:

```bash
./bot
```

# :information_source: useful links

* [How to get the API key U manually](https://udocs.gitbook.io/utopia-api/utopia-api/how-to-enable-api-access)
* [Forum thread](https://talk.u.is/viewtopic.php?pid=5253)
* [UDocs](https://udocs.gitbook.io/utopia-api/)
