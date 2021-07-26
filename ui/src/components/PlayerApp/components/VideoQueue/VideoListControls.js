import React from "react"
import { Button } from 'antd';
import { useDispatch, useSelector } from 'react-redux';
import './VideoListControls.less';
import { updateQueue, nextVideo } from '../../../../store/room/room.actions';
import useDeviceDetect from '../../../common/useDeviceDetect';

const shuffleArray = (array) => {
  for (let i = array.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    const temp = array[i];
    array[i] = array[j];
    array[j] = temp;
  }
};

export function VideoControlComponent() {
  const dispatch = useDispatch();
  const queue = useSelector((state) => state.room.queue);
  const mobile = useDeviceDetect()

  const clearQueue = (e) => {
    e.currentTarget.blur() 
    dispatch(updateQueue([]));
  };

  const skipQueue = (e) => {
    e.currentTarget.blur() 
    dispatch(nextVideo());
  };

  const shuffleQueue = (e) => {
    e.currentTarget.blur() 
    const videoList = [...queue];
    shuffleArray(videoList);
    dispatch(updateQueue(videoList));
  };

  return (
    <div className="list-controls">
      <Button onClick={skipQueue} className="list-controls-item"> {queue.length > 0 ? "Skip" : "Play "} </Button>
      <Button onClick={clearQueue} className="list-controls-item">{ mobile ? "Clear" : "Clear Queue"}</Button>
      <Button onClick={shuffleQueue} className="list-controls-item">Shuffle</Button>
    </div>
  );
}

export const VideoControls = (VideoControlComponent);
