import { send } from '@giantmachines/redux-websocket';
import {
  PROGRESS_UPDATE, EVNT_SEEK_TO_USER, EVNT_PLAYING, EVNT_PAUSING,
} from '../event.types';
import { GetWatcher } from '../user';
import store from '../index';

export const play = () => (dispatch) => {
  const evnt = { action: EVNT_PLAYING, watcher: GetWatcher() };
  dispatch(send(evnt));
};

export const pause = () => (dispatch) => {
  const evnt = { action: EVNT_PAUSING, watcher: GetWatcher() };
  dispatch(send(evnt));
};

export const seekToUser = (seek) => (dispatch) => {
  dispatch({
    type: EVNT_SEEK_TO_USER,
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
  const evnt = { action: 'HANDLE_FINSH', watcher: GetWatcher() };
  dispatch(send(evnt));
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
