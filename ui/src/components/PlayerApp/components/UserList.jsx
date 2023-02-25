import React, { useContext } from 'react';
import {
  List, Button, Progress, Row, Col,
} from 'antd';

import { SyncOutlined } from '@ant-design/icons';
import { RoomContext } from '../../../context/RoomContext';


function UserList() {
  const [room, {seekToUser}] = useContext(RoomContext);

  const watchers = room.watchers;
  const host = room.host;

  const secondsToDate = (seconds) => {
    const time = new Date(seconds * 1000).toISOString().substr(11, 8);
    const res = time.substring(0, 2);
    if (res === '00') {
      return time.substring(3, time.length);
    }
    return time;
  };
  // console.log("room",props.room.watchers)
  return (
  // <Card type="inner" title="Users Progress" className="list">
  //   <div className="container .sc2 userlist">
  /* {JSON.stringify(watchers)} */
    <List
      size="small"
      itemLayout="horizontal"
      dataSource={watchers}
      renderItem={(item) => (
        <List.Item className={item.id === host ? 'userListActive' : null}>
          <Row style={{ width: '100%', padding: '5px' }}>
            <Col flex="100px" style={{ textAlign: 'left', paddingRight: '10px' }}>
              {item.username}
            </Col>
            <Col flex="auto">
              <div style={{ display: 'inline-block', width: '90%', paddingRight: '20px' }}>
                <Progress percent={(item.seek.progress_percent) * 100} showInfo={false} size="small" />
              </div>
            </Col>
            <Col>
              {secondsToDate(item.seek.progress_seconds)}
              <Button icon={<SyncOutlined />} onClick={() => { seekToUser(item.seek); }} style={{ marginLeft: '10px' }} />
            </Col>
          </Row>
        </List.Item>
      )}
    />
  );
}

export default UserList;
