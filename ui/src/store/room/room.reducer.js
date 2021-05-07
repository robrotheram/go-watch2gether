import {openNotificationWithIcon} from "../../components/common/notification"
import { JOIN_SUCCESSFUL, ROOM_ERROR, CLEAR_ERROR, GET_META_SUCCESSFUL, LEAVE_SUCCESSFUL, REJOIN_SUCCESSFUL} from './room.types';
    const INITIAL_STATE = {
      "id": "",
      "owner": "",
      "host":"",
      "controls":false,
      "auto_skip":false,
      "queue":[],
      "watchers":[],
      "error": "",
      "active": true,
    }
export const roomReducer = (state = INITIAL_STATE, action) => {
        
        switch (action.type) {
            case JOIN_SUCCESSFUL:
               return {
                 ...state, id: action.room, error: "", active: true,
               };
            case LEAVE_SUCCESSFUL:
              openNotificationWithIcon("error", "You have been diconnected from the room")   
              return {
                ...state, id: "", error: "", active: false,
              };
            case GET_META_SUCCESSFUL:
              return {
                ...state, 
                error: "", 
                id: action.payload.id, 
                name: action.payload.name, 
                owner: action.payload.owner, 
                host: action.payload.host, 
                controls: action.payload.settings.controls,
                auto_skip: action.payload.settings.auto_skip,
                queue: action.payload.queue,
                watchers: action.payload.watchers,
              };
            case REJOIN_SUCCESSFUL:
              return {
                ...state, 
                error: "", 
                host: action.payload.host, 
                controls: action.payload.controls,
                queue: action.payload.queue,
                watchers: action.payload.watchers,
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

            
            case "LOCAL_QUEUE_UPDATE":
              return {
                ...state, queue: action.queue
              }
            case ROOM_ERROR:
               return {
                  ...state, error: action.error,
               };
            
            case CLEAR_ERROR:
                return {
                   ...state, error: "",
                };
             default: return state;
        }
    };


const process_websocket_event = (state, data) => {
  //console.log("room_reducer action", data.action, data)
  switch(data.action){
    case "ROOM_EXIT":
      openNotificationWithIcon("success", "Room has closed")   
      return {
        ...state, active:false, room:""
      };
    case "UPDATE_QUEUE":             
      openNotificationWithIcon("success", "Queue Updated by "+data.watcher.username)   
      return {
        ...state, queue: data.queue
      };
    case "UPDATE_HOST":                
      return {
        ...state, host: data.host
      };

    case "USER_UPADTE":
      break
      // return {
      //   ...state, watchers: data.watchers,
      // };
    case "UPDATE_CONTROLS":
      return {
        ...state, controls: data.settings.controls,
      };

      case "UPDATE_SKIP":
        return {
          ...state, auto_skip: data.settings.auto_skip,
        };


    case "ON_PROGRESS_UPDATE":
      // let userList = [...state.watchers];
      // userList = userList.map(user => {
      //   if(user.name === data.user) {
      //     return {...user, seek: data.seek};
      //   }
      //   return {...user};
      // });
      //console.log("watchers", data.watchers)
      return {
        ...state, watchers: data.watchers,
      };
    
    default:
      return state;
  }
}



// function NewUser (usr) {
//   return {
//     id: usr.id,
//     name: usr.name,
//     seek: 0.0,
//     current_video: {url:""}
//   }
// }

// function isHost (watchers, user){
//   let u = watchers.filter(u => u.name === user)
//   console.log(u, watchers, user)
//   if (u.length !== 1){
//     return false
//   }
  
//   return u[0].is_host
// }