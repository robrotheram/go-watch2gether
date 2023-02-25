import React, { useContext} from 'react';
import { Layout, Menu, Empty } from 'antd';
import {
  json,
  Outlet,
  useNavigate,
} from 'react-router-dom';
import { LogoutOutlined } from '@ant-design/icons';
import { PageFooter } from '../common/PageFooter';
import './index.less';
import logo from '../WelcomePage/logo.jpg';
import GuildIcon from '../common/icons/GuildIcon';
import UserIcon from '../common/icons/UserIcon';
import { UserContext } from '../../context/UserContext';
import { RoomContextProvider } from '../../context/RoomContext';


const { Sider } = Layout;

const RoomPage = () => {
  const navigate = useNavigate();
  const [user, loading] = useContext(UserContext);
  const logout = () => {
    window.location.replace('/auth/logout');
  };
  if (user === null) {
    return(<Layout style={{ minHeight: '100vh' }}/>)
  }
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider collapsed={true}>
        <div className="logo">
          <img src={logo} alt="logo" width="100%" />
        </div>
        <Menu
          theme="dark"
          defaultSelectedKeys={['1']}
          mode="inline"
          className="guildList"
        >
          {
            user.guilds && user.guilds.map((guild) => (
              <Menu.Item
                key={guild.id}
                className="guildMenu"
                icon={<GuildIcon guild={guild} />}
                onClick={() => navigate(`/app/room/${guild.id}`)}
              >
                {guild.name}

              </Menu.Item>
            ))
          }
        </Menu>
        <Menu theme="dark" defaultSelectedKeys={['1']} defaultOpenKeys={['SubMenu']} mode="inline">
          <Menu.SubMenu
            key="sub1"
            icon={<UserIcon username={user.username} icon={user.avatar_icon} />}
            popupOffset={[0, -3]}
          >
            <Menu.Item style={{background:"#141414"}} onClick={()=> logout()}>
                <LogoutOutlined />
                {' '}
                Logout
            </Menu.Item>
          </Menu.SubMenu>
        </Menu>
      </Sider>
      <Layout className="site-layout">
      <RoomContextProvider>
        <Outlet />
      </RoomContextProvider>
        <PageFooter style={{
          textAlign: 'center', position: 'fixed', bottom: '0px', left: '85px', right: '0px', height: '50px', padding: '15px 50px 28px 50px',
        }}
        />
      </Layout>
    </Layout>
  );
};
export default RoomPage
