import React, {useState} from 'react';
import { Modal } from 'antd';
import { Form, Button } from 'antd';
import { Switch } from 'antd';
import {connect} from 'react-redux'
import {updateControls} from '../store/room/room.actions'

const layout = {
    labelCol: { span: 8 },
    wrapperCol: { span: 16 },
};

function SettingsModal (props) {
    const {isModalVisible, handleOk, handleCancel } = props
    const [controls, setControls] = useState(props.controls);

    const submitForm = () =>{
        updateControls(controls)
        handleOk();

    }
    const cancelForm = () =>{
        setControls(props.controls)
        handleCancel();

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
                label="Enable Player Controls"
                name="controls"
            >
                <Switch checked={controls} onChange={()=>(setControls(!controls))} />
            </Form.Item>
            <Form.Item
                label="Dark Mode"
                name="controls"
            >
                <Switch defaultChecked checked={true}/>
            </Form.Item>
            </Form>
      </Modal>
    )
}
const mapStateToProps  = (state) =>{
    return state.room
  } 
export default connect(mapStateToProps, {updateControls})(SettingsModal)
  
