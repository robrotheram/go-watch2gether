import { uid } from 'uid';
import axios from 'axios';
import {API_URL} from '../'
import store from '../index'

export * from './video.types'
export * from './video.actions'
export * from './video.reducer'

export const validURL = (str) => {
    var pattern = new RegExp('^(https?:\\/\\/)?'+ // protocol
      '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|'+ // domain name
      '((\\d{1,3}\\.){3}\\d{1,3}))'+ // OR ip (v4) address
      '(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*'+ // port and path
      '(\\?[;&a-z\\d%_.~+=-]*)?'+ // query string
      '(\\#[-a-z\\d_]*)?$','i'); // fragment locator
    return !!pattern.test(str) && !str.includes("list=");
  }

export const getTitle = async (url) => {
    const result = await axios(API_URL+"scrape?url="+encodeURI(url),);
    return (result.data.Title);
};

export const createVideoItem = async (url, username) => {
  let title = await getTitle(url)
  console.log("VIDEO GET URL", title)
  return {
      "url":url, 
      "title": title,
      "user":   username, 
      "uid": uid(16)
  }
}