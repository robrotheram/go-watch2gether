import React from "react"
import {
  Card, List, Space, Button, Typography,
} from 'antd';
import {
  VerticalAlignTopOutlined, ArrowUpOutlined, ArrowDownOutlined, DeleteOutlined,
} from '@ant-design/icons';

import { useDispatch, useSelector } from 'react-redux';
import { updateQueue } from '../../../../store/room/room.actions';
import { VideoItem } from './VideoItem';

const { Paragraph, Text } = Typography;

export const VideoListComponent = () => {
  const dispatch = useDispatch();
  const queue = useSelector((state) => state.room.queue);
  const current_video = useSelector((state) => state.video);

  const deleteVideo = (item) => {
    const videoList = [...queue];
    const i = videoList.indexOf(item);
    videoList.splice(i, 1);
    dispatch(updateQueue(videoList));
  };

  const voteUp = (item) => {
    const videoList = [...queue];
    const i = videoList.indexOf(item);
    const z = videoList[i - 1];
    videoList[i - 1] = videoList[i];
    videoList[i] = z;
    dispatch(updateQueue(videoList));
  };

  const moveToTop = (item) => {
    let videoList = [...queue];
    const i = videoList.indexOf(item);
    videoList.splice(i, 1);
    videoList = [item, ...videoList];
    dispatch(updateQueue(videoList));
  };

  const voteDown = (item) => {
    const videoList = [...queue];
    const i = videoList.indexOf(item);
    const z = videoList[i + 1];
    videoList[i + 1] = videoList[i];
    videoList[i] = z;
    dispatch(updateQueue(videoList));
  };

  return (
    <div style={{ height: '100%' }}>
      {current_video !== undefined && current_video.url !== ''
        ? <VideoItem video={current_video} playing />
        : null}
      <Card
        type="inner"
        className="list video"
        title={(
          <Paragraph>
            There are 
            {' '}
            <Text strong>{queue.length}</Text>
            {' '}
            videos in the queue
          </Paragraph>
)}
        style={
            current_video !== undefined && current_video.url !== ''
              ? { height: 'calc( 100% - 150px )' }
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

export const VideoList = VideoListComponent;
