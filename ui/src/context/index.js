import { uid } from 'uid';
import axios from 'axios';
import { API_URL } from "./config"

export const validURL = (str) => {
  const pattern = new RegExp('^(https?:\\/\\/)?' // protocol
      + '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|' // domain name
      + '((\\d{1,3}\\.){3}\\d{1,3}))' // OR ip (v4) address
      + '(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*' // port and path
      + '(\\?[;&a-z\\d%_.~+=-]*)?' // query string
      + '(\\#[-a-z\\d_]*)?$', 'i'); // fragment locator
  return !!pattern.test(str);
};

export const getTitle = async (url) => {
  const result = await axios(`${API_URL}scrape?url=${encodeURI(url)}`);
  return (result.data.Title);
};

export const createVideoItem = async (url, username) => {
  const title = await getTitle(url);
  return {
    url,
    title,
    user: username,
    uid: uid(16),
  };
};


export const addVideosToQueue = async (room_id,url) => {
  await axios.post(`${API_URL}room/${room_id}/videos`, {"url":url}).catch(err => console.log(err))
};
