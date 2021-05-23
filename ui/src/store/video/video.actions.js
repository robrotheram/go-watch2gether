import { send } from '@giantmachines/redux-websocket';
import {PROGRESS_UPDATE, SEEK_TO_USER} from '../room/room.types';
import { GetWatcher } from '../user';
import { PLAYING } from './video.types';
import store from '../index'

export const play = () => {
    return dispatch => {
        let evnt = {action: PLAYING, watcher: GetWatcher()}
        dispatch(send(evnt))
    }
}

export const pause = () => {
    return dispatch => {
        let evnt = {action: "PAUSING", watcher:GetWatcher()}
        dispatch(send(evnt))
    }
}

export const seekToUser = (seek) => {
    console.log("SEEK", seek)
    return dispatch => {
        dispatch( {
            type: SEEK_TO_USER,
            seek: seek
        })
    }
}

export const seekToHost = () => {
    return dispatch => {
        let room = store.getState().room
        let hosts = room.watchers.filter(w => w.id === room.host)
        if (hosts.length === 1) {
            console.log("SEEK_TO_HOST", hosts[0].seek)
            seekToUser(hosts[0].seek)
        }
    }
}

export const handleFinish = () => {
    return dispatch => {
        let evnt = {action: "HANDLE_FINSH", watcher:  GetWatcher()}
        dispatch(send(evnt))  
        dispatch( {
            type: PROGRESS_UPDATE,
            seek:  {
                "progress_percent": 1,
                "progress_seconds": GetWatcher().seek.progress_seconds
            },
        })
    }
}

export const updateSeek = (percent, seconds) => {
    return dispatch => {
        dispatch( {
            type: PROGRESS_UPDATE,
            seek:  {
                "progress_percent": percent,
                "progress_seconds": seconds
            },
        })
    }
}