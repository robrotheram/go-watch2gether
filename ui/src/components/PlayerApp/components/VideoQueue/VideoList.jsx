import React, { useContext } from "react"
import {
  Card, List, Space, Button, Typography,
} from 'antd';
import {
  VerticalAlignTopOutlined, ArrowUpOutlined, ArrowDownOutlined, DeleteOutlined,
} from '@ant-design/icons';

import { VideoItem } from './VideoItem';
import { RoomContext } from "../../../../context/RoomContext";

const { Paragraph, Text } = Typography;

export const VideoList = () => {
  const [room, {updateQueue}] = useContext(RoomContext)
  const queue = room.queue
  const video = room.current_video

  const deleteVideo = (item) => {
    const videoList = [...queue];
    const i = videoList.indexOf(item);
    videoList.splice(i, 1);
    updateQueue(videoList);
  };

  const voteUp = (item) => {
    const videoList = [...queue];
    const i = videoList.indexOf(item);
    const z = videoList[i - 1];
    videoList[i - 1] = videoList[i];
    videoList[i] = z;
    updateQueue(videoList);
  };

  const moveToTop = (item) => {
    let videoList = [...queue];
    const i = videoList.indexOf(item);
    videoList.splice(i, 1);
    videoList = [item, ...videoList];
    updateQueue(videoList);
  };

  const voteDown = (item) => {
    const videoList = [...queue];
    const i = videoList.indexOf(item);
    const z = videoList[i + 1];
    videoList[i + 1] = videoList[i];
    videoList[i] = z;
    updateQueue(videoList);
  };

  return (
    <div style={{ height: '100%' }}>
      
      {video !== undefined && video.url !== ''
        ? <Card style={{margin:"10px 0px"}}><VideoItem video={video} playing /></Card>
        : null}
      <Card
        type="inner"
        className="list video"
        title={(
          <Paragraph style={{marginBottom:"0px"}}>
            There are 
            {' '}
            <Text strong>{queue.length}</Text>
            {' '}
            videos in the queue
          </Paragraph>
)}
        style={
            video !== undefined && video.url !== ''
              ? { height: 'calc( 100% - 170px )' }
              : { height: 'calc( 100% - 55px )', marginTop: '10px' }
}
      >
        <div className="videoQueue">
          <List
            size="small"
            itemLayout="horizontal"
            dataSource={queue}
            renderItem={(item) => (
              <VideoItem video={item} loading={item.loading}>
                <Space style={{ float: 'right' }}>
                  { queue.indexOf(item) !== 0 ? <Button onClick={() => moveToTop(item)} icon={<VerticalAlignTopOutlined />} /> : null}
                  { queue.indexOf(item) !== 0 ? <Button onClick={() => voteUp(item)} icon={<ArrowUpOutlined />} /> : null}
                  { queue.indexOf(item) !== queue.length - 1 ? <Button onClick={() => voteDown(item)} icon={<ArrowDownOutlined />} /> : null}
                  <Button icon={<DeleteOutlined onClick={() => deleteVideo(item)} />} />
                </Space>
              </VideoItem>
            )}
          />
        </div>
      </Card>
    </div>
  );
};
