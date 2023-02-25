import React, { useContext, useEffect, useRef, useState } from 'react';
import ReactPlayer from 'react-player';
import { Empty, Card, Spin, Alert } from 'antd';
import { UserContext } from '../../../context/UserContext';
import { RoomContext } from '../../../context/RoomContext';

const Player = () => {
  const inputRange = useRef(null);
  const [room, {updateSeek, handleFinish, pause, play}] = useContext(RoomContext)
  const [user] = useContext(UserContext)
  const queue = room.queue
  const video = room.current_video
  const playing = room.playing;
  const url = video.url;
  const seek = room.seek_to_user;
  const [seeking, setSeeking] = useState(false);

  useEffect(() => {
    setSeeking(true);
    try {
      console.log('SeekTO', seek);
      inputRange.current.seekTo(parseFloat(seek.progress_seconds));
    } catch {
        return
    }
    setSeeking(false);
  }, [seek]);

  useEffect(() => {
    setSeeking(false);
  }, [url]);

  const handleProgress = (state) => {
    if (!seeking) {
      updateSeek(state.played, state.playedSeconds);
    }
  };
  const handleEnded = () => {
    console.log('UPDATE_SEEK');
    setSeeking(true);
    handleFinish();
  };

  return (
    <Card className="videoPlayer" style={{ height: '100%', width: '100%' }}>
      {/* <pre style={{color:"white"}}>{JSON.stringify(seek, null, 2)}</pre> */}
      {url !== '' ? (
        <ReactPlayer
          ref={inputRange}
          width="100%"
          height="100%"

          playing={playing}
          url={url}
          controls
          onPause={() => { pause(); }}
          onPlay={() => { play(); }}
          onProgress={handleProgress}
          onEnded={handleEnded}
        />
      )
        : <Empty style={{ height: '100%', width: '100%', paddingTop: '100px' }} />}
    </Card>
  );
};

export default Player;
