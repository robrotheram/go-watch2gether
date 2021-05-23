import {AUTH_LOGIN, GET_META_SUCCESSFUL, PROGRESS_UPDATE, EVNT_SEEK_TO_USER, EVNT_PLAYING} from '../event.types';

const INITIAL_STATE = {
  id: "",
  user_id: "",
  url: "",
  title: "",
  current_seek: {
    progress_seconds : 0,
    progress_percent: 0
  },
  seek_to_user: {
    progress_seconds : 0,
    progress_percent: 0
  },
  playing: false
}
export const videoReducer = (state = INITIAL_STATE, action) => {
  switch (action.type) {

    case AUTH_LOGIN: 
      return {
        ...state,
        user_id: action.id
      }
    case GET_META_SUCCESSFUL:
      let video = action.payload.current_video
      return {
        ...state,
        id: video.id,
        url: video.url,
        title: video.title,
        playing: action.payload.playing
      };

    case PROGRESS_UPDATE:
      return {
        ...state, current_seek: action.seek,
      };

    case EVNT_SEEK_TO_USER:
      return {
        ...state, seek_to_user: action.seek,
      };




    case "REDUX_WEBSOCKET::MESSAGE":
      try {
        let data = JSON.parse(action.payload.message)
        return process_websocket_event(state, data)
      } catch (e) {
        console.log("Parse Error", action.payload.message, e)
      }
      return state;
    case "SEEK_TO_HOST":
      return {
        ...state, seek: action.seek,
      };
    default: return state;
  }
};


const process_websocket_event = (state, data) => {
  let video = data.current_video
  switch (data.action) {
    case EVNT_PLAYING:
      if (state.current_seek.progress_percent < 1) {
        return {...state, playing: data.playing,}
      }
      break;
    case EVNT_SEEK_TO_USER :
      let user = data.watchers.filter(w => w.id === state.user_id)
      console.log("video-user", user) 
      if (user.length > 0){
        state.seek_to_user = user[0].seek
      }
      return state

    default:
      let _state = {
        ...state,
        id: video.id,
        url: video.url,
        title: video.title,
      };
      if (state.current_seek.progress_percent < 1 && state.url === data.current_video.url ) {
        _state.playing = data.playing
      }
      return _state
  }

 
}