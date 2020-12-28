import axios from 'axios'
import { CLEAR_ERROR, JOIN_SUCCESSFUL,ROOM_ERROR, LEAVE_SUCCESSFUL, GET_META_SUCCESSFUL, SEEK_TO_HOST, REJOIN_SUCCESSFUL, PROGRESS_UPDATE } from './room.types';
import store, {API_URL, WS_URL, history} from '../index'
import { connect } from '@giantmachines/redux-websocket';
import { send } from '@giantmachines/redux-websocket';

export const join = (roomid, room, user) => {
    return dispatch => {
        axios.post(API_URL+`room/join`, {"id": roomid, "name":room, "username":user}).then(res => {
            console.log("Action", res)
            dispatch( {
                type: JOIN_SUCCESSFUL,
                room: res.data.room_id,
                user: res.data.user
            })
            dispatch(getMeta(res.data.room_id))
            console.log("WS_URL",WS_URL+"room/"+room+"/ws")
            store.dispatch(Connect(res.data.room_id, res.data.user.id))
            history.push('/room/'+res.data.room_id);
        }).catch(e => {
            console.log(e)
            history.push('/');
            if (e.response === undefined){
                dispatch( {
                    type: ROOM_ERROR,
                    error: "Oh Dear something happened when user: "+user+" tried to joining room: "+room+" Unable to contact server",
                })
            }else{
                dispatch( {
                    type: ROOM_ERROR,
                    error: "Oh Dear something happened when user: "+user+" tried to joining room: "+room+" Error: "+e.response.data,
                })
            }
        })
    }
}

export const Connect = (room, user) => {
    return connect(WS_URL+"room/"+room+"/ws"+"?token="+user)
}

export const reJoin = (room) => {
    return dispatch => {
        axios.get(API_URL+`room/`+room).then(res => {
            store.dispatch(Connect(room, store.getState().room.user.name))
            dispatch( {
                type: REJOIN_SUCCESSFUL,
                payload: res.data,
            })
        }).catch(e => {
            console.log(e)
            history.push('/');
            dispatch( {
                type: ROOM_ERROR,
                error: e.response.data,
            })
        })
    }
}

export const leave = (room, user) => {
    return dispatch => {
        axios.post(API_URL+`room/leave`, {
            "id":store.getState().room.id, 
            "name":store.getState().room.name, 
            "username":store.getState().user.name
        }).then(res => {
            console.log("Action", res)
            dispatch( {
                type: LEAVE_SUCCESSFUL,
                room: room,
                user: user
            })
            history.push('/');
        }).catch(e => {
            dispatch( {
                type: ROOM_ERROR,
                error: e.response.data,
            })
            history.push('/');
        })
    }
}


export const getMeta = (id) => {
    console.log("GETTING META", id)
    return dispatch => {
        axios.get(API_URL+`room/`+id).then(res => {
            console.log("Action", res)
            dispatch( {
                type: GET_META_SUCCESSFUL,
                payload: res.data,
            })
            
        }).catch(e => {
            console.log(e)
            dispatch( {
                type: ROOM_ERROR,
                error: "Connection Issue to the server",
            })
        })
    }
}

export const updateSeek = (seek) => {
    store.dispatch( {
        type: PROGRESS_UPDATE,
        seek: seek,
    })
    // let evnt = {action: "ON_PROGRESS_UPDATE", watcher:  store.getState().user, seek:seek}
    // store.dispatch(send(evnt))   
}

export const updateSettings = (cntrls, auto_skip) => {
    let evnt = {action: "UPDATE_SETTINGS", watcher:  store.getState().user, settings: {controls:cntrls, auto_skip:auto_skip}}
    store.dispatch(send(evnt))  
}


export const updateHost = (host) => {
    let evnt = {action: "UPDATE_HOST", watcher:  store.getState().user, host:host}
    store.dispatch(send(evnt))   
}

export async function isAlive() {
    let user = store.getState().user;
    let video = store.getState().video

    let evnt = {
        action: "USER_UPDATE",
        watcher: {
            id: user.id,
            seek: user.seek,
            video_id: video.id

        } 
    }
    return store.dispatch(send(evnt))   
}

export const handleFinish = () => {
    let evnt = {action: "HANDLE_FINSH", watcher:  store.getState().user}
    store.dispatch(send(evnt))   
}


export const sinkToME = (seek) => {
    let evnt = {action: "SEEK_TO_ME", watcher:  store.getState().user, seek:seek}
    store.dispatch(send(evnt))   
}

export const sinkToHost = () => {
    //console.log("SINGK")
    axios.get(API_URL+`room/`+store.getState().room.name).then(res => {
        console.log("SINasdasdasdsadasdGK")
        store.dispatch( {
            type: SEEK_TO_HOST,
            seek: res.data.seek,
        })
    });
}

export const clearError = () => {
    console.log("Action", CLEAR_ERROR)
    return {
        type: CLEAR_ERROR,
    };
}

export const play = () => {
    let evnt = {action: "PLAYING", watcher: store.getState().user}
    store.dispatch(send(evnt))
}

export const pause = () => {
    let evnt = {action: "PAUSING", watcher:store.getState().user}
    store.dispatch(send(evnt))
}

export const updateQueue = (queue) => {
    let evnt = {action: "UPDATE_QUEUE", queue: queue, watcher:store.getState().user}
    store.dispatch(send(evnt))
}

export const updateLocalQueue = (queue) => {
    store.dispatch({
        type: "LOCAL_QUEUE_UPDATE",
        queue: queue,
        watcher:store.getState().user
    })
}


export const nextVideo = () => {
    let evnt = {action: "NEXT_VIDEO", watcher:store.getState().user}
    store.dispatch(send(evnt))
}
