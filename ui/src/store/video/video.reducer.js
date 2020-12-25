import {GET_META_SUCCESSFUL} from '../room/room.types';

const INITIAL_STATE = {
  id: "",
  url: "",
  title: "",
  seek: "",
  playing: false
}
export const videoReducer = (state = INITIAL_STATE, action) => {
  switch (action.type) {

    case GET_META_SUCCESSFUL:
      let video = action.payload.current_video
      return {
        ...state,
        id: video.id,
        url: video.url,
        title: video.title,
        seek: action.payload.seek,
        playing: action.payload.playing,
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
  console.log("video reducer action", data.action, data)
  switch (data.action) {
    case "CHANGE_VIDEO":
      return {
        ...state, id: data.current_video.id, title: data.current_video.title, url: data.current_video.url
      };
    case "PLAYING":
      if (state.playing !== data.playing) {
        if (state.seek < 1) {
          //openNotificationWithIcon("success", "User: " + data.user.name + " started video")
        }
      }
      return {
        ...state, playing: true,
      };
    case "PAUSING":
      if (state.playing !== data.playing) {
        if (state.seek < 1) {
          //openNotificationWithIcon("success", "User: " + data.user.name + " has paused video")
        }
      }
      return {
        ...state, playing: false,
      };
    case "SEEK_TO_USER":
      return {
        ...state, seek: data.seek,
      };
    default:
      return state;
  }
}