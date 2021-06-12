import React from "react"
import { Button } from 'antd';
import { useDispatch, useSelector } from 'react-redux';
import './VideoListControls.less';
import { updateQueue, nextVideo } from '../../../../store/room/room.actions';

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

  const clearQueue = () => {
    dispatch(updateQueue([]));
  };

  const skipQueue = () => {
    dispatch(nextVideo());
  };

  const shuffleQueue = () => {
    const videoList = [...queue];
    shuffleArray(videoList);
    dispatch(updateQueue(videoList));
  };

  return (
    <div className="list-controls">
      <Button onClick={skipQueue} className="list-controls-item"> Skip </Button>
      <Button onClick={clearQueue} className="list-controls-item">Clear Queue</Button>
      <Button onClick={shuffleQueue} className="list-controls-item">Shuffle</Button>
    </div>
  );
}

export const VideoControls = (VideoControlComponent);
