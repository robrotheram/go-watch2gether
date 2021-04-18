import React, {useState, useEffect} from 'react';
import { Modal } from 'antd';
import { Form, Button } from 'antd';
import { Switch } from 'antd';
import {connect} from 'react-redux'
import {updateSettings, updateQueue} from '../../../store/room/room.actions'
import {openNotificationWithIconKey} from "../../common/notification"
import {
    SettingOutlined
  } from '@ant-design/icons';


const layout = {
    labelCol: { span: 12 },
    wrapperCol: { span: 12 },
};

const mapStateToProps  = (state) =>{
    return state.room
  } 

const SettingsModal = connect(mapStateToProps, {updateSettings})( (props) => {
    const {isModalVisible, handleOk, handleCancel, queue } = props

    const [controls, setControls] = useState(props.controls);
    const [auto_skip, setAutoSkip] = useState(props.auto_skip);

    useEffect(() => {
        setControls(props.controls);
        setAutoSkip(props.auto_skip)
    }, [props.controls,props.auto_skip ])
  

    const submitForm = () =>{
        updateSettings(controls, auto_skip)
        handleOk();

    }
    const cancelForm = () =>{
        setControls(props.controls)
        handleCancel();
    }

    const handleDarkMode = () => {
        openNotificationWithIconKey("warning", "Only Dark Mode for you!","darkmode")   
        let videoList = [...queue]; 
        videoList.unshift({url:"https://www.youtube.com/watch?v=dQw4w9WgXcQ", "user": "Watch2Gether"})
        updateQueue(videoList)
    }

    return (
        <Modal
            title="Room Settings"
            visible={isModalVisible}
            onCancel={cancelForm}
            footer={[
                <Button key="back" onClick={cancelForm}>
                  Cancel
                </Button>,
                <Button key="submit" type="primary" onClick={submitForm}>
                  Submit
                </Button>,
              ]}
        >
        <Form {...layout} name="basic">
            <Form.Item
                label="Enable Player Sink To Me Button"
                name="controls"
            >
                <Switch checked={controls} onChange={()=>(setControls(!controls))} />
            </Form.Item>
            <Form.Item
                label="Auto skip when everyone finishes"
                name="auto skip"
            >
                <Switch checked={auto_skip} onChange={()=>(setAutoSkip(!auto_skip))} />
            </Form.Item>
            <Form.Item
                label="Dark Mode"
                name="controls"
            >
                <Switch defaultChecked checked={true} onClick={handleDarkMode}/>
            </Form.Item>
            </Form>
      </Modal>
    )
}
)

const Settings = () => {

    const [isSettingsModalVisible, setIsSettingModalVisible] = useState(false);
    const showSettingsModal = () => {setIsSettingModalVisible(true);};
    const handleSettingsOk = () => {setIsSettingModalVisible(false); };
    const handleSettingsCancel = () => {setIsSettingModalVisible(false);};
  
    return (
      <div style={{"width":"100%"}}>
        <Button style={{"width":"100%"}} type="primary" onClick={() => setIsSettingModalVisible(true) } icon={<SettingOutlined />} key="1">Player Settings</Button>
        <SettingsModal isModalVisible={isSettingsModalVisible} showModal={showSettingsModal} handleOk={handleSettingsOk} handleCancel={handleSettingsCancel}/>
      </div>
    )
  }
export default Settings