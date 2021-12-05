import './App.less';
import React from "react"
import { Layout } from 'antd';

import { VideoControls, VideoList } from './components/VideoQueue';

import VideoPlayer from './components/VideoPlayer'

import { connect } from 'react-redux'
import { join, leave, isAlive, reJoin } from '../../store/room/room.actions'
import { history } from '../../store'
import { withRouter } from "react-router";
import Controls from './components/Controls';

class App extends React.Component {

  state = {
    timer: null,
  };

  update = () => {
    try {
      isAlive();
    } catch (e) {
      console.log("APP_ERROR", e)
    }
  }
  componentDidMount() {
    window.onbeforeunload = (function () { leave() })
    this.startApp();
  }

  componentWillUnmount() {
    clearInterval(this.state.timer);
  }

  startTimer = () => {
    console.log("APP", "Starting Watcher")
    let timer = setInterval(this.update, 2000);
    this.setState({ timer });
  }

  startApp = () => {
    if (this.props.active) {
      this.startTimer();
      return;
    } else {
      if (!this.props.active && this.props.user !== "" && this.props.name !== "" && this.props.name === this.props.match.params.id) {
        this.props.reJoin(this.props.name)
        this.startTimer();
        return
      }
      const id = this.props.match.params.id;
      history.push("/?room=" + id)
    }
  }

  render() {
    return (
      <Layout className="dark-theme">
        <Controls />
        <div className="contentWrapper">
          <div className="queueWrapper">
            <VideoControls />
            <VideoList />
          </div>
          <div className="playerWrapper">
            <VideoPlayer />
          </div>
        </div>
      </Layout>
    );
  }
}

const mapStateToProps = (state) => {

  return state.room
}
export default withRouter(connect(mapStateToProps, { join, leave, isAlive, reJoin })(App))

