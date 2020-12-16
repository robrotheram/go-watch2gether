import React from 'react';
import { Modal } from 'antd';
import { Input, Button } from 'antd';
import {connect} from 'react-redux'
import {updateControls} from '../store/room/room.actions'
import {openNotificationWithIconKey} from "./notification"
import {
    CopyOutlined
  } from '@ant-design/icons';


function ShareModal (props) {
    const {isModalVisible, handleCancel } = props

    const cancelForm = () =>{
        handleCancel();
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
  
