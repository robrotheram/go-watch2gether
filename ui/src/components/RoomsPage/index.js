import React, { useEffect, useState } from "react"
import { Layout, Menu } from 'antd';

import PlayerApp from "../PlayerApp";
import {
  Switch,
  Route
} from "react-router-dom";

import { PageFooter } from "../common/PageFooter";

import "./index.less"
import logo from "../WelcomePage/logo.jpg"
import {connect} from 'react-redux'
import {history} from "../../store/index"
import GuildIcon from "../common/icons/GuildIcon"
import UserIcon from "../common/icons/UserIcon"

import { LogoutOutlined } from '@ant-design/icons';
import {join} from "../../store/room/room.actions"

const { Sider } = Layout;

const RoomPage = ({guilds, icon, username, id, join}) => {

  const [collapsed, setCollapsed] = useState(true)

  const onCollapse = c => {
    console.log(c);
    setCollapsed(c)
  };

  const handleClick = (guild) => {
    //active
    join(guild.id, guild.name, username, false )
  }

  useEffect(()=>{
    if(id === ""){
      history.push("/app")
    }
  },[id])

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider collapsed={collapsed} onCollapse={onCollapse} >
        <div className="logo">
          <img src={logo} alt="logo" width="100%"/>
        </div>
        <Menu 
          theme="dark" 
          defaultSelectedKeys={['1']} 
          mode="inline" 
          className="guildList"
        >
          {
            guilds.map(guild => {
              return(
                <Menu.Item 
                  key={guild.id} 
                  className="guildMenu" 
                  icon={<GuildIcon guild={guild}/>}
                  onClick={() => handleClick(guild)}
                  >
                  {guild.name}
                  
                </Menu.Item>
              )
            })
          }
        </Menu>
        <Menu theme="dark" defaultSelectedKeys={['1']}  defaultOpenKeys={['SubMenu']} mode="inline">
          <Menu.SubMenu 
            key="sub1" 
            icon={<UserIcon username={username} icon={icon}/>} 
            popupOffset={[0,-3]}
          >
          <Menu.Item>
            <a  href={"/auth/logout"} rel="noopener noreferrer">
            <LogoutOutlined /> Logout
            </a>
          </Menu.Item>
          </Menu.SubMenu>
        </Menu>
      </Sider>
      <Layout className="site-layout">
        <Switch>
        <Route path="/app/room/:id">
        <PlayerApp/>
        </Route>
        </Switch>
        <PageFooter style={{ textAlign: 'center', position: "fixed", bottom:"0px", left:"85px", right:"0px", height:"50px", padding: "15px 50px 28px 50px" }}/>
      </Layout>
    </Layout>
  );
}


const mapStateToProps  = (state) =>{
  return {
    ...state.user,
    id: state.room.id
  } 
} 
export default connect(mapStateToProps, {join})(RoomPage)