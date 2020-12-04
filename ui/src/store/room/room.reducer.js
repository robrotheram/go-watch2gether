import {openNotificationWithIcon} from "../../components/notification"
import { JOIN_SUCCESSFUL, ROOM_ERROR, CLEAR_ERROR, GET_META_SUCCESSFUL, UPDATE_SEEK, SEEK_TO_HOST, LEAVE_SUCCESSFUL, REJOIN_SUCCESSFUL} from './room.types';
    const INITIAL_STATE = {
      "name":"",
      "user":"",
      "host":"",
      "current_video":"",
      "seek":0,
      "controls":false,
      "playing":false,
      "queue":[],
      "users":[],
      "error": "",
      "active": false,
    }


export const roomReducer = (state = INITIAL_STATE, action) => {
        // console.log("room_reducer", action)
        switch (action.type) {
            case JOIN_SUCCESSFUL:
               return {
                 ...state, name: action.room, user: action.user, error: "", active: true,
               };
            case LEAVE_SUCCESSFUL:
              return {
                ...state, name: "", user: "", error: "", active: false,
              };
            case GET_META_SUCCESSFUL:
              return {
                ...state, 
                error: "", 
                host: action.payload.host, 
                isHost: action.payload.host===state.user, 
                current_video: action.payload.current_video,
                seek: action.payload.seek,
                controls: action.payload.controls,
                playing: action.payload.playing,
                queue: action.payload.queue,
                users: action.payload.users,
              };
            case REJOIN_SUCCESSFUL:
              return {
                ...state, 
                error: "", 
                host: action.payload.host, 
                isHost: action.payload.host===state.user, 
                current_video: action.payload.current_video,
                seek: action.payload.seek,
                controls: action.payload.controls,
                playing: action.payload.playing,
                queue: action.payload.queue,
                users: action.payload.users,
                active: true
              };
            case "REDUX_WEBSOCKET::MESSAGE":
              try{
                let data = JSON.parse(action.payload.message)
                return process_websocket_event(state, data)
              } catch(e){
                console.log("Parse Error", action.payload.message,  e)
              }
              return state;
            case "REDUX_WEBSOCKET::CLOSED" :
              return {
                ...state, active: false,
              };
            case UPDATE_SEEK: 
              return {
                ...state, seek: action.seek,
              };
            case ROOM_ERROR:
               return {
                  ...state, error: action.error,
               };
            case SEEK_TO_HOST:
              return {
                ...state, seek: action.seek,
             };
            case CLEAR_ERROR:
                return {
                   ...state, error: "",
                };
             default: return state;
        }
    };


const process_websocket_event = (state, data) => {
  // console.log("Procees_Event", data)
  switch(data.action){
    case "ROOM_EXIT":
      openNotificationWithIcon("success", "Room has closed")   
      return {
        ...state, active:false, room:""
      };
    case "UPDATE_QUEUE":             
      openNotificationWithIcon("success", "Queue Updated by "+data.user)   
      return {
        ...state, queue: data.queue
      };
    case "UPDATE_HOST":                
      return {
        ...state, host: data.host, isHost: data.host===state.user
      };
    case "PLAYING": 
      if (state.playing !== data.playing){
        openNotificationWithIcon("success", "User: "+data.user+" started video")
      }
      return {
        ...state, playing: true,
      };
    case "PAUSING": 
      if (state.playing !== data.playing){
        openNotificationWithIcon("success", "User: "+data.user+" has paused video")
      }
      return {
        ...state, playing: false,
      };
    case "USER_UPDATED":
      return {
        ...state, users: data.users,
      };
    case "UPDATE_CONTROLS":
      return {
        ...state, controls: data.controls,
      };
    case "ON_PROGRESS_UPDATE":
      let userList = [...state.users];
      userList = userList.map(user => {
        if(user.name === data.user) {
          return {...user, seek: data.seek};
        }
        return {...user};
      });
      return {
        ...state, users: userList,
      };
    case "SEEK_TO_USER":
        return {
          ...state, seek: data.seek,
        };
    default:
      return state;
  }
}