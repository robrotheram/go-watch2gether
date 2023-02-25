import React, { useContext } from "react"
import { Button } from 'antd';
import './VideoListControls.less';
import { RoomContext } from "../../../../context/RoomContext";

const shuffleArray = (array) => {
  for (let i = array.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    const temp = array[i];
    array[i] = array[j];
    array[j] = temp;
  }
};

export const VideoControls = () => {
  const [room, {nextVideo, updateQueue}] = useContext(RoomContext)
  const queue = room.queue;

  const clearQueue = (e) => {
    e.currentTarget.blur() 
    updateQueue([]);
    nextVideo();
  };

  const skipQueue = (e) => {
    e.currentTarget.blur() 
    nextVideo();
  };

  const shuffleQueue = (e) => {
    e.currentTarget.blur() 
    const videoList = [...queue];
    shuffleArray(videoList);
    updateQueue(videoList);
  };

  return (
    <div className="list-controls">
      <Button onClick={skipQueue} className="list-controls-item"> {queue.length > 0 ? "Skip" : "Play "} </Button>
      <Button onClick={shuffleQueue} className="list-controls-item">Shuffle</Button>
      <Button onClick={clearQueue} className="list-controls-item">Clear All</Button>
    </div>
  );
}
