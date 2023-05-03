const getRoomId = () => {
    const room = window.location.pathname.replace("/app/","")
    return room
}

export async function getController() {

    const response = await fetch(`/api/channel/${getRoomId()}`);
    const jsonData = await response.json();
    if (!response.ok){
        throw jsonData.message
    }
    return jsonData
}

export async function addVideoController(video) {
    const response = await fetch(`/api/channel/${getRoomId()}/add`, {
        method: "PUT", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(video),
    });
    const jsonData = await response.json();
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
    return jsonData
}

export async function playVideoController(video) {
    const response = await fetch(`/api/channel/${getRoomId()}/play`, {
        method: "POST", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
    });
    const jsonData = await response.json();
    return jsonData
}

export async function pauseVideoController(video) {
    const response = await fetch(`/api/channel/${getRoomId()}/pause`, {
        method: "POST", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
    });
    const jsonData = await response.json();
    return jsonData
}

export async function skipVideoController(video) {
    const response = await fetch(`/api/channel/${getRoomId()}/skip`, {
        method: "POST", // or 'PUT'
        headers: {
          "Content-Type": "application/json",
        },
    });
    const jsonData = await response.json();
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
    return jsonData
}


export async function getGuilds() {
    try{
        const response = await fetch(`/api/guilds`);
        const jsonData = await response.json();
        return jsonData
    }catch(e){
        return []
        console.log(e)
    }
    
}

export async function getUser() {
    const response = await fetch(`/auth/user`);
    const jsonData = await response.json();
    if (!response.ok){
        throw jsonData.message
    }
    return jsonData
}
