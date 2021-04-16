import './App.less';
import React from "react"
import { Layout, Row, Col, Divider } from 'antd';
import Navigation from '../common/Nav'
import {PageFooter} from '../common/PageFooter'

import { VideoControls, VideoList } from './components/VideoQueue';
import UserList from './components/UserList'
import VideoPlayer from './components/VideoPlayer'

import {connect} from 'react-redux'
import {join, leave, isAlive, reJoin} from '../../store/room/room.actions'
import {history} from '../../store'
import { withRouter } from "react-router";

import DrawerForm from './components/UserDrawer'
import Controls from './components/Controls';

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
      <Controls/>
      <div style={{
        bottom: "60px",
        "width": "100%", 
        "overflow": "hidden",
        padding:"0px 10px 0px 10px", 
        display:"flex",
        height: "calc(100vh - 130px)"
      }}>  

      <div style={{"height":"100%", "width":"400px" , "padding":"5px"}}>
        <VideoControls/>
         <VideoList/>
      </div>
      <div style={{"height":"100%", "width":"100%" , "padding":"10px", "marginTop":"5px"}}>
      <VideoPlayer/>
      </div>
    

      

      </div>
    </Layout>
    );
  }
}

const mapStateToProps  = (state) =>{
  
  return state.room
} 
export default withRouter(connect(mapStateToProps, {join, leave, isAlive, reJoin})(App))

