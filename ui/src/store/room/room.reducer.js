import { openNotificationWithIcon } from '../../components/common/notification';
import {
  JOIN_SUCCESSFUL, ROOM_ERROR, CLEAR_ERROR, GET_META_SUCCESSFUL, LEAVE_SUCCESSFUL, REJOIN_SUCCESSFUL, EVNT_PLAYING, EVNT_PAUSING, EVNT_NEXT_VIDEO, EVNT_UPDATE_QUEUE,
} from '../event.types';

const INITIAL_STATE = {
  id: '',
  owner: '',
  host: '',
  controls: false,
  auto_skip: false,
  queue: [],
  watchers: [],
  error: '',
  active: true,
};
export const roomReducer = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case JOIN_SUCCESSFUL:
      return {
        ...state, id: action.room, error: '', active: true,
      };
    case LEAVE_SUCCESSFUL:
      openNotificationWithIcon('error', 'You have been diconnected from the room');
      return {
        ...state, id: '', error: '', active: false,
      };
    case GET_META_SUCCESSFUL:
      return {
        ...state,
        error: '',
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
        error: '',
        host: action.payload.host,
        controls: action.payload.controls,
        queue: action.payload.queue,
        watchers: action.payload.watchers,
        active: true,
      };
    case 'REDUX_WEBSOCKET::MESSAGE':
      try {
        const data = JSON.parse(action.payload.message);
        HandleNotification(data);
        return process_websocket_event(state, data);
      } catch (e) {
        console.log('Parse Error', action.payload.message, e);
      }
      return state;
    case 'REDUX_WEBSOCKET::CLOSED':
      return {
        ...state, active: false,
      };

    case 'LOCAL_QUEUE_UPDATE':
      return {
        ...state, queue: action.queue,
      };
    case ROOM_ERROR:
      return {
        ...state, error: action.error,
      };

    case CLEAR_ERROR:
      return {
        ...state, error: '',
      };
    default: return state;
  }
};

const process_websocket_event = (state, data) => {
  switch (data.action) {
    case 'ROOM_EXIT':
      openNotificationWithIcon('success', 'Room has closed');
      return {
        ...state, active: false, room: '',
      };
    default:
      return {
        ...state,
        host: data.host,
        controls: data.settings.controls,
        auto_skip: data.settings.auto_skip,
        queue: data.queue,
        watchers: data.watchers,
      };
  }
};

const HandleNotification = (data) => {
  switch (data.action) {
    case EVNT_PLAYING:
      openNotificationWithIcon('info', `${data.watcher.username} has started the video`);
      break;
    case EVNT_PAUSING:
      openNotificationWithIcon('info', `${data.watcher.username} has stopped the video`);
      break;
    case EVNT_NEXT_VIDEO:
      openNotificationWithIcon('info', `${data.watcher.username} has changed the vidoe the video`);
      break;
    case EVNT_UPDATE_QUEUE:
      openNotificationWithIcon('info', `${data.watcher.username} has updated the queue`);
      break;
    default:
      break;
  }
};
