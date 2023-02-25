import axios from "axios";
import { API_URL } from "../../store";

export const addVideosToQueue = async (room_id,url) => {
    await axios.post(`${API_URL}room/${room_id}/videos`, {"url":url}).catch(err => console.log(err))
  };