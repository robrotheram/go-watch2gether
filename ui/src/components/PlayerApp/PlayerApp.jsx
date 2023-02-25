import './App.less';
import React, { useContext, useEffect, useRef } from "react"
import { Alert, Layout, Spin } from 'antd';
import { VideoControls, VideoList } from './components/VideoQueue';
import VideoPlayer from './components/VideoPlayer'
import Controls from './components/Controls';
import { RoomContext } from '../../context/RoomContext';
import { UserContext } from '../../context/UserContext';

export const PlayerApp = () => {
    const [room, {updateUser}] = useContext(RoomContext)
    const [user] = useContext(UserContext)
    const countRef = useRef({});
    countRef.current = {...user, seek: room.current_seek};

    useEffect(() => {
      let timeout = setInterval(() => {
        updateUser(countRef.current)
      }, 2000);
      return () => clearInterval(timeout)
    }, [])

    if (!room.active){
      <Spin tip="Loading...">
        <Alert
          message="Room Loading"
          description="Please wait as we load the room"
          type="info"
        />
      </Spin>
    }

    return (
      <Layout className="dark-theme">
        <Controls />
        <div className="contentWrapper">
          <div className="queueWrapper">
            <VideoControls/>
            <VideoList />
          </div>
          <div className="playerWrapper">
            <VideoPlayer />
          </div>
        </div>
      </Layout>
    );
  }