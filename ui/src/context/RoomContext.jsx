import axios from "axios";
import React, { createContext, useState, useEffect, useContext } from "react";
import { useParams } from "react-router";
import { useWebSocket } from "react-use-websocket/dist/lib/use-websocket";
import { API_URL, WS_URL } from "./config";
import { CLEAR_ERROR, EVENT_PAUSING, EVENT_PLAYING, EVENT_SEEK_TO_USER, EVENT_UPDATE_QUEUE, EVENT_USER_UPDATE, GET_META_SUCCESSFUL, JOIN_SUCCESSFUL, LEAVE_SUCCESSFUL, LOCAL_QUEUE_UPDATE, PROGRESS_UPDATE, REJOIN_SUCCESSFUL, ROOM_ERROR } from "./event.types";
import { UserContext } from "./UserContext";


// create context
const RoomContext = createContext();
const defaultRoom = {
  id: '',
  owner: '',
  host: '',
  controls: false,
  auto_skip: false,
  queue: [],
  watchers: [],
  error: '',
  active: false,
  current_video: {
    id: '',
    user_id: '',
    url: '',
    title: '',
    seek_to_user: {
      progress_seconds: 0,
      progress_percent: 0,
    },
    playing: false,
  },
  current_seek: {
    progress_seconds: 0,
    progress_percent: 0,
  }
}

const RoomContextProvider = ({ children }) => {
  const { id } = useParams();
  const [room, setRoom] = useState(defaultRoom);
  const [user, loading] = useContext(UserContext);
  const [socketUrl, setSocketUrl] = useState(``);
  const { sendJsonMessage, lastJsonMessage, readyState } = useWebSocket(socketUrl)

  useEffect(() => {
    setRoom({ ...room, active: false });
    if (id !== undefined) {
      console.log("JOIN", id, user)
      axios.post(`${API_URL}room/join`, {
        id: id,
        name: id,
        username: user.username,
        anonymous: false,
      }).then((res) => {
        console.log("JOIN RES", res.data)
        setSocketUrl(`${WS_URL}room/${res.data.room_id}/ws?token=${res.data.user.id}`)
        // setRoom({ ...room, id: res.data.room_id, active: true });
      })
    }
  }, [id]);

  useEffect(() => {
    if (lastJsonMessage !== null) {
      switch (lastJsonMessage.action) {
        case JOIN_SUCCESSFUL:
          setRoom({ ...room, id: action.room, error: '', active: true });
          break;
        case LEAVE_SUCCESSFUL:
          // openNotificationWithIcon('error', 'You have been diconnected from the room');
          setRoom({ ...room, id: '', error: '', active: false });
          break;
        
        case LOCAL_QUEUE_UPDATE:
          setRoom({ ...room, queue: lastJsonMessage.queue });
          break;
        case ROOM_ERROR:
          setRoom({ ...room, error: lastJsonMessage.error });
          break;

        case CLEAR_ERROR:
          setRoom({ ...room, error: '' });
          break;

        case EVENT_UPDATE_QUEUE:
          setRoom({ ...room, ...lastJsonMessage, error: '' });
          break;

        case EVENT_PAUSING:
          setRoom({ ...room, ...lastJsonMessage, error: '' });
          break;

        case EVENT_PLAYING:
          setRoom({ ...room, ...lastJsonMessage, error: '' });
          break;

        case EVENT_SEEK_TO_USER:
          const watcher = lastJsonMessage.watchers.filter((w) => w.id === user.id);
          console.log('video-user', watcher);
          if (watcher.length > 0) {
            room.seek_to_user = watcher[0].seek;
          }
          setRoom(room);

        case EVENT_USER_UPDATE:
          setRoom({ ...room, ...lastJsonMessage, error: '' });
          break;
        default:
        //setRoom(lastJsonMessage)
      }
    }
  }, [lastJsonMessage]);

  const GetWatcher = () => {
    return { ...user, seek: room.current_seek }
  }
  const actions = {
    updateSeek: (percent, seconds) => {
      setRoom({
        ...room, current_seek: {
          progress_percent: percent,
          progress_seconds: seconds,
        }
      });
    },

    updateQueue: (queue) => {
      queue.forEach((v) => { delete v.key; });
      sendJsonMessage({ action: 'UPDATE_QUEUE', queue, watcher: GetWatcher() });
    },

    nextVideo: () => {
      const EVENT = { action: 'NEXT_VIDEO', watcher: GetWatcher() };
      sendJsonMessage(EVENT)
    },

    updateUser: (u) => {
      const EVENT = { action: EVENT_USER_UPDATE, watcher: u };
      sendJsonMessage(EVENT)
    },

    play: () => {
      const EVENT = { action: EVENT_PLAYING, watcher: GetWatcher() };
      sendJsonMessage(EVENT);
    },

    pause: () => {
      const EVENT = { action: EVENT_PAUSING, watcher: GetWatcher() };
      sendJsonMessage(EVENT);
    },

    handleFinish: () => {
      const EVENT = { action: 'HANDLE_FINSH', watcher: GetWatcher() };
      sendJsonMessage(EVENT);
      sendJsonMessage({
        type: PROGRESS_UPDATE,
        seek: {
          progress_percent: 1,
          progress_seconds: GetWatcher().seek.progress_seconds,
        },
      });
    },

    updateSettings: (cntrls, auto_skip) => {
      const EVENT = { action: 'UPDATE_SETTINGS', watcher: GetWatcher(), settings: { controls: cntrls, auto_skip } };
      sendJsonMessage(EVENT);
    },

    forceSinkToMe: () => {
      const EVENT = { action: EVENT_SEEK_TO_USER, watcher: GetWatcher() };
      sendJsonMessage(EVENT);
    },

    seekToUser: (seek) => {
      setRoom({...room, seek_to_user: seek });
    },

    seekToHost: () => {
      const hosts = room.watchers.filter((w) => w.id === room.host);
      if (hosts.length === 1) {
        setRoom({...room, seek_to_user: hosts[0].seek });
      }
    },

    updateHost: (host) => {
      const EVENT = { action: 'UPDATE_HOST', watcher: GetWatcher(), host };
      sendJsonMessage(EVENT);
    },

    pageChange: () =>{
      setRoom({...room, active: false});
    }

  }


  return (
    // the Provider gives access to the context to its children
    <RoomContext.Provider value={[
      room, actions
    ]}>
      {children}
    </RoomContext.Provider>
  );
};

export { RoomContext, RoomContextProvider };