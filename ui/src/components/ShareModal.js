import React, {useState} from 'react';
import { Modal } from 'antd';
import { Input, Button } from 'antd';
import { Switch } from 'antd';
import {connect} from 'react-redux'
import {updateControls, updateQueue} from '../store/room/room.actions'
import {openNotificationWithIconKey} from "./notification"
import {
    CopyOutlined
  } from '@ant-design/icons';

const layout = {
    labelCol: { span: 12 },
    wrapperCol: { span: 12 },
};

function ShareModal (props) {
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
            title="Share Room"
            visible={isModalVisible}
            onCancel={cancelForm}
            footer={[
                <Button key="back" onClick={cancelForm}>
                  Cancel
                </Button>
              ]}
        >
             <Input.Group compact>
      <Input disabled style={{ width: '90%' }} defaultValue={window.location.href} />
      <Button style={{ width: '10%' }} type="primary" icon={<CopyOutlined />}  onClick={() => {
          navigator.clipboard.writeText(window.location.href)
          openNotificationWithIconKey("success", "Link coppied")
          
          }}/>
    </Input.Group>
      </Modal>
    )
}
const mapStateToProps  = (state) =>{
    return state.room
  } 
export default connect(mapStateToProps, {updateControls})(ShareModal)
  
