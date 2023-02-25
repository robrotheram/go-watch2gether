import React, { useContext, useState } from 'react';
import {
  Drawer, Button, Col, Row,
} from 'antd';
import {
  TeamOutlined,
  SyncOutlined,
} from '@ant-design/icons';
import UserList from './UserList';

import useDeviceDetect from '../../common/useDeviceDetect';
import { UserContext } from '../../../context/UserContext';
import { RoomContext } from '../../../context/RoomContext';

const DrawerForm = () => {
  const isMoble = useDeviceDetect();
  const drawerWidth = isMoble ? '100%' : '400px';
  const [room, { forceSinkToMe, seekToHost }] = useContext(RoomContext)
  const [user] = useContext(UserContext)
  const controls = room.controls;
  const isHost = () => {
    return user.id === user.host;
  }
  const [visible, setVisible] = useState(false);

  const showDrawer = () => {
    setVisible(true);
  };

  const onClose = () => {
    setVisible(false);
  };

  return (
    <>
      <Button type="primary" onClick={showDrawer} style={{ height: '33px', margin: '0px 0px 0px 5px' }}>
        <TeamOutlined />
        {!isMoble ? "Watchers":null}
      </Button>
      <Drawer
        title="Watchers Progress"
        width={drawerWidth}
        onClose={onClose}
        open={visible}
        maskClosable
        bodyStyle={{ padding: 0 }}
      >

        <Row>
          { !isHost()
            ? (
              <Col flex="auto" style={{ padding: '5px 5px' }}>
                <Button
                  style={{ width: '100%' }}
                  type="primary"
                  icon={<SyncOutlined />}
                  key="3"
                  onClick={() => seekToHost()}
                >
                  Sync to host
                </Button>
              </Col>
            )
            : null}
          { controls || isHost
            ? (
              <Col flex="auto" style={{ padding: '5px 5px' }}>
                <Button
                  style={{ width: '100%' }}
                  type="primary"
                  icon={<SyncOutlined />}
                  key="2"
                  onClick={() => forceSinkToMe()}
                >
                  Sync everyone to me
                </Button>
              </Col>
            )
            : null}
        </Row>
        {/* <pre style={{color:"white"}}>{JSON.stringify(room, null, 2)}</pre> */}
        <UserList />
      </Drawer>
    </>
  );
};

export default DrawerForm;
