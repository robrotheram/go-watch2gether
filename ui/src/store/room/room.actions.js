import axios from 'axios'
import { CLEAR_ERROR, JOIN_SUCCESSFUL,ROOM_ERROR, LEAVE_SUCCESSFUL, GET_META_SUCCESSFUL, SEEK_TO_HOST, REJOIN_SUCCESSFUL, PROGRESS_UPDATE } from './room.types';
import store, {API_URL, WS_URL, history} from '../index'
import { connect } from '@giantmachines/redux-websocket';
import { send } from '@giantmachines/redux-websocket';

export const join = (room, user) => {
    return dispatch => {
        axios.post(API_URL+`room/`+room+`/join`, {"name":room, "username":user}).then(res => {
            console.log("Action", res)
            dispatch( {
                type: JOIN_SUCCESSFUL,
                room: room,
                user: user
            })
            dispatch(getMeta())
            console.log("WS_URL",WS_URL+"room/"+room+"/ws")
            store.dispatch(Connect(room))
            history.push('/room/'+room);
        }).catch(e => {
            console.log(e)
            history.push('/');
            dispatch( {
                type: ROOM_ERROR,
                error: "Oh Dear something happened when user: "+user+" tried to joining room"+room+" Error: "+e.response.data,
            })
        })
    }
}

export const Connect = (room) => {
    return connect(WS_URL+"room/"+room+"/ws")
}

export const reJoin = (room) => {
    return dispatch => {
        axios.get(API_URL+`room/`+room).then(res => {
            store.dispatch(Connect(room))
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
        axios.post(API_URL+`room/`+store.getState().room.name+`/leave`, {"name":store.getState().room.name, "username":store.getState().room.user.name}).then(res => {
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


export const getMeta = () => {
    return dispatch => {
        axios.get(API_URL+`room/`+store.getState().room.name).then(res => {
            console.log("Action", res)
            dispatch( {
                type: GET_META_SUCCESSFUL,
                payload: res.data,
            })
            
        }).catch(e => {
            console.log(e)
            dispatch( {
                type: ROOM_ERROR,
                error: e.response.data,
            })
        })
    }
}

export const updateSeek = (seek) => {
    store.dispatch( {
        type: PROGRESS_UPDATE,
        seek: seek,
    })
    // let evnt = {action: "ON_PROGRESS_UPDATE", user:store.getState().room.user, seek:seek}
    // store.dispatch(send(evnt))   
}

export const updateControls = (cntrls) => {
    let evnt = {action: "UPDATE_CONTROLS", user:store.getState().room.user, controls:cntrls}
    store.dispatch(send(evnt))   
}

export const updateHost = (host) => {
    let evnt = {action: "UPDATE_HOST", user:store.getState().room.user, host:host}
    store.dispatch(send(evnt))   
}

export async function isAlive() {
    let user= store.getState().room.user;
    user.current_video = store.getState().room.current_video;
    let evnt = {
        action: "USER_UPADTE", user:user, 
    }
    console.log("USER_UPADTE", evnt ); 
    return store.dispatch(send(evnt))   
}

export const handleFinish = () => {
    let evnt = {action: "HANDLE_FINSH", user:store.getState().room.user}
    store.dispatch(send(evnt))   
}


export const sinkToME = (seek) => {
    let evnt = {action: "SEEK_TO_ME", user:store.getState().room.user, seek:seek}
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
    let evnt = {action: "PLAYING", user: store.getState().room.user}
    store.dispatch(send(evnt))
}

export const pause = () => {
    let evnt = {action: "PAUSING", user: store.getState().room.user}
    store.dispatch(send(evnt))
}

export const updateQueue = (queue) => {
    let evnt = {action: "UPDATE_QUEUE", queue: queue, user:store.getState().room.user}
    store.dispatch(send(evnt))
}

export const nextVideo = () => {
    let evnt = {action: "NEXT_VIDEO", user:store.getState().room.user}
    store.dispatch(send(evnt))
}
