
import {Button, Space} from "antd"
import { StarOutlined } from '@ant-design/icons';
import {useDispatch, useSelector} from 'react-redux'

import {useState} from "react"
import { uid } from 'uid';
import axios from 'axios';
import "./VideoListControls.less"
import {API_URL} from '../../../../store'
import {updateQueue, nextVideo, updateLocalQueue} from '../../../../store/room/room.actions'


const shuffleArray = array => {
    for (let i = array.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      const temp = array[i];
      array[i] = array[j];
      array[j] = temp;
    }
  }

  
export function VideoControlComponent () {
  
    const [loading, setLoading] = useState(false);
    const dispatch = useDispatch()
    const queue = useSelector(state => state.room.queue);
    const url = useSelector(state => state.video.url);
    const user = useSelector(state => state.user);

    const getTitle = async (url) => {
        const result = await axios(API_URL+"scrape?url="+encodeURI(url),);
        return (result.data.Title);
    };


    const createVideoItem = async (url) => {
        let title = await getTitle(url)
        console.log("VIDEO GET URL", title)
        return {
            "url":url, 
            "title": title,
            "user":user.username, 
            "uid": uid(16)
        }
    }
        
    const clearQueue = () => {
        dispatch(updateQueue([]));
    }

    const skipQueue = () => {
        dispatch(nextVideo());
    }

    const shuffleQueue = () => {
        let videoList = [...queue];
        shuffleArray(videoList)
        dispatch(updateQueue(videoList));
    }

    return (
        <div className="list-controls">
            <Button onClick={skipQueue} className="list-controls-item"> Skip </Button>
            <Button onClick={clearQueue} className="list-controls-item">Clear Queue</Button>
            <Button onClick={shuffleQueue} className="list-controls-item">Shuffle</Button>
        </div>
    )
}

export const VideoControls = (VideoControlComponent)