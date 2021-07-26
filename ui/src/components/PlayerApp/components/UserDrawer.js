import React, { useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import {
  Drawer, Button, Col, Row,
} from 'antd';
import {
  TeamOutlined,
  SyncOutlined,
} from '@ant-design/icons';
import UserList from './UserList';

import { forceSinkToMe } from '../../../store/room/room.actions';
import { seekToHost } from '../../../store/video/video.actions';
import useDeviceDetect from '../../common/useDeviceDetect';

const DrawerForm = () => {
  const isMoble = useDeviceDetect();
  const drawerWidth = isMoble ? '100%' : '540px';

  const dispatch = useDispatch();
  const controls = useSelector((state) => state.room.controls);
  const isHost = useSelector((state) => state.user.isHost);

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
        visible={visible}
        maskClosable
        bodyStyle={{ padding: 0 }}
      >

        <Row>
          { !isHost
            ? (
              <Col flex="auto" style={{ padding: '5px 5px' }}>
                <Button
                  style={{ width: '100%' }}
                  type="primary"
                  icon={<SyncOutlined />}
                  key="3"
                  onClick={() => dispatch(seekToHost())}
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
                  onClick={() => dispatch(forceSinkToMe())}
                >
                  Sync everyone to me
                </Button>
              </Col>
            )
            : null}
        </Row>
        <UserList />
      </Drawer>
    </>
  );
};

export default DrawerForm;
