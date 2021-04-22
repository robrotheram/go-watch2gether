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


import { Row, Col, Divider } from 'antd';

import {updateQueue, nextVideo, updateLocalQueue} from '../../../store/room/room.actions'
import {leave, sinkToHost, sinkToME} from '../../../store/room/room.actions'
import PlaylistDrawer from './playlists/PlaylistDrawer'
import {createVideoItem, validURL} from '../../../store/video'
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

    
    const addToQueue = async () => {
        if (validURL(newurl)){
            let videoList = [...queue]; 

            videoList.push({url:"", title:"", loading:true})
            updateLocalQueue(videoList)
            videoList = [...queue].filter(i => !i.loading); 
            
            videoList.push(await createVideoItem(newurl, user.username));
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
