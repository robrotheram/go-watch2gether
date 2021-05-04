import React from "react"
import { Layout, Menu, Breadcrumb } from 'antd';
import {
  DesktopOutlined,
  PieChartOutlined,
  FileOutlined,
  TeamOutlined,
  AppstoreOutlined,
} from '@ant-design/icons';
import PlayerApp from "../PlayerApp";
import {
  Switch,
  Route,
  Redirect
} from "react-router-dom";

import { PageFooter } from "../common/PageFooter";
import  Navigation  from "../common/Nav";
import "./index.less"
import logo from "../WelcomePage/logo.jpg"
import {connect} from 'react-redux'
import {history} from "../../store/index"
import GuildIcon from "../common/icons/GuildIcon"
import UserIcon from "../common/icons/UserIcon"

import { LogoutOutlined } from '@ant-design/icons';
import {join} from "../../store/room/room.actions"

const { Sider } = Layout;
const { SubMenu } = Menu;

class RoomPage extends React.Component {
  state = {
    collapsed: true,
  };

  onCollapse = collapsed => {
    console.log(collapsed);
    this.setState({ collapsed });
  };

  handleClick = (guild) => {
    //history.push("/app/room/"+id)
    const {username} = this.props
    this.props.join(guild.id, guild.name, username, false )
  }

  render() {
    const { collapsed } = this.state;
    const {guilds} =this.props
    return (
      <Layout style={{ minHeight: '100vh' }}>
        <Sider collapsed={collapsed} onCollapse={this.onCollapse} >
          <div className="logo">
            <img src={logo} width="100%"/>
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
                    onClick={() => this.handleClick(guild)}
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
              icon={<UserIcon user={this.props}/>} 
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
}

const mapStateToProps  = (state) =>{
  return state.user
} 
export default connect(mapStateToProps, {join})(RoomPage)