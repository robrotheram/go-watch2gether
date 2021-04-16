import React, {useState} from 'react';
import { Modal } from 'antd';
import { Input, Button } from 'antd';
import {connect} from 'react-redux'
import {openNotificationWithIconKey} from "../../common/notification"
import {
    CopyOutlined,
    ShareAltOutlined
  } from '@ant-design/icons';


const mapStateToProps  = (state) =>{
  return state.room
} 
const ShareModal = connect(mapStateToProps, {})((props) => {
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
});

const Share = () => {

  const [isShareModalVisible, setIsShareModalVisible] = useState(false);
  const showShareModal = () => {setIsShareModalVisible(true);};
  const handleShareOk = () => {setIsShareModalVisible(false); };
  const handleShareCancel = () => {setIsShareModalVisible(false);};

  return (
    <div style={{"width":"100%"}}>
      <Button style={{"width":"100%"}} type="primary" icon={<ShareAltOutlined />} onClick={() => setIsShareModalVisible(true) }>Share </Button>
      <ShareModal isModalVisible={isShareModalVisible} showModal={showShareModal} handleOk={handleShareOk} handleCancel={handleShareCancel}/>
    </div>
  )
}
export default Share