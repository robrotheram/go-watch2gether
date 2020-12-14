import React, {useState} from 'react';
import { Modal } from 'antd';
import { Form, Button } from 'antd';
import { Switch } from 'antd';
import {connect} from 'react-redux'
import {updateControls, updateQueue} from '../store/room/room.actions'
import {openNotificationWithIconKey} from "./notification"

const layout = {
    labelCol: { span: 12 },
    wrapperCol: { span: 12 },
};

function SettingsModal (props) {
    const {isModalVisible, handleOk, handleCancel, queue } = props
    const [controls, setControls] = useState(props.controls);

    const submitForm = () =>{
        updateControls(controls)
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
                label="Dark Mode"
                name="controls"
            >
                <Switch defaultChecked checked={true} onClick={handleDarkMode}/>
            </Form.Item>
            </Form>
      </Modal>
    )
}
const mapStateToProps  = (state) =>{
    return state.room
  } 
export default connect(mapStateToProps, {updateControls})(SettingsModal)
  
