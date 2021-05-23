import {AUTH_LOGIN, GET_META_SUCCESSFUL, JOIN_SUCCESSFUL, PROGRESS_UPDATE } from '../event.types';

const INITIAL_STATE = {
  id: "",
  room: "",
  auth : false,
  username: "",
  icon: "",
  guilds: [],
  seek: {
      progress_percent: 0.0,
      progress_seconds: 0
  } ,
  video_id: "",
  isHost: false,
  playing: false
}


export const userReducer = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case AUTH_LOGIN:
      return {
        ...state, id: action.id, auth: action.auth, username: action.username, icon: action.icon, guilds: action.guilds
      };
    case JOIN_SUCCESSFUL:
      return {
        ...state, room: action.room,
      };

    case PROGRESS_UPDATE:
      return {
        ...state, seek: action.seek,
      };

    case GET_META_SUCCESSFUL:
      let video = action.payload.current_video
      return {
        ...state,
        video_id: video.id,
        playing: action.payload.playing,
        isHost: isHost(state, action.payload.host)
      };

    case "REDUX_WEBSOCKET::MESSAGE":
      try {
        let data = JSON.parse(action.payload.message)
        return process_websocket_event(state, data)
      } catch (e) {
        console.log("Parse Error", action.payload.message, e)
      }
      return state;

    default: return state;
  }
};


const process_websocket_event = (state, data) => {
 // console.log("video reducer action", data.action, data)
  switch (data.action) {
    case "UPDATE_HOST":                
      return {
        ...state, isHost: isHost(state, data.host)
      };
    default:
      return state;
  }
}

function isHost (user, host){
  return user.id === host
}