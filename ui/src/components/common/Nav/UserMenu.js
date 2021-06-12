import { Avatar, Menu } from 'antd';
import { UserOutlined, LogoutOutlined } from '@ant-design/icons';
import { connect } from 'react-redux';

import React, { useState } from 'react';
import { BASE_URL } from '../../../store';

const UserAvatar = (props) => {
  if (props.icon !== undefined && props.icon !== '') {
    return (<Avatar shape="square" size={30} src={props.icon} />);
  }
  return (
    <Avatar shape="square" size={30} icon={<UserOutlined />} />
  );
};

const UserMenu = (props) => {
  const [current, setCurrent] = useState('');

  const handleClick = (e) => {
    console.log('click ', e);
    setCurrent(e.key);
  };

  return (
    <Menu
      onClick={handleClick}
      selectedKeys={[current]}
      mode="horizontal"
      style={{ background: 'transparent', lineHeight: '60px' }}
    >
      <Menu.SubMenu key="SubMenu" icon={<UserAvatar icon={props.icon} />} style={{ margin: '0px', padding: '0px' }}>
        <Menu.Item>
          <a href={`${BASE_URL}/auth/logout`} rel="noopener noreferrer">
            <LogoutOutlined />
            {' '}
            Logout
          </a>
        </Menu.Item>
      </Menu.SubMenu>
    </Menu>
  );
};

const mapStateToProps = (state) => state.user;
export default connect(mapStateToProps)(UserMenu);
