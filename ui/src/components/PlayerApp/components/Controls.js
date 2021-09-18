import React, { useState, useEffect } from 'react';
import {
  Button, Space, Card, Input, Row, Col,
} from 'antd';

import { useDispatch, useSelector } from 'react-redux';
import { VideoCameraOutlined } from '@ant-design/icons';
import DrawerForm from './UserDrawer';
import { openNotificationWithIcon } from '../../common/notification';
import PlaylistDrawer from './playlists/PlaylistDrawer';
import { addVideosToQueue, validURL } from '../../../store/video';
import Share from './ShareModal';
import Settings from './SettingsModal';

const Controls = () => {
  const dispatch = useDispatch();
  const user = useSelector((state) => state.user);
  const room_id = useSelector((state) => state.room.id);
  const title = useSelector((state) => state.video.title);
  const [newurl, setURL] = useState('');

  useEffect(() => {
    if (title === '') {
      document.title = 'Watch2gether';
    } else {
      document.title = `Watch2gether | Playing:${title}`;
    }
  }, [title]);

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      dispatch(addToQueue());
    }
  };

  const addToQueue = async () => {
    if (validURL(newurl)) {
      addVideosToQueue(room_id, newurl)
      setURL('');
    } else {
      openNotificationWithIcon('error', 'Invalid URL');
    }
  };

  return (
    <Card className="contolPanel">
      <Row style={{ width: '100%', paddingTop: '10px' }}>
        <Col>
          <Space style={{ marginTop: '1px' }}>
            <PlaylistDrawer />
          </Space>
        </Col>
        <Col flex="auto">
          <Input className="videoInput" defaultValue="mysite" value={newurl} onChange={(e) => setURL(e.target.value)} onKeyDown={handleKeyDown} addonAfter={(<Button type="primary" onClick={addToQueue} icon={<VideoCameraOutlined />}>Add Video</Button>)} />
        </Col>
        <Col>
          <Space style={{ marginTop: '1px' }} size={4}>
            <DrawerForm />
            {user.isHost ? <Settings /> : null}
            <Share />
          </Space>
        </Col>
      </Row>
    </Card>
  );
};

export default Controls;
