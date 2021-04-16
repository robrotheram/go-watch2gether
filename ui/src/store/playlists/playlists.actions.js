import axios from 'axios'
import {API_URL, history} from '../index'

export const getPlaylists = (room_id) => {
    return dispatch => {
        axios.get(API_URL+"room/"+room_id+"/playlist").then(res => {
            console.log("auth", res.data)
            dispatch({
                type: "GET_PLAYLISTS",
                playlists: res.data
            })
        }).catch(e => {
            console.log("Unathoiriz")
        });
    }
}
export const createPlaylists = (room_id, playlist) => {
    return dispatch => {
        axios.put(API_URL+"room/"+room_id+"/playlist", playlist).then(res => {
            console.log("auth", res.data)
            dispatch({
                type: "UPDATE_COMPLETE",
            })
            dispatch(getPlaylists(room_id))
        }).catch(e => {
            console.log("Unathoiriz")
        });
    }
}

export const updatePlaylists = (room_id, playlist) => {
    return dispatch => {
        axios.post(API_URL+"room/"+room_id+"/playlist/"+playlist.id, playlist).then(res => {
            console.log("auth", res.data)
            dispatch({
                type: "UPDATE_COMPLETE",
            })
            dispatch(getPlaylists(room_id))

        }).catch(e => {
            console.log("Unathoiriz")
        });
    }
}
export const deletePlaylists = (room_id, id) => {
    return dispatch => {
        axios.delete(API_URL+"room/"+room_id+"/playlist/"+id).then(res => {
            console.log("auth", res.data)
            dispatch({
                type: "DELETE_COMPLETE",
            })
            dispatch(getPlaylists(room_id))

        }).catch(e => {
            console.log("Unathoiriz")
        });
    }
}

export const loadPlaylists = (room_id, id) => {
    return dispatch => {
        axios.get(API_URL+"room/"+room_id+"/playlist/"+id+"/load").then(res => {
            console.log("auth", res.data)
            dispatch({
                type: "LOAD_COMPLETE",
            })
        }).catch(e => {
            console.log("Unathoiriz")
        });
    }
}