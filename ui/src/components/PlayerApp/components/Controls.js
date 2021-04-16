import React, { useState, useEffect } from 'react';
import { Button, Layout, Space, Card } from 'antd';
import { Input} from "antd"
import {
  ArrowLeftOutlined,
  SyncOutlined,
  SettingOutlined,
  ShareAltOutlined
  
} from '@ant-design/icons';
import {connect} from 'react-redux'

import SettingsModal from './SettingsModal'
import DrawerForm from './UserDrawer';
import { PlaySquareOutlined, VideoCameraOutlined } from '@ant-design/icons';
import { openNotificationWithIcon } from "../../common/notification"

import { uid } from 'uid';
import axios from 'axios';

import { Row, Col, Divider } from 'antd';
import {API_URL} from '../../../store'
import {updateQueue, nextVideo, updateLocalQueue} from '../../../store/room/room.actions'
import {leave, sinkToHost, sinkToME} from '../../../store/room/room.actions'
import PlaylistDrawer from './playlists/PlaylistDrawer'

const { Header} = Layout;

const Controls = (props) => {

  const { host, controls, name } = props.room
  
  const {isHost} = props.user
  const {title} = props.video

  const [newurl, setURL] = useState("");
  const [loading, setLoading] = useState(false);
  const { queue } = props.room
  const { url } = props.video


  useEffect(() => {
    if (title === ""){
      document.title = "Watch2gether"
    }else {
      document.title = "Watch2gether | Playing:"+title
    }
    
  }, [title]);


  
    const handleKeyDown = (e) => {
        if (e.key === 'Enter') {
            addToQueue()
          }
    }
   
    const user = props.user

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
        if (validURL(newurl)){
            let videoList = [...queue]; 

            videoList.push({url:"", title:"", loading:true})
            updateLocalQueue(videoList)
            videoList = [...queue].filter(i => !i.loading); 
            
            videoList.push(await createVideoItem(newurl));
            updateQueue(videoList)            
            setURL("")
        } else {
            openNotificationWithIcon('error', "Invalid URL")
        }
    }    

    return (
        <Card className="contolPanel" style={{"height": "81px"}}>
          <Row style={{width:"100%", paddingTop:"10px"}}>
            <Col>
              <Space>
                <PlaylistDrawer/>
              </Space>
            </Col>
            <Col flex="auto">
                <Input className="videoInput" defaultValue="mysite" value={newurl} onChange={e => setURL(e.target.value)}  onKeyDown={handleKeyDown} addonAfter={( <Button type="primary" onClick={addToQueue} icon={<VideoCameraOutlined />}>Add Video</Button>)}/>
            </Col>
            <Col>
              <Space>
                <DrawerForm/>
              </Space>
            </Col>
          </Row>
          {/* <Space style={{"float":"right", "width":"100%"}}>
            
            
             
          </Space> */}
          
          
          
        </Card>
    )
}

const mapStateToProps  = (state) =>{
  return state
} 
export default connect(mapStateToProps, {leave, sinkToHost, sinkToME })(Controls)
