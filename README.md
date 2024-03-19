# Watch2gether
Yet another video sync website it currently support Youtube, DailyMotion Vimeo Soundcloud and Videos hosted on your own fileserver that you totally legally own😉

## Why? 
This project spun out of my quest to find something that would allow my group of friends watch vidoes together from the random collection of YouTube clips to long movies. All other services did not support playing mp4 files or had features that were not required for example chat. This is little webapp does not need API keys to be set up and users do not need to create any accouts. Just create a room share a link and people can join. 

## Features
### Listen In discord
![Screenshot 1](docs/discord.png)
Bot only supports slash commands
Full list of commands 
```
- /join : Join Bot to a voice channel 
- /leave : Disconnect Bot from channel 

- /pause : Pause Video 
- /skip : Skip to next video in the Queue 
- /play : Play Video 
- /stop : Stop Video 
- /list : List videos in the Queue 
- /add <video> : Add Video to Queue 

- /playlist load <name>
- /now : Current Status of what is playing 

```

### Supported Media files
Watch2gether has been tests with the following URL's other sites may work in the UI but not via the bot. 

- Youtube: Single Vidoe and Playlists
- Podcast-RSS: Currently Supports playing the lastest from the feed
- Icecast streams: Radio m3u streams can be streemed into a room 
- MP4 Video files
- MP3 Files: Note can be used also for podcast if you know the eppisode file




### Playlists
![Screenshot 1](docs/playlists.png)
Watch2gether can save custom playlists without having them public in youtube. 



## Installing
This application packaged as a docker container. 
You can run it with 

```
docker run -d -p 8080:8080 robrotheram/watch2gether
```
### Replacing Values in `app.sample.env` with Discord Bot and OAuth Client Credentials

Follow these steps to replace the placeholder values in the `app.sample.env` file with the actual credentials obtained from creating a Discord bot and OAuth2 client:

1. **Create a Discord Bot:**

   - Go to the [Discord Developer Portal](https://discord.com/developers/applications).
   - Click on the "New Application" button.
   - Enter a name for your application.
   - Navigate to the "Bot" tab on the left sidebar.
   - Click on "Add Bot" and confirm.
   - Copy the generated Bot Token.

2. **Create an OAuth2 Client:**

   - In the [Discord Developer Portal](https://discord.com/developers/applications), go to your application settings.
   - Navigate to the "OAuth2" tab on the left sidebar.
   - Under "Redirects", add a redirect URL (e.g., `http://localhost:8080/auth/callback`).
   - Copy the generated OAuth2 Client ID.
   - Copy the generated OAuth2 Client Secret.

3. **Replace Values in `app.sample.env`:**

   - Copy the `app.sample.env` file and rename it to `app.env`.
   - Replace the placeholder values (`CHANGEME`) with the corresponding values obtained from creating the Discord bot and OAuth2 client:
     ```plaintext
     DISCORD_TOKEN=YOUR_DISCORD_BOT_TOKEN
     DISCORD_CLIENT_ID=YOUR_OAUTH2_CLIENT_ID
     DISCORD_CLIENT_SECRET=YOUR_OAUTH2_CLIENT_SECRET
     SESSION_SECRET=ANY_RANDOM_STRING_FOR_SESSION_SECURITY
     ```
     Ensure to replace `YOUR_DISCORD_BOT_TOKEN`, `YOUR_OAUTH2_CLIENT_ID`, and `YOUR_OAUTH2_CLIENT_SECRET` with the actual values obtained from the Discord Developer Portal.

4. **Save Changes:**

   - Save the `app.env` file.

#### Configuring `docker-compose.yml`

If you are using Docker, you can configure your `docker-compose.yml` file to include the necessary environmental variables. Here's an example:

```yaml
    environment:
      - DISCORD_TOKEN=YOUR_DISCORD_BOT_TOKEN
      - DISCORD_CLIENT_ID=YOUR_OAUTH2_CLIENT_ID
      - DISCORD_CLIENT_SECRET=YOUR_OAUTH2_CLIENT_SECRET
      - SESSION_SECRET=ANY_RANDOM_STRING_FOR_SESSION_SECURITY
```

### Reverse Proxy

For running behind a proxy you will need to forward websoctes as well as http. Below is a sample nginx configuration

```
	server {
        server_name watch2gether.<YOUR DOMAIN>;
        listen 80;
        location / {
        proxy_set_header        Host $host;
        proxy_set_header        X-Real-IP $remote_addr;
        proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header        X-Forwarded-Proto $scheme;
		proxy_set_header 		Upgrade $http_upgrade;
    	proxy_set_header        Connection "upgrade";
        proxy_pass          	http://127.0.0.1:8080;
        proxy_read_timeout  90;
        }
    }
```


## Building
This project uses a server wiritten in go with a react frontend. 
Built with:

go 1.15+ 

Yarn v1.22.10

UI Framework:

Antd v4.9+

There is a handy make file that will build the server, ui and container. 
```
make build
```

# Screenshots
![Screenshot 1](docs/homepage.png)
![Screenshot 2](docs/login.png)
