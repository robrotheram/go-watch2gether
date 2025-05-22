import { DiscordGuilds, Media, Playlist, Settings, State, User } from "@/types";


export async function getController(id:string):Promise<State> {
    const response = await fetch(`/api/channel/${id}`);
    return response.json();
}


export async function addVideoController(id:string, video:string):Promise<Media> {
    const response = await fetch(`/api/channel/${id}/add`, {
            method: "PUT",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(video),
        });
    return response.json();    
}

export async function updateQueueController(id:string, video:Media[]):Promise<Media[]> {
    const response = await fetch(`/api/channel/${id}/queue`, {
        method: "POST", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(video),
    });
    return await response.json();
}

export async function playVideoController(id:string) {
    const response = await fetch(`/api/channel/${id}/play`, {
        method: "POST", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
    });
    return await response.json();
}

export async function pauseVideoController(id:string) {
    const response = await fetch(`/api/channel/${id}/pause`, {
        method: "POST", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
    });
    return await response.json();
}

export async function seekVideoController(id:string, seek:number) {
    const response = await fetch(`/api/channel/${id}/seek`, {
        method: "POST", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(Math.round(seek))
    });
    return await response.json();
}





export async function skipVideoController(id:string) {
    const response = await fetch(`/api/channel/${id}/skip`, {
        method: "POST", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
    });
    return await response.json();
}


export async function loopVideoController(id:string) {
    const response = await fetch(`/api/channel/${id}/loop`, {
        method: "POST", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
    });
    return await response.json();
}

export async function shuffleVideoController(id:string) {
    const response = await fetch(`/api/channel/${id}/shuffle`, {
        method: "POST", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
    });
    return await response.json();
}

export async function clearVideoController(id:string) {
    await fetch(`/api/channel/${id}/clear`, {
        method: "POST", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
    });
}

export async function getSettings():Promise<Settings> {
    const response = await fetch(`/api/settings`);
    return response.json();
}

export async function getGuilds():Promise<DiscordGuilds> {
        const response = await fetch(`/api/guilds`);
        return response.json();
}

export async function getUser():Promise<User> { 
    const response = await fetch(`/auth/user`);
    return response.json();
}

export async function getChannelPlayers(id:string) {
    const response = await fetch(`/api/channel/${id}/players`);
   return await response.json();
}

export async function getChannelPlaylists(id:string):Promise<Playlist[]> {
    
    const response = await fetch(`/api/channel/${id}/playlist`);
    return response.json();
}


export async function addVideoToPlaylist(video:string, playlistID:string):Promise<Media> {
    const response = await fetch(`/api/playist/${playlistID}/add`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(video),
    });
    return await response.json();
}

export async function updatePlaylist(playlist:Playlist) {
    const response = await fetch(`/api/playist/${playlist.id}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body:  JSON.stringify(playlist),
    });
    return await response.json();
}

export async function deletePlaylist(playlist:Playlist) {
    const response = await fetch(`/api/playist/${playlist.id}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
    });
    return await response.json();
}

export async function createPlaylist(id:string) {
    const playist = {
        "name":"new-playlists",
        "channel": id
    }
    const response = await fetch(`/api/playist`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body:  JSON.stringify(playist),
    });
    return await response.json();
}

export async function loadFromPlaylist(id:string, playlistID:string) {
    const response = await fetch(`/api/channel/${id}/add/playlist/${playlistID}`, {
        method: "PUT", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
    });
    return await response.json();
}

export const formatTime = (seconds:number) => {
    let iso = new Date(seconds / 1000000).toISOString()
    return iso.substring(11, iso.length - 5);
}