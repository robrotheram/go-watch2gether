import React, { useState, useEffect } from 'react';
import { Button, Space, Card } from 'antd';
import { Input} from "antd"
import {connect} from 'react-redux'
import DrawerForm from './UserDrawer';
import {  VideoCameraOutlined } from '@ant-design/icons';
import { openNotificationWithIcon } from "../../common/notification"


import { Row, Col } from 'antd';

import {updateQueue, updateLocalQueue} from '../../../store/room/room.actions'
import {leave} from '../../../store/room/room.actions'
import PlaylistDrawer from './playlists/PlaylistDrawer'
import {createVideoItem, validURL} from '../../../store/video'
import Share from "./ShareModal"
import Settings from './SettingsModal'

const Controls = (props) => {
  
  const {isHost} = props.user
  const {title} = props.video

  const [newurl, setURL] = useState("");
  const { queue } = props.room


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
            <Space style={{"marginTop":"1px"}}>
                <PlaylistDrawer/>
              </Space>
            </Col>
            <Col flex="auto">
                <Input className="videoInput" defaultValue="mysite" value={newurl} onChange={e => setURL(e.target.value)}  onKeyDown={handleKeyDown} addonAfter={( <Button type="primary" onClick={addToQueue} icon={<VideoCameraOutlined />}>Add Video</Button>)}/>
            </Col>
            <Col>
              <Space style={{"marginTop":"1px"}} size={4}>
                <DrawerForm/>
                {isHost ?<Settings/>: null}
                <Share/>
              </Space>
            </Col>
          </Row>
        </Card>
    )
}

const mapStateToProps  = (state) =>{
  return state
} 
export default connect(mapStateToProps, {leave})(Controls)
