

import React, { useState } from 'react';
import { Button, Layout, Space } from 'antd';
import {
  ArrowLeftOutlined,
  SyncOutlined,
  SettingOutlined,
  ShareAltOutlined
  
} from '@ant-design/icons';
import {connect} from 'react-redux'
import {leave, sinkToHost, sinkToME} from '../store/room/room.actions'
import SettingsModal from './SettingsModal'
import ShareModal from './ShareModal';

const { Header} = Layout;

function Navigation (props) {
  const { name, isHost, controls} = props

  const [isSettingsModalVisible, setIsSettingModalVisible] = useState(false);
  const showSettingsModal = () => {setIsSettingModalVisible(true);};
  const handleSettingsOk = () => {setIsSettingModalVisible(false); };
  const handleSettingsCancel = () => {setIsSettingModalVisible(false);};

  const [isShareModalVisible, setIsShareModalVisible] = useState(false);
  const showShareModal = () => {setIsShareModalVisible(true);};
  const handleShareOk = () => {setIsShareModalVisible(false); };
  const handleShareCancel = () => {setIsShareModalVisible(false);};

    return (
        <Header style={{"display":"block ruby", "zIndex": "1000", "position":"fixed", "left":0, "right":0, "top":0}}>
          <Button style={{"display": "inline-block"}} type="link" size="large" icon={<ArrowLeftOutlined />} style={{color: "white"}} onClick={() => { props.leave()}}/> 
          <div className="logo" style={{"display": "inline-block"}}>
            <h1 style={{"color":"white"}}>Watch2Gether</h1>
          </div>
          <Space style={{"float":"right"}}>
         
                { !isHost ? <Button type="primary" icon={<SyncOutlined />} key="3" onClick={() => sinkToHost()}>Sync to host</Button> : null}
           
                { controls || isHost ? <Button type="primary" icon={<SyncOutlined />} key="2" onClick={() => sinkToME()}>Sync everyone to me</Button>: null}
      
                {isHost ?<Button type="primary" onClick={() => setIsSettingModalVisible(true) } icon={<SettingOutlined />} key="1"></Button>: null}

                <Button type="primary" icon={<ShareAltOutlined />} onClick={() => setIsShareModalVisible(true) } />
          </Space>
          <SettingsModal isModalVisible={isSettingsModalVisible} showModal={showSettingsModal} handleOk={handleSettingsOk} handleCancel={handleSettingsCancel}/>
          <ShareModal isModalVisible={isShareModalVisible} showModal={showShareModal} handleOk={handleShareOk} handleCancel={handleShareCancel}/>
        </Header>
    )
}
const mapStateToProps  = (state) =>{
    
  return state.room
} 
export default connect(mapStateToProps, {leave, sinkToHost, sinkToME })(Navigation)