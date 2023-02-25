import React, { useState, useEffect, useContext } from 'react';
import {
  Button, Space, Card, Input, Row, Col,
} from 'antd';

import { VideoCameraOutlined } from '@ant-design/icons';
import DrawerForm from './UserDrawer';
import { openNotificationWithIcon } from '../../common/notification';
import Share from './ShareModal';
import { Settings } from './SettingsModal';
import { addVideosToQueue, validURL } from '../../../context';
import { UserContext } from '../../../context/UserContext';
import { RoomContext } from '../../../context/RoomContext';

const Controls = () => {
  const [user] = useContext(UserContext);
  const [room] = useContext(RoomContext);
  const title = ""
  
  const [newurl, setURL] = useState('');
  const isHost = () => {
    return user.id === room.host;
  }

  useEffect(() => {
    if (title === '') {
      document.title = 'Watch2gether';
    } else {
      document.title = `Watch2gether | Playing:${title}`;
    }
  }, [title]);

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      addToQueue();
    }
  };

  const addToQueue = async () => {
    if (validURL(newurl)) {
      addVideosToQueue(room.id, newurl)
      setURL('');
    } else {
      openNotificationWithIcon('error', 'Invalid URL');
    }
  };

  return (
    <Card className="contolPanel">
      <Row style={{ width: '100%' }}>
        <Col>
          <Space style={{ marginTop: '1px' }}>
            {/* <PlaylistDrawer /> */}
          </Space>
        </Col>
        <Col flex="auto">
        <Input.Group className="videoInput">
      <Input style={{ width: 'calc(100% - 120px)' }} value={newurl} onChange={(e) => setURL(e.target.value)} onKeyDown={handleKeyDown} />
      <Button type="primary" onClick={addToQueue} icon={<VideoCameraOutlined />} style={{borderRadius:"0px 5px 5px 0px"}}>Add Video</Button>
    </Input.Group>
        </Col>
        <Col>
          <Space style={{ marginTop: '1px' }} size={4}>
            <DrawerForm />
            {isHost() ? <Settings /> : null} 
            <Share />
          </Space>
        </Col>
      </Row>
    </Card>
  );
};

export default Controls;
