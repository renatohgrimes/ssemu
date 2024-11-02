# SSEMU

SuperSonic EMUlator

A backend emulator for [S4 League](https://s4league.fandom.com/wiki/S4_League) European Client from 2008.

## Disclaimer

The goal of this project was preservation, and was created **for eductational purposes** to demonstrate how a game server works.

Also, I wanted to learn more about Go, networking and OpenTelemetry.

I am not trying to compete against any official or private servers.

Many features are not implemented. You can find more feature complete emulators on the internet.

You can login, explore and play old tutorials.

No game files or resources are available in this source code.

You need a copy of a specific version of the game client. It is *hella* easy to find on the internet.

Remember, the game version is almost 20 years old. There is no garantee that it is safe against remote code execution.

This project was made just for fun, learning and to remember good times.

No copyright infringement intended.

- compatible with european client, version 0.8.19.26302, from September 30th 2008

## Running

Download the zip for your system from this repository and extract.

Copy all your client files inside `bin/client` directory.

Run the application.

- If you want to use OpenTelemetry infrastructure, you can check the docker-compose.yml file inside `deployments` directory or run `make compose`.

## Client Patching

Once you have a client copy, open the patcher website at `localhost:8000/web/patcher.html`.

Provide the `S4Client.exe` to the website and press patch.

Download the executable file into your game client folder and run the patched executable.

Also, make sure to change the IP1 parameter of version.ini file to 127.0.0.1:

```
[SERVER]
IP1 = 127.0.0.1
```

You need to do this step to be able to run the game properly.

## Account Creation

To create an account, use the following suffix: 00

Like so: `myuser00`

It will close your client and then you can login without the suffix like so: `myuser`

## Building

In linux environment, you can run `make release` and all files will be created. You will need `docker`.

## Developing

Run `make devenv` and point your IDE output and working directory to `bin` directory. Make sure to paste your client files inside `bin/client` directory.

### Testing

You can run `make tests` to run tests for this project.

### Linting

You can run `make lint` to run linting for this project.

## Credits & References

[Neowiz MUCA (former Pentavision)](https://muca.world/)

[wtfblub](https://github.com/wtfblub)
- [Pre-Seasons Emulator](https://github.com/WAZAAAAA0/FagNet)
- [Season 8 Emulator](https://gitlab.com/NetspherePirates/NetspherePirates)

[S4LeagueOpenSource](https://github.com/S4LeagueOpenSource)
- [Pre-Seasons Emulator](https://github.com/S4LeagueOpenSource/GodNet)

[Asiro](https://www.elitepvpers.com/forum/members/1494245-asiro.html)
- [Resource Files](https://www.elitepvpers.com/forum/s4-league-hacks-bots-cheats-exploits/395464-release-s4-league-resource-files.html#post3630606)

[0Harakiry](https://www.youtube.com/@0Harakiry)
- [Prototype Trailer](https://www.youtube.com/watch?v=dz5G7UVnk-4)
- [G-STAR 2007 GameShow Trailer](https://www.youtube.com/watch?v=8yiJEarb_A0)

Archive
- [Emulator Archive](https://archive.org/details/s4lserver-emuarchive)
- [Game Client Archive](https://archive.org/details/s4lgameclientarchives)
- [Hella Hell's YouTube Channel](https://www.youtube.com/@HellaHell)
