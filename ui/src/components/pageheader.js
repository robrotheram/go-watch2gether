import React, { useState } from 'react';
import { Button, PageHeader } from 'antd';
import { SyncOutlined, SettingOutlined } from '@ant-design/icons';
import {connect} from 'react-redux'
import {leave, sinkToHost, sinkToME} from '../store/room/room.actions'
import SettingsModal from './SettingsModal'

function Pageheader (props) {
    const { name, isHost, controls} = props

    const [isModalVisible, setIsModalVisible] = useState(false);
    const showModal = () => {setIsModalVisible(true);};
    const handleOk = () => {setIsModalVisible(false); };
    const handleCancel = () => {setIsModalVisible(false);};

    // const currentlyPlaying = () => {
    //     if (queue[0] === undefined){
    //         return ""
    //     }
    //     return (<span>Currently Playing: <a href={queue[0].url}>{queue[0].url}</a></span>)
    // }

    const getActionButtons = () => {
        let buttons = []
        if (!isHost){
            buttons.push(
                <Button type="primary" icon={<SyncOutlined />} key="3" onClick={() => sinkToHost()}>Sync to host</Button>
            );
        }
        if (controls || isHost){
            buttons.push(
                <Button type="primary" icon={<SyncOutlined />} key="2" onClick={() => sinkToME()}>Sync everyone to me</Button>
            );
        }
        if (isHost){
            buttons.push(
                <Button type="primary" onClick={() => setIsModalVisible(true) } icon={<SettingOutlined />} key="1"></Button>
            );
        } 

        return buttons
    }
    
    return(
        <PageHeader
            ghost={false}
            onBack={() => { props.leave()} }
            title={"Room: "+name}
            // subTitle={currentlyPlaying()}
            extra={getActionButtons()}
        >
            <SettingsModal isModalVisible={isModalVisible} showModal={showModal} handleOk={handleOk} handleCancel={handleCancel}/>
        </PageHeader>
    )
}

const mapStateToProps  = (state) =>{
    
    return state.room
  } 
  export default connect(mapStateToProps, {leave, sinkToHost, sinkToME })(Pageheader)
  