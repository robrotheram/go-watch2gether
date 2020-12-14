import {openNotificationWithIcon} from "../../components/notification"
import { JOIN_SUCCESSFUL, ROOM_ERROR, CLEAR_ERROR, GET_META_SUCCESSFUL, UPDATE_SEEK, SEEK_TO_HOST, LEAVE_SUCCESSFUL, REJOIN_SUCCESSFUL, PROGRESS_UPDATE} from './room.types';
    const INITIAL_STATE = {
      "name":"",
      "user":{
        "name":"",
        "seek":0,
        is_host: false
      },
      "host":"",
      "current_video":{},
      "seek":0,
      "controls":false,
      "playing":false,
      "queue":[],
      "users":[],
      "error": "",
      "active": false,
    }

function NewUser (name) {
  return {
    name: name,
    seek: 0.0,
    current_video: {url:""}
  }
}

function isHost (users, user){
  let u = users.filter(u => u.name === user)
  console.log(u, users, user)
  if (u.length !== 1){
    return false
  }
  
  return u[0].is_host
}

export const roomReducer = (state = INITIAL_STATE, action) => {
        // console.log("room_reducer", action)
        switch (action.type) {
            case JOIN_SUCCESSFUL:
               return {
                 ...state, name: action.room, user: NewUser(action.user), error: "", active: true,
               };
            case LEAVE_SUCCESSFUL:
              return {
                ...state, name: "", user: {}, error: "", active: false,
              };
            case GET_META_SUCCESSFUL:
              console.log("USER_META", action.payload.users.filter(u => u.name === state.user.name))
              return {
                ...state, 
                error: "", 
                host: action.payload.host, 
                isHost: isHost(action.payload.users, state.user.name), 
                current_video: action.payload.current_video,
                seek: action.payload.seek,
                controls: action.payload.controls,
                playing: action.payload.playing,
                queue: action.payload.queue,
                users: action.payload.users,
              };
            case REJOIN_SUCCESSFUL:
              console.log("USER_META", action.payload.users.filter(u => u.name === state.user.name))
              return {
                ...state, 
                error: "", 
                host: action.payload.host, 
                isHost: isHost(action.payload.users, state.user.name), 
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
            case PROGRESS_UPDATE: 
              let user = state.user
              user.seek = action.seek
              return {
                ...state, user: user,
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
      openNotificationWithIcon("success", "Queue Updated by "+data.user.name)   
      return {
        ...state, queue: data.queue
      };
    case "UPDATE_HOST":                
      return {
        ...state, host: data.host, isHost: data.host===state.user.name
      };
    case "CHANGE_VIDEO":
      return {
        ...state, queue: data.queue, current_video: data.current_video
      };
    case "PLAYING": 
      if (state.user.seek == 1){
        return state
      }
      if (state.playing !== data.playing){
        if (state.seek < 1){
          openNotificationWithIcon("success", "User: "+data.user.name+" started video")
        }
      }
      return {
        ...state, playing: true,
      };
    case "PAUSING": 
      if (state.user.seek == 1){
        return state
      }
      if (state.playing !== data.playing){
        if (state.seek < 1){
          openNotificationWithIcon("success", "User: "+data.user.name+" has paused video")
        }
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
      // let userList = [...state.users];
      // userList = userList.map(user => {
      //   if(user.name === data.user) {
      //     return {...user, seek: data.seek};
      //   }
      //   return {...user};
      // });
      return {
        ...state, users: data.users,
      };
    case "SEEK_TO_USER":
        return {
          ...state, seek: data.seek,
        };
    default:
      return state;
  }
}