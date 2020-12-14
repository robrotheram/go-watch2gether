import './App.less';
import React from "react"
import { Layout, Row, Col, Divider } from 'antd';
import {Navigation} from './components/Nav'
import Pageheader from './components/pageheader'
import {PageFooter} from './components/PageFooter'

import { VideoControls, VideoList } from './components/VideoQueue';
import UserList from './components/UserList'
import VideoPlayer from './components/VideoPlayer'

import {connect} from 'react-redux'
import {join, leave, isAlive, reJoin} from './store/room/room.actions'
import {history} from './store'
import { withRouter } from "react-router";
import { Typography } from 'antd';

const { Title, Link } = Typography;


const { Content } = Layout;



class App extends React.Component {
  state = {
    timer: null,
  };

  update = () =>{
    try {
      isAlive();
    }catch(e){
      console.log("APP_ERROR", e)
    }
  }
  componentDidMount() {
    window.onbeforeunload = (function(){leave()})
    this.startApp();
  }

  componentWillUnmount() {
    clearInterval(this.state.timer);
  }

  startTimer = () => {
    console.log("APP", "Starting Watcher")
    let timer = setInterval(this.update, 1000);
    this.setState({timer});
  }

  startApp = () => {
    if (this.props.active){
      this.startTimer();
      return;
    }else{
      if (!this.props.active && this.props.user !== "" && this.props.name !== "" && this.props.name === this.props.match.params.id){
        this.props.reJoin(this.props.name)
        this.startTimer();
        return
      }
      const id = this.props.match.params.id;
      history.push("/?room="+id)
    }
  }

  render(){
    return (
      <Layout className="dark-theme">
        <Navigation/>
        
      <Content style={{ padding: '78px 0px', "width":"1550px",  "margin": "0 auto"}}>
      <Pageheader/>
      <Divider/>
        <Row gutter={[16, 16]}>
          <Col span={18} push={6}>
           
          {/* {this.props.current_video.url !== undefined ? 
            <Typography>
              <Title level={3}>Currently Playing: <Link href={this.props.current_video.url} target="_blank">
              {this.props.current_video.url}
              </Link>
              </Title>
            </Typography>
          : null } */}
            <VideoPlayer/>
            <Divider/>
            <UserList/>
          </Col>
          <Col span={6} pull={18}>
              <VideoControls/>
              <VideoList/>
          </Col>
        </Row>
      </Content>
      <PageFooter/>
    </Layout>
    );
  }
}

const mapStateToProps  = (state) =>{
  
  return state.room
} 
export default withRouter(connect(mapStateToProps, {join, leave, isAlive, reJoin})(App))

