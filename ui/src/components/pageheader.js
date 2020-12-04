import React, { useState } from 'react';
import { Button, PageHeader } from 'antd';
import { SyncOutlined, SettingOutlined } from '@ant-design/icons';
import {connect} from 'react-redux'
import {leave, sinkToHost, sinkToME} from '../store/room/room.actions'
import SettingsModal from './SettingsModal'

function Pageheader (props) {
    const { name, isHost, host, queue} = props

    const [isModalVisible, setIsModalVisible] = useState(false);
    const showModal = () => {setIsModalVisible(true);};
    const handleOk = () => {setIsModalVisible(false); };
    const handleCancel = () => {setIsModalVisible(false);};

    const currentlyPlaying = () => {
        if (queue[0] === undefined){
            return ""
        }
        return (<span>Currently Playing: <a href={queue[0].url}>{queue[0].url}</a></span>)
    }

    return(
        <PageHeader
            ghost={false}
            onBack={() => { props.leave()} }
            title={name}
            subTitle={currentlyPlaying()}
            extra={isHost ? [
                <Button type="primary" icon={<SyncOutlined />} key="3" onClick={() => sinkToHost()}>Sync to host</Button>,
                <Button type="primary" icon={<SyncOutlined />} key="2" onClick={() => sinkToME()}>Sync everyone to me</Button>,
                <Button type="primary" onClick={() => setIsModalVisible(true) } icon={<SettingOutlined />} key="1"></Button>
            ]:
            [
                <Button type="primary" icon={<SyncOutlined />} onClick={() => sinkToHost()} key="3">Sync to host</Button>,
                <Button type="primary" icon={<SyncOutlined />} onClick={() => sinkToME()} key="2">Sync everyone to me</Button>,
            ]
        }><SettingsModal isModalVisible={isModalVisible} showModal={showModal} handleOk={handleOk} handleCancel={handleCancel}/>
        </PageHeader>
    )
}

const mapStateToProps  = (state) =>{
    
    return state.room
  } 
  export default connect(mapStateToProps, {leave, sinkToHost, sinkToME })(Pageheader)
  