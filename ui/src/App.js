import './App.less';
import { useEffect } from "react"
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

const { Content } = Layout;



function App(props) {
  
  const update = () =>{
    try {
      isAlive();
    }catch(e){
      console.log("APP", e)
    }
  }

  useEffect(() => {
    if (props.active){
      console.log("APP", "Starting Watcher")
      const intervalId = setInterval(() => update(), 1000);
      return () => clearInterval(intervalId);
    }else{
      if (!props.active && props.user !== "" && props.name !== ""){
        props.reJoin(props.name)
        console.log("APP", "Starting Watcher")
        const intervalId = setInterval(() => update(), 1000);
        return () => clearInterval(intervalId);
      }else{
        const id = props.match.params.id;
        history.push("/?room="+id)
      }      
    }
  }, [props, props.location.search]);

  useEffect(() => {
    window.onbeforeunload = (function(){leave()})
  });

  return (
    <Layout className="dark-theme">
      <Navigation/>
      
    <Content style={{ padding: '78px 0px', "width":"1550px",  "margin": "0 auto"}}>
    <Pageheader/>
    <Divider/>
      <Row gutter={[16, 16]}>
        <Col span={18} push={6}>
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

const mapStateToProps  = (state) =>{
  
  return state.room
} 
export default withRouter(connect(mapStateToProps, {join, leave, isAlive, reJoin})(App))

