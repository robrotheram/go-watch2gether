import axios from 'axios';
import { API_URL } from '../index';
import { GetUsername } from '../user';

export const getPlaylists = (room_id) => (dispatch) => {
  axios.get(`${API_URL}room/${room_id}/playlist`).then((res) => {
    console.log('auth', res.data);
    dispatch({
      type: 'GET_PLAYLISTS',
      playlists: res.data,
    });
  }).catch((e) => {
    console.warn('Unathoiriz', e);
  });
};
export const createPlaylists = (room_id, playlist) => {
  playlist.username = GetUsername();
  return (dispatch) => {
    axios.put(`${API_URL}room/${room_id}/playlist`, playlist).then((res) => {
      console.log('auth', res.data);
      dispatch({
        type: 'UPDATE_COMPLETE',
      });
      dispatch(getPlaylists(room_id));
    }).catch((e) => {
      console.warn('Unathoiriz', e);
    });
  };
};

export const updatePlaylists = (room_id, playlist) => {
  playlist.username = GetUsername();
  return (dispatch) => {
    axios.post(`${API_URL}room/${room_id}/playlist/${playlist.id}`, playlist).then((res) => {
      console.log('auth', res.data);
      dispatch({
        type: 'UPDATE_COMPLETE',
      });
      dispatch(getPlaylists(room_id));
    }).catch((e) => {
      console.warn('Unathoiriz', e);
    });
  };
};
export const deletePlaylists = (room_id, id) => (dispatch) => {
  axios.delete(`${API_URL}room/${room_id}/playlist/${id}`).then((res) => {
    console.log('auth', res.data);
    dispatch({
      type: 'DELETE_COMPLETE',
    });
    dispatch(getPlaylists(room_id));
  }).catch((e) => {
    console.warn('Unathoiriz', e);
  });
};

export const loadPlaylists = (room_id, id) => (dispatch) => {
  axios.get(`${API_URL}room/${room_id}/playlist/${id}/load`).then((res) => {
    console.log('auth', res.data);
    dispatch({
      type: 'LOAD_COMPLETE',
    });
  }).catch((e) => {
    console.warn('Unathoiriz', e);
  });
};
