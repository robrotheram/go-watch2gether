
import { VideoList }from "./RandomVideo"
import {Button, Space} from "antd"
import { StarOutlined } from '@ant-design/icons';
import {connect} from 'react-redux'

import {useState} from "react"
import { uid } from 'uid';
import axios from 'axios';

import {API_URL} from '../../../../store'
import {updateQueue, nextVideo, updateLocalQueue} from '../../../../store/room/room.actions'

export function VideoControlComponent (props) {
  
    const [loading, setLoading] = useState(false);
    const { queue } = props.room
    const { url } = props.video

    const user = props.user

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
        updateQueue([]);
    }

    const skipQueue = () => {
        nextVideo();
    }


    const randomVideo = async() => {
        if (loading){ return; }
        setLoading(true);
        let videoList = [...queue]; 

        videoList.push({url:"", title:"", loading:true})
        updateLocalQueue(videoList)
        videoList = [...queue].filter(i => !i.loading); 

        for (var i=1; i < 100; i += 2){
            let randomElement = VideoList[Math.floor(Math.random() * VideoList.length)];
            if (videoList.filter(e => e.url === randomElement || url === randomElement).length === 0) {
                videoList.push(await createVideoItem(randomElement, user.username));
                break;
            }
        }
        updateQueue(videoList)
        setLoading(false);
    }
      
    return (
        <div>
            <Space size="small" style={{width:"100%", marginTop: "10px", marginBottom: "10px"}}>
                <Button onClick={skipQueue} style={{width:"100%"}}> Skip To Next</Button>
                <Button onClick={clearQueue} style={{width:"100%"}}>Clear Queue</Button>
                <Button icon={<StarOutlined />}onClick={randomVideo} style={{width:"100%"}}>Add Random</Button>
            </Space>
            
        </div>
    )
}

const mapStateToProps  = (state) =>{
    return state
  } 
export const VideoControls = connect(mapStateToProps, {updateQueue})(VideoControlComponent)