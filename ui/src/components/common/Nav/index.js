import React, { useEffect } from 'react';
import { Layout, Space } from 'antd';
import UserMenu from './UserMenu';

const { Header } = Layout;

function Navigation() {
  
  const title  = useSelector(state => state.video.title);
  const dispatch = useDispatch()

  useEffect(() => {
    if (title === '') {
      document.title = 'Watch2gether';
    } else {
      document.title = `Watch2gether | Playing:${title}`;
    }
  }, [title]);
  
  return (
    <Header style={{
      display: 'block ruby', background: '#1f1f1f', zIndex: '1000', left: 0, right: 0, top: 0, padding: '0px 20px 0px 20px',
    }}
    >
      {title !== '' ? (
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
export default (Navigation);
