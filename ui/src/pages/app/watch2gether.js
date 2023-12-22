export const getRoomId = () => {
    const regex = /\/app\/(\d+)/;
    const match = window.location.pathname.match(regex);
    if (match ==undefined) {
        return ""
    }
    const id = match[1];
    return id
}

export async function getController() {
    const response = await fetch(`/api/channel/${getRoomId()}`);
    const jsonData = await response.json();
    if (!response.ok){
        throw jsonData.message
    }
    return jsonData
}

export async function createController() {
    const response = await fetch(`/api/channel/${getRoomId()}`,{
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        }
    });
    const jsonData = await response.json();
    if (!response.ok){
        throw jsonData.message
    }
    return jsonData
}


export function getSocket() {
    let room = getRoomId()
    let socket = ((window.location.protocol === "https:") ? "wss://" : "ws://") + window.location.host + "/api/channel/"+room+"/ws";
    return socket
}


export async function addVideoController(video) {
    let response;
    try{
        response = await fetch(`/api/channel/${getRoomId()}/add`, {
            method: "PUT", // or 'PUT'
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(video),
        });
    }catch(e){
        throw e
    }

    const jsonData = await response.json();
    if (!response.ok){
        throw (jsonData.message)
    }
    return jsonData
    
}

export async function updateQueueController(video) {
    const response = await fetch(`/api/channel/${getRoomId()}/queue`, {
        method: "POST", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(video),
    });
    const jsonData = await response.json();
    if (!response.ok){
        throw jsonData.message
    }
    return jsonData
}

export async function playVideoController(video) {
    const response = await fetch(`/api/channel/${getRoomId()}/play`, {
        method: "POST", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
    });
    // const jsonData = await response.json();
    // if (!response.ok){
    //     throw jsonData.message
    // }
    // return jsonData
}

export async function pauseVideoController(video) {
    const response = await fetch(`/api/channel/${getRoomId()}/pause`, {
        method: "POST", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
    });
    // const jsonData = await response.json();
    // if (!response.ok){
    //     throw jsonData.message
    // }
    // return jsonData
}





export async function skipVideoController(video) {
    const response = await fetch(`/api/channel/${getRoomId()}/skip`, {
        method: "POST", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
    });
    const jsonData = await response.json();
    if (!response.ok){
        throw jsonData.message
    }
    return jsonData
}


export async function loopVideoController(video) {
    const response = await fetch(`/api/channel/${getRoomId()}/loop`, {
        method: "POST", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
    });
    const jsonData = await response.json();
    if (!response.ok){
        throw jsonData.message
    }
    return jsonData
}

export async function shuffleVideoController(video) {
    const response = await fetch(`/api/channel/${getRoomId()}/shuffle`, {
        method: "POST", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
    });
    const jsonData = await response.json();
    if (!response.ok){
        throw jsonData.message
    }
    return jsonData
}


export async function getSettings() {
    const response = await fetch(`/api/settings`);
    const jsonData = await response.json();
    if (!response.ok){
        throw jsonData.message
    }
    return jsonData
}

export async function getGuilds() {
    
        const response = await fetch(`/api/guilds`);
        const jsonData = await response.json();
        if (!response.ok){
            throw jsonData.message
        }
        return jsonData
    
}

export async function getUser() {
    const response = await fetch(`/auth/user`);
    const jsonData = await response.json();
    if (!response.ok){
        throw jsonData.message
    }
    return jsonData
}


export async function getChannelPlaylists() {
    
    const response = await fetch(`/api/channel/${getRoomId()}/playlist`);
    const jsonData = await response.json();
    if (!response.ok){
        throw jsonData.messages
    }
    return jsonData.sort((a, b) => { return a.name.localeCompare(b.name)});
}

export async function updatePlaylist(playlist) {
    const response = await fetch(`/api/playist/${playlist.id}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body:  JSON.stringify(playlist),
    });
    const jsonData = await response.json();
    if (!response.ok){
        throw jsonData.message
    }
    return jsonData
}

export async function deletePlaylist(playlist) {
    const response = await fetch(`/api/playist/${playlist.id}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
    });
    if (!response.ok){
        const jsonData = await response.json();
        throw jsonData.message
    }
}

export async function createPlaylist() {
    const playist = {
        "name":"new-playlists",
        "channel": getRoomId()
    }
    const response = await fetch(`/api/playist`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body:  JSON.stringify(playist),
    });
    const jsonData = await response.json();
    if (!response.ok){
        throw jsonData.message
    }
    return jsonData
}

export async function loadFromPlaylist(playlistID) {
    const response = await fetch(`/api/channel/${getRoomId()}/add/playlist/${playlistID}`, {
        method: "PUT", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
    });
    const jsonData = await response.json();
    if (!response.ok){
        throw jsonData.message
    }
    return jsonData
}