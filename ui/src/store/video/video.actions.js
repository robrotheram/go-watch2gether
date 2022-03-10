import { send } from '@giantmachines/redux-websocket';
import {
  PROGRESS_UPDATE, EVENT_SEEK_TO_USER, EVENT_PLAYING, EVENT_PAUSING,
} from '../event.types';
import { GetWatcher } from '../user';
import store from '../index';

export const play = () => (dispatch) => {
  const EVENT = { action: EVENT_PLAYING, watcher: GetWatcher() };
  dispatch(send(EVENT));
};

export const pause = () => (dispatch) => {
  const EVENT = { action: EVENT_PAUSING, watcher: GetWatcher() };
  dispatch(send(EVENT));
};

export const seekToUser = (seek) => (dispatch) => {
  dispatch({
    type: EVENT_SEEK_TO_USER,
    seek,
  });
};

export const seekToHost = () => (dispatch) => {
  const { room } = store.getState();
  const hosts = room.watchers.filter((w) => w.id === room.host);
  if (hosts.length === 1) {
    dispatch(seekToUser(hosts[0].seek));
  }
};

export const handleFinish = () => (dispatch) => {
  const EVENT = { action: 'HANDLE_FINSH', watcher: GetWatcher() };
  dispatch(send(EVENT));
  dispatch({
    type: PROGRESS_UPDATE,
    seek: {
      progress_percent: 1,
      progress_seconds: GetWatcher().seek.progress_seconds,
    },
  });
};

export const updateSeek = (percent, seconds) => (dispatch) => {
  dispatch({
    type: PROGRESS_UPDATE,
    seek: {
      progress_percent: percent,
      progress_seconds: seconds,
    },
  });
};
