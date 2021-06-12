import axios from 'axios';
import { BASE_URL, API_URL, history } from '../index';
import { join } from '../room/room.actions';
import { AUTH_LOGIN } from '../event.types';

export const checklogin = (room) => (dispatch) => {
  axios.get(`${BASE_URL}/auth/user`).then((res) => {
    console.log('auth', res.data);
    dispatch({
      type: AUTH_LOGIN,
      auth: true,
      id: res.data.user.id,
      username: res.data.user.username,
      icon: res.data.user.avatar_icon,
      guilds: res.data.guilds,
    });
    if (room === undefined || room === '') {
      history.push('/app');
    }
    axios.get(`${API_URL}room/${room}`).then((data) => {
      dispatch(join(room, data.data.name, res.data.user.username, false));
    }).catch(() => {
      history.push('/app');
    });
  }).catch((e) => {
    console.warn('Unathoiriz', e);
    dispatch({
      type: AUTH_LOGIN,
      auth: false,
      username: '',
      guilds: [],
    });
  });
};
