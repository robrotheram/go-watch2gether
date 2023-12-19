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

There is also a Docker-compose file avalible. 

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
