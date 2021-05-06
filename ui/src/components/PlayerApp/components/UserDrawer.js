import React, {useState}from 'react'
import {connect} from 'react-redux'
import { Drawer, Space, Button, Col, Row, Input, Select, DatePicker, Divider } from 'antd';
import { TeamOutlined } from '@ant-design/icons';
import UserList from './UserList';
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
        <Button type="primary" onClick={showDrawer}  style={{"height":"33px", "margin": "0px 0px 0px 5px"}}>
          <TeamOutlined />Watchers
        </Button>
        <Drawer
          title="Watchers Progress"
          width={460}
          onClose={onClose}
          visible={visible}
          maskClosable={true}
          bodyStyle={{ padding: 0 }}
        >

            <Row>
            <Col  flex="auto" style={{padding:"5px 5px"}}>
              { !isHost ? <Button style={{"width":"100%"}} type="primary" icon={<SyncOutlined />} key="3" onClick={() => props.sinkToHost()}>Sync to host</Button> : null}
            </Col>
            <Col  flex="auto" style={{padding:"5px 5px"}}>
            { controls || isHost ? <Button style={{"width":"100%"}} type="primary" icon={<SyncOutlined />} key="2" onClick={() => sinkToME()}>Sync everyone to me</Button>: null}
            </Col>
            </Row> 
          


          

            
       

        <UserList/>
        </Drawer>
      </>
    );
  
}


const mapStateToProps  = (state) =>{
  return state
} 
export default connect(mapStateToProps, {leave, sinkToHost, sinkToME })(DrawerForm)
