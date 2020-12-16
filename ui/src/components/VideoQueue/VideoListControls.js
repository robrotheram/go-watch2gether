
import { VideoList }from "./RandomVideo"
import {Button, Space, Input} from "antd"
import { StarOutlined, VideoCameraOutlined } from '@ant-design/icons';
import { openNotificationWithIcon } from "../notification"
import {connect} from 'react-redux'
import {updateQueue, nextVideo, updateLocalQueue} from '../../store/room/room.actions'
import {useState} from "react"
import { uid } from 'uid';
import axios from 'axios';
import {API_URL} from '../../store'

export function VideoControlComponent (props) {
    const [url, setURL] = useState("");
    const [loading, setLoading] = useState(false);
    const { queue, user } = props

    const validURL = (str) => {
        var pattern = new RegExp('^(https?:\\/\\/)?'+ // protocol
          '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|'+ // domain name
          '((\\d{1,3}\\.){3}\\d{1,3}))'+ // OR ip (v4) address
          '(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*'+ // port and path
          '(\\?[;&a-z\\d%_.~+=-]*)?'+ // query string
          '(\\#[-a-z\\d_]*)?$','i'); // fragment locator
        return !!pattern.test(str) && !str.includes("list=");
      }

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
            "user":user.name, 
            "uid": uid(16)
        }
    }

    
    const addToQueue = async () => {
        if (validURL(url)){
            let videoList = [...queue]; 

            videoList.push({url:"", title:"", loading:true})
            updateLocalQueue(videoList)
            videoList = [...queue].filter(i => !i.loading); 
            
            videoList.push(await createVideoItem(url));
            updateQueue(videoList)            
            setURL("")
        } else {
            openNotificationWithIcon('error', "Invalid URL")
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
            if (videoList.filter(e => e.url === randomElement).length === 0) {
                videoList.push(await createVideoItem(randomElement));
                break;
            }
        }
        updateQueue(videoList)
        setLoading(false);
    }
      
    return (
        <div>
            <Input className="videoInput" defaultValue="mysite" value={url} onChange={e => setURL(e.target.value)} addonAfter={( <Button type="primary" onClick={addToQueue} icon={<VideoCameraOutlined />}>Add Video</Button>)}/>   
            <Space size="small" style={{width:"100%", marginTop: "10px", marginBottom: "10px"}}>
                <Button onClick={skipQueue} style={{width:"100%"}}> Skip To Next</Button>
                <Button onClick={clearQueue} style={{width:"100%"}}>Clear Queue</Button>
                <Button icon={<StarOutlined />}onClick={randomVideo} style={{width:"100%"}}>Add Random</Button>
            </Space>
            
        </div>
    )
}

const mapStateToProps  = (state) =>{
    return state.room
  } 
export const VideoControls = connect(mapStateToProps, {updateQueue})(VideoControlComponent)