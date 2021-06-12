import React, { useEffect } from 'react';
import { Layout, Space } from 'antd';

import { connect } from 'react-redux';
import { leave, sinkToME } from '../../../store/room/room.actions';
import UserMenu from './UserMenu';

const { Header } = Layout;

function Navigation(props) {
  const { host, controls, name } = props.room;
  const { isHost } = props.user;
  const { title } = props.video;

  useEffect(() => {
    if (title === '') {
      document.title = 'Watch2gether';
    } else {
      document.title = `Watch2gether | Playing:${title}`;
    }
  }, [title]);
  console.log(props);
  return (
    <Header style={{
      display: 'block ruby', background: '#1f1f1f', zIndex: '1000', left: 0, right: 0, top: 0, padding: '0px 20px 0px 20px',
    }}
    >
      {props.video !== undefined && props.video.title !== '' ? (
        <h1 style={{ color: 'white' }}>
          {' '}
          Current Playing:
          { props.video.title}
        </h1>
      ) : null}
      <Space style={{ float: 'right' }}>
        <UserMenu />
      </Space>
    </Header>
  );
}

const mapStateToProps = (state) => state;

export default connect(mapStateToProps, { leave })(Navigation);
