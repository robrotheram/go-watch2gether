import React, {useState}from 'react'
import {connect} from 'react-redux'
import { Drawer, Space, Button, Col, Row, Input, Select, DatePicker, Divider } from 'antd';
import { InfoCircleOutlined } from '@ant-design/icons';
import UserList from './UserList';
import Share from './ShareModal'
import Settings from './SettingsModal'
import {
  SyncOutlined
} from '@ant-design/icons';
import {leave, sinkToHost, sinkToME} from '../../../store/room/room.actions'


const { Option } = Select;

const DrawerForm = (props) => {

  const { host, controls, name } = props.room
  const {isHost} = props.user
  

  const [visible, setVisible] = useState(false);
  
  const showDrawer = () => {
    setVisible(true);
  };

  const onClose = () => {
    setVisible(false);
  };

    return (
      <>
        <Button type="primary" onClick={showDrawer}  style={{"height":"33px", "margin": "0px 10px"}}>
          <InfoCircleOutlined /> Info
        </Button>
        <Drawer
          title="More Info"
          width={460}
          onClose={onClose}
          visible={visible}
          maskClosable={true}
          bodyStyle={{ padding: 0 }}
        >

          <Space size="small" style={{width:"100%", padding:"10px", marginBottom:"-20px" }}>
            { !isHost ? <Button style={{"width":"100%"}} type="primary" icon={<SyncOutlined />} key="3" onClick={() => props.sinkToHost()}>Sync to host</Button> : null}
            { controls || isHost ? <Button style={{"width":"100%"}} type="primary" icon={<SyncOutlined />} key="2" onClick={() => sinkToME()}>Sync everyone to me</Button>: null}
            {isHost ?<Settings/>: null}
            <Share/>
          </Space>
          <Divider>User Progress</Divider>


          

            
       

        <UserList/>
        </Drawer>
      </>
    );
  
}


const mapStateToProps  = (state) =>{
  return state
} 
export default connect(mapStateToProps, {leave, sinkToHost, sinkToME })(DrawerForm)
