import { GET_META_SUCCESSFUL, JOIN_SUCCESSFUL, PROGRESS_UPDATE } from '../room/room.types';
import {openNotificationWithIcon} from "../../components/notification"
const INITIAL_STATE = {
  id: "",
  name: "",
  seek: 0.0,
  video_id: "",
  isHost: false,
  playing: false
}


export const userReducer = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case JOIN_SUCCESSFUL:
      return {
        ...state, id: action.user.id, name: action.user.name,
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
    case "CHANGE_VIDEO":
      return {
        ...state, video_id: data.current_video.id
      };
    case "UPDATE_HOST":                
      return {
        ...state, isHost: isHost(state, data.host)
      };
    case "PLAYING":
      if (state.playing !== data.playing) {
        if (state.seek < 1) {
          openNotificationWithIcon("success", "User: " + data.watcher.name + " started video")
          return {
            ...state, playing: true,
          };
        }
      }
      return state
    case "PAUSING":
      if (state.playing !== data.playing) {
        if (state.seek < 1) {
          openNotificationWithIcon("success", "User: " + data.watcher.name + " has paused video")
          return {
            ...state, playing: false,
          };
        }
      }
      return state
    default:
      return state;
  }
}

function isHost (user, host){
  return user.id === host
}