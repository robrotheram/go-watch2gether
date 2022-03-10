import axios from 'axios';
import { connect, send } from '@giantmachines/redux-websocket';
import {
  CLEAR_ERROR, JOIN_SUCCESSFUL, ROOM_ERROR, LEAVE_SUCCESSFUL, GET_META_SUCCESSFUL, REJOIN_SUCCESSFUL, PROGRESS_UPDATE, EVENT_SEEK_TO_USER,
} from '../event.types';
import store, { API_URL, WS_URL, history } from '../index';

import { GetUsername, GetWatcher } from '../user';

export const join = (roomid, room, user, anonymous) => (dispatch) => {
  axios.post(`${API_URL}room/join`, {
    id: roomid, name: room, username: user, anonymous,
  }).then((res) => {
    console.log('Action', res);
    dispatch({
      type: JOIN_SUCCESSFUL,
      room: res.data.room_id,
      user: res.data.user,
    });
    dispatch(getMeta(res.data.room_id));
    console.log('WS_URL', `${WS_URL}room/${room}/ws`);
    store.dispatch(Connect(res.data.room_id, res.data.user.id));
    history.push(`/app/room/${res.data.room_id}`);
  }).catch((e) => {
    console.log(e);
    history.push('/');
    if (e.response === undefined) {
      dispatch({
        type: ROOM_ERROR,
        error: `Oh Dear something happened when user: ${user} tried to joining room: ${room} Unable to contact server`,
      });
    } else {
      dispatch({
        type: ROOM_ERROR,
        error: e.response.data,
      });
    }
  });
};

export const Connect = (room, user) => connect(`${WS_URL}room/${room}/ws?token=${user}`);

export const reJoin = (room) => (dispatch) => {
  axios.get(`${API_URL}room/${room}`).then((res) => {
    store.dispatch(Connect(room, GetUsername()));
    dispatch({
      type: REJOIN_SUCCESSFUL,
      payload: res.data,
    });
  }).catch((e) => {
    console.log(e);
    history.push('/');
    dispatch({
      historytype: ROOM_ERROR,
      error: e.response.data,
    });
  });
};

export const leave = (room, user) => (dispatch) => {
  axios.post(`${API_URL}room/leave`, {
    id: store.getState().room.id,
    name: store.getState().room.name,
    username: GetUsername(),
  }).then((res) => {
    console.log('Action', res);
    dispatch({
      type: LEAVE_SUCCESSFUL,
      room,
      user,
    });
    history.push('/');
  }).catch((e) => {
    dispatch({
      type: ROOM_ERROR,
      error: e.response.data,
    });
    history.push('/');
  });
};

export const getMeta = (id) => {
  console.log('GETTING META', id);
  return (dispatch) => {
    axios.get(`${API_URL}room/${id}`).then((res) => {
      console.log('Action', res);
      store.dispatch({
        type: GET_META_SUCCESSFUL,
        payload: res.data,
      });
    }).catch((e) => {
      console.log(e);
      dispatch({
        type: ROOM_ERROR,
        error: 'Connection Issue to the server',
      });
    });
  };
};

export const updateSeek = (percent, seconds) => {
  const seek = {
    progress_percent: percent,
    progress_seconds: seconds,
  };

  store.dispatch({
    type: PROGRESS_UPDATE,
    seek,
  });
  // let EVENT = {action: "ON_PROGRESS_UPDATE", watcher:  GetWatcher(), seek:seek}
  // store.dispatch(send(EVENT))
};

export const updateSettings = (cntrls, auto_skip) => {
  const EVENT = { action: 'UPDATE_SETTINGS', watcher: GetWatcher(), settings: { controls: cntrls, auto_skip } };
  store.dispatch(send(EVENT));
};

export const forceSinkToMe = () => (dispatch) => {
  const EVENT = { action: EVENT_SEEK_TO_USER, watcher: GetWatcher() };
  dispatch(send(EVENT));
};

export const updateHost = (host) => {
  const EVENT = { action: 'UPDATE_HOST', watcher: GetWatcher(), host };
  store.dispatch(send(EVENT));
};

export async function isAlive() {
  const EVENT = {
    action: 'USER_UPDATE',
    watcher: GetWatcher(),
  };
  return store.dispatch(send(EVENT));
}

export const sinkToME = (seek) => {
  const EVENT = { action: 'SEEK_TO_ME', watcher: GetWatcher(), seek };
  store.dispatch(send(EVENT));
};

export const clearError = () => {
  console.log('Action', CLEAR_ERROR);
  return {
    type: CLEAR_ERROR,
  };
};

export const updateQueue = (queue) => (dispatch) => {
  queue.forEach((v) => { delete v.key; });
  const EVENT = { action: 'UPDATE_QUEUE', queue, watcher: GetWatcher() };
  dispatch(send(EVENT));
};

export const updateLocalQueue = (queue) => {
  store.dispatch({
    type: 'LOCAL_QUEUE_UPDATE',
    queue,
    watcher: GetWatcher(),
  });
};

export const nextVideo = () => (dispatch) => {
  const EVENT = { action: 'NEXT_VIDEO', watcher: GetWatcher() };
  dispatch(send(EVENT));
};
